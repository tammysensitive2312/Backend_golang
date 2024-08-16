package use_cases

import (
	"Backend_golang_project/internal/domain/dto/request"
	"Backend_golang_project/internal/domain/entities"
	"Backend_golang_project/internal/pkg"
	"Backend_golang_project/internal/repositories"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type IProjectService interface {
	Create(ctx context.Context, request request.CreateProjectRequest) (*entities.Project, error)
	Delete(ctx context.Context, name string) error
	Update(ctx context.Context, id int, request request.UpdateProjectRequest) (*entities.Project, error)
	GetById(ctx context.Context, id int) (*entities.Project, error)
	GetProjectList(ctx context.Context, page int, pageSize int) (*pkg.Pagination, error)
}

type ProjectService struct {
	projectRepository repositories.IProjectRepository
}

func (p ProjectService) GetProjectList(ctx context.Context, page int, pageSize int) (*pkg.Pagination, error) {
	pagination, err := p.projectRepository.GetList(ctx, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get project list: %w", err)
	}
	return pagination, nil
}

func (p ProjectService) Create(ctx context.Context, request request.CreateProjectRequest) (*entities.Project, error) {
	entity := request.ToProjectEntity()

	if entity.ProjectEndedAt != nil && entity.ProjectEndedAt.Before(entity.ProjectStartedAt) {
		return nil, fmt.Errorf("project ended time is invalid")
	}
	data, err := p.projectRepository.Create(ctx, entity)
	if err != nil {
		log.Error("Error in service.Create with error: ", err)
		return nil, err
	}

	log.Info("Project created successfully", entity)
	return data, nil
}

func (p ProjectService) Delete(ctx context.Context, name string) error {
	// Thực hiện xóa project qua repository
	deleted := p.projectRepository.Delete(ctx, name)

	if deleted {
		// Nếu xóa thành công, ghi log thành công và trả về nil (không có lỗi)
		log.WithFields(log.Fields{
			"name": name,
		}).Info("Project deleted successfully")
		return nil
	} else {
		err := fmt.Errorf("failed to delete project with name: %s", name)
		log.WithFields(log.Fields{
			"name":  name,
			"error": err,
		}).Error("Delete project failed")
		return err
	}
}

func (p ProjectService) Update(ctx context.Context, id int, req request.UpdateProjectRequest) (*entities.Project, error) {
	// 1. Lấy project hiện tại từ repository
	existingProject, err := p.projectRepository.GetById(ctx, id)
	if err != nil {
		log.Error("Error fetching existing project: ", err)
		return nil, fmt.Errorf("error fetching existing project: %w", err)
	}

	// 2. Cập nhật thông tin project
	existingProject.Name = req.Name
	existingProject.Category = req.Category
	existingProject.ProjectSpend = req.ProjectSpend
	existingProject.ProjectVariance = req.ProjectVariance
	existingProject.RevenueRecognised = req.RevenueRecognised
	existingProject.ProjectStartedAt = req.ProjectStartedAt
	existingProject.ProjectEndedAt = req.ProjectEndedAt

	updatedProject, err := p.projectRepository.Update(ctx, existingProject)
	if err != nil {
		log.Error("Error updating project: ", err)
		return nil, fmt.Errorf("error updating project: %w", err)
	}

	log.Info("Project updated successfully", updatedProject)
	return updatedProject, nil
}

func (p ProjectService) GetById(ctx context.Context, id int) (*entities.Project, error) {
	getById, err := p.projectRepository.GetById(ctx, id)
	if err != nil {
		log.Error("Failed to get getByIdproject by ID:", err)
		return nil, err
	}
	return getById, nil
}

func NewProjectService(repository repositories.IProjectRepository) IProjectService {
	return &ProjectService{
		projectRepository: repository,
	}
}
