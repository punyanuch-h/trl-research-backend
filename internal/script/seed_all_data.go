package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"trl-research-backend/internal/database"
	"trl-research-backend/internal/models"
)

func main() {
	// 🔥 Initialize Firebase
	database.InitFirebase("trl-research-service-account.json")
	defer database.CloseFirebase()
	ctx := context.Background()
	client := database.FirestoreClient

	fmt.Println("🌱 Starting Firestore seeding process...")
	fmt.Println(strings.Repeat("=", 60))

	// =============================
	// 1️⃣ Admins
	// =============================
	admins := []models.AdminInfo{
		{
			AdminPrefix:           "Dr.",
			AdminAcademicPosition: "Assistant Professor",
			AdminFirstName:        "Ann",
			AdminLastName:         "Smith",
			AdminDepartment:       "Computer Science",
			AdminPhoneNumber:      "+66-81-234-5678",
			AdminEmail:            "admin@example.com",
			AdminPassword:         "password123",
		},
	}

	for _, admin := range admins {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(admin.AdminPassword), bcrypt.DefaultCost)
		admin.AdminPassword = string(hashed)
		admin.CreatedAt = time.Now()
		admin.UpdatedAt = time.Now()

		docRef := client.Collection("admin_info").Doc(admin.AdminEmail)
		_, err := docRef.Set(ctx, admin)
		if err != nil {
			log.Printf("❌ Failed to seed admin %s: %v\n", admin.AdminEmail, err)
		} else {
			fmt.Printf("✅ Admin seeded: %s\n", admin.AdminEmail)
		}
	}

	// =============================
	// 2️⃣ Researchers
	// =============================
	researchers := []models.ResearcherInfo{
		{
			AdminID:                    "A-00001",
			ResearcherPrefix:           "Dr.",
			ResearcherAcademicPosition: "Research Fellow",
			ResearcherFirstName:        "Pair",
			ResearcherLastName:         "Brown",
			ResearcherDepartment:       "Software Engineering",
			ResearcherPhoneNumber:      "+66-83-111-2222",
			ResearcherEmail:            "researcher@example.com",
			ResearcherPassword:         "password123",
		},
	}

	for _, r := range researchers {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(r.ResearcherPassword), bcrypt.DefaultCost)
		r.ResearcherPassword = string(hashed)
		r.CreatedAt = time.Now()
		r.UpdatedAt = time.Now()

		docRef := client.Collection("researchers").Doc(r.ResearcherEmail)
		_, err := docRef.Set(ctx, r)
		if err != nil {
			log.Printf("❌ Failed to seed researcher %s: %v\n", r.ResearcherEmail, err)
		} else {
			fmt.Printf("✅ Researcher seeded: %s\n", r.ResearcherEmail)
		}
	}

	// =============================
	// 3️⃣ Coordinators
	// =============================
	coordinators := []map[string]interface{}{
		{
			"coordinator_id":    "C-0001",
			"coordinator_email": "coordinator1@university.edu",
			"coordinator_name":  "Dr. Michael Chen",
			"coordinator_phone": "+66-91-111-1111",
			"department":        "Research Development",
			"created_at":        time.Now(),
			"updated_at":        time.Now(),
		},
	}

	for _, c := range coordinators {
		docRef := client.Collection("coordinators").Doc(c["coordinator_email"].(string))
		_, err := docRef.Set(ctx, c)
		if err != nil {
			log.Printf("❌ Failed to seed coordinator %v\n", err)
		} else {
			fmt.Printf("✅ Coordinator seeded: %s\n", c["coordinator_email"])
		}
	}

	// =============================
	// 4️⃣ Cases
	// =============================
	cases := []map[string]interface{}{
		{
			"case_id":          "CS-0001",
			"researcher_id":    "RS-0001",
			"case_title":       "AI-powered Diagnosis",
			"case_type":        "Software",
			"case_description": "Developing ML models for early disease detection.",
			"status":           true,
			"created_at":       time.Now(),
			"updated_at":       time.Now(),
		},
	}

	for _, c := range cases {
		docRef := client.Collection("cases").Doc(c["case_id"].(string))
		_, err := docRef.Set(ctx, c)
		if err != nil {
			log.Printf("❌ Failed to seed case %v\n", err)
		} else {
			fmt.Printf("✅ Case seeded: %s\n", c["case_id"])
		}
	}

	// =============================
	// 5️⃣ Appointments
	// =============================
	appointments := []map[string]interface{}{
		{
			"appointment_id": "AP-0001",
			"case_id":        "CS-0001",
			"date":           time.Now().AddDate(0, 0, 7),
			"status":         "Scheduled",
			"location":       "Conference Room A",
			"note":           "Discuss initial progress",
			"summary":        "Introductory meeting with researcher",
			"created_at":     time.Now(),
			"updated_at":     time.Now(),
		},
	}

	for _, a := range appointments {
		docRef := client.Collection("appointments").Doc(a["appointment_id"].(string))
		_, err := docRef.Set(ctx, a)
		if err != nil {
			log.Printf("❌ Failed to seed appointment %v\n", err)
		} else {
			fmt.Printf("✅ Appointment seeded: %s\n", a["appointment_id"])
		}
	}

	// =============================
	// 6️⃣ Assessment TRL
	// =============================
	assessments := []map[string]interface{}{
		{
			"case_id":          "CS-0001",
			"trl_level_result": 3,
			"rq1_answer":       true,
			"rq2_answer":       false,
			"rq3_answer":       true,
			"created_at":       time.Now(),
			"updated_at":       time.Now(),
		},
	}

	for _, a := range assessments {
		docRef := client.Collection("assessment_trl").Doc(a["case_id"].(string))
		_, err := docRef.Set(ctx, a)
		if err != nil {
			log.Printf("❌ Failed to seed assessment %v\n", err)
		} else {
			fmt.Printf("✅ Assessment TRL seeded for case: %s\n", a["case_id"])
		}
	}

	// =============================
	// 7️⃣ Intellectual Property
	// =============================
	ips := []map[string]interface{}{
		{
			"case_id":              "CS-0001",
			"ip_types":             "Patent",
			"ip_protection_status": "Application Filed",
			"ip_request_number":    "US2024001234A1",
			"created_at":           time.Now(),
			"updated_at":           time.Now(),
		},
	}

	for _, ip := range ips {
		docRef := client.Collection("intellectual_property").Doc(ip["case_id"].(string))
		_, err := docRef.Set(ctx, ip)
		if err != nil {
			log.Printf("❌ Failed to seed IP %v\n", err)
		} else {
			fmt.Printf("✅ Intellectual Property seeded for case: %s\n", ip["case_id"])
		}
	}

	// =============================
	// 8️⃣ Supporters
	// =============================
	supporters := []map[string]interface{}{
		{
			"case_id":      "CS-0001",
			"need_test":    true,
			"need_funding": false,
			"need_partners": true,
			"need_guidelines": true,
			"need_account": false,
			"created_at":   time.Now(),
			"updated_at":   time.Now(),
		},
	}

	for _, s := range supporters {
		docRef := client.Collection("supporters").Doc(s["case_id"].(string))
		_, err := docRef.Set(ctx, s)
		if err != nil {
			log.Printf("❌ Failed to seed supporter %v\n", err)
		} else {
			fmt.Printf("✅ Supporter seeded for case: %s\n", s["case_id"])
		}
	}

	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("🎉 Firestore seeding completed successfully!")
}
