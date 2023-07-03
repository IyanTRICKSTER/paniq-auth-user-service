package tests

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"paniq-auth-user-service/pkg/contracts"
	"paniq-auth-user-service/pkg/domains/notification"
	"testing"
)

var notificationUsecase contracts.INotificationService

func init() {

	//Load .env file
	if err := godotenv.Load("../../../../.env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

}

func TestNotifyWithEmail(t *testing.T) {

	notificationUsecase = notification.NewUsecase(
		os.Getenv("SMTP_HOST"),
		os.Getenv("SMTP_PORT"),
		os.Getenv("SMTP_USERNAME"),
		os.Getenv("SMTP_PASSWORD"),
	)

	notificationUsecase.NotifyWithEmail(
		"iyan@gmail.com",
		"ucok@gmail.com",
		"test",
		"hello world")
}
