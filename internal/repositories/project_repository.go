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
	GetList(ctx context.Context, page int, pageSize int) (*Pagination, error)
}

type ProjectRepository struct {
	base
}

func (p ProjectRepository) GetList(ctx context.Context, page int, pageSize int) (*Pagination, error) {
	var projectList []entities.Project
	var totalRecords int64

	// đếm tổng số bản ghi
	if err := p.db.WithContext(ctx).Model(&entities.Project{}).Count(&totalRecords).Error; err != nil {
		return nil, err
	}
	// tính số trang
	totalPages := int(totalRecords) / pageSize
	if int(totalRecords)%pageSize != 0 {
		totalPages++
	}

	// lấy dữ liệu cho trang hiện tại
	offset := (page - 1) * pageSize
	if err := p.db.WithContext(ctx).Limit(pageSize).Offset(offset).Find(&projectList).Error; err != nil {
		return nil, err
	}

	pagination := &Pagination{
		Projects:     projectList,
		TotalPages:   totalPages,
		TotalRecords: int(totalRecords),
		CurrentPage:  page,
	}

	return pagination, nil
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

// Update The function performs the following steps:
// 1. Checks if a project with the given ID exists in the database.
//   - If no project is found, it returns nil and an error indicating that the project was not found.
//   - If an error occurs during the check, it returns nil and the error.
//
// 2. Updates the project in the database using the provided Project entity.
//   - If an error occurs during the update, it returns nil and the error.
//   - If no rows are affected (i.e., the project was not updated), it returns nil and an error indicating that the update might have failed.
//
// 3. Fetches the updated project from the database.
//   - If an error occurs during the fetch, it returns nil and the error.
//
// 4. Returns the updated Project entity and nil.
func (p ProjectRepository) Update(ctx context.Context, pj *entities.Project) (*entities.Project, error) {
	// existing check
	existingProject := entities.Project{}
	result := p.db.WithContext(ctx).First(&existingProject, pj.ID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("project not found")
		}
		return nil, fmt.Errorf("error checking existing project: %w", result.Error)
	}

	result = p.db.WithContext(ctx).Model(pj).Updates(pj)
	if result.Error != nil {
		return nil, fmt.Errorf("error updating project: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("no rows affected, update might have failed")
	}

	updatedProject := &entities.Project{}
	if err := p.db.WithContext(ctx).First(updatedProject, pj.ID).Error; err != nil {
		return nil, fmt.Errorf("error fetching updated project: %w", err)
	}

	return updatedProject, nil
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
