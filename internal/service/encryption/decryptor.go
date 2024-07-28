package encryption

import (
	"crypto"
	"crypto/rsa"
	"errors"
	"os"
)

//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=Decryptor
type Decryptor interface {
	// Decrypt decrypts specified bytes with RSA OAEP.
	Decrypt(ciphertext []byte) ([]byte, error)
}

type decryptor struct {
	privateKey *rsa.PrivateKey
}

// NewDecryptor initiates new instance of Encryptor.
func NewDecryptor(keyFile string) (Decryptor, error) {
	pemKey, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}

	privateKey, err := parseRSAPrivateKeyFromPEMStr(pemKey)
	if err != nil {
		return nil, err
	}

	return &decryptor{
		privateKey: privateKey,
	}, nil
}

// Decrypt decrypts specified bytes with RSA OAEP.
func (d *decryptor) Decrypt(ciphertext []byte) ([]byte, error) {
	if d.privateKey == nil {
		return nil, errors.New("RSA private key not specified")
	}

	return d.privateKey.Decrypt(nil, ciphertext, &rsa.OAEPOptions{Hash: crypto.SHA256})
}
