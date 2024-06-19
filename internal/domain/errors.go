package domain

import "errors"

var (
	InvalidProof     = errors.New("invalid proof")
	InvalidChallenge = errors.New("invalid challenge")
)
