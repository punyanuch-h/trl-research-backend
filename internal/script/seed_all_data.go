package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"

	"trl-research-backend/internal/config"
	"trl-research-backend/internal/database"
	"trl-research-backend/internal/models"
)

func main() {
	// Load config
	config.LoadConfig()

	// 🔥 Initialize Postgres
	database.InitPostgres()
	defer database.ClosePostgres()

	db := database.DB

	fmt.Println("🌱 Starting Postgres seeding process...")
	fmt.Println(strings.Repeat("=", 60))

	// =============================
	// 🧹 Clear all existing data first (Optional, but good for seeding)
	// =============================
	db.Exec("TRUNCATE TABLE admins, researchers, coordinators, cases, appointments, assessment_trls, intellectual_properties, supporters, file_metadatas RESTART IDENTITY CASCADE")

	now := time.Now()

	// =============================
	// 1️⃣ Admins
	// =============================
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	admins := []models.AdminInfo{
		{AdminID: "AD-00001", AdminPrefix: "ดร.", AdminAcademicPosition: "ผู้ช่วยศาสตราจารย์", AdminFirstName: "ทิพาจินต์", AdminLastName: "ไทยพิสุทธิกุล", AdminDepartment: "Computer Science", AdminPhoneNumber: "+66-81-234-5678", AdminEmail: "admin1@example.com", AdminPassword: string(hashedPassword), CaseID: "CS-00001", CreatedAt: now, UpdatedAt: now},
		{AdminID: "AD-00002", AdminPrefix: "ศ.ดร.", AdminAcademicPosition: "ศาสตราจารย์", AdminFirstName: "ปรีมน", AdminLastName: "ปุณณกิติเกษม", AdminDepartment: "AI Research", AdminPhoneNumber: "+66-82-234-5678", AdminEmail: "admin2@example.com", AdminPassword: string(hashedPassword), CaseID: "CS-00002", CreatedAt: now, UpdatedAt: now},
	}

	for _, admin := range admins {
		if err := db.Create(&admin).Error; err != nil {
			log.Printf("❌ Failed to seed admin %s: %v\n", admin.AdminEmail, err)
		} else {
			fmt.Printf("✅ Admin seeded: %s\n", admin.AdminEmail)
		}
	}

	// =============================
	// 2️⃣ Researchers
	// =============================
	researchers := []models.ResearcherInfo{
		{ResearcherID: "RS-00001", AdminID: "AD-00001", Prefix: "ดร.", AcademicPosition: "Research Fellow", FirstName: "ปุณยนุช", LastName: "หทัยวสีวงศ์", Department: "Software Engineering", PhoneNumber: "+66-83-111-2222", Email: "researcher1@example.com", Password: string(hashedPassword), CreatedAt: now, UpdatedAt: now},
		{ResearcherID: "RS-00002", AdminID: "AD-00002", Prefix: "ผศ. ดร.", AcademicPosition: "Postdoc", FirstName: "ศุภิสรา", LastName: "งามชัยพิสิฐ", Department: "Bioinformatics", PhoneNumber: "+66-84-222-3333", Email: "researcher2@example.com", Password: string(hashedPassword), CreatedAt: now, UpdatedAt: now},
		{ResearcherID: "RS-00003", AdminID: "AD-00001", Prefix: "นางสาว", AcademicPosition: "Assistant", FirstName: "สิทธิดา", LastName: "ศรีธนกฤตาธิการ", Department: "Electronics", PhoneNumber: "+66-85-333-4444", Email: "researcher3@example.com", Password: string(hashedPassword), CreatedAt: now, UpdatedAt: now},
	}

	for _, r := range researchers {
		if err := db.Create(&r).Error; err != nil {
			log.Printf("❌ Failed to seed researcher %s: %v\n", r.Email, err)
		} else {
			fmt.Printf("✅ Researcher seeded: %s\n", r.Email)
		}
	}

	// =============================
	// 3️⃣ Coordinators
	// =============================
	coordinators := []models.CoordinatorInfo{
		{CoordinatorID: "C-00001", CoordinatorEmail: "admin1@example.com", CoordinatorName: "ดร. ทิพาจินต์ ไทยพิสุทธิกุล", CoordinatorPhone: "+66-91-111-1111", Department: "Research Development", CaseID: "CS-00001", CreatedAt: now, UpdatedAt: now},
		{CoordinatorID: "C-00002", CoordinatorEmail: "admin2@example.com", CoordinatorName: "ศ.ดร. ปรีมน ปุณณกิติเกษม", CoordinatorPhone: "+66-92-111-1111", Department: "Innovation Center", CaseID: "CS-00002", CreatedAt: now, UpdatedAt: now},
		{CoordinatorID: "C-00003", CoordinatorEmail: "coordinator3@university.edu", CoordinatorName: "ดร. วราภรณ์ อินทรกุล", CoordinatorPhone: "+66-93-111-1111", Department: "AI Lab", CaseID: "CS-00003", CreatedAt: now, UpdatedAt: now},
		{CoordinatorID: "C-00004", CoordinatorEmail: "coordinator4@university.edu", CoordinatorName: "ดร. วราภรณ์ อินทรกุล", CoordinatorPhone: "+66-93-111-1111", Department: "AI Lab", CaseID: "CS-00004", CreatedAt: now, UpdatedAt: now},
		{CoordinatorID: "C-00005", CoordinatorEmail: "coordinator5@university.edu", CoordinatorName: "ดร. วราภรณ์ อินทรกุล", CoordinatorPhone: "+66-93-111-1111", Department: "AI Lab", CaseID: "CS-00005", CreatedAt: now, UpdatedAt: now},
	}
	for _, c := range coordinators {
		if err := db.Create(&c).Error; err != nil {
			log.Printf("❌ Failed to seed coordinator %v\n", err)
		} else {
			fmt.Printf("✅ Coordinator seeded: %s\n", c.CoordinatorEmail)
		}
	}

	// =============================
	// 4️⃣ Cases
	// =============================
	emptyJSON, _ := json.Marshal([]string{})
	cases := []models.CaseInfo{
		{CaseID: "CS-00001", CoordinatorEmail: "admin1@example.com", TrlScore: "5", TrlSuggestion: "Excellent progress", Status: "approved", IsUrgent: false, CaseTitle: "AI-powered Diagnosis", CaseType: "Software", CaseDescription: "ML model for early disease detection", CaseKeywords: "AI, ML, Medical", ResearcherID: "RS-00001", CaseAttachments: datatypes.JSON(emptyJSON), CreatedAt: now, UpdatedAt: now},
		{CaseID: "CS-00002", CoordinatorEmail: "admin2@example.com", TrlScore: "6", TrlSuggestion: "Ready for pilot testing", Status: "approved", IsUrgent: false, CaseTitle: "Robotics Arm Control", CaseType: "Hardware", CaseDescription: "Design for precise robot movement", CaseKeywords: "Robot, Control, Sensor", ResearcherID: "RS-00002", CaseAttachments: datatypes.JSON(emptyJSON), CreatedAt: now, UpdatedAt: now},
		{CaseID: "CS-00003", CoordinatorEmail: "coordinator3@university.edu", TrlScore: "2", TrlSuggestion: "Need prototype validation", Status: "pending", IsUrgent: false, CaseTitle: "Smart Irrigation", CaseType: "IoT", CaseDescription: "Water system for agriculture", CaseKeywords: "IoT, Sensor", ResearcherID: "RS-00003", CaseAttachments: datatypes.JSON(emptyJSON), CreatedAt: now, UpdatedAt: now},
		{CaseID: "CS-00004", CoordinatorEmail: "coordinator3@university.edu", TrlScore: "7", TrlSuggestion: "Improve prototype stability", Status: "pending", IsUrgent: true, UrgentReason: "ต้องการขอทุนภายในเดือนมิถุนายน 2026", CaseTitle: "Nanotech Coating", CaseType: "Material", CaseDescription: "Durable coating for surfaces", CaseKeywords: "Nano, Surface", ResearcherID: "RS-00002", CaseAttachments: datatypes.JSON(emptyJSON), CreatedAt: now, UpdatedAt: now},
		{CaseID: "CS-00005", CoordinatorEmail: "admin2@example.com", TrlScore: "1", TrlSuggestion: "In concept phase", Status: "pending", IsUrgent: false, CaseTitle: "Green Battery", CaseType: "Energy", CaseDescription: "New eco battery", CaseKeywords: "Energy, Battery", ResearcherID: "RS-00003", CaseAttachments: datatypes.JSON(emptyJSON), CreatedAt: now, UpdatedAt: now},
	}
	for _, c := range cases {
		if err := db.Create(&c).Error; err != nil {
			log.Printf("❌ Failed to seed case %v\n", err)
		} else {
			fmt.Printf("✅ Case seeded: %s\n", c.CaseID)
		}
	}

	// =============================
	// 5️⃣ Appointments
	// =============================
	appointments := []models.Appointment{
		{AppointmentID: "AP-00001", CaseID: "CS-00001", Date: now.AddDate(0, 0, 7), Status: "attended", Location: "Conference Room A", Note: "Discuss progress", Summary: "Kickoff meeting", CreatedAt: now, UpdatedAt: now},
		{AppointmentID: "AP-00002", CaseID: "CS-00001", Date: now.AddDate(0, 0, 14), Status: "absent", Location: "Conference Room A", Note: "Follow-up", Summary: "Researcher sick", CreatedAt: now, UpdatedAt: now},
		{AppointmentID: "AP-00003", CaseID: "CS-00002", Date: now.AddDate(0, 0, 10), Status: "pending", Location: "Conference Room B", Note: "Prototype review", Summary: "Awaiting confirmation", CreatedAt: now, UpdatedAt: now},
		{AppointmentID: "AP-00004", CaseID: "CS-00003", Date: now.AddDate(0, 0, 12), Status: "attended", Location: "Meeting Room 2", Note: "Test field setup", Summary: "Completed", CreatedAt: now, UpdatedAt: now},
		{AppointmentID: "AP-00005", CaseID: "CS-00004", Date: now.AddDate(0, 0, 20), Status: "pending", Location: "Zoom", Note: "Online sync", Summary: "Progress update", CreatedAt: now, UpdatedAt: now},
	}
	for _, a := range appointments {
		if err := db.Create(&a).Error; err != nil {
			log.Printf("❌ Failed to seed appointment %v\n", err)
		} else {
			fmt.Printf("✅ Appointment seeded: %s\n", a.AppointmentID)
		}
	}

	// =============================
	// 6️⃣ Assessment TRL
	// =============================
	checkboxQuestionList := [][]string{
		{"สมมุติฐานมีทฤษฎีทางวิทยาศาสตร์หรือคณิตศาสตร์รองรับ", "สมมุติฐานเป็นไปตามงานวิจัยที่เกี่ยวข้อง", "ผู้วิจัยมีการพัฒนาแนวคิดหรือสมการเพื่อสนับสนุนสมมุติฐาน"},
		{"สมมุติฐานผ่านการตรวจสอบโดยผู้เชี่ยวชาญ และยืนยันหลักการทางวิทยาศาสตร์พื้นฐาน", "สมมุติฐานแสดงแนวทางที่เป็นไปได้พร้อม ระบุส่วนประกอบสำคัญของเทคโนโลยี", "สมมุติฐานมีการประเมินหรือคาดการณ์ประสิทธิภาพเบื้องต้นขององค์ประกอบหลัก", "มีการศึกษาเบื้องต้นยืนยันความเป็นไปได้ของการจำลอง กระบวนการอย่างง่าย (การศึกษาโดยไม่มีการทดลองในห้องปฏิบัติการ", "สมมุติฐานมีการทดสอบแนวคิด (Proof of Concept) ด้วยข้อมูลสังเคราะห์"},
	}

	pickRandomSubset := func(options []string) datatypes.JSON {
		if len(options) == 0 {
			res, _ := json.Marshal([]string{})
			return datatypes.JSON(res)
		}
		count := rand.Intn(len(options)) + 1
		rand.Shuffle(len(options), func(i, j int) { options[i], options[j] = options[j], options[i] })
		res, _ := json.Marshal(options[:count])
		return datatypes.JSON(res)
	}

	for i := 1; i <= 5; i++ {
		a := models.AssessmentTrl{
			ID:             fmt.Sprintf("AS-0000%d", i),
			CaseID:         fmt.Sprintf("CS-0000%d", i),
			TrlLevelResult: i,
			Rq1Answer:      true,
			Rq1Attachments: datatypes.JSON(emptyJSON),
			Rq2Answer:      true,
			Rq2Attachments: datatypes.JSON(emptyJSON),
			Cq1Answer:      pickRandomSubset(checkboxQuestionList[0]),
			Cq1Attachments: datatypes.JSON(emptyJSON),
			Cq2Answer:      pickRandomSubset(checkboxQuestionList[1]),
			Cq2Attachments: datatypes.JSON(emptyJSON),
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		if err := db.Create(&a).Error; err != nil {
			log.Printf("❌ Failed to seed assessment %v\n", err)
		} else {
			fmt.Printf("✅ Assessment TRL seeded for case: %s\n", a.CaseID)
		}
	}

	// =============================
	// 7️⃣ Intellectual Property
	// =============================
	ipTypes := []string{"สิทธิบัตร", "อนุสิทธิบัตร", "สิทธิบัตรออกแบบผลิตภัณฑ์", "ลิขสิทธิ์", "เครื่องหมายการค้า"}
	for i := 1; i <= 5; i++ {
		ip := models.IntellectualProperty{
			ID:                 fmt.Sprintf("IP-0000%d", i),
			CaseID:             fmt.Sprintf("CS-0000%d", i),
			IPTypes:            ipTypes[i-1],
			IPProtectionStatus: "Application Filed",
			IPRequestNumber:    fmt.Sprintf("TH2025%04dA1", i),
			IPAttachments:      datatypes.JSON(emptyJSON),
			CreatedAt:          now,
			UpdatedAt:          now,
		}
		if err := db.Create(&ip).Error; err != nil {
			log.Printf("❌ Failed to seed IP %v\n", err)
		} else {
			fmt.Printf("✅ Intellectual Property seeded for case: %s\n", ip.CaseID)
		}
	}

	// =============================
	// 8️⃣ Supporters
	// =============================
	for i := 1; i <= 5; i++ {
		s := models.Supporter{
			SupporterID:                     fmt.Sprintf("SP-0000%d", i),
			CaseID:                          fmt.Sprintf("CS-0000%d", i),
			SupportResearch:                 i%2 == 0,
			SupportSiEIC:                    true,
			NeedProtectIntellectualProperty: i%2 != 0,
			NeedActivities:                  true,
			NeedTest:                        true,
			NeedCapital:                     i%2 == 0,
			NeedPartners:                    true,
			Need:                            "Require collaboration and mentorship",
			CreatedAt:                       now,
			UpdatedAt:                       now,
		}
		if err := db.Create(&s).Error; err != nil {
			log.Printf("❌ Failed to seed supporter %v\n", err)
		} else {
			fmt.Printf("✅ Supporter seeded for case: %s\n", s.CaseID)
		}
	}

	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("🎉 Postgres seeding completed successfully!")
}
