package services

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

func SendAlertTelegram(message string) error {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	token := os.Getenv("TELEGRAM_TOKEN")
	chat_id := os.Getenv("TELEGRAM_CHAT_ID")

	escapedMessage := url.QueryEscape(message)
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s", token, chat_id, escapedMessage)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error sending message to Telegram: %v\n", err)
		return err
	}
	defer resp.Body.Close()
	fmt.Println("Message sent to Telegram successfully")
	return nil
}
