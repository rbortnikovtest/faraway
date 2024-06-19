package repository

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

// KeyRepository provides access to the private key.
// TODO: Keep the private key in a secure location.
type KeyRepository struct {
	key *rsa.PrivateKey
}

func NewKeyRepository() (*KeyRepository, error) {
	k, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("generate key: %w", err)
	}

	return &KeyRepository{key: k}, nil
}

func (k *KeyRepository) GetPrivateKey() (*rsa.PrivateKey, error) {
	return k.key, nil
}
