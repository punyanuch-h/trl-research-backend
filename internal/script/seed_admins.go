package main

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"

	"trl-research-backend/internal/database"
	"trl-research-backend/internal/models"
	"trl-research-backend/internal/repository"
)

func main() {
	// ✅ 1. Initialize Firestore connection
	database.InitFirebase("trl-research-service-account.json")
	defer database.CloseFirebase()

	// ✅ 2. Initialize repositories
	adminRepo := repository.NewAdminRepo(database.FirestoreClient)
	researcherRepo := repository.NewResearcherRepo(database.FirestoreClient)

	// ✅ 3. Create default admin
	admin := models.AdminInfo{
		AdminPrefix:           "Dr.",
		AdminAcademicPosition: "Professor",
		AdminFirstName:        "System",
		AdminLastName:         "Admin",
		AdminDepartment:       "IT Department",
		AdminPhoneNumber:      "+66-81-111-1111",
		AdminEmail:            "admin@example.com",
		AdminPassword:         "password123", // will hash below
		CaseID:                "",
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	// ✅ Hash password before saving
	hashedAdminPwd, err := bcrypt.GenerateFromPassword([]byte(admin.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("❌ Failed to hash admin password: %v", err)
	}
	admin.AdminPassword = string(hashedAdminPwd)

	if err := adminRepo.CreateAdmin(&admin); err != nil {
		log.Fatalf("❌ Failed to create admin: %v", err)
	}
	fmt.Println("✅ Admin created successfully:")
	fmt.Printf("   📧 Email: %s\n", admin.AdminEmail)
	fmt.Printf("   🔑 Password: password123\n\n")

	// ✅ 4. Create default researcher
	researcher := models.ResearcherInfo{
		ResearcherPrefix:     "Mr.",
		ResearcherFirstName:  "Default",
		ResearcherLastName:   "Researcher",
		ResearcherDepartment: "AI Lab",
		ResearcherPhoneNumber: "+66-82-222-2222",
		ResearcherEmail:      "researcher@example.com",
		ResearcherPassword:   "password123",
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	hashedResPwd, err := bcrypt.GenerateFromPassword([]byte(researcher.ResearcherPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("❌ Failed to hash researcher password: %v", err)
	}
	researcher.ResearcherPassword = string(hashedResPwd)

	if err := researcherRepo.CreateResearcher(&researcher); err != nil {
		log.Fatalf("❌ Failed to create researcher: %v", err)
	}
	fmt.Println("✅ Researcher created successfully:")
	fmt.Printf("   📧 Email: %s\n", researcher.ResearcherEmail)
	fmt.Printf("   🔑 Password: password123\n\n")

	fmt.Println("🎉 Seeding completed successfully!")
}
