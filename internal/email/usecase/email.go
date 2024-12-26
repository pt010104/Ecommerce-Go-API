package usecase

import (
	"context"
	"fmt"
	"os"

	"github.com/pt010104/api-golang/internal/resources"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func (uc implUsecase) SendVerificationEmail(userEmail string, verificationToken string) error {
	ctx := context.Background()
	client := getClient()

	service, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("unable to retrieve Gmail client: %v", err)
	}

	message := createMessage(os.Getenv("EMAIL_USERNAME"), userEmail, "Verify Your Email", verificationToken)

	_, err = service.Users.Messages.Send("me", message).Do()
	if err != nil {
		return fmt.Errorf("unable to send email: %v", err)
	}

	return nil
}

func (uc implUsecase) SendResetPasswordEmail(userEmail string, verificationToken string) error {
	ctx := context.Background()
	client := getClient()

	service, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("unable to retrieve Gmail client: %v", err)
	}

	message := createResetPassMessage(os.Getenv("EMAIL_USERNAME"), userEmail, "Reset your password", verificationToken)

	_, err = service.Users.Messages.Send("me", message).Do()
	if err != nil {
		return fmt.Errorf("unable to send email: %v", err)
	}

	return nil
}

func (uc implUsecase) SendEmail(ctx context.Context, data resources.EmailData) error {
	client := getClient()

	service, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("unable to retrieve Gmail client: %v", err)
	}

	emailFrom := os.Getenv("EMAIL_USERNAME")
	message := createEmailMessage(emailFrom, data)

	_, err = service.Users.Messages.Send("me", message).Do()
	if err != nil {
		return fmt.Errorf("unable to send email: %v", err)
	}

	return nil
}
