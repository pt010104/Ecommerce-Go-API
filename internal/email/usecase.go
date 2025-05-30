package email

import (
	"context"

	"github.com/pt010104/api-golang/internal/resources"
)

type UseCase interface {
	SendEmail(ctx context.Context, data resources.EmailData) error
	SendVerificationEmail(userEmail string, verificationToken string) error
	SendResetPasswordEmail(userEmail string, verificationToken string) error
}
