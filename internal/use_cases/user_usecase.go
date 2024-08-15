package use_cases

import (
	dto "Backend_golang_project/internal/domain/dto/user"
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
}

type UserService struct {
	userRepository repositories.IUserRepository
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
		log.Error("Failed to create user:", err)
		return nil, err
	}

	log.Info("User created successfully")
	return data, nil
}

func (u UserService) GetUserByID(ctx context.Context, ID int) (*entities.User, error) {
	user, err := u.userRepository.GetUserById(ctx, ID)
	if err != nil {
		log.Error("Failed to get user by ID:", err)
		return nil, err
	}

	return user, nil
}

func NewUserService(userRepository repositories.IUserRepository) IUserService {
	return &UserService{
		userRepository: userRepository,
	}
}
