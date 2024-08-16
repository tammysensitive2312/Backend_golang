package handlers

import (
	dto "Backend_golang_project/internal/domain/dto/request"
	"Backend_golang_project/internal/domain/dto/response"
	"Backend_golang_project/internal/pkg"
	"Backend_golang_project/internal/use_cases"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

type UserHandler struct {
	userService use_cases.IUserService
}

func NewUserHandler(userService use_cases.IUserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetUserById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		pkg.AbortErrorHandler(ctx, id)
		return
	}

	user, err := h.userService.GetUserByID(ctx, id)
	if err != nil {
		pkg.GetErrorResponse(500)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"response": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"projects": user.Projects,
		},
	})
}

// CreateNewUser vẫn còn bug khi tạo response chứa json project thì chưa validate được các trường dữ liệu của
// các trường dữ liệu của project
func (h *UserHandler) CreateNewUser(ctx *gin.Context) {
	var request dto.CreateUserRequest

	if err := ctx.BindJSON(&request); err != nil {
		pkg.AbortErrorHandleCustomMessage(ctx, pkg.CannotBindJson, err.Error())
		return
	}

	// 500
	newUser, err := h.userService.Create(ctx, &request)
	if err != nil {
		pkg.AbortErrorHandleCustomMessage(ctx, pkg.CannotCreateNewUser, err.Error())
		return
	}

	// trả vể mã 201
	pkg.SuccessfulHandle(ctx, newUser)
}

func (h *UserHandler) LoginUser(ctx *gin.Context) {
	var request dto.LoginRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Request body contains malformed JSON (syntax error)"})

		case errors.As(err, &unmarshalTypeError):
			log.WithFields(log.Fields{
				"field": unmarshalTypeError.Field,
				"value": unmarshalTypeError.Value,
				"type":  unmarshalTypeError.Type.String(),
			}).Error("Failed to unmarshal JSON field")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Request body contains an invalid value for the " + unmarshalTypeError.Field + " field"})

		case errors.Is(err, io.EOF):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Request body must not be empty"})

		default:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request payload: %v", err)})
		}
		return
	}

	accessToken, refreshToken, err := h.userService.Login(ctx, request)

	if err != nil || len(accessToken) == 0 || len(refreshToken) == 0 {
		pkg.AbortErrorHandleCustomMessage(ctx, pkg.InvalidLogin, err.Error())
		return
	}
	pkg.SuccessfulHandle(ctx, response.ToLoginResponse(accessToken, refreshToken))
}

func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.AbortErrorHandleCustomMessage(c, pkg.CannotBindJson, err.Error())
		return
	}

	newAccessToken, err := h.userService.RefreshToken(c, &req)
	if err != nil {
		pkg.AbortErrorHandleCustomMessage(c, 500, err.Error())
		return
	}

	pkg.SuccessfulHandle(c, newAccessToken)
}
