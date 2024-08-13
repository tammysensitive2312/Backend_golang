package repositories

import (
	"Backend_golang_project/internal/domain/entities"
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IProjectRepository interface {
	Create(ctx context.Context, pj *entities.Project) (*entities.Project, error)
	Update(ctx context.Context, pj *entities.Project) (*entities.Project, error)
	Delete(ctx context.Context, name string) bool
	GetById(ctx context.Context, id int) (*entities.Project, error)
}

type ProjectRepository struct {
	base
}

// Create creates a new project in the database.
//
// The function accepts a context and a pointer to a Project entity.
// It uses the provided context to ensure that the operation is performed within the specified deadline.
// The Project entity contains the necessary information to create a new project.
//
// The function uses the GORM library to interact with the database.
// It creates a new record in the database using the provided Project entity.
// If an error occurs during the creation process, the function logs the error and returns nil, along with the error.
// Otherwise, it returns the created Project entity and nil.
func (p ProjectRepository) Create(ctx context.Context, pj *entities.Project) (*entities.Project, error) {
	if err := p.db.WithContext(ctx).Create(pj).Error; err != nil {
		log.Error("Cannot create project with err: ", err)
		return nil, err
	}
	return pj, nil
}

func (p ProjectRepository) Update(ctx context.Context, pj *entities.Project) (*entities.Project, error) {
	panic("err")
}

func (p ProjectRepository) Delete(ctx context.Context, name string) bool {
	var project entities.Project
	if err := p.db.WithContext(ctx).First(&project, name).Error; err != nil {
		log.Error("Having error when find project: ", name, err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Errorf("project %v not found", project)
			return false
		}
	} else {
		if err := p.db.WithContext(ctx).Delete(&project).Error; err != nil {
			log.Error("failed to delete project: ", name, err)
			return false
		}
	}
	return true
}

func (p ProjectRepository) GetById(ctx context.Context, id int) (*entities.Project, error) {
	var project entities.Project
	if err := p.db.WithContext(ctx).First(&project, id).Error; err != nil {
		log.Error("Having error when find project: ", id, err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("project %v not found", project)
		}
		return nil, fmt.Errorf("error retrieving project: %w", err)
	}
	return &project, nil
}

// NewProjectRepository constructor
func NewProjectRepository(db *gorm.DB) IProjectRepository {
	return &ProjectRepository{
		base: base{db: db}}
}
