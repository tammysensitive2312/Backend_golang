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

// Delete removes a project from the database by its name.
//
// The function accepts a context and a project name as parameters.
// It uses the provided context to ensure that the operation is performed within the specified deadline.
//
// The function first attempts to find a project with the given name in the database.
// If no project is found, it logs an error message and returns false.
// If an error occurs during the search process, it logs the error message and returns false.
//
// If a project is found, the function proceeds to delete the project from the database.
// If an error occurs during the deletion process, it logs an error message and returns false.
// If no project is deleted (i.e., RowsAffected is 0), it logs an error message and returns false.
//
// If the project is successfully deleted, it logs an informational message and returns true.
func (p ProjectRepository) Delete(ctx context.Context, name string) bool {
	var project entities.Project
	err := p.db.WithContext(ctx).Where("name = ?", name).First(&project)
	if err.Error != nil {
		if errors.Is(err.Error, gorm.ErrRecordNotFound) {
			log.Errorf("Project with name '%s' not found", name)
			return false
		}
		log.Errorf("Error when finding project with name '%s': %v", name, err.Error)
		return false
	}

	result := p.db.WithContext(ctx).Delete(&project)
	if result.Error != nil {
		log.Errorf("Failed to delete project with name '%s': %v", name, result.Error)
		return false
	}

	if result.RowsAffected == 0 {
		log.Errorf("No project deleted with name '%s'", name)
		return false
	}

	log.Infof("Successfully deleted project with name '%s'", name)
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
