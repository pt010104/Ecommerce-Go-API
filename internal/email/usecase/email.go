package usecase

import (
	"context"
	"fmt"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"os"
)

func (uc implUsecase) SendVerificationEmail(userEmail string, verificationToken string) error {
	ctx := context.Background()
	client := getClient()

	// Create Gmail service
	service, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("unable to retrieve Gmail client: %v", err)
	}

	// Create the email message
	message := createMessage(os.Getenv("EMAIL_USERNAME"), userEmail, "Verify Your Email", verificationToken)

	// Send the email
	_, err = service.Users.Messages.Send("me", message).Do()
	if err != nil {
		return fmt.Errorf("unable to send email: %v", err)
	}

	return nil
}
