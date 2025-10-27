package main

import (
	"fmt"
	"trl-research-backend/internal/utils"
)

func main() {
	fmt.Println("🔑 Testing KeyProvider...")

	// Test key provider initialization
	kp, err := utils.NewEnvKeyProvider()
	if err != nil {
		fmt.Printf("❌ KeyProvider initialization failed: %v\n", err)
		return
	}

	fmt.Println("✅ KeyProvider initialized successfully!")

	// Test private key retrieval
	priv, err := kp.GetPrivateKey("v1")
	if err != nil {
		fmt.Printf("❌ Private key retrieval failed: %v\n", err)
		return
	}
	fmt.Printf("✅ Private key retrieved successfully (bits=%d)\n", priv.N.BitLen())

	// Test public key retrieval
	pub, err := kp.GetPublicKey("v1")
	if err != nil {
		fmt.Printf("❌ Public key retrieval failed: %v\n", err)
		return
	}
	fmt.Printf("✅ Public key retrieved successfully (bits=%d)\n", pub.N.BitLen())

	// Test JWT generation
	token, err := utils.GenerateJWT(
		"test-admin-id",
		"test@example.com",
		"admin",
		"",
		"",
		"test-issuer",
		"test-audience",
		"v1",
		24,
		*kp,
	)
	if err != nil {
		fmt.Printf("❌ JWT generation failed: %v\n", err)
		return
	}
	fmt.Printf("✅ JWT generated successfully: %s...\n", token[:50])

	// Test JWT validation
	claims, err := utils.ValidateJWT(
		token,
		"test-issuer",
		"test-audience",
		*kp,
	)
	if err != nil {
		fmt.Printf("❌ JWT validation failed: %v\n", err)
		return
	}
	fmt.Printf("✅ JWT validated successfully for: %s\n", claims.UserEmail)

	fmt.Println("\n🎉 All key tests passed!")
}
