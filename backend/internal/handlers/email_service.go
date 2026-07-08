package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"strings"
)

// EmailMessage represents an email to send
type EmailMessage struct {
	To      string
	Subject string
	Body    string
	IsHTML  bool
}

// SendEmail sends an email using SMTP or Mailgun (gracefully skips if neither is configured)
func SendEmail(msg EmailMessage) error {
	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost != "" {
		return sendViaSMTP(msg)
	}
	mailgunKey := os.Getenv("MAILGUN_API_KEY")
	if mailgunKey != "" {
		return sendViaMailgun(msg)
	}
	log.Printf("[Email] No email service configured. Would send to %s: %s", msg.To, msg.Subject)
	return nil
}

func sendViaSMTP(msg EmailMessage) error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM")
	if from == "" {
		from = user
	}
	if port == "" {
		port = "587"
	}

	auth := smtp.PlainAuth("", user, password, host)

	contentType := "text/plain"
	if msg.IsHTML {
		contentType = "text/html"
	}

	headers := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: %s; charset=UTF-8\r\n\r\n",
		from, msg.To, msg.Subject, contentType,
	)
	body := headers + msg.Body

	addr := host + ":" + port
	err := smtp.SendMail(addr, auth, from, []string{msg.To}, []byte(body))
	if err != nil {
		log.Printf("[SMTP] Failed to send email to %s: %v", msg.To, err)
		return err
	}
	log.Printf("[SMTP] Email sent to %s: %s", msg.To, msg.Subject)
	return nil
}

func sendViaMailgun(msg EmailMessage) error {
	apiKey := os.Getenv("MAILGUN_API_KEY")
	domain := os.Getenv("MAILGUN_DOMAIN")
	from := os.Getenv("SMTP_FROM")
	if from == "" {
		from = "noreply@" + domain
	}

	apiURL := fmt.Sprintf("https://api.mailgun.net/v3/%s/messages", domain)

	formData := url.Values{}
	formData.Set("from", from)
	formData.Set("to", msg.To)
	formData.Set("subject", msg.Subject)
	if msg.IsHTML {
		formData.Set("html", msg.Body)
	} else {
		formData.Set("text", msg.Body)
	}

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return err
	}
	req.SetBasicAuth("api", apiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[Mailgun] Failed to send email to %s: %v", msg.To, err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		log.Printf("[Mailgun] Non-2xx response for %s: %d", msg.To, resp.StatusCode)
		return fmt.Errorf("mailgun returned status %d", resp.StatusCode)
	}
	log.Printf("[Mailgun] Email sent to %s: %s", msg.To, msg.Subject)
	return nil
}

// ---- Email Template Functions ----

// SendWelcomeEmail sends a welcome email after registration
func SendWelcomeEmail(toEmail, fullName string) {
	subject := "Welcome to Yojana Portal!"
	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<body style="font-family: Arial, sans-serif; color: #333;">
  <h2 style="color:#1a56db;">Welcome to Yojana Portal, %s!</h2>
  <p>Thank you for registering with us. Your account has been successfully created.</p>
  <p>You can now:</p>
  <ul>
    <li>Explore government schemes and check your eligibility</li>
    <li>Browse and apply to government and private sector jobs</li>
    <li>Get AI-powered career recommendations</li>
  </ul>
  <p>If you have any questions, feel free to contact our support team.</p>
  <br/>
  <p>Best regards,<br/>The Yojana Portal Team</p>
</body>
</html>`, fullName)

	go func() {
		if err := SendEmail(EmailMessage{
			To:      toEmail,
			Subject: subject,
			Body:    body,
			IsHTML:  true,
		}); err != nil {
			log.Printf("[Email] Failed to send welcome email to %s: %v", toEmail, err)
		}
	}()
}

// SendNewsletterWelcomeEmail sends a confirmation to new newsletter subscribers
func SendNewsletterWelcomeEmail(toEmail, name string) {
	subject := "You're subscribed to Yojana Portal Newsletter!"
	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<body style="font-family: Arial, sans-serif; color: #333;">
  <h2 style="color:#1a56db;">Thank you for subscribing%s!</h2>
  <p>You will now receive updates on:</p>
  <ul>
    <li>New government schemes and benefits</li>
    <li>Job openings and career opportunities</li>
    <li>Important deadlines and announcements</li>
  </ul>
  <p>To unsubscribe, simply reply to any newsletter email with "UNSUBSCRIBE" in the subject.</p>
  <br/>
  <p>Best regards,<br/>The Yojana Portal Team</p>
</body>
</html>`, func() string {
		if name != "" {
			return ", " + name
		}
		return ""
	}())

	go func() {
		if err := SendEmail(EmailMessage{
			To:      toEmail,
			Subject: subject,
			Body:    body,
			IsHTML:  true,
		}); err != nil {
			log.Printf("[Email] Failed to send newsletter welcome email to %s: %v", toEmail, err)
		}
	}()
}

// SendApplicationStatusEmail notifies user of scheme application status change
func SendApplicationStatusEmail(toEmail, fullName, schemeName, status string) {
	subject := fmt.Sprintf("Application Status Update: %s", schemeName)
	statusText := "updated"
	statusColor := "#333"
	switch status {
	case "approved":
		statusText = "approved ✅"
		statusColor = "#16a34a"
	case "rejected":
		statusText = "rejected ❌"
		statusColor = "#dc2626"
	case "pending":
		statusText = "under review ⏳"
		statusColor = "#d97706"
	}

	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<body style="font-family: Arial, sans-serif; color: #333;">
  <h2 style="color:#1a56db;">Application Status Update</h2>
  <p>Dear %s,</p>
  <p>Your application for <strong>%s</strong> has been <span style="color:%s;font-weight:bold;">%s</span>.</p>
  <p>Please log in to your Yojana Portal account for more details.</p>
  <br/>
  <p>Best regards,<br/>The Yojana Portal Team</p>
</body>
</html>`, fullName, schemeName, statusColor, statusText)

	go func() {
		if err := SendEmail(EmailMessage{
			To:      toEmail,
			Subject: subject,
			Body:    body,
			IsHTML:  true,
		}); err != nil {
			log.Printf("[Email] Failed to send status email to %s: %v", toEmail, err)
		}
	}()
}
