package handlers

import (
	"Backend_golang_project/internal/domain/dto"
	"Backend_golang_project/internal/use_cases"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type ProjectHandler struct {
	service use_cases.IProjectService
}

func NewProjectHandler(service use_cases.IProjectService) *ProjectHandler {
	return &ProjectHandler{
		service: service,
	}
}

func (h *ProjectHandler) Create(ctx *gin.Context) {
	var request dto.CreateProjectRequest

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

	newProject, err := h.service.Create(ctx, request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":                 newProject.ID,
		"name":               newProject.Name,
		"category":           newProject.Category,
		"project_spend":      newProject.ProjectSpend,
		"project_variance":   newProject.ProjectVariance,
		"revenue_recognised": newProject.RevenueRecognised,
		"project_started_at": newProject.ProjectStartedAt,
		"created_at":         newProject.CreatedAt,
		"updated_at":         newProject.UpdatedAt,
	})
}

func (h *ProjectHandler) GetById(ctx *gin.Context) {

}

func (h *ProjectHandler) Delete(ctx *gin.Context) {

}

func (h *ProjectHandler) Update(ctx *gin.Context) {

}
