package use_cases

import (
	"Backend_golang_project/infrastructure/config"
	"Backend_golang_project/infrastructure/middleware/jwt"
	dto "Backend_golang_project/internal/domain/dto/request"
	"Backend_golang_project/internal/domain/entities"
	"Backend_golang_project/internal/pkg"
	"Backend_golang_project/internal/repositories"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

type IUserService interface {
	Create(ctx context.Context, request *dto.CreateUserRequest) (*entities.User, error)
	GetUserByID(ctx context.Context, ID int) (*entities.User, error)
	Login(cxt context.Context, req dto.LoginRequest) (string, string, error)
	RefreshToken(req *dto.RefreshTokenRequest) (string, error)
	ExportToS3(ctx context.Context, filename string) error
}

type UserService struct {
	config         *config.Config
	userRepository repositories.IUserRepository
	s3repository   repositories.S3RepositoryInterface
}

func (u UserService) ExportToS3(ctx context.Context, filename string) error {
	bucket := u.config.S3Config.Bucket
	// Tạo file tạm thời để lưu dữ liệu CSV
	tempFile, err := os.CreateTemp("", "data-*.csv")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	writer := csv.NewWriter(tempFile)
	defer writer.Flush()

	// Ghi header của CSV (tùy thuộc vào cấu trúc dữ liệu của bạn)
	writer.Write([]string{"ID", "username", "email", "password", "created_at", "updated_at"}) // Thay đổi cho phù hợp với bảng của bạn

	batchSize := 1000
	offset := 0

	for {
		// Lấy dữ liệu từ cơ sở dữ liệu theo từng batch
		data, err := u.userRepository.GetUsersBatch(ctx, batchSize, offset) // Cần implement phương thức này trong userRepository
		if err != nil {
			return fmt.Errorf("failed to fetch data from database: %v", err)
		}

		// Nếu không còn dữ liệu, thoát khỏi vòng lặp
		if len(data) == 0 {
			break
		}

		// Ghi từng record vào CSV
		for _, user := range data {
			row := []string{
				strconv.Itoa(user.ID), // Thay bằng các trường dữ liệu thực tế của bạn
				user.Username,
				user.Email,
				user.Password,
				user.CreatedAt.Format("2006-01-02T15:04:05"),
				user.UpdatedAt.Format("2006-01-02T15:04:05"),
			}
			if err := writer.Write(row); err != nil {
				return fmt.Errorf("failed to write to csv: %v", err)
			}
		}

		// Tăng offset để lấy batch tiếp theo
		offset += batchSize
	}

	// Đảm bảo rằng tất cả các dữ liệu đã được ghi vào file
	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("error flushing csv writer: %v", err)
	}

	// Reset con trỏ file về đầu để upload
	if _, err := tempFile.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek file: %v", err)
	}

	// Upload file CSV lên S3
	if err := u.s3repository.UploadFile(bucket, filename, tempFile); err != nil {
		return fmt.Errorf("failed to upload file to s3: %v", err)
	}

	return nil
}

func (u UserService) RefreshToken(req *dto.RefreshTokenRequest) (string, error) {
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

func NewUserService(config *config.Config, userRepository repositories.IUserRepository, s3repository repositories.S3RepositoryInterface) IUserService {
	return &UserService{
		config:         config,
		userRepository: userRepository,
		s3repository:   s3repository,
	}
}
