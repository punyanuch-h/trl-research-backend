package send_email

import (
	"fmt"
	"trl-research-backend/internal/models"
)

// Data structures for email templates
type Recipient struct {
	Name  string
	Email string
}

type Professor struct {
	Name             string
	Email            string
	AcademicPosition string
}

type AppointmentDetails struct {
	ResearchTitle string
	Date          string
	Time          string
	Location      string
	Detail        string
	Researchers   []Recipient
	Coordinators  []Recipient
}

type EmailResult struct {
	Email   string
	Name    string
	Success bool
	Error   error
}

type SendEmailResults struct {
	Status       string
	TotalSent    int
	TotalFailed  int
	TotalSkipped int
	Success      []EmailResult
	Failed       []EmailResult
}

// GatherAppointmentEmailData helper function to gather and format information from models
func GatherAppointmentEmailData(ap *models.Appointments) (AppointmentDetails, Professor) {
	details := AppointmentDetails{
		Date:     ap.Date.Format("02 January 2006"),
		Time:     ap.Date.Format("15:04"),
		Location: ap.Location,
		Detail:   ap.Detail,
	}

	if ap.Case != nil {
		details.ResearchTitle = ap.Case.Title

		// Add Researcher
		if ap.Case.Researcher != nil {
			r := ap.Case.Researcher
			name := fmt.Sprintf("%s %s %s", r.Prefix, r.FirstName, r.LastName)
			details.Researchers = append(details.Researchers, Recipient{
				Name:  name,
				Email: r.Email,
			})
		}

		// Add Coordinator
		if ap.Case.Coordinator != nil {
			c := ap.Case.Coordinator
			name := fmt.Sprintf("%s %s %s", c.Prefix, c.FirstName, c.LastName)
			details.Coordinators = append(details.Coordinators, Recipient{
				Name:  name,
				Email: c.Email,
			})
		}
	}

	prof := Professor{}
	if ap.Case != nil && ap.Case.Coordinator != nil {
		c := ap.Case.Coordinator
		prof = Professor{
			Name:             fmt.Sprintf("%s %s %s", c.Prefix, c.FirstName, c.LastName),
			Email:            c.Email,
			AcademicPosition: c.AcademicPosition,
		}
	}

	return details, prof
}

// GenerateEmailSubject creates a standard subject for new appointment emails
func GenerateEmailSubject(title string) string {
	return fmt.Sprintf("New Appointment Scheduled: %s", title)
}

// TemplateCreate for new appointments
func TemplateCreate(recipient Recipient, details AppointmentDetails, isResearcher bool, coordinatorName string) string {
	content := fmt.Sprintf("Dear %s,\n\n", recipient.Name)
	content += fmt.Sprintf("A new appointment has been scheduled for the research project: \"%s\"\n\n", details.ResearchTitle)
	content += "Appointment Details:\n"
	content += fmt.Sprintf("- Date: %s\n", details.Date)
	content += fmt.Sprintf("- Time: %s\n", details.Time)
	content += fmt.Sprintf("- Location: %s\n", details.Location)
	
	if details.Detail != "" {
		content += fmt.Sprintf("- Note: %s\n", details.Detail)
	}

	if isResearcher && coordinatorName != "" {
		content += fmt.Sprintf("\nYour coordinator for this session will be: %s\n", coordinatorName)
	}

	content += "\nPlease be on time.\n\nBest regards,\nREA System"
	return content
}

// TemplateReminder for reminders
func TemplateReminder(details AppointmentDetails, professor Professor, researcher Recipient, coordinator *Recipient) string {
	content := fmt.Sprintf("Reminder: Upcoming Appointment for \"%s\"\n\n", details.ResearchTitle)
	content += "This is a reminder for your upcoming appointment:\n"
	content += fmt.Sprintf("- Date: %s\n", details.Date)
	content += fmt.Sprintf("- Time: %s\n", details.Time)
	content += fmt.Sprintf("- Location: %s\n\n", details.Location)

	content += "Participants:\n"
	content += fmt.Sprintf("- Professor: %s (%s)\n", professor.Name, professor.Email)
	content += fmt.Sprintf("- Researcher: %s\n", researcher.Name)
	if coordinator != nil {
		content += fmt.Sprintf("- Coordinator: %s\n", coordinator.Name)
	}

	content += "\nThank you.\nREA System"
	return content
}
