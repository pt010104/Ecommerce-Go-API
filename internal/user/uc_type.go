package user

type UseCaseType struct {
	UserName string
	Password string
	Email    string
}

type SignInType struct {
	Email     string
	Password  string
	SessionID string
}

type ForgetPasswordRequest struct {
	Email string
}
type ResetPassWordReq struct {
	NewPassword string
}
