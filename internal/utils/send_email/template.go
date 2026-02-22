package send_email

import (
	"fmt"
	"html"
	"time"
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
	Summary       string
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
	// Convert time to Asia/Bangkok (UTC+7)
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		loc = time.FixedZone("ICT", 7*60*60)
	}
	localTime := ap.Date.In(loc)

	details := AppointmentDetails{
		Date:     localTime.Format("02 January 2006"),
		Time:     localTime.Format("15:04"),
		Location: ap.Location,
		Detail:   ap.Detail,
		Summary:  ap.Summary,
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

// getHTMLTemplate wraps content in a beautiful, formal HTML structure
func getHTMLTemplate(title string, bodyContent string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; line-height: 1.6; color: #333; background-color: #f4f4f4; margin: 0; padding: 0; }
        .container { max-width: 600px; margin: 30px auto; background-color: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 4px 15px rgba(0,0,0,0.1); }
        .header { background-color: #00C1D6; color: #ffffff; padding: 25px; text-align: center; }
        .header h1 { margin: 0; font-size: 22px; font-weight: 600; letter-spacing: 0.5px; }
        .content { padding: 30px; border-left: 1px solid #eee; border-right: 1px solid #eee; }
        .footer { background-color: #f8f9fa; padding: 10px; text-align: center; font-size: 12px; color: #888; border: 1px solid #eee; border-bottom-left-radius: 8px; border-bottom-right-radius: 8px; }
        .info-table { width: 100%%; border-collapse: collapse; margin: 20px 0; }
        .info-table td { padding: 12px 5px; border-bottom: 1px solid #f0f0f0; vertical-align: top; }
        .label { font-weight: 600; width: 30%%; color: #555; }
        .value { color: #222; }
        .highlight { color: #d9534f; font-weight: bold; font-size: 0.9em; background-color: #fdf0f0; padding: 2px 6px; border-radius: 4px; }
        p { margin-bottom: 15px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>%s</h1>
        </div>
        <div class="content">
            %s
        </div>
        <div class="footer">
            <p>This is an automated message from the TRL Research Administration System.<br>Please do not reply to this email.</p>
            <p style="color: #ccc; font-size: 10px;">Ref: %s</p>
        </div>
    </div>
</body>
</html>`, title, bodyContent, time.Now().Format("20060102-150405"))
}

// TemplateCreate for new appointments
func TemplateCreate(recipient Recipient, details AppointmentDetails, isResearcher bool, coordinatorName string) string {
	body := fmt.Sprintf("<p>Dear <strong>%s</strong>,</p>", html.EscapeString(recipient.Name))
	body += fmt.Sprintf("<p>We are writing to inform you that a new appointment has been scheduled for the research project: <strong>%s</strong>.</p>", html.EscapeString(details.ResearchTitle))

	body += "<table class='info-table'>"
	body += fmt.Sprintf("<tr><td class='label'>Date:</td><td class='value'>%s</td></tr>", html.EscapeString(details.Date))
	body += fmt.Sprintf("<tr><td class='label'>Time:</td><td class='value'>%s</td></tr>", html.EscapeString(details.Time))
	body += fmt.Sprintf("<tr><td class='label'>Location:</td><td class='value'>%s</td></tr>", html.EscapeString(details.Location))

	if details.Detail != "" {
		body += fmt.Sprintf("<tr><td class='label'>Note:</td><td class='value'>%s</td></tr>", html.EscapeString(details.Detail))
	}

	if isResearcher && coordinatorName != "" {
		body += fmt.Sprintf("<tr><td class='label'>Coordinator:</td><td class='value'>%s</td></tr>", html.EscapeString(coordinatorName))
	}
	body += "</table>"

	body += "<p>We kindly request your punctual attendance.</p>"
	body += "<p>Sincerely,<br>TRL Research Administration System</p>"

	return getHTMLTemplate("New Appointment Scheduled", body)
}

// TemplateUpdate handles notifications when appointment details are modified
func TemplateUpdate(recipient Recipient, details AppointmentDetails, appointmentID string, updatedAt string, changedFields map[string]bool) string {
	body := fmt.Sprintf("<p>Dear <strong>%s</strong>,</p>", html.EscapeString(recipient.Name))
	body += "<p><strong style='color: #d9534f;'>[URGENT]</strong> The details for your appointment have been updated.</p>"

	body += "<table class='info-table'>"
	body += fmt.Sprintf("<tr><td class='label'>Project:</td><td class='value'>%s</td></tr>", html.EscapeString(details.ResearchTitle))
	body += fmt.Sprintf("<tr><td class='label'>Last Updated:</td><td class='value'>%s</td></tr>", html.EscapeString(updatedAt))
	body += "</table>"

	body += "<h3>Updated Details</h3>"
	body += "<table class='info-table'>"

	// Date
	dateVal := html.EscapeString(details.Date)
	if changedFields["date"] || changedFields["time"] {
		dateVal += " <span class='highlight'>(UPDATED)</span>"
	}
	body += fmt.Sprintf("<tr><td class='label'>Date:</td><td class='value'>%s</td></tr>", dateVal)

	// Time
	timeVal := html.EscapeString(details.Time)
	if changedFields["date"] || changedFields["time"] {
		timeVal += " <span class='highlight'>(UPDATED)</span>"
	}
	body += fmt.Sprintf("<tr><td class='label'>Time:</td><td class='value'>%s</td></tr>", timeVal)

	// Location
	locVal := html.EscapeString(details.Location)
	if changedFields["location"] {
		locVal += " <span class='highlight'>(UPDATED)</span>"
	}
	body += fmt.Sprintf("<tr><td class='label'>Location:</td><td class='value'>%s</td></tr>", locVal)

	// Detail
	detVal := html.EscapeString(details.Detail)
	if changedFields["detail"] {
		detVal += " <span class='highlight'>(UPDATED)</span>"
	}
	body += fmt.Sprintf("<tr><td class='label'>Note:</td><td class='value'>%s</td></tr>", detVal)

	if details.Summary != "" {
		sumVal := html.EscapeString(details.Summary)
		if changedFields["summary"] {
			sumVal += " <span class='highlight'>(UPDATED)</span>"
		}
		body += fmt.Sprintf("<tr><td class='label'>Summary:</td><td class='value'>%s</td></tr>", sumVal)
	}
	body += "</table>"

	body += "<p>Please review these changes carefully and adjust your schedule accordingly.</p>"
	body += "<p>Sincerely,<br>TRL Research Administration System</p>"

	return getHTMLTemplate("Appointment Update Notification", body)
}

// TemplateReminder for reminders
func TemplateReminder(details AppointmentDetails, professor Professor, researcher Recipient, coordinator *Recipient) string {
	body := fmt.Sprintf("<p>This is a reminder for your upcoming appointment regarding the project: <strong>%s</strong>.</p>", html.EscapeString(details.ResearchTitle))

	body += "<table class='info-table'>"
	body += fmt.Sprintf("<tr><td class='label'>Date:</td><td class='value'>%s</td></tr>", html.EscapeString(details.Date))
	body += fmt.Sprintf("<tr><td class='label'>Time:</td><td class='value'>%s</td></tr>", html.EscapeString(details.Time))
	body += fmt.Sprintf("<tr><td class='label'>Location:</td><td class='value'>%s</td></tr>", html.EscapeString(details.Location))
	body += "</table>"

	body += "<h3>Participants</h3>"
	body += "<table class='info-table'>"
	body += fmt.Sprintf("<tr><td class='label'>Professor:</td><td class='value'>%s (%s)</td></tr>", html.EscapeString(professor.Name), html.EscapeString(professor.Email))
	body += fmt.Sprintf("<tr><td class='label'>Researcher:</td><td class='value'>%s</td></tr>", html.EscapeString(researcher.Name))
	if coordinator != nil {
		body += fmt.Sprintf("<tr><td class='label'>Coordinator:</td><td class='value'>%s</td></tr>", html.EscapeString(coordinator.Name))
	}
	body += "</table>"

	body += "<p>Sincerely,<br>TRL Research Administration System</p>"

	return getHTMLTemplate("Appointment Reminder", body)
}

// TemplateForgetPassword creates a formal email for temporary password reset
func TemplateForgetPassword(tempPass string) string {
	body := "<p>Dear User,</p>"
	body += "<p>We received a request to reset your password for the TRL Research Administration System.</p>"
	body += "<p>Your temporary password is:</p>"
	body += fmt.Sprintf("<div style='text-align: center; margin: 30px 0;'><div style='background-color: #f8f9fa; padding: 15px; font-size: 18px; font-weight: bold; letter-spacing: 1px; color: #00C1D6; border-radius: 6px; border: 1px solid #e9ecef; display: inline-block; max-width: 100%%; word-break: break-all; overflow-wrap: break-word;'>%s</div></div>", html.EscapeString(tempPass))
	body += "<p>This temporary password will expire in <strong>10 minutes</strong>. Please log in and change your password immediately.</p>"
	body += "<p>If you did not request a password reset, please contact the administrator immediately.</p>"
	body += "<p>Sincerely,<br>TRL Research Administration System</p>"

	return getHTMLTemplate("Password Reset Request", body)
}
