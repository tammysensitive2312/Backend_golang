package pkg

import "Backend_golang_project/internal/domain/entities"

type Pagination struct {
	Projects     []entities.Project `json:"projects"`
	TotalPages   int                `json:"total_pages"`
	TotalRecords int                `json:"total_records"`
	CurrentPage  int                `json:"current_page"`
}
