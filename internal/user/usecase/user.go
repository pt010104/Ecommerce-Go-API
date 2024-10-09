package usecase

import (
	"context"
	"time"

	"crypto/rand"
	"encoding/base64"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/user"
	"github.com/pt010104/api-golang/pkg/jwt"
	"github.com/pt010104/api-golang/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"os"

	"golang.org/x/crypto/bcrypt"
	"log"
)

func generateRandomString(size int) (string, error) {
	// Create a byte slice with the required size
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// Encode the byte slice to base64 to get a safe string
	return base64.URLEncoding.EncodeToString(b), nil
}
func (uc implUsecase) CreateUser(ctx context.Context, uct user.UseCaseType) (models.User, error) {
	err := uc.validateDataCreateUser(ctx, uct.Email)
	if err != nil {
		uc.l.Errorf(ctx, "error during validate data: %v", err)
		return models.User{}, err
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

func (uc implUsecase) SignIn(ctx context.Context, sit user.SignInType) (models.User, string, string, error) {
	u, err := uc.repo.GetUser(ctx, user.GetUserOption{
		Email: sit.Email,
	})
	if err != nil {
		uc.l.Errorf(ctx, "error during finding matching user: %v", err)
		return models.User{}, "", "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(sit.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			uc.l.Warnf(ctx, "password mismatch for user: %v", u.Email)
			return models.User{}, "", "", err
		}
		uc.l.Errorf(ctx, "error comparing passwords: %v", err)
		return models.User{}, "", "", err
	}

	sessionId, err2 := generateRandomString(32)
	if err2 != nil {
		uc.l.Errorf(ctx, "user.usecase.SignIn.generateRandomstring: %v", err2)
		return models.User{}, "", "", err2

	}

	e1 := uc.repo.DeleteRecord(ctx, u.ID.Hex(), sit.SessionID)
	if e1 != nil {
		uc.l.Errorf(ctx, "user.usecase.signin.deleterecord : ", e1)
		println("delete")
		return models.User{}, "", "", e1
	}
	kt, err := uc.repo.CreateKeyToken(ctx, u.ID, sessionId)
	if err != nil {
		uc.l.Errorf(ctx, "error during finding matching user: %v", err)
		return models.User{}, "", "", err
	}

	payload := jwt.Payload{
		UserID:    u.ID.Hex(),
		Refresh:   false,
		SessionID: sessionId,
	}

	expirationTime := time.Hour * 24
	token, err := jwt.Sign(payload, expirationTime, kt.SecretKey)
	if err != nil {
		uc.l.Errorf(ctx, "error signing token: %v", err)
		return models.User{}, "", "", err
	}

	return u, sessionId, token, nil
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
	uc.repo.DeleteRecord(ctx, sc.UserID, sc.SessionID)
	return u, nil
}
func (uc implUsecase) LogOut(ctx context.Context, sc models.Scope) {
	err := uc.repo.DeleteRecord(ctx, sc.UserID, sc.SessionID)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.LogOut.repo.DeleteRecord")
	}

}
func (uc implUsecase) ResetPassWord(ctx context.Context, UserId string, newPass string) error {
	u, err := uc.repo.GetUser(ctx, user.GetUserOption{
		ID: UserId,
	})
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.ResetPassword.GetUser:", err)
		return err
	}

	newHashesPass, err := uc.hashPassword(newPass)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.user.ResetPassword.hashnewpass:", err)
		return err
	}
	updateData := bson.M{
		"password": newHashesPass,
	}
	err = uc.repo.UpdateRecord(ctx, u.ID.String(), updateData)
	if err != nil {
		log.Fatalf("Error updating record: %v", err)
		return err
	}

	return nil
}
func (uc implUsecase) MartJWTasUsed(ctx context.Context, JWT string) error {
	updateData := bson.M{
		"is_used": true,
	}
	err := uc.repo.UpdateRequestTokenRecord(ctx, JWT, updateData)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.MarkJWTASused:", err)
		return err
	}
	return nil
}

func (uc implUsecase) IsJWTresetVaLID(ctx context.Context, JWT string) (bool, error) {

	valid, err := uc.repo.IsJWTresetVaLID(ctx, JWT)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.IsJwtValid", err)
	}
	return valid, nil
}
