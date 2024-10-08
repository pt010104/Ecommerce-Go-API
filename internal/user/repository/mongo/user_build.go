package mongo

import (
	"context"
	"time"

	"crypto/rand"
	"encoding/base64"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
func (impl implRepo) buildUserModel(context context.Context, opt user.CreateUserOption) (models.User, error) {
	u := models.User{
		ID:       primitive.NewObjectID(),
		Email:    opt.Email,
		UserName: opt.UserName,
		Password: opt.Password,
	}
	return u, nil

}
func (impl implRepo) buildKeyTokenModel(context context.Context, userId primitive.ObjectID, sessionID string) (models.KeyToken, error) {
	secretKey, err := generateRandomString(32)
	if err != nil {
		return models.KeyToken{}, err
	}

	refreshToken, err := generateRandomString(64)
	if err != nil {
		return models.KeyToken{}, err
	}

	u := models.KeyToken{
		ID:           primitive.NewObjectID(),
		UserID:       userId,
		SecretKey:    secretKey,
		RefreshToken: refreshToken,
		SessionID:    sessionID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	return u, nil

}

func (impl implRepo) buildRequestTokenModel(context context.Context, userId primitive.ObjectID, token string) (models.RequestToken, error) {

	u := models.RequestToken{
		ID:        primitive.NewObjectID(),
		Token:     token,
		UserID:    userId,
		CreatedAt: time.Now(),
		UpdateAt:  time.Now(),
		Is_Used:   false,
	}
	return u, nil
}
