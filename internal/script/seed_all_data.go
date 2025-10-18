package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"trl-research-backend/internal/database"
	"trl-research-backend/internal/models"

	"golang.org/x/crypto/bcrypt"
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
			AdminID:               "A-00001",
			AdminPrefix:           "Dr.",
			AdminAcademicPosition: "Assistant Professor",
			AdminFirstName:        "Ann",
			AdminLastName:         "Smith",
			AdminDepartment:       "Computer Science",
			AdminPhoneNumber:      "+66-81-234-5678",
			AdminEmail:            "admin@example.com",
			AdminPassword:         "password123",
			CaseID:                "CS-00001",
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
			ResearcherID:               "RS-00001",
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
	coordinators := []models.CoordinatorInfo{
		{
			CoordinatorID:    "C-00001",
			CoordinatorEmail: "coordinator1@university.edu",
			CoordinatorName:  "Dr. Michael Chen",
			CoordinatorPhone: "+66-91-111-1111",
			Department:       "Research Development",
			CaseID:           "CS-00001",
		},
	}

	for _, c := range coordinators {
		c.CreatedAt = time.Now()
		c.UpdatedAt = time.Now()
		docRef := client.Collection("coordinators").Doc(c.CoordinatorEmail)
		_, err := docRef.Set(ctx, c)
		if err != nil {
			log.Printf("❌ Failed to seed coordinator %v\n", err)
		} else {
			fmt.Printf("✅ Coordinator seeded: %s\n", c.CoordinatorEmail)
		}
	}

	// =============================
	// 4️⃣ Cases
	// =============================
	cases := []models.CaseInfo{
		{
			CaseID:           "CS-00001",
			CoordinatorEmail: "coordinator1@university.edu",
			TrlScore:         "3",
			TrlSuggestion:    "Focus on prototype development",
			Status:           true,
			IsUrgent:         false,
			UrgentReason:     "",
			UrgentFeedback:   "",
			CaseTitle:        "AI-powered Diagnosis",
			CaseType:         "Software",
			CaseDescription:  "Developing ML models for early disease detection.",
			CaseKeywords:     "AI, Machine Learning, Medical Diagnosis",
			ResearcherID:     "RS-00001",
		},
	}

	for _, c := range cases {
		c.CreatedAt = time.Now()
		c.UpdatedAt = time.Now()
		docRef := client.Collection("cases").Doc(c.CaseID)
		_, err := docRef.Set(ctx, c)
		if err != nil {
			log.Printf("❌ Failed to seed case %v\n", err)
		} else {
			fmt.Printf("✅ Case seeded: %s\n", c.CaseID)
		}
	}

	// =============================
	// 5️⃣ Appointments
	// =============================
	appointments := []models.Appointment{
		{
			AppointmentID: "AP-00001",
			CaseID:        "CS-00001",
			Date:          time.Now().AddDate(0, 0, 7),
			Status:        "Scheduled",
			Location:      "Conference Room A",
			Note:          "Discuss initial progress",
			Summary:       "Introductory meeting with researcher",
		},
	}

	for _, a := range appointments {
		a.CreatedAt = time.Now()
		a.UpdatedAt = time.Now()
		docRef := client.Collection("appointments").Doc(a.AppointmentID)
		_, err := docRef.Set(ctx, a)
		if err != nil {
			log.Printf("❌ Failed to seed appointment %v\n", err)
		} else {
			fmt.Printf("✅ Appointment seeded: %s\n", a.AppointmentID)
		}
	}

	// =============================
	// 6️⃣ Assessment TRL
	// =============================
	assessments := []models.AssessmentTrl{
		{
			ID:             "AT-00001",
			CaseID:         "CS-00001",
			TrlLevelResult: 3,
			Rq1Answer:      true,
			Rq2Answer:      false,
			Rq3Answer:      true,
			Rq4Answer:      false,
			Rq5Answer:      true,
			Rq6Answer:      false,
			Rq7Answer:      true,
			Cq1Answer:      []string{"Option A", "Option B"},
			Cq2Answer:      []string{"Option C"},
			Cq3Answer:      []string{"Option D", "Option E"},
			Cq4Answer:      []string{"Option F"},
			Cq5Answer:      []string{"Option G", "Option H"},
			Cq6Answer:      []string{"Option I"},
			Cq7Answer:      []string{"Option J", "Option K"},
			Cq8Answer:      []string{"Option L"},
			Cq9Answer:      []string{"Option M", "Option N"},
		},
	}

	for _, a := range assessments {
		a.CreatedAt = time.Now()
		a.UpdatedAt = time.Now()
		docRef := client.Collection("assessment_trl").Doc(a.ID)
		_, err := docRef.Set(ctx, a)
		if err != nil {
			log.Printf("❌ Failed to seed assessment %v\n", err)
		} else {
			fmt.Printf("✅ Assessment TRL seeded for case: %s\n", a.CaseID)
		}
	}

	// =============================
	// 7️⃣ Intellectual Property
	// =============================
	ips := []models.IntellectualProperty{
		{
			ID:                 "IP-00001",
			CaseID:             "CS-00001",
			IPTypes:            "Patent",
			IPProtectionStatus: "Application Filed",
			IPRequestNumber:    "US2024001234A1",
		},
	}

	for _, ip := range ips {
		ip.CreatedAt = time.Now()
		ip.UpdatedAt = time.Now()
		docRef := client.Collection("intellectual_property").Doc(ip.ID)
		_, err := docRef.Set(ctx, ip)
		if err != nil {
			log.Printf("❌ Failed to seed IP %v\n", err)
		} else {
			fmt.Printf("✅ Intellectual Property seeded for case: %s\n", ip.CaseID)
		}
	}

	// =============================
	// 8️⃣ Supporters
	// =============================
	supporters := []models.Supporter{
		{
			SupporterID:                     "SP-00001",
			CaseID:                          "CS-00001",
			SupportResearch:                 true,
			SupportVDC:                      false,
			SupportSiEIC:                    true,
			NeedProtectIntellectualProperty: true,
			NeedCoDevelopers:                false,
			NeedActivities:                  true,
			NeedTest:                        true,
			NeedCapital:                     false,
			NeedPartners:                    true,
			NeedGuidelines:                  true,
			NeedCertification:               false,
			NeedAccount:                     false,
			Need:                            "Additional funding and technical support",
			AdditionalDocuments:             "Research proposal and preliminary results",
		},
	}

	for _, s := range supporters {
		s.CreatedAt = time.Now()
		s.UpdatedAt = time.Now()
		docRef := client.Collection("supporters").Doc(s.SupporterID)
		_, err := docRef.Set(ctx, s)
		if err != nil {
			log.Printf("❌ Failed to seed supporter %v\n", err)
		} else {
			fmt.Printf("✅ Supporter seeded for case: %s\n", s.CaseID)
		}
	}

	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("🎉 Firestore seeding completed successfully!")
}
