package main

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"

	"go.uber.org/zap"

	"faraway/internal/service"
	"faraway/internal/transport"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	var (
		host = os.Getenv("SERVER_HOST")
		port = os.Getenv("SERVER_PORT")
	)

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "8181"
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	var challenge, proof string
	for {
		select {
		case <-ctx.Done():
			stop()
		default:
			v := url.Values{}
			v.Set("challenge", challenge)
			v.Set("proof", proof)

			u := url.URL{
				Scheme:   "http",
				Host:     net.JoinHostPort(host, port),
				Path:     "/quote",
				RawQuery: v.Encode(),
			}

			resp, err := http.Get(u.String())
			if err != nil {
				logger.Fatal("failed to do get request", zap.Error(err))
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				logger.Error("failed to read body", zap.Error(err))
				continue
			}
			_ = resp.Body.Close()

			switch resp.StatusCode {
			case http.StatusOK:
				var q transport.GetQuoteResponse
				if err := json.Unmarshal(body, &q); err != nil {
					logger.Error("failed to unmarshall quote response", zap.Error(err))
					continue
				}

				logger.Info("got a quote", zap.String("quote", q.Quote))
				// reset proof and a challenge
				challenge, proof = "", ""
			case http.StatusTooManyRequests:
				var c transport.ChallengeResponse
				if err := json.Unmarshal(body, &c); err != nil {
					logger.Error("failed to unmarshall challenge response", zap.Error(err))
					continue
				}

				logger.Info("got a challenge", zap.String("challenge", c.Challenge),
					zap.Int("difficulty", c.Difficulty))

				p, err := service.SolveChallenge(c.Challenge, c.Difficulty)
				if err != nil {
					logger.Error("failed to solve the challenge", zap.Error(err))
					continue
				}

				proof = p
				challenge = c.Challenge
				logger.Info("solved the challenge", zap.String("challenge", challenge),
					zap.String("proof", proof))
			case http.StatusBadRequest:
				logger.Info("got bad request response", zap.String("challenge", challenge),
					zap.String("proof", proof), zap.String("body", string(body)))
				// reset proof and a challenge
				challenge, proof = "", ""
			default:
				logger.Info("got unexpected response", zap.Int("status", resp.StatusCode),
					zap.String("challenge", challenge), zap.String("proof", proof),
					zap.String("body", string(body)))
			}
		}
	}
}
