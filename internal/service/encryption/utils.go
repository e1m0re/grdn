package encryption

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

// parseRSAPublicKeyFromPEMStr parses public key from PEM.
func parseRSAPublicKeyFromPEMStr(data []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	return x509.ParsePKCS1PublicKey(block.Bytes)
}

// parseRSAPrivateKeyFromPEMStr parses public key from PEM.
func parseRSAPrivateKeyFromPEMStr(data []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}
