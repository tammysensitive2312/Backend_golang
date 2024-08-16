package use_cases

import (
	"Backend_golang_project/infrastructure/config"
	"Backend_golang_project/infrastructure/middleware/jwt"
	dto "Backend_golang_project/internal/domain/dto/request"
	"Backend_golang_project/internal/domain/entities"
	"Backend_golang_project/internal/pkg"
	"Backend_golang_project/internal/repositories"
	"context"
	"github.com/go-playground/validator/v10"
	"time"

	log "github.com/sirupsen/logrus"
)

type IUserService interface {
	Create(ctx context.Context, request *dto.CreateUserRequest) (*entities.User, error)
	GetUserByID(ctx context.Context, ID int) (*entities.User, error)
	Login(cxt context.Context, req dto.LoginRequest) (string, string, error)
	RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (string, error)
}

type UserService struct {
	config         *config.Config
	userRepository repositories.IUserRepository
}

func (u UserService) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (string, error) {
	claims, err := jwt.ClaimRefreshToken(req.RefreshToken, u.config)
	if err != nil {
		log.Error("Invalid refresh token")
		return "", err
	}

	newAccessToken, _, err := jwt.GenerateJwtToken(u.config, claims.ID)
	if err != nil {
		log.Error("Failed to generate new access token")
		return "", err
	}
	return newAccessToken, nil
}

func (u UserService) Create(ctx context.Context, request *dto.CreateUserRequest) (*entities.User, error) {
	v := validator.New()
	newUser, err := request.ToUserEntity(v)
	if err != nil {
		log.Error(err)
	}

	hashPassword, err := pkg.HashPassword(request.Password)
	if err != nil {
		log.Error("Hash password fail")
		return nil, err
	}
	newUser.Password = hashPassword

	now := time.Now()
	newUser.CreatedAt = now
	newUser.UpdatedAt = now

	data, err := u.userRepository.CreateUser(ctx, newUser)
	if err != nil {
		log.Error("Failed to create response:", err)
		return nil, err
	}

	log.Info("User created successfully")
	return data, nil
}

func (u UserService) GetUserByID(ctx context.Context, ID int) (*entities.User, error) {
	user, err := u.userRepository.GetUserById(ctx, ID)
	if err != nil {
		log.Error("Failed to get response by ID:", err)
		return nil, err
	}

	return user, nil
}

func (u *UserService) Login(ctx context.Context, req dto.LoginRequest) (string, string, error) {
	user, err := u.userRepository.GetUserByEmail(ctx, req.Email)
	if err != nil || user == nil {
		log.Error("Invalid login")
		return "", "", err
	}
	ok := pkg.CheckPasswordHash(req.Password, user.Password)
	if !ok {
		log.Error("Invalid password")
		return "", "", err
	}

	accessToken, refreshToken, err := jwt.GenerateJwtToken(u.config, user.ID)
	if err != nil {
		log.Error("Cannot generate token ", err)
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func NewUserService(config *config.Config, userRepository repositories.IUserRepository) IUserService {
	return &UserService{
		config:         config,
		userRepository: userRepository,
	}
}
