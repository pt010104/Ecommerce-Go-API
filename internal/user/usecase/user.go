package usecase

import (
	"context"
	"sync"
	"time"

	"os"

	"crypto/rand"
	"encoding/base64"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/user"
	"github.com/pt010104/api-golang/pkg/jwt"
	"github.com/pt010104/api-golang/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func GenerateRandomString(n int) (string, error) {
	// Create a byte slice of size n
	bytes := make([]byte, n)

	// Read random bytes into the slice
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Encode the byte slice to a base64 string
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (uc implUsecase) CreateUser(ctx context.Context, uct user.CreateUserInput) (models.User, error) {
	err := uc.validateDataUser(ctx, uct.Email, uct.Password)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.CreateUser.repo.Create.validateDataUser: %v", err)
		return models.User{}, err
	}

	u, err := uc.repo.GetUser(ctx, user.GetUserOption{
		Email: uct.Email,
	})
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.CreateUser.repo.GetUser: %v", err)
		return models.User{}, err
	}

	if u.ID != primitive.NilObjectID {
		uc.l.Errorf(ctx, "user.usecase.CreateUser: user already exists")
		return models.User{}, user.ErrEmailExisted
	}

	hashedPass, err := uc.hashPassword(uct.Password)
	if err != nil {
		uc.l.Errorf(ctx, "error during  hashing pass : %v", err)
		return models.User{}, err
	}

	nu, err := uc.repo.CreateUser(ctx, user.CreateUserOption{
		Email:    uct.Email,
		Password: hashedPass,
		UserName: uct.UserName,
	})
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.CreateUser.repo.Create: %v", err)
		return models.User{}, err
	}
	return nu, nil

}

func (uc implUsecase) SignIn(ctx context.Context, sit user.SignInType) (user.SignInOutput, error) {
	u, err := uc.repo.GetUser(ctx, user.GetUserOption{
		Email: sit.Email,
	})
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.SignIn.GetUser: %v", err)
		return user.SignInOutput{}, err
	}

	if !u.IsVerified {
		uc.l.Errorf(ctx, "user.usecase.SignIn: user is not verified")
		return user.SignInOutput{}, user.ErrUserNotVerified
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(sit.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			uc.l.Errorf(ctx, "user.usecase.SignIn.CompareHashAndPassword: %v", err)
			return user.SignInOutput{}, err
		}
		uc.l.Errorf(ctx, "error comparing passwords: %v", err)
		return user.SignInOutput{}, err
	}

	var wg sync.WaitGroup
	var wgErr error
	var sessionId, secretKey, refreshToken string

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = uc.repo.DeleteKeyToken(ctx, u.ID.Hex(), sit.SessionID)
		if err != nil {
			uc.l.Errorf(ctx, "user.usecase.SignIn.DeleteKeyToken : ", err)
			wgErr = err
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		sessionId, secretKey, refreshToken, err = uc.createKeyToken(ctx)
		if err != nil {
			uc.l.Errorf(ctx, "user.usecase.SignIn.createKeyToken: %v", err)
			wgErr = err
			return
		}
	}()

	wg.Wait()

	if wgErr != nil {
		return user.SignInOutput{}, wgErr
	}

	kt, err := uc.repo.CreateKeyToken(ctx, user.CreateKeyTokenOption{
		UserID:       u.ID,
		SessionID:    sessionId,
		SecretKey:    secretKey,
		RefrestToken: refreshToken,
	})
	if err != nil {
		uc.l.Errorf(ctx, "error during finding matching user: %v", err)
		return user.SignInOutput{}, err
	}

	payload := jwt.Payload{
		UserID:    u.ID.Hex(),
		Refresh:   false,
		SessionID: sessionId,
	}

	expirationTime := time.Hour * 24
	t, err := jwt.Sign(payload, expirationTime, kt.SecretKey)
	if err != nil {
		uc.l.Errorf(ctx, "error signing token: %v", err)
		return user.SignInOutput{}, err
	}

	token := user.Token{
		AccessToken:  t,
		RefreshToken: refreshToken,
	}

	return user.SignInOutput{
		User:      u,
		Token:     token,
		SessionID: sessionId,
	}, nil
}
func (uc implUsecase) ForgetPasswordRequest(ctx context.Context, email string) (token string, err error) {
	u, err := uc.repo.GetUser(ctx, user.GetUserOption{
		Email: email,
	})
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.ForgetPasswordRequest: %v", err)
		return "", err
	}
	payload := jwt.Payload{
		UserID:  u.ID.Hex(),
		Refresh: false,
		Type:    "reset-request",
	}
	expirationTime := time.Hour * 1
	token, err = jwt.Sign(payload, expirationTime, os.Getenv("SUPER_SECRET_KEY"))
	if err != nil {
		uc.l.Errorf(ctx, "error signing token: %v", err)
		return "", err
	}

	wg := sync.WaitGroup{}
	var wgErr error

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = uc.emailUC.SendVerificationEmail(u.Email, token)
		if err != nil {
			uc.l.Errorf(ctx, "user.usecase.ForgetPasswordRequest: %v", err)
			wgErr = err
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err = uc.repo.CreateRequestToken(ctx, u.ID, token)
		if err != nil {
			uc.l.Errorf(ctx, "user.usecase.Forgetpasswordrequest.CreateRequestToken: ", err)
			wgErr = err
			return
		}
	}()

	wg.Wait()
	if wgErr != nil {
		return "", wgErr
	}

	return token, nil

}
func (uc implUsecase) VerifyRequest(ctx context.Context, email string) (token string, err error) {
	u, err := uc.repo.GetUser(ctx, user.GetUserOption{
		Email: email,
	})
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.Verify: %v", err)
		return "", err
	}
	payload := jwt.Payload{
		UserID:  u.ID.Hex(),
		Refresh: false,
		Type:    "verify",
	}
	expirationTime := time.Hour * 1
	token, err = jwt.Sign(payload, expirationTime, os.Getenv("SUPER_SECRET_KEY"))
	if err != nil {
		uc.l.Errorf(ctx, "error signing token user.usecase.verify: %v", err)
		return "", err
	}
	wg := sync.WaitGroup{}
	var wgErr error

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = uc.emailUC.SendVerificationEmail(u.Email, token)
		if err != nil {
			uc.l.Errorf(ctx, "user.usecase.ForgetPasswordRequest: %v", err)
			wgErr = err
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err = uc.repo.CreateRequestToken(ctx, u.ID, token)
		if err != nil {
			uc.l.Errorf(ctx, "user.usecase.Forgetpasswordrequest.CreateRequestToken: ", err)
			wgErr = err
			return
		}
	}()

	wg.Wait()
	if wgErr != nil {
		return "", wgErr
	}

	return token, nil

}

func (uc implUsecase) Detail(ctx context.Context, sc models.Scope, id string) (models.User, error) {
	u, err := uc.repo.DetailUser(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			uc.l.Warnf(ctx, "user.usecase.Detail.repo.DetailUser: %v", err)
			return models.User{}, user.ErrUserNotFound
		}
		uc.l.Errorf(ctx, "user.usecase.Detail.repo.DetailUser: %v", err)
		return models.User{}, err
	}
	uc.repo.DeleteKeyToken(ctx, sc.UserID, sc.SessionID)
	return u, nil
}

func (uc implUsecase) LogOut(ctx context.Context, sc models.Scope) {
	err := uc.repo.DeleteKeyToken(ctx, sc.UserID, sc.SessionID)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.LogOut.repo.DeleteKeyToken")
	}

}
func (uc implUsecase) ResetPassWord(ctx context.Context, input user.ResetPasswordInput) error {
	u, err := uc.repo.GetUser(ctx, user.GetUserOption{
		ID: input.UserId,
	})
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.ResetPassword.GetUser:", user.ErrUserNotFound)
		return err
	}

	rt, err := uc.repo.DetailRequestToken(ctx, input.Token)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.ResetPassword.DetailRequestToken: %v", err)
		return err
	}

	if rt.IsUsed {
		uc.l.Errorf(ctx, "user.usecase.ResetPassword: token is already used")
		return user.ErrTokenUsed
	}

	newHashesPass, err := uc.hashPassword(input.NewPass)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.user.ResetPassword.hashnewpass:", err)
		return err
	}

	opt := user.UpdateUserOption{
		Model:    u,
		Password: newHashesPass,
	}
	_, err = uc.repo.UpdateUser(ctx, opt)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.ResetPassword.UpdateUser: %v", err)
		return err
	}

	isUsed := true
	err = uc.repo.UpdateRequestToken(ctx, user.UpdateRequestTokenOption{
		Token:  input.Token,
		IsUsed: &isUsed,
	})
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.ResetPassword.UpdateRequestToken: %v", err)
		return err
	}

	return nil
}
func (uc implUsecase) VerifyUser(ctx context.Context, input user.VerifyUserInput) error {
	u, err := uc.repo.DetailUser(ctx, input.UserId)
	if err != nil {
		uc.l.Errorf(ctx, " user.usecase.Verify.getuser:", err)
		return err
	}
	rt, err := uc.repo.DetailRequestToken(ctx, input.Token)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.ResetPassword.DetailRequestToken: %v", err)
		return err
	}

	if rt.IsUsed {
		uc.l.Errorf(ctx, "user.usecase.ResetPassword: token is already used")
		return user.ErrTokenUsed
	}
	opt := user.UpdateUserOption{
		Model:      u,
		IsVerified: true,
	}
	_, err = uc.repo.UpdateUser(ctx, opt)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.VerifyUser.UpdateUser: %v", err)
		return err
	}
	isUsed := true
	err = uc.repo.UpdateRequestToken(ctx, user.UpdateRequestTokenOption{
		Token:  input.Token,
		IsUsed: &isUsed,
	})
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.ResetPassword.UpdateRequestToken: %v", err)
		return err
	}

	return nil

}
func (uc implUsecase) DistributeNewToken(ctx context.Context, input user.DistributeNewTokenInput) (output user.DistributeNewTokenOutPut, er error) {
	u, err := uc.repo.DetailUser(ctx, input.UserId)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.DistributeNewToken.DetailUser :", err)
		return user.DistributeNewTokenOutPut{}, err
	}
	kt, err := uc.repo.DetailKeyToken(ctx, u.ID.Hex(), input.SessionID)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.DistributeNewToken.DetailKeyToken :", err)
		return user.DistributeNewTokenOutPut{}, err
	}
	if kt.RefreshToken != input.RefreshToken {
		uc.l.Errorf(ctx, "user.usecase.DistributeNewToken.Kt.RefreshToken!=input :")
		return user.DistributeNewTokenOutPut{}, user.ErrRefreshTokenIsNotValid
	}

	if time.Since(kt.UpdatedAt) > time.Hour {
		uc.l.Errorf(ctx, "user.usecase.DistributeNewToken.RefreshToken expired")
		return user.DistributeNewTokenOutPut{}, user.ErrRefreshTokenIsExpired
	}
	newRefreshToken, err1 := GenerateRandomString(32)
	if err1 != nil {
		return user.DistributeNewTokenOutPut{}, err1
	}
	uc.repo.UpdateKeyToken(ctx, user.UpdateKeyTokenInput{
		ID:           kt.ID,
		RefreshToken: newRefreshToken,
		UpdatedAt:    time.Now(),
	})
	payload := jwt.Payload{
		UserID:    u.ID.Hex(),
		Refresh:   false,
		SessionID: kt.SessionID,
	}

	expirationTime := time.Hour * 24
	t, err := jwt.Sign(payload, expirationTime, kt.SecretKey)
	if err != nil {
		uc.l.Errorf(ctx, "error signing token: %v", err)
		return user.DistributeNewTokenOutPut{}, err
	}
	token := user.Token{
		AccessToken:  t,
		RefreshToken: newRefreshToken,
	}
	return user.DistributeNewTokenOutPut{

		JWT:          token.AccessToken,
		RefreshToken: token.RefreshToken,
		UserID:       u.ID.String(),
	}, nil

}
