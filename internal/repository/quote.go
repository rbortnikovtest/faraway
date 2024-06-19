package repository

import (
	"crypto/rand"
	"encoding/csv"
	"io"
	"math/big"
	"os"

	"faraway/internal/domain"
)

type QuoteRepository struct {
	quotes []domain.Quote
}

func NewQuoteRepository(filePath string) (*QuoteRepository, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var (
		reader = csv.NewReader(f)
		quotes []domain.Quote
	)
	for {
		r, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return &QuoteRepository{quotes: quotes}, err
		}

		quotes = append(quotes, domain.Quote{
			Quote:  r[1],
			Author: r[0],
		})
	}
}

func (q *QuoteRepository) GetRandomQuote() (domain.Quote, error) {
	i, err := rand.Int(rand.Reader, big.NewInt(int64(len(q.quotes))))
	if err != nil {
		return domain.Quote{}, err
	}
	return q.quotes[i.Int64()], nil
}
