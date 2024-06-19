package transport

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"faraway/internal/domain"
)

type Puzzler interface {
	GenerateChallenge(id string) (string, int, error)
	VerifyProof(id, challenge, proof string) error
}

type QuoteProvider interface {
	GetRandomQuote() (domain.Quote, error)
}

type Handler struct {
	l *zap.Logger
	p Puzzler
	q QuoteProvider
}

func NewHandler(l *zap.Logger, p Puzzler, q QuoteProvider) *Handler {
	return &Handler{l: l, p: p, q: q}
}

func (h *Handler) GetQuote(w http.ResponseWriter, r *http.Request) {
	q, err := h.q.GetRandomQuote()
	if err != nil {
		h.l.Error("failed to get random quote", zap.Error(err))
		h.writeResponse(w, ErrorResponse{Message: "internal error"}, http.StatusInternalServerError)
		return
	}

	h.writeResponse(w, GetQuoteResponse{
		Quote:  q.Quote,
		Author: q.Author,
	}, http.StatusOK)
}

func (h *Handler) writeResponse(w http.ResponseWriter, data any, status int) {
	p, err := json.Marshal(data)
	if err != nil {
		h.l.Error("failed to marshal data", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(p); err != nil {
		h.l.Error("failed to write challenge response", zap.Error(err))
		return
	}
}

type GetQuoteResponse struct {
	Quote  string `json:"quote"`
	Author string `json:"author"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
