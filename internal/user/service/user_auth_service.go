package service

import (
	"context"
	"enigmanations/cats-social/internal/user"
	"enigmanations/cats-social/internal/user/errs"
	"enigmanations/cats-social/internal/user/repository"
	"enigmanations/cats-social/internal/user/request"
	"enigmanations/cats-social/internal/user/response"
	"enigmanations/cats-social/pkg/bcrypt"
	"enigmanations/cats-social/pkg/jwt"
	"fmt"

	emailverifier "github.com/AfterShip/email-verifier"
)

type UserAuthService interface {
	Login(req *request.UserLoginRequest) (*response.UserLoginResponse, error)
}

type userAuthService struct {
	userDB  repository.UserRepository
	authDB  repository.UserAuthRepository
	Context context.Context
}

func NewUserAuthService(userDB repository.UserRepository, authDB repository.UserAuthRepository, ctx context.Context) UserAuthService {
	return &userAuthService{userDB: userDB, authDB: authDB, Context: ctx}
}

func (service *userAuthService) Login(req *request.UserLoginRequest) (*response.UserLoginResponse, error) {
	var userCredential *user.User

	// Check email
	if req.Email != "" {
		var verifier = emailverifier.NewVerifier()
		ret, err := verifier.Verify(req.Email)
		if err != nil {
			return nil, errs.UserErrEmailInvalidFormat
		}
		if !ret.Syntax.Valid {
			return nil, errs.UserErrEmailInvalidFormat
		}

		// Get user
		userCredentialFound, err := service.userDB.GetByEmail(service.Context, req.Email)
		if err != nil {
			return nil, errs.UserErrNotFound
		}

		userCredential = userCredentialFound
	}

	// Check password
	if req.Password != "" {
		if !bcrypt.CheckPasswordHash(req.Password, userCredential.Password) {
			return nil, errs.WrongPassword
		}
	}

	// Create or get session
	var userSession *user.UserSession
	userSessionFound, _ := service.authDB.GetIfExists(service.Context, userCredential.Id)
	if userSessionFound != nil {
		userSession = userSessionFound
	} else {
		userSessionCreated, err := service.authDB.Save(service.Context, userCredential)
		if err != nil {
			return nil, err
		}
		userSession = userSessionCreated
	}

	// Generate access token
	accessToken, err := jwt.GenerateAccessToken(uint64(userSession.UserId), userCredential)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return response.UserToUserLoginResponse(*userCredential, accessToken), nil
}
