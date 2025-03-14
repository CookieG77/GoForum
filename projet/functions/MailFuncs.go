package functions

import (
	"github.com/go-gomail/gomail"
	"os"
	"strconv"
	"sync"
)

// wg (WaitGroup) is a struct that waits for a collection of goroutines to finish
var wg sync.WaitGroup

var dialer *gomail.Dialer

// initialized is used to check if the SMTP service has been initialized or not.
var initialized bool = false

// InitMail initializes the mailer.
// If the SMTP server configuration file is not found, the function will log an error and return.
func InitMail(SMTPServerConfigFile string) {
	// Load the SMTP server configuration from the .env file
	smtpServer := os.Getenv("SMTP_HOST")
	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	if err != nil || smtpServer == "" || smtpUser == "" || smtpPassword == "" {
		ErrorPrintln("SMTP configuration not found in the .env file")
		return
	}
	dialer = gomail.NewDialer(smtpServer, smtpPort, smtpUser, smtpPassword)
	initialized = true
	SuccessPrintf(
		"SMTP server connected\n\t- host : \"%s\"\n\t- port : \"%s\"\n\t- user : \"%s\"\n",
		smtpServer,
		strconv.Itoa(smtpPort),
		smtpUser,
	)
}

// SendMail sends an email to the specified address.
// If the mailer has not been initialized, the function will log an error and return.
func SendMail(to string, subject string, content string) {
	if !initialized {
		ErrorPrintln("Mail Service not initialized, check the SMTP server configuration file")
		return
	}

	wg.Add(1)
	go sendMail(to, subject, content, &wg)
	wg.Wait()
}

// sendMail sends an email to the specified address
// This function is used by SendMail to send the email in a goroutine
func sendMail(to string, subject string, content string, wg *sync.WaitGroup) {
	defer wg.Done()

	m := gomail.NewMessage()
	m.SetHeader("From", dialer.Username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)

	if err := dialer.DialAndSend(m); err != nil {
		ErrorPrintf("Could not send mail to %s -> %v\n", to, err)
	} else {
		InfoPrintf("Mail sent to %s\n", to)
	}
}

func SendMailWithAttachments(to string, subject string, content string, attachments ...string) {
	if !initialized {
		ErrorPrintln("Mail Service not initialized, check the SMTP server configuration file")
		return
	}

	wg.Add(1)
	go sendMailWithAttachments(to, subject, content, attachments, &wg)
	wg.Wait()
}

func sendMailWithAttachments(to string, subject string, content string, attachments []string, wg *sync.WaitGroup) {
	defer wg.Done()

	m := gomail.NewMessage()
	m.SetHeader("From", dialer.Username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)

	for _, attachment := range attachments {
		m.Attach(attachment)
	}

	if err := dialer.DialAndSend(m); err != nil {
		ErrorPrintf("Could not send mail to %s -> %v\n", to, err)
	} else {
		InfoPrintf("Mail with one or more attachments sent to %s\n", to)
	}
}
