package usecase

import (
	"context"
	"sync"
	"time"

	"os"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/user"
	"github.com/pt010104/api-golang/pkg/jwt"
	"github.com/pt010104/api-golang/pkg/mongo"

	"golang.org/x/crypto/bcrypt"
)

func (uc implUsecase) CreateUser(ctx context.Context, uct user.CreateUserInput) (models.User, error) {
	err := uc.validateDataUser(ctx, uct.Email, uct.Password)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.CreateUser.repo.Create.validateDataUser: %v", err)
		return models.User{}, err
	}

	_, err = uc.repo.GetUser(ctx, user.GetUserOption{
		Email: uct.Email,
	})
	if err != nil {
		if err != mongo.ErrNoDocuments {
			uc.l.Errorf(ctx, "error during finding matching user: %v", err)
			return models.User{}, user.ErrEmailExisted
		}
	}

	hashedPass, err := uc.hashPassword(uct.Password)
	if err != nil {
		uc.l.Errorf(ctx, "error during  hashing pass : %v", err)
		return models.User{}, err
	}

	u, err := uc.repo.CreateUser(ctx, user.CreateUserOption{
		Email:    uct.Email,
		Password: hashedPass,
		UserName: uct.UserName,
	})
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.CreateUser.repo.Create: %v", err)
		return models.User{}, err
	}
	return u, nil

}

func (uc implUsecase) SignIn(ctx context.Context, sit user.SignInType) (user.SignInOutput, error) {
	u, err := uc.repo.GetUser(ctx, user.GetUserOption{
		Email: sit.Email,
	})
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.SignIn.GetUser: %v", err)
		return user.SignInOutput{}, err
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
	err1 := uc.emailUC.SendVerificationEmail(u.Email, token)
	if err1 != nil {
		uc.l.Errorf(ctx, "user.usecase.ForgetPasswordRequest: %v", err1)
	}
	_, err2 := uc.repo.CreateRequestToken(ctx, u.ID, token)
	if err2 != nil {
		uc.l.Errorf(ctx, "user.usecase.Forgetpasswordrequest.CreateRequestToken: ", err2)
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
