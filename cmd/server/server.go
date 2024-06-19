package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"faraway/internal/repository"
	"faraway/internal/service"
	"faraway/internal/transport"
)

// TODO: move to config
const (
	difficulty = 3
	quotesFile = "./assets/quotes.csv"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	var (
		host = os.Getenv("HOST")
		port = os.Getenv("PORT")
	)

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "8181"
	}

	keyRepo, err := repository.NewKeyRepository()
	if err != nil {
		logger.Fatal("failed to init key repository", zap.Error(err))
	}

	quoteRepo, err := repository.NewQuoteRepository(quotesFile)
	if err != nil {
		logger.Fatal("failed to init quote repository", zap.Error(err))
	}

	puzzler := service.NewPuzzler(keyRepo, difficulty)
	handler := transport.NewHandler(logger, puzzler, quoteRepo)

	router := mux.NewRouter()
	router.Use(handler.ChallengeMiddleware)
	router.HandleFunc("/quote", handler.GetQuote).Methods(http.MethodGet)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	server := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: router,
	}

	go func() {
		logger.Info("starting server", zap.String("host", host), zap.String("port", port))
		if err := server.ListenAndServe(); err != nil {
			logger.Fatal("failed to serve", zap.Error(err))
		}
	}()

	<-ctx.Done()
	stop()

	sdCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(sdCtx); err != nil {
		logger.Error("failed to shutdown server", zap.Error(err))
	}
}
