package transport

import (
	"net"
	"net/http"

	"go.uber.org/zap"
)

func (h *Handler) ChallengeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		challenge := r.URL.Query().Get("challenge")
		proof := r.URL.Query().Get("proof")

		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			h.l.Error("failed to parse host", zap.Error(err))
			h.writeResponse(w, ErrorResponse{Message: "internal error"}, http.StatusInternalServerError)
			return
		}

		if proof == "" {
			c, difficulty, err := h.p.GenerateChallenge(host)
			if err != nil {
				h.l.Error("failed to generate challenge", zap.Error(err))
				h.writeResponse(w, ErrorResponse{Message: "internal error"}, http.StatusInternalServerError)
				return
			}

			h.writeResponse(w, ChallengeResponse{Challenge: c, Difficulty: difficulty}, http.StatusTooManyRequests)
			return
		}

		if err := h.p.VerifyProof(host, challenge, proof); err != nil {
			h.writeResponse(w, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type ChallengeResponse struct {
	Challenge  string `json:"challenge"`
	Difficulty int    `json:"difficulty"`
}
