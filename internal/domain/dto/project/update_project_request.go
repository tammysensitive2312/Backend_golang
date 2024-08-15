package project

import "time"

type UpdateProjectRequest struct {
	Name              string     `json:"name"`
	Category          string     `json:"category"`
	ProjectSpend      int        `json:"project_spend"`
	ProjectVariance   int        `json:"project_variance"`
	RevenueRecognised int        `json:"revenue_recognised"`
	ProjectStartedAt  time.Time  `json:"project_started_at"`
	ProjectEndedAt    *time.Time `json:"project_ended_at,omitempty"`
}
