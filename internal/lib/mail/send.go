package mail

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"

	"github.com/prabhatlabs/internals/lib/env"
)

func SendMail(recipient string, subject string, body string) {
	serverAddr := getServerAddr()
	auth := getAuth()
	senderEmail := env.Vars.EMAIL_USER

	var msg bytes.Buffer

	// Every header field must end with \r\n
	// it's the standard, idk why!
	msg.WriteString(fmt.Sprintf("From: %s\r\n", senderEmail))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", recipient))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/html; charset=\"utf-8\"\r\n")

	// Add an empty line to separate headers from the message body
	msg.WriteString("\r\n")

	msg.WriteString(body)

	// Supports multiple recipients via the slice, might implement later!
	err := smtp.SendMail(serverAddr, auth, senderEmail, []string{recipient}, msg.Bytes())
	if err != nil {
		log.Fatalf("failed to send mail: %v", err)
	}

	log.Printf("sent mail to %s", recipient)
}
