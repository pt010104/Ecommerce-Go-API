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
