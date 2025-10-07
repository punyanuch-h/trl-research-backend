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

	// Create repositories
	adminRepo := repository.NewAdminRepo(database.FirestoreClient)
	caseRepo := repository.NewCaseRepo(database.FirestoreClient)
	researcherRepo := repository.NewResearcherRepo(database.FirestoreClient)
	coordinatorRepo := repository.NewCoordinatorRepo(database.FirestoreClient)
	appointmentRepo := repository.NewAppointmentRepo(database.FirestoreClient)
	assessmentTrlRepo := repository.NewAssessmentTrlRepo(database.FirestoreClient)
	intellectualPropertyRepo := repository.NewIntellectualPropertyRepo(database.FirestoreClient)
	supporterRepo := repository.NewSupporterRepo(database.FirestoreClient)

	fmt.Println("🌱 Starting comprehensive data seeding...")
	fmt.Println(strings.Repeat("=", 60))

	// 1. Seed Admins
	fmt.Println("\n👥 Seeding Admins...")
	admins := []models.AdminInfo{
		{
			AdminPrefix:           "Dr.",
			AdminAcademicPosition: "Assistant Professor",
			AdminFirstName:        "Ann",
			AdminLastName:         "Smith",
			AdminDepartment:       "Computer Science",
			AdminPhoneNumber:      "+66-81-234-5678",
			AdminEmail:            "ann.smith@university.edu",
			AdminPassword:         "password123",
			CaseID:                "",
		},
		{
			AdminPrefix:           "Prof.",
			AdminAcademicPosition: "Professor",
			AdminFirstName:        "Mint",
			AdminLastName:         "Johnson",
			AdminDepartment:       "Information Technology",
			AdminPhoneNumber:      "+66-82-345-6789",
			AdminEmail:            "mint.johnson@university.edu",
			AdminPassword:         "password123",
			CaseID:                "",
		},
		{
			AdminPrefix:           "Dr.",
			AdminAcademicPosition: "Associate Professor",
			AdminFirstName:        "Pair",
			AdminLastName:         "Brown",
			AdminDepartment:       "Software Engineering",
			AdminPhoneNumber:      "+66-83-456-7890",
			AdminEmail:            "pair.brown@university.edu",
			AdminPassword:         "password123",
			CaseID:                "",
		},
		{
			AdminPrefix:           "Dr.",
			AdminAcademicPosition: "Assistant Professor",
			AdminFirstName:        "Sarah",
			AdminLastName:         "Wilson",
			AdminDepartment:       "Data Science",
			AdminPhoneNumber:      "+66-84-567-8901",
			AdminEmail:            "sarah.wilson@university.edu",
			AdminPassword:         "password123",
			CaseID:                "",
		},
		{
			AdminPrefix:           "Prof.",
			AdminAcademicPosition: "Professor",
			AdminFirstName:        "David",
			AdminLastName:         "Lee",
			AdminDepartment:       "Artificial Intelligence",
			AdminPhoneNumber:      "+66-85-678-9012",
			AdminEmail:            "david.lee@university.edu",
			AdminPassword:         "password123",
			CaseID:                "",
		},
	}

	seedAdmins(adminRepo, admins)

	// 2. Seed Coordinators
	fmt.Println("\n🎯 Seeding Coordinators...")
	coordinators := []models.CoordinatorInfo{
		{
			CoordinatorEmail: "coordinator1@university.edu",
			CoordinatorName:  "Dr. Michael Chen",
			CoordinatorPhone: "+66-91-111-1111",
			Department:       "Research Development",
		},
		{
			CoordinatorEmail: "coordinator2@university.edu",
			CoordinatorName:  "Prof. Lisa Anderson",
			CoordinatorPhone: "+66-92-222-2222",
			Department:       "Innovation Center",
		},
		{
			CoordinatorEmail: "coordinator3@university.edu",
			CoordinatorName:  "Dr. James Taylor",
			CoordinatorPhone: "+66-93-333-3333",
			Department:       "Technology Transfer",
		},
		{
			CoordinatorEmail: "coordinator4@university.edu",
			CoordinatorName:  "Prof. Maria Garcia",
			CoordinatorPhone: "+66-94-444-4444",
			Department:       "Research Support",
		},
		{
			CoordinatorEmail: "coordinator5@university.edu",
			CoordinatorName:  "Dr. Robert Kim",
			CoordinatorPhone: "+66-95-555-5555",
			Department:       "Commercialization",
		},
	}

	seedCoordinators(coordinatorRepo, coordinators)

	// 3. Seed Researchers
	fmt.Println("\n🔬 Seeding Researchers...")
	researchers := []models.ResearcherInfo{
		{
			AdminID:                    "SI-00001", // Will be updated after admin creation
			ResearcherPrefix:           "Dr.",
			ResearcherAcademicPosition: "Research Fellow",
			ResearcherFirstName:        "Alex",
			ResearcherLastName:         "Thompson",
			ResearcherDepartment:       "Computer Science",
			ResearcherPhoneNumber:      "+66-71-111-1111",
			ResearcherEmail:            "alex.thompson@university.edu",
		},
		{
			AdminID:                    "SI-00002",
			ResearcherPrefix:           "Dr.",
			ResearcherAcademicPosition: "Senior Researcher",
			ResearcherFirstName:        "Emma",
			ResearcherLastName:         "Davis",
			ResearcherDepartment:       "Information Technology",
			ResearcherPhoneNumber:      "+66-72-222-2222",
			ResearcherEmail:            "emma.davis@university.edu",
		},
		{
			AdminID:                    "SI-00003",
			ResearcherPrefix:           "Prof.",
			ResearcherAcademicPosition: "Principal Investigator",
			ResearcherFirstName:        "John",
			ResearcherLastName:         "Miller",
			ResearcherDepartment:       "Software Engineering",
			ResearcherPhoneNumber:      "+66-73-333-3333",
			ResearcherEmail:            "john.miller@university.edu",
		},
		{
			AdminID:                    "SI-00004",
			ResearcherPrefix:           "Dr.",
			ResearcherAcademicPosition: "Research Associate",
			ResearcherFirstName:        "Sophia",
			ResearcherLastName:         "White",
			ResearcherDepartment:       "Data Science",
			ResearcherPhoneNumber:      "+66-74-444-4444",
			ResearcherEmail:            "sophia.white@university.edu",
		},
		{
			AdminID:                    "SI-00005",
			ResearcherPrefix:           "Dr.",
			ResearcherAcademicPosition: "Research Scientist",
			ResearcherFirstName:        "Daniel",
			ResearcherLastName:         "Clark",
			ResearcherDepartment:       "Artificial Intelligence",
			ResearcherPhoneNumber:      "+66-75-555-5555",
			ResearcherEmail:            "daniel.clark@university.edu",
		},
	}

	seedResearchers(researcherRepo, researchers)

	// 4. Seed Cases
	fmt.Println("\n📋 Seeding Cases...")
	cases := []models.CaseInfo{
		{
			ResearcherID:     "RS-00001", // Will be updated after researcher creation
			CoordinatorEmail: "coordinator1@university.edu",
			TrlScore:         "TRL-3",
			Status:           true,
			IsUrgent:         false,
			UrgentReason:     "",
			UrgentFeedback:   "",
			CaseTitle:        "AI-Powered Medical Diagnosis System",
			CaseType:         "Software Development",
			CaseDescription:  "Development of an AI system for early detection of medical conditions using machine learning algorithms.",
			CaseKeywords:     "AI, Machine Learning, Medical Diagnosis, Healthcare",
		},
		{
			ResearcherID:     "RS-00002",
			CoordinatorEmail: "coordinator2@university.edu",
			TrlScore:         "TRL-5",
			Status:           true,
			IsUrgent:         true,
			UrgentReason:     "Patent deadline approaching",
			UrgentFeedback:   "Priority review required",
			CaseTitle:        "Blockchain-based Supply Chain Management",
			CaseType:         "Technology Innovation",
			CaseDescription:  "Implementation of blockchain technology for transparent and secure supply chain tracking.",
			CaseKeywords:     "Blockchain, Supply Chain, Transparency, Security",
		},
		{
			ResearcherID:     "RS-00003",
			CoordinatorEmail: "coordinator3@university.edu",
			TrlScore:         "TRL-4",
			Status:           false,
			IsUrgent:         false,
			UrgentReason:     "",
			UrgentFeedback:   "",
			CaseTitle:        "IoT Smart City Infrastructure",
			CaseType:         "Infrastructure Development",
			CaseDescription:  "Development of IoT sensors and systems for smart city infrastructure management.",
			CaseKeywords:     "IoT, Smart City, Infrastructure, Sensors",
		},
		{
			ResearcherID:     "RS-00004",
			CoordinatorEmail: "coordinator4@university.edu",
			TrlScore:         "TRL-6",
			Status:           true,
			IsUrgent:         false,
			UrgentReason:     "",
			UrgentFeedback:   "",
			CaseTitle:        "Renewable Energy Optimization Platform",
			CaseType:         "Energy Technology",
			CaseDescription:  "Platform for optimizing renewable energy distribution and consumption patterns.",
			CaseKeywords:     "Renewable Energy, Optimization, Platform, Sustainability",
		},
		{
			ResearcherID:     "RS-00005",
			CoordinatorEmail: "coordinator5@university.edu",
			TrlScore:         "TRL-2",
			Status:           true,
			IsUrgent:         false,
			UrgentReason:     "",
			UrgentFeedback:   "",
			CaseTitle:        "Quantum Computing Algorithm Development",
			CaseType:         "Advanced Computing",
			CaseDescription:  "Research and development of quantum algorithms for complex computational problems.",
			CaseKeywords:     "Quantum Computing, Algorithms, Research, Advanced Computing",
		},
	}

	seedCases(caseRepo, cases)

	// 5. Seed Appointments
	fmt.Println("\n📅 Seeding Appointments...")
	appointments := []models.Appointment{
		{
			CaseID:   "CS-00001",                  // Will be updated after case creation
			Date:     time.Now().AddDate(0, 0, 7), // 1 week from now
			Status:   "Scheduled",
			Location: "Conference Room A",
			Note:     "Initial project discussion",
			Summary:  "Meeting to discuss project requirements and timeline",
		},
		{
			CaseID:   "CS-00002",
			Date:     time.Now().AddDate(0, 0, 14), // 2 weeks from now
			Status:   "Confirmed",
			Location: "Virtual Meeting",
			Note:     "Progress review meeting",
			Summary:  "Review of current progress and next steps",
		},
		{
			CaseID:   "CS-00003",
			Date:     time.Now().AddDate(0, 0, 21), // 3 weeks from now
			Status:   "Pending",
			Location: "Lab 101",
			Note:     "Technical demonstration",
			Summary:  "Demonstration of prototype functionality",
		},
		{
			CaseID:   "CS-00004",
			Date:     time.Now().AddDate(0, 0, 28), // 4 weeks from now
			Status:   "Scheduled",
			Location: "Innovation Center",
			Note:     "Final presentation",
			Summary:  "Final project presentation and evaluation",
		},
		{
			CaseID:   "CS-00005",
			Date:     time.Now().AddDate(0, 0, 35), // 5 weeks from now
			Status:   "Tentative",
			Location: "Conference Room B",
			Note:     "Follow-up meeting",
			Summary:  "Follow-up discussion on project outcomes",
		},
	}

	seedAppointments(appointmentRepo, appointments)

	// 6. Seed Assessment TRL
	fmt.Println("\n📊 Seeding Assessment TRL...")
	assessments := []models.AssessmentTrl{
		{
			CaseID:         "CS-00001",
			TrlLevelResult: 3,
			Rq1Answer:      true,
			Rq2Answer:      false,
			Rq3Answer:      true,
			Rq4Answer:      true,
			Rq5Answer:      false,
			Rq6Answer:      true,
			Rq7Answer:      false,
			Cq1Answer:      "Basic research completed",
			Cq2Answer:      "Proof of concept demonstrated",
			Cq3Answer:      "Laboratory validation in progress",
			Cq4Answer:      "Component testing underway",
			Cq5Answer:      "System integration planned",
			Cq6Answer:      "Field testing not yet started",
			Cq7Answer:      "Commercial validation pending",
			Cq8Answer:      "Market analysis in progress",
			Cq9Answer:      "Technology transfer discussions ongoing",
		},
		{
			CaseID:         "CS-00002",
			TrlLevelResult: 5,
			Rq1Answer:      true,
			Rq2Answer:      true,
			Rq3Answer:      true,
			Rq4Answer:      true,
			Rq5Answer:      true,
			Rq6Answer:      false,
			Rq7Answer:      false,
			Cq1Answer:      "Research and development completed",
			Cq2Answer:      "Proof of concept validated",
			Cq3Answer:      "Laboratory validation completed",
			Cq4Answer:      "Component testing completed",
			Cq5Answer:      "System integration in progress",
			Cq6Answer:      "Field testing planned",
			Cq7Answer:      "Commercial validation pending",
			Cq8Answer:      "Market analysis completed",
			Cq9Answer:      "Technology transfer in progress",
		},
		{
			CaseID:         "CS-00003",
			TrlLevelResult: 4,
			Rq1Answer:      true,
			Rq2Answer:      true,
			Rq3Answer:      true,
			Rq4Answer:      true,
			Rq5Answer:      false,
			Rq6Answer:      false,
			Rq7Answer:      false,
			Cq1Answer:      "Research completed",
			Cq2Answer:      "Proof of concept validated",
			Cq3Answer:      "Laboratory validation completed",
			Cq4Answer:      "Component testing in progress",
			Cq5Answer:      "System integration planned",
			Cq6Answer:      "Field testing not started",
			Cq7Answer:      "Commercial validation pending",
			Cq8Answer:      "Market analysis in progress",
			Cq9Answer:      "Technology transfer discussions",
		},
		{
			CaseID:         "CS-00004",
			TrlLevelResult: 6,
			Rq1Answer:      true,
			Rq2Answer:      true,
			Rq3Answer:      true,
			Rq4Answer:      true,
			Rq5Answer:      true,
			Rq6Answer:      true,
			Rq7Answer:      false,
			Cq1Answer:      "Research and development completed",
			Cq2Answer:      "Proof of concept validated",
			Cq3Answer:      "Laboratory validation completed",
			Cq4Answer:      "Component testing completed",
			Cq5Answer:      "System integration completed",
			Cq6Answer:      "Field testing in progress",
			Cq7Answer:      "Commercial validation pending",
			Cq8Answer:      "Market analysis completed",
			Cq9Answer:      "Technology transfer in progress",
		},
		{
			CaseID:         "CS-00005",
			TrlLevelResult: 2,
			Rq1Answer:      true,
			Rq2Answer:      false,
			Rq3Answer:      false,
			Rq4Answer:      false,
			Rq5Answer:      false,
			Rq6Answer:      false,
			Rq7Answer:      false,
			Cq1Answer:      "Basic research in progress",
			Cq2Answer:      "Proof of concept under development",
			Cq3Answer:      "Laboratory validation not started",
			Cq4Answer:      "Component testing not started",
			Cq5Answer:      "System integration not planned",
			Cq6Answer:      "Field testing not planned",
			Cq7Answer:      "Commercial validation not planned",
			Cq8Answer:      "Market analysis not started",
			Cq9Answer:      "Technology transfer not discussed",
		},
	}

	seedAssessmentTrl(assessmentTrlRepo, assessments)

	// 7. Seed Intellectual Property
	fmt.Println("\n🏛️ Seeding Intellectual Property...")
	intellectualProperties := []models.IntellectualProperty{
		{
			CaseID:             "CS-00001",
			IPTypes:            "Patent",
			IPProtectionStatus: "Application Filed",
			IPRequestNumber:    "US2024001234A1",
		},
		{
			CaseID:             "CS-00002",
			IPTypes:            "Patent",
			IPProtectionStatus: "Granted",
			IPRequestNumber:    "US2023005678B2",
		},
		{
			CaseID:             "CS-00003",
			IPTypes:            "Trademark",
			IPProtectionStatus: "Application Filed",
			IPRequestNumber:    "TM2024009012",
		},
		{
			CaseID:             "CS-00004",
			IPTypes:            "Patent",
			IPProtectionStatus: "Under Review",
			IPRequestNumber:    "US2024003456A1",
		},
		{
			CaseID:             "CS-00005",
			IPTypes:            "Copyright",
			IPProtectionStatus: "Registered",
			IPRequestNumber:    "CR2024007890",
		},
	}

	seedIntellectualProperty(intellectualPropertyRepo, intellectualProperties)

	// 8. Seed Supporters
	fmt.Println("\n🤝 Seeding Supporters...")
	supporters := []models.Supporter{
		{
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
			NeedGuidelines:                  false,
			NeedCertification:               true,
			NeedAccount:                     false,
			Need:                            "Technical expertise and testing facilities",
			AdditionalDocuments:             "Technical specifications and test results",
		},
		{
			CaseID:                          "CS-00002",
			SupportResearch:                 true,
			SupportVDC:                      true,
			SupportSiEIC:                    true,
			NeedProtectIntellectualProperty: true,
			NeedCoDevelopers:                true,
			NeedActivities:                  true,
			NeedTest:                        true,
			NeedCapital:                     true,
			NeedPartners:                    true,
			NeedGuidelines:                  true,
			NeedCertification:               true,
			NeedAccount:                     true,
			Need:                            "Comprehensive support for commercialization",
			AdditionalDocuments:             "Business plan and market analysis",
		},
		{
			CaseID:                          "CS-00003",
			SupportResearch:                 true,
			SupportVDC:                      false,
			SupportSiEIC:                    false,
			NeedProtectIntellectualProperty: false,
			NeedCoDevelopers:                true,
			NeedActivities:                  false,
			NeedTest:                        true,
			NeedCapital:                     false,
			NeedPartners:                    false,
			NeedGuidelines:                  true,
			NeedCertification:               false,
			NeedAccount:                     false,
			Need:                            "Technical development support",
			AdditionalDocuments:             "Development roadmap",
		},
		{
			CaseID:                          "CS-00004",
			SupportResearch:                 true,
			SupportVDC:                      true,
			SupportSiEIC:                    true,
			NeedProtectIntellectualProperty: true,
			NeedCoDevelopers:                false,
			NeedActivities:                  true,
			NeedTest:                        false,
			NeedCapital:                     true,
			NeedPartners:                    true,
			NeedGuidelines:                  false,
			NeedCertification:               true,
			NeedAccount:                     true,
			Need:                            "Market entry and funding support",
			AdditionalDocuments:             "Financial projections and investor materials",
		},
		{
			CaseID:                          "CS-00005",
			SupportResearch:                 true,
			SupportVDC:                      false,
			SupportSiEIC:                    false,
			NeedProtectIntellectualProperty: false,
			NeedCoDevelopers:                true,
			NeedActivities:                  true,
			NeedTest:                        false,
			NeedCapital:                     false,
			NeedPartners:                    true,
			NeedGuidelines:                  true,
			NeedCertification:               false,
			NeedAccount:                     false,
			Need:                            "Research collaboration and guidance",
			AdditionalDocuments:             "Research proposal and literature review",
		},
	}

	seedSupporters(supporterRepo, supporters)

	fmt.Println("\n🎉 Comprehensive data seeding completed successfully!")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("📊 Summary:")
	fmt.Println("   👥 Admins: 5 records")
	fmt.Println("   🎯 Coordinators: 5 records")
	fmt.Println("   🔬 Researchers: 5 records")
	fmt.Println("   📋 Cases: 5 records")
	fmt.Println("   📅 Appointments: 5 records")
	fmt.Println("   📊 Assessment TRL: 5 records")
	fmt.Println("   🏛️ Intellectual Property: 5 records")
	fmt.Println("   🤝 Supporters: 5 records")
	fmt.Println("   📝 Total: 40 records created")
}

// Helper functions for seeding each entity type

func seedAdmins(repo *repository.AdminRepo, admins []models.AdminInfo) {
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
		err = repo.CreateAdmin(&admin)
		if err != nil {
			log.Printf("❌ Error creating admin %s: %v", admin.AdminEmail, err)
			continue
		}

		fmt.Printf("✅ Admin %d: %s %s %s (%s)\n", i+1, admin.AdminPrefix, admin.AdminFirstName, admin.AdminLastName, admin.AdminID)
	}
}

func seedCoordinators(repo *repository.CoordinatorRepo, coordinators []models.CoordinatorInfo) {
	for i, coordinator := range coordinators {
		// Set timestamps
		now := time.Now()
		coordinator.CreatedAt = now
		coordinator.UpdatedAt = now

		// Create coordinator
		err := repo.CreateCoordinator(&coordinator)
		if err != nil {
			log.Printf("❌ Error creating coordinator %s: %v", coordinator.CoordinatorEmail, err)
			continue
		}

		fmt.Printf("✅ Coordinator %d: %s (%s)\n", i+1, coordinator.CoordinatorName, coordinator.CoordinatorID)
	}
}

func seedResearchers(repo *repository.ResearcherRepo, researchers []models.ResearcherInfo) {
	for i, researcher := range researchers {
		// Set timestamps
		now := time.Now()
		researcher.CreatedAt = now
		researcher.UpdatedAt = now

		// Create researcher
		err := repo.CreateResearcher(&researcher)
		if err != nil {
			log.Printf("❌ Error creating researcher %s: %v", researcher.ResearcherEmail, err)
			continue
		}

		fmt.Printf("✅ Researcher %d: %s %s %s (%s)\n", i+1, researcher.ResearcherPrefix, researcher.ResearcherFirstName, researcher.ResearcherLastName, researcher.ResearcherID)
	}
}

func seedCases(repo *repository.CaseRepo, cases []models.CaseInfo) {
	for i, caseInfo := range cases {
		// Set timestamps
		now := time.Now()
		caseInfo.CreatedAt = now
		caseInfo.UpdatedAt = now

		// Create case
		err := repo.CreateCase(&caseInfo)
		if err != nil {
			log.Printf("❌ Error creating case %s: %v", caseInfo.CaseTitle, err)
			continue
		}

		fmt.Printf("✅ Case %d: %s (%s)\n", i+1, caseInfo.CaseTitle, caseInfo.CaseID)
	}
}

func seedAppointments(repo *repository.AppointmentRepo, appointments []models.Appointment) {
	for i, appointment := range appointments {
		// Set timestamps
		now := time.Now()
		appointment.CreatedAt = now
		appointment.UpdatedAt = now

		// Create appointment
		err := repo.CreateAppointment(&appointment)
		if err != nil {
			log.Printf("❌ Error creating appointment for case %s: %v", appointment.CaseID, err)
			continue
		}

		fmt.Printf("✅ Appointment %d: %s - %s (%s)\n", i+1, appointment.Date.Format("2006-01-02"), appointment.Status, appointment.AppointmentID)
	}
}

func seedAssessmentTrl(repo *repository.AssessmentTrlRepo, assessments []models.AssessmentTrl) {
	for i, assessment := range assessments {
		// Set timestamps
		now := time.Now()
		assessment.CreatedAt = now
		assessment.UpdatedAt = now

		// Create assessment
		err := repo.CreateAssessmentTrl(&assessment)
		if err != nil {
			log.Printf("❌ Error creating assessment for case %s: %v", assessment.CaseID, err)
			continue
		}

		fmt.Printf("✅ Assessment %d: TRL Level %d (%s)\n", i+1, assessment.TrlLevelResult, assessment.ID)
	}
}

func seedIntellectualProperty(repo *repository.IntellectualPropertyRepo, ips []models.IntellectualProperty) {
	for i, ip := range ips {
		// Set timestamps
		now := time.Now()
		ip.CreatedAt = now
		ip.UpdatedAt = now

		// Create intellectual property
		err := repo.CreateIP(&ip) // Changed from CreateIntellectualProperty to CreateIP
		if err != nil {
			log.Printf("❌ Error creating IP for case %s: %v", ip.CaseID, err)
			continue
		}

		fmt.Printf("✅ IP %d: %s - %s (%s)\n", i+1, ip.IPTypes, ip.IPProtectionStatus, ip.ID)
	}
}

func seedSupporters(repo *repository.SupporterRepo, supporters []models.Supporter) {
	for i, supporter := range supporters {
		// Set timestamps
		now := time.Now()
		supporter.CreatedAt = now
		supporter.UpdatedAt = now

		// Create supporter
		err := repo.CreateSupporter(&supporter)
		if err != nil {
			log.Printf("❌ Error creating supporter for case %s: %v", supporter.CaseID, err)
			continue
		}

		fmt.Printf("✅ Supporter %d: %s (%s)\n", i+1, supporter.Need, supporter.SupporterID)
	}
}
