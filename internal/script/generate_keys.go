package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func main() {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Printf("Error generating private key: %v\n", err)
		return
	}

	// Encode private key to PKCS#8 format
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		fmt.Printf("Error marshaling private key: %v\n", err)
		return
	}

	privateKeyPEM := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	privateKeyFile, err := os.Create("private_key_v1.pem")
	if err != nil {
		fmt.Printf("Error creating private key file: %v\n", err)
		return
	}
	defer privateKeyFile.Close()

	err = pem.Encode(privateKeyFile, privateKeyPEM)
	if err != nil {
		fmt.Printf("Error encoding private key: %v\n", err)
		return
	}

	// Encode public key to PKIX format
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		fmt.Printf("Error marshaling public key: %v\n", err)
		return
	}

	publicKeyPEM := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	publicKeyFile, err := os.Create("public_key_v1.pem")
	if err != nil {
		fmt.Printf("Error creating public key file: %v\n", err)
		return
	}
	defer publicKeyFile.Close()

	err = pem.Encode(publicKeyFile, publicKeyPEM)
	if err != nil {
		fmt.Printf("Error encoding public key: %v\n", err)
		return
	}

	fmt.Println("✅ RSA key pair generated successfully!")
	fmt.Println("📁 private_key_v1.pem (PKCS#8 format)")
	fmt.Println("📁 public_key_v1.pem (PKIX format)")
}
