package use_cases

import (
	"Backend_golang_project/internal/domain/dto"
	"Backend_golang_project/internal/domain/entities"
	"Backend_golang_project/internal/repositories"
	"context"
	log "github.com/sirupsen/logrus"
)

type IProjectService interface {
	Create(ctx context.Context, request dto.CreateProjectRequest) (*entities.Project, error)
	Delete(ctx context.Context, name string)
	Update()
	GetById(ctx context.Context, id int) (*entities.Project, error)
}

type ProjectService struct {
	projectRepository repositories.IProjectRepository
}

func (p ProjectService) Create(ctx context.Context, request dto.CreateProjectRequest) (*entities.Project, error) {
	project := request.ToProjectEntity()
	data, err := p.projectRepository.Create(ctx, project)
	if err != nil {
		log.Error("Error in service.Create with error: ", err)
		return nil, err
	}

	log.Info("Project created successfully", project)
	return data, nil
}

func (p ProjectService) Delete(ctx context.Context, name string) {
	if p.projectRepository.Delete(ctx, name) {
		log.Info("Project deleted successfully", name)
	}
	log.Error("Delete project failed")
	return
}

func (p ProjectService) Update() {
	//TODO implement me
	panic("implement me")
}

func (p ProjectService) GetById(ctx context.Context, id int) (*entities.Project, error) {
	project, err := p.projectRepository.GetById(ctx, id)
	if err != nil {
		log.Error("Failed to get project by ID:", err)
		return nil, err
	}
	return project, nil
}

func NewProjectService(repository repositories.IProjectRepository) IProjectService {
	return &ProjectService{
		projectRepository: repository,
	}
}
