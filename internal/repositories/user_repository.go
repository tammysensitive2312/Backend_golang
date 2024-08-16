package repositories

import (
	"Backend_golang_project/internal/domain/entities"
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, user *entities.User) (*entities.User, error)
	GetUserById(ctx context.Context, ID int) (*entities.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
}

type UserRepository struct {
	base
}

func (u UserRepository) CreateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	if err := u.db.WithContext(ctx).Create(user).Error; err != nil {
		log.Error("Cannot create response with err:", err.Error())
		return nil, err
	}
	return user, nil
}

func (u UserRepository) GetUserById(ctx context.Context, ID int) (*entities.User, error) {
	var user entities.User

	if err := u.db.WithContext(ctx).Preload("Projects").First(&user, ID).Error; err != nil {
		log.Error("Can not find response with ID: ", ID, err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("response %v not found", user)
		}
		return nil, fmt.Errorf("error retrieving response: %w", err)
	}
	return &user, nil
}

func (u UserRepository) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User

	if err := u.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		log.Error("Can not find response with email: ", email, err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("response %v not found", user)
		}
		return nil, fmt.Errorf("error retrieving response: %w", err)
	}
	return &user, nil
}

// NewUserRepository constructor
func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{base: base{db: db}}
}
