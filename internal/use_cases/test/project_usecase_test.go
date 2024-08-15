package test

import (
	"Backend_golang_project/internal/domain/dto/project"
	"Backend_golang_project/internal/domain/entities"
	"Backend_golang_project/internal/pkg"
	"Backend_golang_project/internal/use_cases"

	//"Backend_golang_project/internal/repositories"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProjectRepository là một mock của IProjectRepository
type MockProjectRepository struct {
	mock.Mock
}

func (m *MockProjectRepository) GetList(ctx context.Context, page int, pageSize int) (*pkg.Pagination, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockProjectRepository) Create(ctx context.Context, pj *entities.Project) (*entities.Project, error) {
	args := m.Called(ctx, pj)
	return args.Get(0).(*entities.Project), args.Error(1)
}

func (m *MockProjectRepository) Update(ctx context.Context, pj *entities.Project) (*entities.Project, error) {
	args := m.Called(ctx, pj)
	return args.Get(0).(*entities.Project), args.Error(1)
}

func (m *MockProjectRepository) Delete(ctx context.Context, name string) bool {
	args := m.Called(ctx, name)
	return args.Bool(0)
}

func (m *MockProjectRepository) GetById(ctx context.Context, id int) (*entities.Project, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.Project), args.Error(1)
}

func TestProjectService_Create(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	service := use_cases.NewProjectService(mockRepo)

	ctx := context.Background()
	request := project.CreateProjectRequest{
		Name: "Test Project",
	}

	expectedProject := &entities.Project{
		ID:   1,
		Name: request.Name,
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*entities.Project")).Return(expectedProject, nil)

	result, err := service.Create(ctx, request)

	assert.NoError(t, err)
	assert.Equal(t, expectedProject, result)
	mockRepo.AssertExpectations(t)
}

func TestProjectService_Delete(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	service := use_cases.NewProjectService(mockRepo)

	ctx := context.Background()
	projectName := "Test Project"

	mockRepo.On("Delete", ctx, projectName).Return(true)

	service.Delete(ctx, projectName)

	mockRepo.AssertExpectations(t)
}

func TestProjectService_GetById(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	service := use_cases.NewProjectService(mockRepo)

	ctx := context.Background()
	projectID := 1

	expectedProject := &entities.Project{
		ID:   projectID,
		Name: "Test Project",
	}

	mockRepo.On("GetById", ctx, projectID).Return(expectedProject, nil)

	result, err := service.GetById(ctx, projectID)

	assert.NoError(t, err)
	assert.Equal(t, expectedProject, result)
	mockRepo.AssertExpectations(t)
}

func TestProjectService_GetById_Error(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	service := use_cases.NewProjectService(mockRepo)

	ctx := context.Background()
	projectID := 1

	expectedError := errors.New("project not found")

	mockRepo.On("GetById", ctx, projectID).Return((*entities.Project)(nil), expectedError)

	result, err := service.GetById(ctx, projectID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}
