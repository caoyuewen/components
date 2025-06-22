package rsaencry

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func LoadRSAPrivateKeyFromFile(filename string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("invalid private key data")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("invalid private key type")
	}

	return rsaPrivateKey, nil
}

// SignRSA 私钥签名
func SignRSA(privateKey *rsa.PrivateKey, message []byte) ([]byte, error) {
	hashed := sha256.Sum256(message)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return nil, err
	}
	return signature, nil
}

// EncryptPrivateKeyF 分段加密
func EncryptPrivateKeyF(privateKey *rsa.PrivateKey, plaintext []byte) ([]byte, error) {
	var ciphertext []byte
	keySize := privateKey.PublicKey.N.BitLen() / 8

	for len(plaintext) > 0 {
		var block []byte
		if len(plaintext) > keySize-11 {
			block, plaintext = plaintext[:keySize-11], plaintext[keySize-11:]
		} else {
			block, plaintext = plaintext, nil
		}

		cipherBlock, err := rsa.EncryptPKCS1v15(rand.Reader, &privateKey.PublicKey, block)
		if err != nil {
			return nil, err
		}
		ciphertext = append(ciphertext, cipherBlock...)
	}

	return ciphertext, nil
}
