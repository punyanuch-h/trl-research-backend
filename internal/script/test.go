package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	hash := "$2a$10$CvSLirnqFeP5..jPrGID4u80JTw3Mmndm0vzDgDd/s8/cV4HDIxAi"
	password := "password123"

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println("❌ Password mismatch:", err)
	} else {
		fmt.Println("✅ Password valid")
	}
}
