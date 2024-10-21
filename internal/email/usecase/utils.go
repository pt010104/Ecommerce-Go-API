package usecase

import (
	"context"
	"encoding/base64"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"net/http"
	"os"
)

func getClient() *http.Client {
	config := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URI"),
		Endpoint:     google.Endpoint,
	}

	token := &oauth2.Token{
		RefreshToken: os.Getenv("GOOGLE_REFRESH_TOKEN"),
	}

	client := config.Client(context.Background(), token)
	return client
}

func createMessage(from, to, subject, token string) *gmail.Message {

	clientBaseURL := os.Getenv("CLIENT_BASE_URL")
	if clientBaseURL == "" {
		clientBaseURL = "http://localhost:3000"
	}

	verificationURL := fmt.Sprintf("%s/verify-email?token=%s", clientBaseURL, token)

	message := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\nPlease verify your email by clicking this link: %s", to, subject, verificationURL))

	var msg gmail.Message
	msg.Raw = base64.URLEncoding.EncodeToString(message)

	return &msg
}
func createResetPassMessage(from, to, subject, token string) *gmail.Message {

	clientBaseURL := os.Getenv("CLIENT_BASE_URL")
	if clientBaseURL == "" {
		clientBaseURL = "http://localhost:3000"
	}

	resetPasswordURL := fmt.Sprintf("%s/reset-password?token=%s", clientBaseURL, token)

	message := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\nPlease reset your password by clicking this link: %s", to, subject, resetPasswordURL))

	var msg gmail.Message
	msg.Raw = base64.URLEncoding.EncodeToString(message)

	return &msg
}
