package service

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"faraway/internal/domain"
)

//go:generate mockgen -source=puzzler.go -destination=puzzler_mock.go -package=service

const challengeTTL = 5 * time.Second

type KeyRepository interface {
	GetPrivateKey() (*rsa.PrivateKey, error)
}

type Puzzler struct {
	kr         KeyRepository
	difficulty int
}

func NewPuzzler(kr KeyRepository, difficulty int) *Puzzler {
	return &Puzzler{kr: kr, difficulty: difficulty}
}

func (p *Puzzler) GenerateChallenge(id string) (string, int, error) {
	key, err := p.kr.GetPrivateKey()
	if err != nil {
		return "", 0, fmt.Errorf("get private key: %w", err)
	}

	prefix := fmt.Sprintf("%s:%d", id, time.Now().Unix())
	signed, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &key.PublicKey, []byte(prefix), nil)
	if err != nil {
		return "", 0, fmt.Errorf("encrypt prefix: %w", err)
	}

	return hex.EncodeToString(signed), p.difficulty, nil
}

func (p *Puzzler) VerifyProof(id, challenge, proof string) error {
	signed, err := hex.DecodeString(challenge)
	if err != nil {
		return fmt.Errorf("decode challenge: %w", err)
	}

	if err := p.verifyChallenge(id, signed); err != nil {
		return fmt.Errorf("malformed challenge: %w", err)
	}

	once, err := hex.DecodeString(proof)
	if err != nil {
		return fmt.Errorf("decode proof: %w", err)
	}

	hash := sha256.Sum256(append(signed, once...))
	for i := range p.difficulty {
		if hash[i] != 0 {
			return domain.InvalidProof
		}
	}

	return nil
}

func SolveChallenge(challenge string, difficulty int) (string, error) {
	prefix, err := hex.DecodeString(challenge)
	if err != nil {
		return "", fmt.Errorf("%w: decode challenge: %w", domain.InvalidChallenge, err)
	}

	var once uint64
	for once = range math.MaxUint64 {
		hash := sha256.Sum256(binary.BigEndian.AppendUint64(prefix, once))
		solved := true
		for i := range difficulty {
			if hash[i] != 0 {
				solved = false
				break
			}
		}

		if solved {
			buf := make([]byte, 8)
			binary.BigEndian.PutUint64(buf, once)
			return hex.EncodeToString(buf), nil
		}
	}

	return "", errors.New("can't solve the challenge")
}

func (p *Puzzler) verifyChallenge(id string, signed []byte) error {
	key, err := p.kr.GetPrivateKey()
	if err != nil {
		return fmt.Errorf("get private key: %w", err)
	}

	prefix, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, key, signed, nil)
	if err != nil {
		return fmt.Errorf("%w: decrypt challenge: %w", domain.InvalidChallenge, err)
	}

	parts := strings.Split(string(prefix), ":")
	if len(parts) != 2 || parts[0] != id {
		return domain.InvalidChallenge
	}

	t, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return fmt.Errorf("%w: parse timestamp: %w", domain.InvalidChallenge, err)
	}

	now := time.Now()
	ct := time.Unix(t, 0)

	if now.Before(ct) {
		return fmt.Errorf("%w: invalid time", domain.InvalidChallenge)
	}

	if now.After(ct.Add(challengeTTL)) {
		return fmt.Errorf("%w: challenge expired", domain.InvalidChallenge)
	}

	return nil
}
