package send_email

import (
    "fmt"
    "net/smtp"
    "strings"
    "sync"
)

type SMTPConfig struct {
    Host     string
    Port     string
    Username string
    Password string
    From     string
}

type EmailService func(to, subject, body string) error

type VerifiedRecipients struct {
    Researchers        []Recipient
    UniqueCoordinators []Recipient
    Duplicates         []Recipient
}

// VerifyRecipients separates researchers from coordinators to avoid duplicate emails
func VerifyRecipients(researchers, coordinators []Recipient) VerifiedRecipients {
    researcherEmails := make(map[string]bool)
    
    for _, r := range researchers {
        researcherEmails[strings.ToLower(r.Email)] = true
    }
    
    var uniqueCoordinators []Recipient
    var duplicates []Recipient
    
    for _, c := range coordinators {
        if researcherEmails[strings.ToLower(c.Email)] {
            duplicates = append(duplicates, c)
        } else {
            uniqueCoordinators = append(uniqueCoordinators, c)
        }
    }
    
    return VerifiedRecipients{
        Researchers:        researchers,
        UniqueCoordinators: uniqueCoordinators,
        Duplicates:         duplicates,
    }
}

// SendEmail sends emails to all recipients in the list
func SendEmail(details AppointmentDetails, emailService EmailService) (*SendEmailResults, error) {
    verified := VerifyRecipients(details.Researchers, details.Coordinators)
    
    results := &SendEmailResults{
        Status:       "completed",
        TotalSkipped: len(verified.Duplicates),
        Success:      []EmailResult{},
        Failed:       []EmailResult{},
    }
    
    var wg sync.WaitGroup
    var mu sync.Mutex
    

    // Get coordinator name for researcher emails
    var coordinatorName string
    if len(details.Coordinators) > 0 {
        coordinatorName = details.Coordinators[0].Name
    }
    
    // Send to Researchers
    for _, researcher := range verified.Researchers {
        wg.Add(1)
        go func(r Recipient) {
            defer wg.Done()
            content := TemplateCreate(r, details, true, coordinatorName)
            subject := GenerateEmailSubject(details.ResearchTitle)
            err := emailService(r.Email, subject, content)
            
            mu.Lock()
            updateResults(results, r.Email, r.Name, err)
            mu.Unlock()
        }(researcher)
    }
    
    // Send to Unique Coordinators
    for _, coordinator := range verified.UniqueCoordinators {
        wg.Add(1)
        go func(c Recipient) {
            defer wg.Done()
            content := TemplateCreate(c, details, false, "")
            subject := GenerateEmailSubject(details.ResearchTitle)
            err := emailService(c.Email, subject, content)
            
            mu.Lock()
            updateResults(results, c.Email, c.Name, err)
            mu.Unlock()
        }(coordinator)
    }
    
    wg.Wait()
    return results, nil
}

func updateResults(results *SendEmailResults, email string, name string, err error) {
    res := EmailResult{
        Email:   email,
        Name:    name,
        Success: err == nil,
        Error:   err,
    }
    if err != nil {
        results.Failed = append(results.Failed, res)
        results.TotalFailed++
    } else {
        results.Success = append(results.Success, res)
        results.TotalSent++
    }
}

// CreateSMTPEmailService creates an email service using SMTP
func CreateSMTPEmailService(config SMTPConfig) EmailService {
    return func(to, subject, body string) error {
        auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)
        
        headers := map[string]string{
            "From":         config.From,
            "To":           to,
            "Subject":      subject,
            "MIME-Version": "1.0",
            "Content-Type": "text/plain; charset=\"utf-8\"",
        }
        
        message := ""
        for k, v := range headers {
            message += fmt.Sprintf("%s: %s\r\n", k, v)
        }
        message += "\r\n" + body
        
        addr := fmt.Sprintf("%s:%s", config.Host, config.Port)
        err := smtp.SendMail(addr, auth, config.From, []string{to}, []byte(message))
        if err != nil {
            return fmt.Errorf("failed to send email to %s: %w", to, err)
        }
        return nil
    }
}

// CreateMockEmailService creates a mock email service for testing
func CreateMockEmailService() EmailService {
    return func(to, subject, body string) error {
        fmt.Printf("📧 Mock Email Sent | To: %s | Subject: %s\n", to, subject)
        return nil
    }
}