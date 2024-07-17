package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"os"
)

type Encryptor interface {
	// Encrypt encrypts specified bytes with RSA OAEP.
	Encrypt(plaintext []byte) ([]byte, error)
}

type encryptor struct {
	publicKey *rsa.PublicKey
}

// NewEncryptor initiates new instance of Encryptor.
func NewEncryptor(keyFile string) (Encryptor, error) {
	pemKey, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}

	publicKey, err := parseRsaPublicKeyFromPemStr(pemKey)
	if err != nil {
		return nil, err
	}

	return &encryptor{
		publicKey: publicKey,
	}, nil
}

// Encrypt encrypts specified bytes with RSA OAEP.
func (e *encryptor) Encrypt(plaintext []byte) ([]byte, error) {
	if e.publicKey == nil {
		return make([]byte, 0), fmt.Errorf("RSA public key not specified")
	}

	return rsa.EncryptOAEP(sha256.New(), rand.Reader, e.publicKey, plaintext, nil)
}
