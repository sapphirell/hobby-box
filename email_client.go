package main

import (
	"log"
	"strings"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

func sendMail() {
	// Set up authentication information.
	auth := sasl.NewPlainClient("", "username", "password")

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	to := []string{"recipient@example.net"}
	msg := strings.NewReader("To: recipient@example.net\r\n" +
		"Subject: discount Gophers!\r\n" +
		"\r\n" +
		"This is the email body.\r\n")
	log.Println("sending...")
	err := smtp.SendMail("localhost:1025", auth, "sender@example.org", to, msg)
	if err != nil {
		log.Fatal(err)
	}
	err1 := smtp.SendMail("localhost:1025", auth, "sender@example.org", to, msg)
	if err1 != nil {
		log.Fatal(err1)
	}
}
