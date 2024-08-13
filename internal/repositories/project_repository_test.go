package repositories_test

import (
	"Backend_golang_project/internal/domain/entities"
	"Backend_golang_project/internal/repositories"
	"context"
	"github.com/glebarez/sqlite"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type ProjectRepositoryTestSuite struct {
	suite.Suite
	mockDB *gorm.DB
	repo   repositories.IProjectRepository
}

func (suite *ProjectRepositoryTestSuite) SetupTest() {
	// Tạo một mock DB sử dụng SQLite in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)

	// Migrate schema
	err = db.AutoMigrate(&entities.Project{})
	suite.Require().NoError(err)

	suite.mockDB = db
	suite.repo = repositories.NewProjectRepository(db)
}

func (suite *ProjectRepositoryTestSuite) TestCreate() {
	ctx := context.Background()
	project := &entities.Project{
		Name:      "Test Project",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdProject, err := suite.repo.Create(ctx, project)

	suite.Require().NoError(err)
	suite.Require().NotNil(createdProject)
	suite.Equal(project.Name, createdProject.Name)
}

func (suite *ProjectRepositoryTestSuite) TestGetById() {
	ctx := context.Background()
	project := &entities.Project{
		Name: "Test Project",
	}

	// Tạo project trước
	createdProject, err := suite.repo.Create(ctx, project)
	suite.Require().NoError(err)

	// Lấy project bằng ID
	foundProject, err := suite.repo.GetById(ctx, createdProject.ID)

	suite.Require().NoError(err)
	suite.Require().NotNil(foundProject)
	suite.Equal(createdProject.ID, foundProject.ID)
	suite.Equal(createdProject.Name, foundProject.Name)
}

func (suite *ProjectRepositoryTestSuite) TestDelete() {
	ctx := context.Background()
	project := &entities.Project{
		Name: "Test Project",
	}

	// Tạo project trước
	createdProject, err := suite.repo.Create(ctx, project)
	suite.Require().NoError(err)

	// Xóa project
	result := suite.repo.Delete(ctx, createdProject.Name)
	suite.True(result)
}

func TestProjectRepositorySuite(t *testing.T) {
	suite.Run(t, new(ProjectRepositoryTestSuite))
}
