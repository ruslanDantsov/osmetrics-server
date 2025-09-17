package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
)

func LoadPublicKey(path string) (*rsa.PublicKey, error) {
	keyData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block with public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not RSA public key")
	}

	return rsaPub, nil
}

func LoadPrivateKey(path string) (*rsa.PrivateKey, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(b)
	if block == nil {
		return nil, errors.New("failed to decode PEM block from private key")
	}

	// поддерживаем PKCS#1 и PKCS#8
	if block.Type == "RSA PRIVATE KEY" {
		priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		return priv, nil
	}

	if block.Type == "PRIVATE KEY" {
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		priv, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("not RSA private key inside PKCS#8")
		}
		return priv, nil
	}

	return nil, fmt.Errorf("unsupported private key type: %s", block.Type)
}

func EncryptPayload(pubKey *rsa.PublicKey, data []byte) ([]byte, error) {
	if pubKey == nil {
		// Ключ не задан, возвращаем оригинальные данные (JSON метрики)
		return data, nil
	}

	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, data)
	if err != nil {
		return nil, err
	}

	// Возвращаем чистый Base64
	encoded := base64.StdEncoding.EncodeToString(encrypted)
	return []byte(encoded), nil
}

func DecryptRSA(privKey *rsa.PrivateKey, cipherData []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, privKey, cipherData)
}
