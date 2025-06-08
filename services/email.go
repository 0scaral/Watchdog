package services

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func SendAlertMail(text string) {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	email := gomail.NewMessage()

	email.SetHeader("From", os.Getenv("EMAIL_SRC"))
	email.SetHeader("To", os.Getenv("EMAIL_DST"))
	email.SetHeader("Subject", "Watchdog Alert")

	email.SetBody("text/plain", text)

	dialer := gomail.NewDialer(os.Getenv("SMTP_SERVER"), 587, os.Getenv("EMAIL_SRC"), os.Getenv("EMAIL_PASSWD"))

	if err := dialer.DialAndSend(email); err != nil {
		fmt.Println("Error sending email:", err)
	} else {
		fmt.Println("Alert email sent successfully")
	}
}
