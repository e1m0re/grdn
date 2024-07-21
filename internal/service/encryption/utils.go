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

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break // fall through
	}
	return nil, errors.New("key type is not RSA")
}

// parseRSAPrivateKeyFromPEMStr parses public key from PEM.
func parseRSAPrivateKeyFromPEMStr(data []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}
