package cron

import (
	"fmt"
	"log"
	"time"
	"trl-research-backend/internal/config"
	"trl-research-backend/internal/repository"
	"trl-research-backend/internal/utils/send_email"

	"github.com/robfig/cron/v3"
)

type ReminderCron struct {
	AppointmentRepo repository.AppointmentRepository
	AdminRepo       repository.AdminRepository
	Cfg             config.Config
	Location        *time.Location
}

func NewReminderCron(apRepo repository.AppointmentRepository, adRepo repository.AdminRepository, cfg config.Config) *ReminderCron {
	return &ReminderCron{
		AppointmentRepo: apRepo,
		AdminRepo:       adRepo,
		Cfg:             cfg,
	}
}

func (c *ReminderCron) Start() {
	// Use Asia/Bangkok as the standard for this project
	jakartaTime, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Printf("⚠️ Could not load Asia/Bangkok location, falling back to local: %v", err)
		jakartaTime = time.Local
	}
	c.Location = jakartaTime
	
	cr := cron.New(cron.WithLocation(jakartaTime))

	// Scheduling: Run every 5 minutes
	_, err = cr.AddFunc("*/5 * * * *", func() {
		log.Println("⏰ Running meeting reminder cron job...")
		c.HandleReminders()
	})

	if err != nil {
		log.Fatalf("❌ Error adding cron job: %v", err)
	}

	cr.Start()
	log.Println("🚀 Reminder cron job started - Running every 5 minutes")
}

func (c *ReminderCron) HandleReminders() {
	// 1. Calculate the target time: Current Time + 1 Hour
	loc := c.Location
	if loc == nil {
		loc = time.Local
	}
	now := time.Now().In(loc)
	targetTime := now.Add(1 * time.Hour)

	// 2. Query Logic: Window range of +/- 5 minutes
	startRange := targetTime.Add(-5 * time.Minute)
	endRange := targetTime.Add(5 * time.Minute)

	// Query database for appointments within range where IsNotify is FALSE
	appointments, err := c.AppointmentRepo.GetUpcomingAppointments(startRange, endRange)
	if err != nil {
		log.Printf("❌ Error fetching upcoming appointments: %v", err)
		return
	}

	if len(appointments) == 0 {
		return
	}

	log.Printf("🔔 Found %d appointments to notify", len(appointments))

	// Fetch all admins for notification
	admins, err := c.AdminRepo.GetAdminAll()
	if err != nil {
		log.Printf("❌ Error fetching admins: %v", err)
	}

	// Initialize email service
	emailService := send_email.CreateSMTPEmailService(c.Cfg.GetSMTPConfig())

	for _, ap := range appointments {
		log.Printf("📨 Processing reminder for appointment %s (Case: %s)", ap.ID, ap.CaseID)

		// Extract emails and handle duplicates using a Set-like map
		recipientEmails := make(map[string]string) // email -> name

		// A. Admins
		for _, admin := range admins {
			name := fmt.Sprintf("%s %s %s", admin.Prefix, admin.FirstName, admin.LastName)
			recipientEmails[admin.Email] = name
		}

		// B. Researcher
		var researcherName string
		if ap.Case != nil && ap.Case.Researcher != nil {
			r := ap.Case.Researcher
			researcherName = fmt.Sprintf("%s %s %s", r.Prefix, r.FirstName, r.LastName)
			recipientEmails[r.Email] = researcherName
		}

		// C. Coordinator
		var coordinatorName string
		var coordinatorRecipient *send_email.Recipient
		if ap.Case != nil && ap.Case.Coordinator != nil {
			co := ap.Case.Coordinator
			coordinatorName = fmt.Sprintf("%s %s %s", co.Prefix, co.FirstName, co.LastName)
			recipientEmails[co.Email] = coordinatorName
			coordinatorRecipient = &send_email.Recipient{Name: coordinatorName, Email: co.Email}
		}

		// Prepare data for template
		details, professor := send_email.GatherAppointmentEmailData(&ap)
		researcherRecipient := send_email.Recipient{Name: researcherName}

		// Generate template content
		content := send_email.TemplateReminder(details, professor, researcherRecipient, coordinatorRecipient)
		subject := fmt.Sprintf("Reminder: Upcoming Appointment for %s", details.ResearchTitle)

		if len(recipientEmails) == 0 {
			log.Printf("⚠️ No recipients found for appointment %s", ap.ID)
			continue
		}

		// Send the email to unique recipients
		allSent := true
		for email, name := range recipientEmails {
			log.Printf("📤 Sending reminder to %s (%s)", name, email)
			err := emailService(email, subject, content)
			if err != nil {
				log.Printf("❌ Failed to send email to %s: %v", email, err)
				allSent = false
			}
		}

		if !allSent {
			log.Printf("⚠️ Not all reminders sent for appointment %s; leaving IsNotify=false for retry", ap.ID)
			continue
		}

		// Update Database: Immediately update IsNotify to TRUE
		err = c.AppointmentRepo.UpdateNotifyStatus(ap.ID, true)
		if err != nil {
			log.Printf("❌ Failed to update IsNotify status for appointment %s: %v", ap.ID, err)
		} else {
			log.Printf("✅ Successfully updated notification status for appointment %s", ap.ID)
		}
	}
}
