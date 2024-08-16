package handlers

import (
	"Backend_golang_project/internal/domain/dto/request"
	"Backend_golang_project/internal/pkg"
	"Backend_golang_project/internal/use_cases"
	"encoding/json"
	"errors"
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
	var projectRequest request.CreateProjectRequest

	if err := ctx.ShouldBindJSON(&projectRequest); err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			pkg.AbortErrorHandleCustomMessage(ctx, pkg.CannotBindJson, "Request body contains malformed JSON (syntax error)")
		case errors.As(err, &unmarshalTypeError):
			log.WithFields(log.Fields{
				"field": unmarshalTypeError.Field,
				"value": unmarshalTypeError.Value,
				"type":  unmarshalTypeError.Type.String(),
			}).Error("Failed to unmarshal JSON field")
			pkg.AbortErrorHandleCustomMessage(ctx, pkg.CannotBindJson, "Request body contains an invalid value for the "+unmarshalTypeError.Field+" field")
		case errors.Is(err, io.EOF):
			pkg.AbortErrorHandleCustomMessage(ctx, pkg.CannotBindJson, "Request body must not be empty")
		default:
			pkg.AbortErrorHandleCustomMessage(ctx, pkg.CannotBindJson, "Invalid request payload")
		}
		return
	}

	newProject, err := h.service.Create(ctx, projectRequest)
	if err != nil {
		pkg.AbortErrorHandleCustomMessage(ctx, http.StatusInternalServerError, "Failed to create project")
		return
	}

	pkg.SuccessfulHandle(ctx, newProject)
}

func (h *ProjectHandler) GetById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		pkg.AbortErrorHandleCustomMessage(ctx, pkg.CannotBindJson, "Invalid ID")
		return
	}

	project, err := h.service.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			pkg.AbortErrorHandler(ctx, pkg.RecordNotFound)
		} else {
			pkg.AbortErrorHandleCustomMessage(ctx, http.StatusInternalServerError, "Failed to retrieve project")
		}
		return
	}
	pkg.SuccessfulHandle(ctx, project)
}

func (h *ProjectHandler) Delete(ctx *gin.Context) {
	name := ctx.Param("name")

	err := h.service.Delete(ctx, name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			pkg.AbortErrorHandler(ctx, pkg.RecordNotFound)
		} else {
			pkg.AbortErrorHandleCustomMessage(ctx, http.StatusInternalServerError, "Failed to delete project")
		}
		return
	}

	pkg.SuccessfulHandle(ctx, gin.H{"message": "Project deleted successfully", "name": name})
}

func (h *ProjectHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pkg.AbortErrorHandleCustomMessage(c, pkg.CannotBindJson, "Invalid project ID")
		return
	}

	var req request.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.AbortErrorHandleCustomMessage(c, pkg.CannotBindJson, err.Error())
		return
	}

	updatedProject, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			pkg.AbortErrorHandler(c, pkg.RecordNotFound)
		default:
			pkg.AbortErrorHandleCustomMessage(c, http.StatusInternalServerError, "Failed to update project")
		}
		return
	}

	pkg.SuccessfulHandle(c, updatedProject)
}

func (h *ProjectHandler) GetProjects(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	pagination, err := h.service.GetProjectList(ctx, page, pageSize)
	if err != nil {
		pkg.AbortErrorHandleCustomMessage(ctx, http.StatusInternalServerError, "Failed to retrieve projects")
		return
	}

	pkg.SuccessfulHandle(ctx, pagination)
}
