package services

import (
	"context"

	"github.com/rafli-lutfi/kanban-app-mongodb/src/models"
	"github.com/rafli-lutfi/kanban-app-mongodb/src/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CategoryService interface {
	GetCategoryByID(ctx context.Context, categoryID primitive.ObjectID) (models.Category, error)
	GetCategories(ctx context.Context, userID primitive.ObjectID) ([]models.Category, error)
	StoreCategory(ctx context.Context, category models.Category) (interface{}, error)
	DeleteCategory(ctx context.Context, categoryID primitive.ObjectID) error
}
type categoryService struct {
	categoryRepository repository.CategoryRepository
	taskRepository     repository.TaskRepository
}

func NewCategoryService(categoryRepository repository.CategoryRepository, taskRepository repository.TaskRepository) *categoryService {
	return &categoryService{categoryRepository, taskRepository}
}

func (s *categoryService) GetCategoryByID(ctx context.Context, categoryID primitive.ObjectID) (models.Category, error) {
	category, err := s.categoryRepository.GetCategoryByID(ctx, categoryID)
	if err == mongo.ErrNoDocuments {
		return models.Category{}, models.ErrRecordNotFound
	}
	if err != nil {
		return models.Category{}, err
	}

	return category, nil
}

func (s *categoryService) GetCategories(ctx context.Context, userID primitive.ObjectID) ([]models.Category, error) {
	categories, err := s.categoryRepository.GetCategories(ctx, userID)
	if err != nil {
		return []models.Category{}, err
	}

	return categories, nil
}

func (s *categoryService) StoreCategory(ctx context.Context, category models.Category) (interface{}, error) {
	categoryID, err := s.categoryRepository.StoreCategory(ctx, category)
	if err != nil {
		return nil, err
	}

	return categoryID, nil
}

func (s *categoryService) UpdateCategory() {}

func (s *categoryService) DeleteCategory(ctx context.Context, categoryID primitive.ObjectID) error {
	_, err := s.categoryRepository.GetCategoryByID(ctx, categoryID)
	if err != nil {
		return err
	}
	// TODO: delete all task according to categoryID
	err = s.taskRepository.DeleteAllTaskByCategory(ctx, categoryID)
	if err != nil {
		return err
	}

	return s.categoryRepository.DeleteCategory(ctx, categoryID)
}
