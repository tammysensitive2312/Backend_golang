package dto

import "database/sql"

type CreateProjectRequest struct {
	Name             string       `json:"name" binding:"required,lte=255"`
	Category         string       `json:"category"`
	ProjectSpend     int          `json:"project_spend"`
	ProjectVariance  int          `json:"project_variance"`
	ProjectStartedAt sql.NullTime `json:"project_started_at"`
	ProjectEndedAt   sql.NullTime `json:"project_ended_at"`
}
