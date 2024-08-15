package handlers

import (
	"Backend_golang_project/internal/domain/dto/project"
	"Backend_golang_project/internal/use_cases"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"io"
	"net/http"
	"strconv"
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
	var request project.CreateProjectRequest

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
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	project, err := h.service.GetById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user", "details": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"project": gin.H{
			"ID":                project.ID,
			"Name":              project.Name,
			"Category":          project.Category,
			"ProjectVariance":   project.ProjectVariance,
			"RevenueRecognised": project.RevenueRecognised,
			"ProjectSpend":      project.ProjectSpend,
			"ProjectStartedAt":  project.ProjectStartedAt,
			"CreatedAt":         project.CreatedAt,
			"UpdatedAt":         project.UpdatedAt,
			"ProjectEndedAt":    project.ProjectEndedAt,
		},
	})
}

func (h *ProjectHandler) Delete(ctx *gin.Context) {
	name := ctx.Param("name")

	err := h.service.Delete(ctx, name)
	if err != nil {
		// Trả về lỗi nếu xóa thất bại
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Trả về phản hồi thành công nếu xóa thành công
	ctx.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully", "name": name})
}

func (h *ProjectHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// 2. Bind JSON request body vào struct UpdateProjectRequest
	var req project.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Gọi service để cập nhật project
	updatedProject, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to update project",
				"details": err,
			})
		}
		return
	}

	c.JSON(http.StatusOK, updatedProject)
}

func (h *ProjectHandler) GetProjects(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	// Gọi service để lấy danh sách project với phân trang
	pagination, err := h.service.GetProjectList(ctx, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve projects", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, pagination)
}
