package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"trl-research-backend/internal/database"
	"trl-research-backend/internal/models"
	"trl-research-backend/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Initialize Firebase
	database.InitFirebase("localServiceAccountKey.json")
	defer database.CloseFirebase()

	// Create admin repository
	adminRepo := repository.NewAdminRepo(database.FirestoreClient)

	// Define admin data
	admins := []models.AdminInfo{
		{
			AdminPrefix:           "Dr.",
			AdminAcademicPosition: "Assistant Professor",
			AdminFirstName:        "Ann",
			AdminLastName:         "Smith",
			AdminDepartment:       "Computer Science",
			AdminPhoneNumber:      "+66-81-234-5678",
			AdminEmail:            "Ann@gmail.com",
			AdminPassword:         "password123", // Will be hashed
			CaseID:                "",
		},
		{
			AdminPrefix:           "Prof.",
			AdminAcademicPosition: "Professor",
			AdminFirstName:        "Mint",
			AdminLastName:         "Johnson",
			AdminDepartment:       "Information Technology",
			AdminPhoneNumber:      "+66-82-345-6789",
			AdminEmail:            "Mint@gmail.com",
			AdminPassword:         "password123", // Will be hashed
			CaseID:                "",
		},
		{
			AdminPrefix:           "Dr.",
			AdminAcademicPosition: "Associate Professor",
			AdminFirstName:        "Pair",
			AdminLastName:         "Brown",
			AdminDepartment:       "Software Engineering",
			AdminPhoneNumber:      "+66-83-456-7890",
			AdminEmail:            "Pair@gmail.com",
			AdminPassword:         "password123", // Will be hashed
			CaseID:                "",
		},
	}

	// Seed admins
	for i, admin := range admins {
		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.AdminPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("❌ Error hashing password for %s: %v", admin.AdminEmail, err)
			continue
		}
		admin.AdminPassword = string(hashedPassword)

		// Set timestamps
		now := time.Now()
		admin.CreatedAt = now
		admin.UpdatedAt = now

		// Create admin
		err = adminRepo.CreateAdmin(&admin)
		if err != nil {
			log.Printf("❌ Error creating admin %s: %v", admin.AdminEmail, err)
			continue
		}

		fmt.Printf("✅ Admin %d created successfully:\n", i+1)
		fmt.Printf("   📧 Email: %s\n", admin.AdminEmail)
		fmt.Printf("   👤 Name: %s %s %s\n", admin.AdminPrefix, admin.AdminFirstName, admin.AdminLastName)
		fmt.Printf("   🏫 Position: %s\n", admin.AdminAcademicPosition)
		fmt.Printf("   🏢 Department: %s\n", admin.AdminDepartment)
		fmt.Printf("   📱 Phone: %s\n", admin.AdminPhoneNumber)
		fmt.Printf("   🔑 Password: password123 (hashed)\n")
		fmt.Printf("   🆔 Admin ID: %s\n", admin.AdminID)
		fmt.Println("   " + strings.Repeat("-", 50))
	}

	fmt.Println("\n🎉 Admin seeding completed!")
	fmt.Println("📝 Note: All admins have the default password 'password123'")
	fmt.Println("🔐 Passwords are securely hashed using bcrypt")
}
