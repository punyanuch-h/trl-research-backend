package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"

	"trl-research-backend/internal/config"
	"trl-research-backend/internal/database"
	"trl-research-backend/internal/models"
)

func ptrInt16(v int16) *int16 {
	return &v
}

func ptrString(s string) *string {
	return &s
}

func main() {
	// ==========================================
	// 0. INIT & CONFIG
	// ==========================================
	config.LoadConfig()
	database.InitPostgres()
	defer database.ClosePostgres()
	db := database.DB

	fmt.Println("🚀 Starting FINAL Seeding Process...")
	fmt.Println("   - Logic: RQ Decision Tree Applied")
	fmt.Println("   - Logic: Thai Prefixes & Academic Positions")
	fmt.Println("   - Logic: 5-Digit IDs")
	fmt.Println(strings.Repeat("=", 60))

	// ==========================================
	// 🧹 CLEANUP (Truncate all tables)
	// ==========================================
	tables := []string{
		"supportments", "intellectual_properties", "assessments",
		"appointments", "cases", "coordinators", "researchers", "admins",
	}
	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table)).Error; err != nil {
			log.Fatalf("❌ Failed to truncate %s: %v", table, err)
		}
	}

	// ==========================================
	// 🛠️ HELPERS & DATA
	// ==========================================
	now := time.Now()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		panic("failed to hash seed password: " + err.Error())
	}

	// Helper: Convert Struct/Slice -> JSONB
	toJSON := func(v interface{}) datatypes.JSON {
		b, err := json.Marshal(v)
		if err != nil {
			panic("failed to marshal JSON seed data: " + err.Error())
		}
		return datatypes.JSON(b)
	}
	emptyJSON := toJSON([]string{})

	checkboxData := [][]string{
		// Level 1 (Index 0)
		{"สมมุติฐานมีทฤษฎีทางวิทยาศาสตร์หรือคณิตศาสตร์รองรับ",
			"สมมุติฐานเป็นไปตามงานวิจัยที่เกี่ยวข้อง",
			"ผู้วิจัยมีการพัฒนาแนวคิดหรือสมการเพื่อสนับสนุนสมมุติฐาน"},
		// Level 2 (Index 1)
		{"สมมุติฐานผ่านการตรวจสอบโดยผู้เชี่ยวชาญ และยืนยันหลักการทางวิทยาศาสตร์พื้นฐาน",
			"สมมุติฐานแสดงแนวทางที่เป็นไปได้พร้อม ระบุส่วนประกอบสำคัญของเทคโนโลยี",
			"สมมุติฐานมีการประเมินหรือคาดการณ์ประสิทธิภาพเบื้องต้นขององค์ประกอบหลัก",
			"มีการศึกษาเบื้องต้นยืนยันความเป็นไปได้ของการจำลอง กระบวนการอย่างง่าย (การศึกษาโดยไม่มีการทดลองในห้องปฏิบัติการ",
			"สมมุติฐานมีการทดสอบแนวคิด (Proof of Concept) ด้วยข้อมูลสังเคราะห์"},
		// Level 3 (Index 2)
		{"สมมุติฐานถูกพิสูจน์ด้วยการทดลองเบื้องต้นแล้ว",
			"การทดลองสามารถคาดการณ์ของส่วนประกอบเทคโนโลยีได้",
			"มีการสร้างตัวชี้วัดประสิทธิภาพเทคโนโลยีหรือระบบ",
			"มีข้อเท็จจริงวิทยาศาสตร์ที่เกี่ยวข้องกับการพัฒนาเทคโนโลยีที่สามารถจำลองทำซ้ำได้",
			"มีการยืนยันคุณสมบัติและประสิทธิภาพของเทคโนโลยีหรือระบบด้วยสมการ หรือตัวแปร",
			"มีหลักฐานงานวิจัยที่เผยแพร่แล้วว่าการรวมเทคโนโลยีและส่วนประกอบของระบบประสบความสำเร็จ",
			"มีการระบุความเสี่ยงและมีการบริหารความเสี่ยงสำหรับงานวิจัย"},
		// Level 4 (Index 3)
		{"มีการสรุปและจัดทำข้อกำหนดของระบบ/การออกแบบ โดยอ้างอิงจากความต้องการจริง",
			"มีการระบุวัสดุ กระบวนการ และเทคนิคที่เกี่ยวข้อง",
			"มีต้นแบบเทคโนโลยีที่ปรับขนาดได้",
			"มีการทดสอบและแสดงประสิทธิภาพของส่วนประกอบและต้นแบบในห้องปฏิบัติการ",
			"มีการจำลองและตรวจสอบความเป็นไปได้ของกระบวนการ",
			"มีส่วนประกอบของระบบครบถ้วนและเพียงพอ",
			"มีการเริ่มศึกษาบูรณาการกับการใช้งานอื่น",
			"มีการระบุปัจจัยต้นทุน",
			"มีการริเริ่มโปรแกรมการจัดการความเสี่ยงอย่างเป็นทางการและบูรณาการกับการจัดการโครงการ"},
		// Level 5 (Index 4)
		{"ต้นแบบถูกพัฒนาและทำงานได้จริง โดยมีการรวมโมดูล/ฟังก์ชันสำคัญ และทดสอบการทำงานภายใต้สภาวะที่ใกล้เคียงหรือเป็นจริง",
			"ส่วนประกอบและส่วนต่อประสานของระบบได้รับการกำหนด ตรวจสอบ และรับรองตามมาตรฐานที่ยอมรับได้",
			"มีการวัดผลกระบวนการที่เที่ยงตรง",
			"มีการระบุปัญหาและประเมินความน่าเชื่อถือด้านคุณภาพ",
			"มีการสรุปกระบวนการออกแบบสำหรับการใช้งานจริง",
			"มีการจัดทำและดำเนินการตามแผนบริหารความเสี่ยง"},
		// Level 6 (Index 5)
		{"มีการทดสอบและสาธิตต้นแบบในสภาพแวดล้อมที่เกี่ยวข้อง/จำลองจริง พร้อมการยืนยันคุณสมบัติทางวิศวกรรมและประสิทธิภาพของระบบ",
			"ส่วนประกอบของสินค้าหรือบริการต้นแบบนั้นสามารถทำงานร่วมกันได้ในการทดสอบการแก้ปัญหาจริง",
			"มีการจัดเตรียมวัสดุ/อุปกรณ์ภายนอกครบถ้วน",
			"มีการรวบรวมข้อมูลด้านการบำรุงรักษาและระบบสนับสนุนที่เชื่อถือได้"},
		// Level 7 (Index 6)
		{"มีการทดสอบและตรวจสอบการปฏิบัติงานของอุปกรณ์/กระบวนการในสภาวะจริง เพื่อหาข้อจำกัด จุดบกพร่อง และยืนยันความถูกต้องกับระบบที่ใช้งานอยู่",
			"มีต้นแบบและส่วนประกอบที่ใกล้เคียงของจริง แสดงให้เห็นถึงความพอดีและฟังก์ชันการทำงานที่สอดคล้องกับการผลิตจริง",
			"มีข้อมูลสนับสนุนด้านความน่าเชื่อถือ การบำรุงรักษา",
			"มีอุปกรณ์และวัสดุที่ใช้ได้จริงในกระบวนการผลิต"},
		// Level 8 (Index 7)
		{"ทุกองค์ประกอบของเทคโนโลยี/ระบบมีความพอดี ฟังก์ชันเข้ากันได้ และเหมาะสมกับสภาพแวดล้อมการทำงานจริง",
			"วัสดุทั้งหมดในการผลิตและพร้อมใช้งาน",
			"มีข้อมูลและเอกสารการบำรุงรักษา/การสนับสนุนที่สมบูรณ์และอยู่ภายใต้การควบคุมการกำหนดค่า"},
		// Level 9 (Index 8)
		{"เทคโนโลยี/ระบบทำงานได้ตามที่กำหนดในเอกสารแนวคิด มีการนำไปปรับใช้ในสภาพแวดล้อมจริง และแสดงศักยภาพได้อย่างสมบูรณ์",
			"มีการทดสอบและประเมินผลการปฏิบัติงานสำเร็จแล้วและจัดทำเป็นเอกสาร",
			"มีการออกแบบโดยคำนึงถึงเป้าหมายด้านต้นทุน",
			"มีการระบุและบรรเทาความเสี่ยงด้านความปลอดภัยและผลข้างเคียง"},
	}

	// ==========================================
	// 👤 1. ADMINS
	// ==========================================
	admins := []models.Admins{
		{ID: "AD-00001", Prefix: "นาย", AcademicPosition: "ศ. ดร.", FirstName: "วิชัย", LastName: "การุณย์", Department: "Research Admin", PhoneNumber: "0811111111", Email: "admin.wichai@uni.ac.th", Password: string(hashedPassword)},
		{ID: "AD-00002", Prefix: "นาง", AcademicPosition: "รศ. ดร.", FirstName: "นภา", LastName: "สุขสวัสดิ์", Department: "Tech Transfer", PhoneNumber: "0811111112", Email: "admin.napa@uni.ac.th", Password: string(hashedPassword)},
		{ID: "AD-00003", Prefix: "นาย", AcademicPosition: "ผศ. ดร.", FirstName: "เกรียงไกร", LastName: "วิทย์", Department: "Patent Office", PhoneNumber: "0811111113", Email: "admin.kriangkrai@uni.ac.th", Password: string(hashedPassword)},
		{ID: "AD-00004", Prefix: "นางสาว", AcademicPosition: "ดร.", FirstName: "สิริ", LastName: "เพ็ญ", Department: "Grant Management", PhoneNumber: "0811111114", Email: "admin.siri@uni.ac.th", Password: string(hashedPassword)},
		{ID: "AD-00005", Prefix: "นาย", AcademicPosition: "อาจารย์", FirstName: "สมศักดิ์", LastName: "ไอที", Department: "System Admin", PhoneNumber: "0811111115", Email: "admin.somsak@uni.ac.th", Password: string(hashedPassword)},
	}
	if err := db.CreateInBatches(&admins, 5).Error; err != nil {
		log.Fatalf("❌ Failed to seed admins: %v", err)
	}
	fmt.Println("✅ Admins seeded")

	// ==========================================
	// 👨‍🔬 2. RESEARCHERS
	// ==========================================
	researchers := []models.Researchers{
		{ID: "RS-00001", Prefix: "นาย", AcademicPosition: "นพ.", FirstName: "สมชาย", LastName: "รักษ์ดี", Department: "Faculty of Medicine", PhoneNumber: "0822222221", Email: "somchai.med@uni.ac.th", Password: string(hashedPassword)},
		{ID: "RS-00002", Prefix: "นางสาว", AcademicPosition: "ผศ. ดร.", FirstName: "กานดา", LastName: "วิศวะ", Department: "Computer Eng.", PhoneNumber: "0822222222", Email: "kanda.eng@uni.ac.th", Password: string(hashedPassword)},
		{ID: "RS-00003", Prefix: "นาง", AcademicPosition: "รศ. ภญ.", FirstName: "มาลี", LastName: "เภสัช", Department: "Pharmacy", PhoneNumber: "0822222223", Email: "malee.pharm@uni.ac.th", Password: string(hashedPassword)},
		{ID: "RS-00004", Prefix: "นาย", AcademicPosition: "ดร.", FirstName: "ประทีป", LastName: "วัสดุ", Department: "Material Science", PhoneNumber: "0822222224", Email: "prateep.mat@uni.ac.th", Password: string(hashedPassword)},
		{ID: "RS-00005", Prefix: "นางสาว", AcademicPosition: "ผศ.", FirstName: "อริสรา", LastName: "อาหาร", Department: "Food Tech", PhoneNumber: "0822222225", Email: "arisara.food@uni.ac.th", Password: string(hashedPassword)},
	}
	db.CreateInBatches(&researchers, 5)
	fmt.Println("✅ Researchers seeded")

	// ==========================================
	// 🤝 3. COORDINATORS
	// ==========================================
	coordinators := []models.Coordinators{
		{ID: "CO-00001", Prefix: "นางสาว", AcademicPosition: "นพ.", FirstName: "สุดารัตน์", LastName: "ประสาน", Department: "Medical Center", PhoneNumber: "0833333331", Email: "coor.suda@uni.ac.th"},
		{ID: "CO-00002", Prefix: "นาย", AcademicPosition: "ดร.", FirstName: "ปิติ", LastName: "บริการ", Department: "Innovation Hub", PhoneNumber: "0833333332", Email: "coor.piti@uni.ac.th"},
		{ID: "CO-00003", Prefix: "นาง", AcademicPosition: "ผศ.", FirstName: "วรรณา", LastName: "ใจดี", Department: "Science Fac.", PhoneNumber: "0833333333", Email: "coor.wanna@uni.ac.th"},
		{ID: "CO-00004", Prefix: "นางสาว", AcademicPosition: "ผศ. ดร.", FirstName: "ดาริน", LastName: "ช่วยงาน", Department: "Agro-Industry", PhoneNumber: "0833333334", Email: "coor.darin@uni.ac.th"},
		{ID: "CO-00005", Prefix: "นาย", AcademicPosition: "รศ. ภญ.", FirstName: "เอก", LastName: "กฎหมาย", Department: "IP Dept", PhoneNumber: "0833333335", Email: "coor.ek@uni.ac.th"},
	}
	if err := db.CreateInBatches(&coordinators, 5).Error; err != nil {
		log.Fatalf("❌ Failed to seed coordinators: %v", err)
	}
	fmt.Println("✅ Coordinators seeded")

	// ==========================================
	// 📁 4. CASES
	// ==========================================

	// TRL Score logic applies to the specific Assessment scenario below
	cases := []models.Cases{
		// Case 1: TRL 8. Finished (Status=True). No Urgent. Has Admin.
		{
			ID: "CS-00001", ResearcherID: "RS-00001", CoordinatorID: "CO-00001",
			Title: "AI Screening for Retinopathy", Type: "Software", Description: "Deep learning model.", Keywords: "AI, Retina",
			TrlScore: ptrInt16(8), Status: true, IsUrgent: false,
			CreatedAt: now.AddDate(0, -6, 0), UpdatedAt: now,
		},
		// Case 2: TRL 4. In Progress. Urgent (Funding). Has Admin.
		{
			ID: "CS-00002", ResearcherID: "RS-00002", CoordinatorID: "CO-00002",
			Title: "Smart PM2.5 Grid Sensor", Type: "Medical Devices", Description: "LoRaWAN sensor.", Keywords: "IoT, Environment",
			TrlScore: ptrInt16(4), Status: false, IsUrgent: false, UrgentReason: "Submit for NRCT funding", UrgentFeedback: "Please revise budget",
			CreatedAt: now.AddDate(0, -2, 0), UpdatedAt: now,
		},
		// Case 3: TRL 2. In Progress. No Urgent. NO Admin (New Case).
		{
			ID: "CS-00003", ResearcherID: "RS-00003", CoordinatorID: "CO-00003",
			Title: "Nano-Curcumin Extraction", Type: "Medical Devices", Description: "New extraction method.", Keywords: "Herb, Nano",
			TrlScore: ptrInt16(2), Status: false, IsUrgent: true, UrgentReason: "Submit for NRCT funding",
			CreatedAt: now.AddDate(0, 0, -10), UpdatedAt: now,
		},
		// Case 4: TRL 5. In Progress. Urgent (Contest). Has Admin.
		{
			ID: "CS-00004", ResearcherID: "RS-00004", CoordinatorID: "CO-00004",
			Title: "Self-Healing Bio-Concrete", Type: "Medicines Vaccines Stem Cells", Description: "Bacteria concrete.", Keywords: "Bio, Construction",
			TrlScore: ptrInt16(4), Status: false, IsUrgent: true, UrgentReason: "Competition deadline",
			CreatedAt: now.AddDate(0, -4, 0), UpdatedAt: now,
		},
		// Case 5: TRL 9. Finished. No Urgent. Has Admin.
		{
			ID: "CS-00005", ResearcherID: "RS-00005", CoordinatorID: "CO-00005",
			Title: "Probiotic Durian Bites", Type: "Plant/Animal Breeds", Description: "Freeze-dried product.", Keywords: "Food, Export",
			TrlScore: ptrInt16(1), Status: true, IsUrgent: false,
			CreatedAt: now.AddDate(-1, 0, 0), UpdatedAt: now,
		},
	}
	if err := db.CreateInBatches(&cases, 5).Error; err != nil {
		log.Fatalf("❌ Failed to seed cases: %v", err)
	}
	fmt.Println("✅ Cases seeded")

	// ==========================================
	// 📅 5. APPOINTMENTS
	// ==========================================
	appointments := []models.Appointments{
		// Case 1
		{ID: "AP-00001", CaseID: "CS-00001", Date: now.AddDate(0, 0, -5), Status: "pending", Location: "Zoom", Detail: "MOU Signing", Summary: "Signed successfully", IsNotify: true},
		// Case 2
		{ID: "AP-00002", CaseID: "CS-00002", Date: now.AddDate(0, 0, 2), Status: "pending", Location: "Innovation Lab", Detail: "Budget Review", Summary: ""},
		// Case 3
		{ID: "AP-00003", CaseID: "CS-00003", Date: now.AddDate(0, -1, 0), Status: "attended", Location: "Research Office", Detail: "Initial Consultation", Summary: "Discussed research plan", IsNotify: true},
		// Case 4
		{ID: "AP-00004", CaseID: "CS-00004", Date: now.AddDate(0, 0, 7), Status: "absent", Location: "Site Visit", Detail: "Witness Stress Test", Summary: ""},
		// Case 5
		{ID: "AP-00005", CaseID: "CS-00005", Date: now.AddDate(0, -2, 0), Status: "attended", Location: "Factory", Detail: "Final Inspection", Summary: "Passed", IsNotify: true},
	}
	if err := db.CreateInBatches(&appointments, 5).Error; err != nil {
		log.Fatalf("❌ Failed to seed appointments: %v", err)
	}
	fmt.Println("✅ Appointments seeded")

	// ==========================================
	// 📝 6. ASSESSMENTS (DECISION TREE LOGIC)
	// ==========================================
	// Checkbox logic: If TRL achieved is X, levels 1 to X are 'passed' (full/partial check), X+1 is empty.
	assessments := []models.Assessments{
		// --------------------------------------------------------
		// Case 1 (Target: TRL 8)
		// --------------------------------------------------------
		{
			ID: "AS-00001", CaseID: "CS-00001", TrlEstimate: 8,
			Rq1Answer: true, Rq2Answer: true, Rq3Answer: true, Rq4Answer: true, Rq5Answer: true,
			Rq6Answer: true, Rq7Answer: true,
			Cq1Answer:             emptyJSON,
			Cq2Answer:             emptyJSON,
			Cq3Answer:             emptyJSON,
			Cq4Answer:             emptyJSON,
			Cq5Answer:             emptyJSON,
			Cq6Answer:             emptyJSON,
			Cq7Answer:             emptyJSON,
			Cq8Answer:             toJSON(checkboxData[7]),
			Cq9Answer:             toJSON([]string{checkboxData[8][0], checkboxData[8][1]}),
			ImprovementSuggestion: "มีการออกแบบโดยคำนึงถึงเป้าหมายด้านต้นทุน, มีการระบุและบรรเทาความเสี่ยงด้านความปลอดภัยและผลข้างเคียง, Ready for market strategy (Level 9).",
		},
		// --------------------------------------------------------
		// Case 2 (Target: TRL 4)
		// --------------------------------------------------------
		{
			ID: "AS-00002", CaseID: "CS-00002", TrlEstimate: 4,
			Rq1Answer: true, Rq2Answer: false, Rq3Answer: false, Rq4Answer: false, Rq5Answer: false,
			Rq6Answer: true, Rq7Answer: true,
			Cq1Answer:             emptyJSON,
			Cq2Answer:             emptyJSON,
			Cq3Answer:             emptyJSON,
			Cq4Answer:             toJSON(checkboxData[3]),
			Cq5Answer:             toJSON([]string{checkboxData[4][0], checkboxData[4][1], checkboxData[4][4], checkboxData[4][5]}),
			Cq6Answer:             emptyJSON,
			Cq7Answer:             emptyJSON,
			Cq8Answer:             emptyJSON,
			Cq9Answer:             emptyJSON,
			ImprovementSuggestion: "มีการวัดผลกระบวนการที่เที่ยงตรง, มีการระบุปัญหาและประเมินความน่าเชื่อถือด้านคุณภาพ, Need to test in relevant environment for TRL 5.",
		},
		// --------------------------------------------------------
		// Case 3 (Target: TRL 2)
		// Path: RQ1(0)->RQ7(0) -> result 2
		// --------------------------------------------------------
		{
			ID: "AS-00003", CaseID: "CS-00003", TrlEstimate: 2,
			Rq1Answer: false, Rq2Answer: false, Rq3Answer: false, Rq4Answer: false, Rq5Answer: false,
			Rq6Answer: false, Rq7Answer: true,
			Cq1Answer:             emptyJSON,
			Cq2Answer:             toJSON(checkboxData[1]),
			Cq3Answer:             toJSON([]string{checkboxData[2][0], checkboxData[2][3], checkboxData[2][4], checkboxData[2][5], checkboxData[2][6]}),
			Cq4Answer:             emptyJSON,
			Cq5Answer:             emptyJSON,
			Cq6Answer:             emptyJSON,
			Cq7Answer:             emptyJSON,
			Cq8Answer:             emptyJSON,
			Cq9Answer:             emptyJSON,
			ImprovementSuggestion: "การทดลองสามารถคาดการณ์ของส่วนประกอบเทคโนโลยีได้, มีการสร้างตัวชี้วัดประสิทธิภาพเทคโนโลยีหรือระบบ, Technology concept formulated. Need experimental proof of concept (TRL 3).",
		},
		// --------------------------------------------------------
		// Case 4 (Target: TRL 5)
		// Path: RQ1(1)->RQ2(0)->RQ6(1) -> result 5
		// --------------------------------------------------------
		{
			ID: "AS-00004", CaseID: "CS-00004", TrlEstimate: 5,
			Rq1Answer: true, Rq2Answer: false, Rq3Answer: false, Rq4Answer: false, Rq5Answer: false,
			Rq6Answer: true, Rq7Answer: false,
			Cq1Answer:             emptyJSON,
			Cq2Answer:             emptyJSON,
			Cq3Answer:             emptyJSON,
			Cq4Answer:             emptyJSON,
			Cq5Answer:             toJSON(checkboxData[4]),
			Cq6Answer:             emptyJSON,
			Cq7Answer:             emptyJSON,
			Cq8Answer:             emptyJSON,
			Cq9Answer:             emptyJSON,
			ImprovementSuggestion: "Next step: System model/prototype (TRL 6).",
		},
		// --------------------------------------------------------
		// Case 5 (Target: TRL 9)
		// Path: RQ1(1)->RQ2(1)->RQ3(1)->RQ4(1)->RQ5(1) -> result 9
		// --------------------------------------------------------
		{
			ID: "AS-00005", CaseID: "CS-00005", TrlEstimate: 9,
			Rq1Answer: true, Rq2Answer: true, Rq3Answer: true, Rq4Answer: true, Rq5Answer: true,
			Rq6Answer: true, Rq7Answer: true,
			Cq1Answer:             emptyJSON,
			Cq2Answer:             emptyJSON,
			Cq3Answer:             emptyJSON,
			Cq4Answer:             emptyJSON,
			Cq5Answer:             emptyJSON,
			Cq6Answer:             emptyJSON,
			Cq7Answer:             emptyJSON,
			Cq8Answer:             emptyJSON,
			Cq9Answer:             toJSON(checkboxData[8]),
			ImprovementSuggestion: "System is fully deployed and operational.",
		},
	}
	if err := db.CreateInBatches(&assessments, 5).Error; err != nil {
		log.Fatalf("❌ Failed to seed assessments: %v", err)
	}
	fmt.Println("✅ Assessments seeded")

	// ==========================================
	// 💡 7. INTELLECTUAL PROPERTIES
	// ==========================================
	ips := []models.IntellectualProperties{
		{ID: "IP-00001", CaseID: "CS-00001", Types: "ลิขสิทธิ์", ProtectionStatus: "Granted", RequestNumber: "CR-2024-001"},
		{ID: "IP-00002", CaseID: "CS-00002", Types: "อนุสิทธิบัตร", ProtectionStatus: "Pending", RequestNumber: "PT-2025-099"},
		{ID: "IP-00003", CaseID: "CS-00003", Types: "สิทธิบัตร", ProtectionStatus: "Drafting", RequestNumber: ""},
		{ID: "IP-00004", CaseID: "CS-00004", Types: "เครื่องหมายการค้า", ProtectionStatus: "Registered", RequestNumber: "TM-88899"},
		{ID: "IP-00005", CaseID: "CS-00005", Types: "สิทธิบัตรออกแบบผลิตภัณฑ์", ProtectionStatus: "Application Filed", RequestNumber: "DS-2024-555"},
	}
	if err := db.CreateInBatches(&ips, 5).Error; err != nil {
		log.Fatalf("❌ Failed to seed ips: %v", err)
	}
	fmt.Println("✅ IPs seeded")

	// ==========================================
	// 🆘 8. SUPPORTMENTS
	// ==========================================
	supportments := []models.Supportments{
		{ID: "SP-00001", CaseID: "CS-00001", NeedPartners: true, NeedCapital: true, Need: "Looking for private hospital partners."},
		{ID: "SP-00002", CaseID: "CS-00002", NeedTest: true, NeedGuidelines: true, Need: "Need standard calibration lab."},
		{ID: "SP-00003", CaseID: "CS-00003", SupportResearch: true, NeedProtectIntellectualProperty: true, Need: "Consultation on patent drafting."},
		{ID: "SP-00004", CaseID: "CS-00004", NeedActivities: true, NeedCoDevelopers: true, Need: "Matching with construction firms."},
		{ID: "SP-00005", CaseID: "CS-00005", SupportSiEIC: true, NeedAccount: true, Need: "Business model canvas workshop."},
	}
	if err := db.CreateInBatches(&supportments, 5).Error; err != nil {
		log.Fatalf("❌ Failed to seed supportments: %v", err)
	}
	fmt.Println("✅ Supportments seeded")

	// ==========================================
	// 🗣️ 9. CHAT LOGS
	// ==========================================
	chatLogs := []models.ChatLogs{
		{ID: "CL-00001", AdminID: nil, ResearcherID: ptrString("RS-00001"), History: datatypes.JSON([]byte(`[{"role": "user", "content": "สวัสดี"}, {"role": "assistant", "content": "สวัสดี! มีอะไรให้ฉันช่วยไหม? 😊"}, {"role": "user", "content": "trl คืออะไร"}, {"role": "assistant", "content": "TRL ย่อมาจาก \"Technology Readiness Level\" หรือ \"ระดับความพร้อมด้านเทคโนโลยี\" เป็นระบบที่ใช้วัดระดับการพัฒนาของเทคโนโลยี ตั้งแต่แนวคิดพื้นฐานไปจนถึงการนำไปใช้งานจริงในสภาพแวดล้อมการปฏิบัติงาน\n\nระดับ TRL มีทั้งหมด 9 ขั้น ได้แก่:\n\n1. **TRL 1** - สังเกตหลักการพื้นฐาน\n2. **TRL 2** - กำหนดแนวคิดทางเทคโนโลยี\n3. **TRL 3** - พิสูจน์แนวคิดในระดับทดลอง\n4. **TRL 4** - ทดสอบในห้องปฏิบัติการ\n5. **TRL 5** - ทดสอบในสภาพแวดล้อมจริงขั้นพื้นฐาน\n6. **TRL 6** - ทดสอบต้นแบบในสภาพแวดล้อมจริง\n7. **TRL 7** - แสดงต้นแบบในสภาพแวดล้อมการปฏิบัติงาน\n8. **TRL 8** - เทคโนโลยีพร้อมใช้งานจริงและผ่านการทดสอบ\n9. **TRL 9** - เทคโนโลยีนำมาใช้งานจริงแล้วในภารกิจจริง\n\nระดับ TRL นี้ช่วยให้ประเมินความก้าวหน้าและความเสี่ยงของการพัฒนาเทคโนโลยีได้อย่างมีระบบ"}]`))},
		{ID: "CL-00002", AdminID: nil, ResearcherID: ptrString("RS-00002"), History: datatypes.JSON([]byte(`[{"role": "user", "content": "hi"}, {"role": "assistant", "content": "Hello! How can I assist you with Technology Readiness Levels (TRL 1-9) today?"}]`))},
		{ID: "CL-00003", AdminID: ptrString("AD-00001"), ResearcherID: nil, History: datatypes.JSON([]byte(`[{"role": "user", "content": "trl คืออะไร"}, {"role": "assistant", "content": "TRL หรือ **Technology Readiness Level** คือ **ระดับความพร้อมของเทคโนโลยี** ที่ใช้ในการประเมินความก้าวหน้าของการพัฒนาเทคโนโลยี ตั้งแต่แนวคิดเริ่มต้นจนถึงการใช้งานจริงในสภาพแวดล้อมปฏิบัติการ  \n\nโดย TRL มีทั้งหมด **9 ระดับ (TRL 1–TRL 9)** ซึ่งแต่ละระดับแสดงถึงขั้นตอนของการพัฒนาเทคโนโลยี เช่น  \n\n- **TRL 1:** การสังเกตและรายงานหลักการทางวิทยาศาสตร์  \n- **TRL 2:** การกำหนดแนวคิดหรือการประยุกต์ใช้เทคโนโลยี  \n- **TRL 3:** การยืนยันแนวคิดในห้องปฏิบัติการ  \n- **TRL 4–6:** การทดสอบต้นแบบในสภาพแวดล้อมต่าง ๆ  \n- **TRL 7–9:** การสาธิตและใช้งานเทคโนโลยีจริงในสภาพแวดล้อมปฏิบัติการ  \n\nสรุปคือ TRL ใช้เพื่อประเมินว่าเทคโนโลยีหนึ่ง ๆ “พร้อมแค่ไหน” สำหรับการใช้งานจริง."}]`))},
		{ID: "CL-00004", AdminID: ptrString("AD-00004"), ResearcherID: nil, History: datatypes.JSON([]byte(`[{"role": "user", "content": "trl 7 ต้องมีอะไรบ้าง"}, {"role": "assistant", "content": "**TRL 7** คือระดับที่เทคโนโลยีอยู่ในขั้นตอน **“การสาธิตต้นแบบในสภาพแวดล้อมจริง”** \n\nสิ่งที่ต้องมีในระดับนี้ ได้แก่  \n\n- **ระบบต้นแบบ (Prototype system)** ที่รวมองค์ประกอบทั้งหมดของเทคโนโลยีครบถ้วนตามการออกแบบสุดท้าย  \n- การ **สาธิตการทำงานในสภาพแวดล้อมจริงหรือใกล้เคียงกับการใช้งานจริง** เพื่อพิสูจน์ว่าเทคโนโลยีทำงานได้ภายใต้เงื่อนไขที่คาดว่าจะพบจริง  \n- มีการ **เก็บข้อมูลผลการทดสอบ** เพื่อยืนยันประสิทธิภาพ ความน่าเชื่อถือ และความเสถียรของระบบ  \n\nกล่าวโดยสรุป TRL 7 คือขั้นที่เทคโนโลยีได้รับการ **พิสูจน์แล้วว่าสามารถทำงานได้จริงในสภาพจริง** แม้ยังไม่ถึงขั้นผลิตเชิงพาณิชย์เต็มรูปแบบ."}]`))},
		{ID: "CL-00005", AdminID: nil, ResearcherID: ptrString("RS-00005"), History: datatypes.JSON([]byte(`[{"role": "user", "content": "trl 5 คือ"}, {"role": "assistant", "content": "**TRL 5** คือระดับที่เรียกว่า  \n**“Component and/or Breadboard Validation in Relevant Environments”** \nหรือ **การทดสอบองค์ประกอบและ/หรืออุปกรณ์ที่ใช้ทดลอง (breadboard) ในสถานการณ์ที่เกี่ยวข้อง** \n\n**รายละเอียดของ TRL 5:**\n- เป็นขั้นที่เทคโนโลยีถูกเชื่อมต่อเข้าด้วยกันในระดับที่มีนัยสำคัญ  \n- องค์ประกอบพื้นฐานของเทคโนโลยีถูกจัดประกอบเข้ากับชิ้นงานบางส่วนและนำไปทดสอบในสถานการณ์จำลองที่เกี่ยวข้องจริง  \n- มีการเปรียบเทียบผลการทดสอบกับสมมติฐานหรือเป้าหมายที่ตั้งไว้  \n- มีการระบุปัญหาที่พบและประเมินความใกล้เคียงของผลกับเป้าหมาย  \n\n**ตัวอย่างการประยุกต์ TRL 5 ในแต่ละด้าน:**\n- **ด้านซอฟต์แวร์:** ทดสอบส่วนประกอบหรือระบบย่อยในสภาพแวดล้อมที่เกี่ยวข้อง  \n- **ด้านอุปกรณ์ทางการแพทย์:** ทำการทดสอบยืนยันความถูกต้องและความปลอดภัย (Validation and Safety Test)  \n- **ด้านยา วัคซีน สเต็มเซลล์:** ทำการวิเคราะห์คุณสมบัติขั้นสูงของตัวอย่างและเริ่มพัฒนาในระดับ GMP  \n- **ด้านพันธุ์พืช/พันธุ์สัตว์:** คัดเลือกประชากรที่มีลักษณะตามเป้าหมายโดยการปลูกหรือเลี้ยงในระดับแปลงสาธิตหรือแปลงควบคุม  \n\nสรุปคือ TRL 5 เป็นขั้นที่เริ่ม **เชื่อมโยงองค์ประกอบของเทคโนโลยีเข้าด้วยกันและทดสอบในสภาพแวดล้อมจำลองที่ใกล้เคียงของจริง** ครับ."}]`))},
	}
	if err := db.CreateInBatches(&chatLogs, 5).Error; err != nil {
		log.Fatalf("❌ Failed to seed chat logs: %v", err)
	}
	fmt.Println("✅ Chat logs seeded")

	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("🎉 FINAL SEEDING COMPLETE!")
}
