package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/iterator"

	"trl-research-backend/internal/database"
	"trl-research-backend/internal/models"
)

// 🧹 clearCollection deletes all documents in the given Firestore collection
func clearCollection(ctx context.Context, client *firestore.Client, collection string) error {
	iter := client.Collection(collection).Documents(ctx)
	batch := client.Batch()
	count := 0

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		batch.Delete(doc.Ref)
		count++

		// Firestore batch limit = 500
		if count%400 == 0 {
			_, err := batch.Commit(ctx)
			if err != nil {
				return err
			}
			batch = client.Batch()
		}
	}

	if count > 0 {
		_, err := batch.Commit(ctx)
		if err != nil {
			return err
		}
	}

	fmt.Printf("🧹 Cleared %d documents from %s\n", count, collection)
	return nil
}

func main() {
	// 🔥 Initialize Firebase
	database.InitFirebase("trl-research-service-account.json")
	defer database.CloseFirebase()
	ctx := context.Background()
	client := database.FirestoreClient

	fmt.Println("🌱 Starting Firestore seeding process...")
	fmt.Println(strings.Repeat("=", 60))

	// =============================
	// 🧹 Clear all existing data first
	// =============================
	collections := []string{
		"admin_info",
		"researchers",
		"coordinators",
		"cases",
		"appointments",
		"assessment_trl",
		"intellectual_properties",
		"supporters",
	}

	for _, col := range collections {
		if err := clearCollection(ctx, client, col); err != nil {
			log.Fatalf("❌ Failed to clear %s: %v\n", col, err)
		} else {
			fmt.Printf("✅ Cleared collection: %s\n", col)
		}
	}

	now := time.Now()

	// =============================
	// 1️⃣ Admins
	// =============================
	admins := []models.AdminInfo{
		{"A-00001", "Dr.", "Assistant Professor", "Ann", "Smith", "Computer Science", "+66-81-234-5678", "admin1@example.com", "password123", "CS-00001", now, now},
		{"A-00002", "Prof.", "Professor", "John", "Doe", "Information Tech", "+66-82-234-5678", "admin2@example.com", "password123", "CS-00002", now, now},
		{"A-00003", "Dr.", "Lecturer", "May", "Tan", "AI Research", "+66-83-234-5678", "admin3@example.com", "password123", "CS-00003", now, now},
		{"A-00004", "Dr.", "Assistant Professor", "Nina", "Park", "Robotics", "+66-84-234-5678", "admin4@example.com", "password123", "CS-00004", now, now},
		{"A-00005", "Prof.", "Dean", "Tom", "Lee", "Innovation", "+66-85-234-5678", "admin5@example.com", "password123", "CS-00005", now, now},
	}

	for _, admin := range admins {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(admin.AdminPassword), bcrypt.DefaultCost)
		admin.AdminPassword = string(hashed)
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
		{"RS-00001", "A-00001", "Dr.", "Research Fellow", "Pair", "Brown", "Software Engineering", "+66-83-111-2222", "researcher1@example.com", "password123", now, now},
		{"RS-00002", "A-00002", "Dr.", "Postdoc", "Kate", "Miller", "Bioinformatics", "+66-84-222-3333", "researcher2@example.com", "password123", now, now},
		{"RS-00003", "A-00003", "Mr.", "Assistant", "Jay", "Wong", "Electronics", "+66-85-333-4444", "researcher3@example.com", "password123", now, now},
		{"RS-00004", "A-00004", "Ms.", "Analyst", "Sue", "Kim", "Chemical", "+66-86-444-5555", "researcher4@example.com", "password123", now, now},
		{"RS-00005", "A-00005", "Dr.", "Scientist", "Beam", "Chan", "AI Systems", "+66-87-555-6666", "researcher5@example.com", "password123", now, now},
	}

	for _, r := range researchers {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(r.ResearcherPassword), bcrypt.DefaultCost)
		r.ResearcherPassword = string(hashed)
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
		{
			CoordinatorID:    "C-00002",
			CoordinatorEmail: "coordinator2@university.edu",
			CoordinatorName:  "Dr. Laila Wong",
			CoordinatorPhone: "+66-92-111-1111",
			Department:       "Innovation Center",
			CaseID:           "CS-00002",
		},
		{
			CoordinatorID:    "C-00003",
			CoordinatorEmail: "coordinator3@university.edu",
			CoordinatorName:  "Dr. Tanawat Lee",
			CoordinatorPhone: "+66-93-111-1111",
			Department:       "AI Lab",
			CaseID:           "CS-00003",
		},
		{
			CoordinatorID:    "C-00004",
			CoordinatorEmail: "coordinator4@university.edu",
			CoordinatorName:  "Dr. Pailin Sae",
			CoordinatorPhone: "+66-94-111-1111",
			Department:       "R&D Hub",
			CaseID:           "CS-00004",
		},
		{
			CoordinatorID:    "C-00005",
			CoordinatorEmail: "coordinator5@university.edu",
			CoordinatorName:  "Dr. Min Cho",
			CoordinatorPhone: "+66-95-111-1111",
			Department:       "Research Admin",
			CaseID:           "CS-00005",
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
		// ✅ 2 approved
		{
			CaseID:           "CS-00001",
			CoordinatorEmail: "coordinator1@university.edu",
			TrlScore:         "5",
			TrlSuggestion:    "Excellent progress",
			Status:           true,
			IsUrgent:         false,
			UrgentReason:     "",
			UrgentFeedback:   "",
			CaseTitle:        "AI-powered Diagnosis",
			CaseType:         "Software",
			CaseDescription:  "ML model for early disease detection",
			CaseKeywords:     "AI, ML, Medical",
			ResearcherID:     "RS-00001",
		},
		{
			CaseID:           "CS-00002",
			CoordinatorEmail: "coordinator2@university.edu",
			TrlScore:         "4",
			TrlSuggestion:    "Ready for pilot testing",
			Status:           true,
			IsUrgent:         false,
			UrgentReason:     "",
			UrgentFeedback:   "",
			CaseTitle:        "Robotics Arm Control",
			CaseType:         "Hardware",
			CaseDescription:  "Design for precise robot movement",
			CaseKeywords:     "Robot, Control, Sensor",
			ResearcherID:     "RS-00002",
		},

		// 🕓 3 in process
		{
			CaseID:           "CS-00003",
			CoordinatorEmail: "coordinator3@university.edu",
			TrlScore:         "2",
			TrlSuggestion:    "Need prototype validation",
			Status:           false,
			IsUrgent:         false,
			UrgentReason:     "",
			UrgentFeedback:   "",
			CaseTitle:        "Smart Irrigation",
			CaseType:         "IoT",
			CaseDescription:  "Water system for agriculture",
			CaseKeywords:     "IoT, Sensor",
			ResearcherID:     "RS-00003",
		},
		{
			CaseID:           "CS-00004",
			CoordinatorEmail: "coordinator4@university.edu",
			TrlScore:         "3",
			TrlSuggestion:    "Improve prototype stability",
			Status:           false,
			IsUrgent:         true,
			UrgentReason:     "ไม่ urgent",
			UrgentFeedback:   "",
			CaseTitle:        "Nanotech Coating",
			CaseType:         "Material",
			CaseDescription:  "Durable coating for surfaces",
			CaseKeywords:     "Nano, Surface",
			ResearcherID:     "RS-00004",
		},
		{
			CaseID:           "CS-00005",
			CoordinatorEmail: "coordinator5@university.edu",
			TrlScore:         "1",
			TrlSuggestion:    "In concept phase",
			Status:           false,
			IsUrgent:         false,
			UrgentReason:     "",
			UrgentFeedback:   "",
			CaseTitle:        "Green Battery",
			CaseType:         "Energy",
			CaseDescription:  "New eco battery",
			CaseKeywords:     "Energy, Battery",
			ResearcherID:     "RS-00005",
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
	// 5️⃣ Appointments (case1 has 2)
	// =============================
	appointments := []models.Appointment{
		{"AP-00001", "CS-00001", now.AddDate(0, 0, 7), "attended", "Conference Room A", "Discuss progress", "Kickoff meeting", now, now},
		{"AP-00002", "CS-00001", now.AddDate(0, 0, 14), "absent", "Conference Room A", "Follow-up", "Researcher sick", now, now},
		{"AP-00003", "CS-00002", now.AddDate(0, 0, 10), "pending", "Conference Room B", "Prototype review", "Awaiting confirmation", now, now},
		{"AP-00004", "CS-00003", now.AddDate(0, 0, 12), "attended", "Meeting Room 2", "Test field setup", "Completed", now, now},
		{"AP-00005", "CS-00004", now.AddDate(0, 0, 20), "pending", "Zoom", "Online sync", "Progress update", now, now},
	}
	for _, a := range appointments {
		docRef := client.Collection("appointments").Doc(a.AppointmentID)
		_, err := docRef.Set(ctx, a)
		if err != nil {
			log.Printf("❌ Failed to seed appointment %v\n", err)
		} else {
			fmt.Printf("✅ Appointment seeded: %s\n", a.AppointmentID)
		}
	}

	// =============================
	// 6️⃣ Assessment TRL (correct CQ answers from checkboxQuestionList)
	// =============================

	// 📋 Master question list (same as your React checkboxQuestionList)
	checkboxQuestionList := [][]string{
		{
			"สมมุติฐานมีทฤษฎีทางวิทยาศาสตร์หรือคณิตศาสตร์รองรับ",
			"สมมุติฐานเป็นไปตามงานวิจัยที่เกี่ยวข้อง",
			"ผู้วิจัยมีการพัฒนาแนวคิดหรือสมการเพื่อสนับสนุนสมมุติฐาน",
		},
		{
			"สมมุติฐานผ่านการตรวจสอบโดยผู้เชี่ยวชาญ และยืนยันหลักการทางวิทยาศาสตร์พื้นฐาน",
			"สมมุติฐานแสดงแนวทางที่เป็นไปได้พร้อม ระบุส่วนประกอบสำคัญของเทคโนโลยี",
			"สมมุติฐานมีการประเมินหรือคาดการณ์ประสิทธิภาพเบื้องต้นขององค์ประกอบหลัก",
			"มีการศึกษาเบื้องต้นยืนยันความเป็นไปได้ของการจำลอง กระบวนการอย่างง่าย (การศึกษาโดยไม่มีการทดลองในห้องปฏิบัติการ",
			"สมมุติฐานมีการทดสอบแนวคิด (Proof of Concept) ด้วยข้อมูลสังเคราะห์",
		},
		{
			"สมมุติฐานถูกพิสูจน์ด้วยการทดลองเบื้องต้นแล้ว",
			"การทดลองสามารถคาดการณ์ของส่วนประกอบเทคโนโลยีได้",
			"มีการสร้างตัวชี้วัดประสิทธิภาพเทคโนโลยีหรือระบบ",
			"มีข้อเท็จจริงวิทยาศาสตร์ที่เกี่ยวข้องกับการพัฒนาเทคโนโลยีที่สามารถจำลองทำซ้ำได้",
			"มีการยืนยันคุณสมบัติและประสิทธิภาพของเทคโนโลยีหรือระบบด้วยสมการ หรือตัวแปร",
			"มีหลักฐานงานวิจัยที่เผยแพร่แล้วว่าการรวมเทคโนโลยีและส่วนประกอบของระบบประสบความสำเร็จ",
			"มีการระบุความเสี่ยงและมีการบริหารความเสี่ยงสำหรับงานวิจัย",
		},
		{
			"มีการสรุปและจัดทำข้อกำหนดของระบบ/การออกแบบ โดยอ้างอิงจากความต้องการจริง",
			"มีการระบุวัสดุ กระบวนการ และเทคนิคที่เกี่ยวข้อง",
			"มีต้นแบบเทคโนโลยีที่ปรับขนาดได้",
			"มีการทดสอบและแสดงประสิทธิภาพของส่วนประกอบและต้นแบบในห้องปฏิบัติการ",
			"มีการจำลองและตรวจสอบความเป็นไปได้ของกระบวนการ",
			"มีส่วนประกอบของระบบครบถ้วนและเพียงพอ",
			"มีการเริ่มศึกษาบูรณาการกับการใช้งานอื่น",
			"มีการระบุปัจจัยต้นทุน",
			"มีการริเริ่มโปรแกรมการจัดการความเสี่ยงอย่างเป็นทางการและบูรณาการกับการจัดการโครงการ",
		},
		{
			"ต้นแบบถูกพัฒนาและทำงานได้จริง โดยมีการรวมโมดูล/ฟังก์ชันสำคัญ และทดสอบการทำงานภายใต้สภาวะที่ใกล้เคียงหรือเป็นจริง",
			"ส่วนประกอบและส่วนต่อประสานของระบบได้รับการกำหนด ตรวจสอบ และรับรองตามมาตรฐานที่ยอมรับได้",
			"มีการวัดผลกระบวนการที่เที่ยงตรง",
			"มีการระบุปัญหาและประเมินความน่าเชื่อถือด้านคุณภาพ",
			"มีการสรุปกระบวนการออกแบบสำหรับการใช้งานจริง",
			"มีการจัดทำและดำเนินการตามแผนบริหารความเสี่ยง",
		},
		{
			"มีการทดสอบและสาธิตต้นแบบในสภาพแวดล้อมที่เกี่ยวข้อง/จำลองจริง พร้อมการยืนยันคุณสมบัติทางวิศวกรรมและประสิทธิภาพของระบบ",
			"ส่วนประกอบของสินค้าหรือบริการต้นแบบนั้นสามารถทำงานร่วมกันได้ในการทดสอบการแก้ปัญหาจริง",
			"มีการจัดเตรียมวัสดุ/อุปกรณ์ภายนอกครบถ้วน",
			"มีการรวบรวมข้อมูลด้านการบำรุงรักษาและระบบสนับสนุนที่เชื่อถือได้",
		},
		{
			"มีการทดสอบและตรวจสอบการปฏิบัติงานของอุปกรณ์/กระบวนการในสภาวะจริง เพื่อหาข้อจำกัด จุดบกพร่อง และยืนยันความถูกต้องกับระบบที่ใช้งานอยู่",
			"มีต้นแบบและส่วนประกอบที่ใกล้เคียงของจริง แสดงให้เห็นถึงความพอดีและฟังก์ชันการทำงานที่สอดคล้องกับการผลิตจริง",
			"มีข้อมูลสนับสนุนด้านความน่าเชื่อถือ การบำรุงรักษา",
			"มีอุปกรณ์และวัสดุที่ใช้ได้จริงในกระบวนการผลิต",
		},
		{
			"ทุกองค์ประกอบของเทคโนโลยี/ระบบมีความพอดี ฟังก์ชันเข้ากันได้ และเหมาะสมกับสภาพแวดล้อมการทำงานจริง",
			"วัสดุทั้งหมดในการผลิตและพร้อมใช้งาน",
			"มีข้อมูลและเอกสารการบำรุงรักษา/การสนับสนุนที่สมบูรณ์และอยู่ภายใต้การควบคุมการกำหนดค่า",
		},
		{
			"เทคโนโลยี/ระบบทำงานได้ตามที่กำหนดในเอกสารแนวคิด มีการนำไปปรับใช้ในสภาพแวดล้อมจริง และแสดงศักยภาพได้อย่างสมบูรณ์",
			"มีการทดสอบและประเมินผลการปฏิบัติงานสำเร็จแล้วและจัดทำเป็นเอกสาร",
			"มีการออกแบบโดยคำนึงถึงเป้าหมายด้านต้นทุน",
			"มีการระบุและบรรเทาความเสี่ยงด้านความปลอดภัยและผลข้างเคียง",
		},
	}

	// ✅ helper to select a random subset of question labels
	pickRandomSubset := func(options []string) []string {
		if len(options) == 0 {
			return []string{}
		}
		count := rand.Intn(len(options)) + 1 // at least 1
		rand.Shuffle(len(options), func(i, j int) { options[i], options[j] = options[j], options[i] })
		return options[:count]
	}

	// ✅ random true/false answers
	trueFalseSet := [][]bool{
		{true, true, false, false, true, false, true},
		{true, false, true, true, false, true, false},
		{false, true, true, false, false, true, true},
		{true, true, true, true, true, true, true},
		{false, false, true, true, false, false, true},
	}

	// ✅ create and insert AssessmentTRL docs
	for i := 1; i <= 5; i++ {
		a := models.AssessmentTrl{
			ID:             fmt.Sprintf("AT-0000%d", i),
			CaseID:         fmt.Sprintf("CS-0000%d", i),
			TrlLevelResult: i,
			Rq1Answer:      trueFalseSet[i-1][0],
			Rq2Answer:      trueFalseSet[i-1][1],
			Rq3Answer:      trueFalseSet[i-1][2],
			Rq4Answer:      trueFalseSet[i-1][3],
			Rq5Answer:      trueFalseSet[i-1][4],
			Rq6Answer:      trueFalseSet[i-1][5],
			Rq7Answer:      trueFalseSet[i-1][6],
			Cq1Answer:      pickRandomSubset(checkboxQuestionList[0]),
			Cq2Answer:      pickRandomSubset(checkboxQuestionList[1]),
			Cq3Answer:      pickRandomSubset(checkboxQuestionList[2]),
			Cq4Answer:      pickRandomSubset(checkboxQuestionList[3]),
			Cq5Answer:      pickRandomSubset(checkboxQuestionList[4]),
			Cq6Answer:      pickRandomSubset(checkboxQuestionList[5]),
			Cq7Answer:      pickRandomSubset(checkboxQuestionList[6]),
			Cq8Answer:      pickRandomSubset(checkboxQuestionList[7]),
			Cq9Answer:      pickRandomSubset(checkboxQuestionList[8]),
			CreatedAt:      now,
			UpdatedAt:      now,
		}

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
	ipTypes := []string{"สิทธิบัตร", "อนุสิทธิบัตร", "สิทธิบัตรออกแบบผลิตภัณฑ์", "ลิขสิทธิ์", "เครื่องหมายการค้า", "ความลับทางการค้า"}

	for i := 1; i <= 5; i++ {
		ip := models.IntellectualProperty{
			ID:                 fmt.Sprintf("IP-0000%d", i),
			CaseID:             fmt.Sprintf("CS-0000%d", i),
			IPTypes:            ipTypes[i-1],
			IPProtectionStatus: "Application Filed",
			IPRequestNumber:    fmt.Sprintf("TH2025%04dA1", i),
			CreatedAt:          now,
			UpdatedAt:          now,
		}
		docRef := client.Collection("intellectual_properties").Doc(ip.ID)
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
	for i := 1; i <= 5; i++ {
		s := models.Supporter{
			SupporterID:                     fmt.Sprintf("SP-0000%d", i),
			CaseID:                          fmt.Sprintf("CS-0000%d", i),
			SupportResearch:                 i%2 == 0,
			SupportVDC:                      i%3 == 0,
			SupportSiEIC:                    true,
			NeedProtectIntellectualProperty: i%2 != 0,
			NeedCoDevelopers:                false,
			NeedActivities:                  true,
			NeedTest:                        true,
			NeedCapital:                     i%2 == 0,
			NeedPartners:                    true,
			NeedGuidelines:                  true,
			NeedCertification:               false,
			NeedAccount:                     false,
			Need:                            "Require collaboration and mentorship",
			AdditionalDocuments:             "Project plan and reference materials",
			CreatedAt:                       now,
			UpdatedAt:                       now,
		}
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
