package encryption

import (
	"crypto"
	"crypto/rsa"
	"fmt"
	"os"
)

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

	privateKey, err := parseRsaPrivateKeyFromPemStr(pemKey)
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
		return make([]byte, 0), fmt.Errorf("RSA private key not specified")
	}

	return d.privateKey.Decrypt(nil, ciphertext, rsa.OAEPOptions{Hash: crypto.SHA256})
}
