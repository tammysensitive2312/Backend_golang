package dto

import (
	"Backend_golang_project/internal/domain/entities"
	"database/sql"
	"time"
)

type CreateProjectRequest struct {
	Name             string       `json:"name" binding:"required,lte=255"`
	Category         string       `json:"category"`
	ProjectSpend     int          `json:"project_spend"`
	ProjectVariance  int          `json:"project_variance"`
	ProjectStartedAt sql.NullTime `json:"project_started_at"`
	ProjectEndedAt   sql.NullTime `json:"project_ended_at"`
}

func (req *CreateProjectRequest) ToProjectEntity() *entities.Project {
	project := &entities.Project{
		Name:              req.Name,
		Category:          req.Category,
		ProjectSpend:      req.ProjectSpend,
		ProjectVariance:   req.ProjectVariance,
		RevenueRecognised: 0,
	}

	if req.ProjectStartedAt.Valid {
		project.ProjectStartedAt = req.ProjectStartedAt.Time
	} else {
		project.ProjectStartedAt = time.Now()
	}

	if req.ProjectEndedAt.Valid {
		project.ProjectEndedAt = &req.ProjectEndedAt.Time
	}

	return project
}
