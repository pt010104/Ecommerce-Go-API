package email

type UseCase interface {
	SendVerificationEmail(userEmail string, verificationToken string) error
}
