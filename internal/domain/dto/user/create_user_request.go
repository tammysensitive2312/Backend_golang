package dto

import (
	"Backend_golang_project/internal/domain/dto/project"
	"Backend_golang_project/internal/domain/entities"
	"github.com/go-playground/validator/v10"
	_ "time"
)

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=20,password_strength"`
	Username string `json:"username" binding:"required"`

	// tag dive dùng để đi sâu vào các cấu trúc dữ liệu lồng nhau
	Projects []project.CreateProjectRequest `json:"projects,omitempty" binding:"dive"`
}

func (req *CreateUserRequest) ToUserEntity(v *validator.Validate) (*entities.User, error) {
	for _, projectReq := range req.Projects {
		if err := v.Struct(projectReq); err != nil {
			return nil, err
		}
	}

	var projects []entities.Project
	for _, projectReq := range req.Projects {
		entity := projectReq.ToProjectEntity()
		projects = append(projects, *entity)
	}

	return &entities.User{
		Email:    req.Email,
		Password: req.Password,
		Username: req.Username,
		Projects: projects,
	}, nil
}
