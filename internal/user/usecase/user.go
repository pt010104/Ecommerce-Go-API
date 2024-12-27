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

	"github.com/pt010104/api-golang/pkg/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func (uc implUsecase) CreateUser(ctx context.Context, uct user.CreateUserInput) (models.User, error) {
	err := uc.validateDataUser(ctx, uct.Email, uct.Password)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.CreateUser.repo.Create.validateDataUser: %v", err)
		return models.User{}, err
	}

	u, err := uc.repo.GetUser(ctx, user.GetUserOption{
		GetFilter: user.GetFilter{
			Email: uct.Email,
		},
	})
	if err != nil {
		if err != mongo.ErrNoDocuments {
			uc.l.Errorf(ctx, "USER.usecase.createuser.getuser:")
			return models.User{}, err
		}
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
		Name:     uct.Name,
	})
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.CreateUser.repo.Create: %v", err)
		return models.User{}, err
	}
	return nu, nil

}

func (uc implUsecase) SignIn(ctx context.Context, sit user.SignInType) (user.SignInOutput, error) {
	if sit.SessionID == sessionIDTest {
		return user.SignInOutput{}, nil
	}

	u, err := uc.repo.GetUser(ctx, user.GetUserOption{
		GetFilter: user.GetFilter{
			Email: sit.Email,
		},
	})
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.SignIn.GetUser: %v", err)
		return user.SignInOutput{}, user.ErrMismatchedHashAndPassword
	}

	if !u.IsVerified {
		uc.l.Debugf(ctx, "userveriry:", u.IsVerified)
		uc.l.Errorf(ctx, "user.usecase.SignIn: user is not verified")
		return user.SignInOutput{}, user.ErrUserNotVerified
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(sit.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			uc.l.Errorf(ctx, "user.usecase.SignIn.CompareHashAndPassword: %v", err)
			return user.SignInOutput{}, user.ErrMismatchedHashAndPassword
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

	expirationTime := accessTokenExpireTime * 365
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
		GetFilter: user.GetFilter{
			Email: email,
		},
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
		err = uc.emailUC.SendResetPasswordEmail(u.Email, token)
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
func (uc implUsecase) VerifyEmail(ctx context.Context, email string) (token string, err error) {
	u, err := uc.repo.GetUser(ctx, user.GetUserOption{
		GetFilter: user.GetFilter{
			Email: email,
		},
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

func (uc implUsecase) Detail(ctx context.Context, sc models.Scope, id string) (user.DetailUserOutput, error) {
	var u models.User
	u, err := uc.redisRepo.DetailUser(ctx, id)
	if err != nil {
		uc.l.Warnf(ctx, "user.usecase.Detail.redis.DetailUser: %v", err)
		u, err = uc.repo.DetailUser(ctx, id)
		if err != nil {
			uc.l.Errorf(ctx, "user.usecase.Detail.DetailUser: %v", err)
			return user.DetailUserOutput{}, err
		}

		go func() {
			redisCtx := context.Background()
			if err := uc.redisRepo.StoreUser(redisCtx, u); err != nil {
				uc.l.Warnf(ctx, "user.usecase.Detail.redis.StoreUser: %v", err)
			}
		}()
	}

	var avatar models.Media
	if u.MediaID != primitive.NilObjectID {
		avatar, err = uc.mediaUC.Detail(ctx, models.Scope{}, u.MediaID.Hex())
		if err != nil {
			uc.l.Errorf(ctx, "user.usecase.Detail.Detail: %v", err)
		}
	}

	return user.DetailUserOutput{
		User:       u,
		Avatar_URL: avatar.URL,
	}, nil
}
func (uc implUsecase) GetModel(ctx context.Context, id string) (models.User, error) {
	u, err := uc.redisRepo.DetailUser(ctx, id)
	if err != nil {
		uc.l.Warnf(ctx, "user.usecase.GetModel.redis.DetailUser: %v", err)
		u, err = uc.repo.DetailUser(ctx, id)
		if err != nil {
			uc.l.Errorf(ctx, "user.usecase.GetModel.DetailUser: %v", err)
			if err == mongo.ErrNoDocuments {
				return models.User{}, user.ErrUserNotFound
			}

			return models.User{}, err
		}

		go func() {
			redisCtx := context.Background()
			if err := uc.redisRepo.StoreUser(redisCtx, u); err != nil {
				uc.l.Warnf(ctx, "user.usecase.GetModel.redis.StoreUser: %v", err)
			}
		}()

		return u, nil
	}

	return u, nil
}

func (uc implUsecase) LogOut(ctx context.Context, sc models.Scope) error {
	err := uc.repo.DeleteKeyToken(ctx, sc.UserID, sc.SessionID)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.LogOut.repo.DeleteKeyToken")
		return err
	}

	return nil
}
func (uc implUsecase) ResetPassWord(ctx context.Context, input user.ResetPasswordInput) error {
	u, err := uc.repo.GetUser(ctx, user.GetUserOption{
		GetFilter: user.GetFilter{
			ID: input.UserId,
		},
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
	_, err = uc.repo.UpdatePatchUser(ctx, opt)
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
		uc.l.Errorf(ctx, "user.usecase.ResetPassword.updateRequestToken: %v", err)
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
	_, err = uc.repo.UpdatePatchUser(ctx, opt)
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
		uc.l.Errorf(ctx, "user.usecase.ResetPassword.updateRequestToken: %v", err)
		return err
	}

	return nil

}
func (uc implUsecase) DistributeNewToken(ctx context.Context, input user.DistributeNewTokenInput) (output user.DistributeNewTokenOutput, er error) {
	if input.UserId == "" || input.SessionID == "" || input.RefreshToken == "" {
		uc.l.Errorf(ctx, "user.usecase.DistributeNewToken: invalid input")
		return user.DistributeNewTokenOutput{}, user.ErrInvalidInput
	}

	if input.SessionID == sessionIDTest {
		return user.DistributeNewTokenOutput{}, nil
	}

	u, err := uc.repo.DetailUser(ctx, input.UserId)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.DistributeNewToken.DetailUser :", err)
		return user.DistributeNewTokenOutput{}, err
	}

	kt, err := uc.repo.DetailKeyToken(ctx, u.ID.Hex(), input.SessionID)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.DistributeNewToken.DetailKeyToken :", err)
		return user.DistributeNewTokenOutput{}, err
	}

	if kt.RefreshToken != input.RefreshToken {
		uc.l.Errorf(ctx, "user.usecase.DistributeNewToken.Kt.RefreshToken!=input :")
		return user.DistributeNewTokenOutput{}, user.ErrRefreshTokenIsNotValid
	}

	if time.Since(kt.UpdatedAt) > refrestTokenExpireTime {
		uc.l.Errorf(ctx, "user.usecase.DistributeNewToken.RefreshToken expired")
		return user.DistributeNewTokenOutput{}, user.ErrRefreshTokenIsExpired
	}

	newRefreshToken, err := util.GenerateRandomString(32)
	if err != nil {
		return user.DistributeNewTokenOutput{}, err
	}

	err = uc.repo.UpdateKeyToken(ctx, user.UpdateKeyTokenInput{
		UserID:       u.ID.Hex(),
		SessionID:    input.SessionID,
		RefreshToken: newRefreshToken,
	})
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.DistributeNewToken.UpdateKeyToken :", err)
		return user.DistributeNewTokenOutput{}, err
	}

	payload := jwt.Payload{
		UserID:    u.ID.Hex(),
		Refresh:   false,
		SessionID: kt.SessionID,
	}
	expirationTime := accessTokenExpireTime
	t, err := jwt.Sign(payload, expirationTime, kt.SecretKey)
	if err != nil {
		uc.l.Errorf(ctx, "error signing token: %v", err)
		return user.DistributeNewTokenOutput{}, err
	}

	token := user.Token{
		AccessToken:  t,
		RefreshToken: newRefreshToken,
	}

	return user.DistributeNewTokenOutput{
		Token: token,
	}, nil

}

func (uc implUsecase) DetailKeyToken(ctx context.Context, userID string, sessionID string) (models.KeyToken, error) {

	key, err := uc.redisRepo.GetSecretKey(ctx, models.Scope{
		UserID:    userID,
		SessionID: sessionID,
	})
	if err != nil {
		uc.l.Warnf(ctx, "user.usecase.DetailKeyToken.GetSecretKey", err)
	}

	if len(key) > 0 {
		return models.KeyToken{SecretKey: string(key)}, nil
	}

	k, err := uc.repo.DetailKeyToken(ctx, userID, sessionID)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.DetailKeyToken.Repo", err)
		return models.KeyToken{}, err
	}

	go func() {
		redisCtx := context.Background()
		cacheErr := uc.redisRepo.StoreSecretKey(models.Scope{
			UserID:    userID,
			SessionID: sessionID,
		}, k.SecretKey, redisCtx)
		if cacheErr != nil {
			uc.l.Warnf(ctx, "user.usecase.DetailKeyToken.CacheUpdate", cacheErr)
		}
	}()

	return k, nil
}
func (uc implUsecase) Update(ctx context.Context, sc models.Scope, input user.UpdateInput) (user.DetailUserOutput, error) {
	var u models.User
	var avatar models.Media
	var wg sync.WaitGroup
	var wgErr error

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		u, err = uc.GetModel(ctx, sc.UserID)
		if err != nil {
			uc.l.Errorf(ctx, "user.usecase.Update.GetModel: %v", err)
			wgErr = err
			return
		}
	}()

	if input.MediaID != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error
			avatar, err = uc.mediaUC.Detail(ctx, sc, input.MediaID)
			if err != nil {
				uc.l.Errorf(ctx, "user.usecase.Update.repo.Detail: %v", err)
				wgErr = err
				return
			}
		}()
	}

	wg.Wait()
	if wgErr != nil {
		return user.DetailUserOutput{}, wgErr
	}

	if input.Email != "" && input.Email != u.Email {
		existingUser, err := uc.repo.GetUser(ctx, user.GetUserOption{
			GetFilter: user.GetFilter{
				Email: input.Email,
			},
		})
		if err != nil {
			if err != mongo.ErrNoDocuments {
				uc.l.Errorf(ctx, "user.usecase.Update.GetUser: %v", err)
				return user.DetailUserOutput{}, err
			}
		}

		if existingUser.ID != primitive.NilObjectID {
			uc.l.Errorf(ctx, "user.usecase.Update: email already exists")
			return user.DetailUserOutput{}, user.ErrEmailExisted
		}
	}

	nu, err := uc.repo.UpdateUser(ctx, user.UpdateUserOption{
		Model:   u,
		Email:   input.Email,
		Name:    input.Name,
		MediaID: input.MediaID,
	})
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.Update.repo.UpdateUser: %v", err)
		return user.DetailUserOutput{}, err
	}

	go func() {
		redisCtx := context.Background()
		err = uc.redisRepo.StoreUser(redisCtx, nu)
		if err != nil {
			uc.l.Warnf(ctx, "user.usecase.Detail.redis.StoreUser: %v", err)
		}
	}()

	return user.DetailUserOutput{
		User:       nu,
		Avatar_URL: avatar.URL,
	}, nil
}

func (uc implUsecase) AddAddress(ctx context.Context, sc models.Scope, input user.AddAddressInput) (user.DetailAddressOutput, error) {
	u, err := uc.GetModel(ctx, sc.UserID)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.AddAddress.GetModel: %v", err)
		return user.DetailAddressOutput{}, err
	}

	newAddress := models.Address{
		ID:       primitive.NewObjectID(),
		Street:   input.Street,
		District: input.District,
		City:     input.City,
		Province: input.Province,
		Phone:    input.Phone,
		Default:  input.Default,
	}

	if input.Default {
		for i := range u.Address {
			u.Address[i].Default = false
		}
	}

	u.Address = append(u.Address, newAddress)

	_, err = uc.repo.UpdateUser(ctx, user.UpdateUserOption{
		Model:   u,
		Address: u.Address,
	})
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.AddAddress.UpdateUser: %v", err)
		return user.DetailAddressOutput{}, err
	}

	go func() {
		redisCtx := context.Background()
		if err := uc.redisRepo.StoreUser(redisCtx, u); err != nil {
			uc.l.Warnf(ctx, "user.usecase.AddAddress.redis.StoreUser: %v", err)
		}
	}()

	return user.DetailAddressOutput{
		Addressess: u.Address,
	}, nil
}

func (uc implUsecase) UpdateAddress(ctx context.Context, sc models.Scope, input user.UpdateAddressInput) (user.DetailAddressOutput, error) {
	u, err := uc.GetModel(ctx, sc.UserID)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.UpdateAddress.GetModel: %v", err)
		return user.DetailAddressOutput{}, err
	}

	addressID, err := primitive.ObjectIDFromHex(input.AddressID)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.UpdateAddress.ObjectIDFromHex: %v", err)
		return user.DetailAddressOutput{}, err
	}

	addressIndex := -1
	for i, addr := range u.Address {
		if addr.ID == addressID {
			addressIndex = i
			break
		}
	}

	if addressIndex == -1 {
		uc.l.Errorf(ctx, "user.usecase.UpdateAddress: address not found")
		return user.DetailAddressOutput{}, user.ErrAddressNotFound
	}

	if input.Default {
		for i := range u.Address {
			u.Address[i].Default = false
		}
	}

	u.Address[addressIndex].Street = input.Street
	u.Address[addressIndex].District = input.District
	u.Address[addressIndex].City = input.City
	u.Address[addressIndex].Province = input.Province
	u.Address[addressIndex].Phone = input.Phone
	u.Address[addressIndex].Default = input.Default

	_, err = uc.repo.UpdateUser(ctx, user.UpdateUserOption{
		Model:   u,
		Address: u.Address,
	})
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.UpdateAddress.UpdateUser: %v", err)
		return user.DetailAddressOutput{}, err
	}

	go func() {
		redisCtx := context.Background()
		if err := uc.redisRepo.StoreUser(redisCtx, u); err != nil {
			uc.l.Warnf(ctx, "user.usecase.UpdateAddress.redis.StoreUser: %v", err)
		}
	}()

	return user.DetailAddressOutput{
		Addressess: u.Address,
	}, nil
}

func (uc implUsecase) ListUsers(ctx context.Context, sc models.Scope, input user.ListUserInput) ([]models.User, error) {
	users, err := uc.repo.ListUser(ctx, user.ListUserOption{
		GetFilter: user.GetFilter{
			IDs: input.IDs,
		},
	})
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.List.repo.ListUser: %v", err)
		return nil, err
	}

	return users, nil
}
