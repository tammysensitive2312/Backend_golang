package project

import (
	"Backend_golang_project/internal/domain/entities"
	"time"
)

type CreateProjectRequest struct {
	Name             string     `json:"name" binding:"required,lte=255"`
	Category         string     `json:"category" binding:"required,valid_category"`
	ProjectSpend     int        `json:"project_spend"`
	ProjectVariance  int        `json:"project_variance"`
	ProjectStartedAt time.Time  `json:"project_started_at" binding:"required,future_date"`
	ProjectEndedAt   *time.Time `json:"project_ended_at"`
}

func (req *CreateProjectRequest) ToProjectEntity() *entities.Project {
	project := &entities.Project{
		Name:              req.Name,
		Category:          req.Category,
		ProjectSpend:      req.ProjectSpend,
		ProjectVariance:   req.ProjectVariance,
		ProjectStartedAt:  req.ProjectStartedAt,
		ProjectEndedAt:    req.ProjectEndedAt,
		RevenueRecognised: 0,
	}
	return project
}
