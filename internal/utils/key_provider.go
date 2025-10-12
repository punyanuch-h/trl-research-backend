package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

// KeyProvider จัดการ private/public key ที่ใช้เซ็นและตรวจสอบ JWT
type KeyProvider struct {
	privateKeys map[string]*rsa.PrivateKey
	publicKeys  map[string]*rsa.PublicKey
}

// NewEnvKeyProvider โหลด key จาก Environment Variables (ไม่ใช้ไฟล์)
func NewEnvKeyProvider() (*KeyProvider, error) {
	// อ่านค่าจาก environment variables
	privKeyB64 := os.Getenv("PRIVATE_KEY_V1_B64")
	pubKeyB64 := os.Getenv("PUBLIC_KEY_V1_B64")

	if privKeyB64 == "" {
		return nil, errors.New("missing PRIVATE_KEY_V1_B64 in environment")
	}
	if pubKeyB64 == "" {
		return nil, errors.New("missing PUBLIC_KEY_V1_B64 in environment")
	}

	// Decode จาก Base64
	privKeyBytes, err := base64.StdEncoding.DecodeString(privKeyB64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key base64: %v", err)
	}

	pubKeyBytes, err := base64.StdEncoding.DecodeString(pubKeyB64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key base64: %v", err)
	}

	// แปลงเป็น object key จริง
	privKey, err := parseStrictPrivateKey(privKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %v", err)
	}

	pubKey, err := parseStrictPublicKey(pubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("invalid public key: %v", err)
	}

	return &KeyProvider{
		privateKeys: map[string]*rsa.PrivateKey{"v1": privKey},
		publicKeys:  map[string]*rsa.PublicKey{"v1": pubKey},
	}, nil
}

// parseStrictPrivateKey รองรับเฉพาะ PKCS#8 (RSA)
func parseStrictPrivateKey(pemBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("cannot decode PEM block (private key)")
	}
	if block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("unexpected PEM type: %s (expected PRIVATE KEY for PKCS#8)", block.Type)
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PKCS#8: %v", err)
	}
	priv, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("private key is not RSA")
	}
	return priv, nil
}

// parseStrictPublicKey รองรับเฉพาะ PKIX (RSA)
func parseStrictPublicKey(pemBytes []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("cannot decode PEM block (public key)")
	}
	if block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("unexpected PEM type: %s (expected PUBLIC KEY for PKIX)", block.Type)
	}
	pubAny, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PKIX public key: %v", err)
	}
	pub, ok := pubAny.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("public key is not RSA")
	}
	return pub, nil
}

// GetPrivateKey ดึง private key ตาม kid
func (f *KeyProvider) GetPrivateKey(kid string) (*rsa.PrivateKey, error) {
	k, ok := f.privateKeys[kid]
	if !ok {
		return nil, fmt.Errorf("private key for kid=%s not found", kid)
	}
	return k, nil
}

// GetPublicKey ดึง public key ตาม kid
func (f *KeyProvider) GetPublicKey(kid string) (*rsa.PublicKey, error) {
	k, ok := f.publicKeys[kid]
	if !ok {
		return nil, fmt.Errorf("public key for kid=%s not found", kid)
	}
	return k, nil
}
