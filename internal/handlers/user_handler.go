package handlers

import (
	dto "Backend_golang_project/internal/domain/dto/user"
	"Backend_golang_project/internal/use_cases"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	userService use_cases.IUserService
}

func NewUserHandler(userService use_cases.IUserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateNewUser vẫn còn bug khi tạo user chứa json project thì chưa validate được các trường dữ liệu của
// các trường dữ liệu của project
func (h *UserHandler) CreateNewUser(ctx *gin.Context) {
	var request dto.CreateUserRequest

	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request payload: %v", err)})
		return
	}

	// 500
	newUser, err := h.userService.Create(ctx, &request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
		return
	}

	// trả vể mã 201
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user": gin.H{
			"id":         newUser.ID,
			"name":       newUser.Username,
			"email":      newUser.Email,
			"created_at": newUser.CreatedAt,
			"updated_at": newUser.UpdatedAt,
		},
	})
}
