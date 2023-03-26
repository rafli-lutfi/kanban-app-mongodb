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
	StoreCategory(ctx context.Context, category models.Category) (interface{}, error)
}
type categoryService struct {
	categoryRepository repository.CategoryRepository
}

func NewCategoryService(categoryRepository repository.CategoryRepository) *categoryService {
	return &categoryService{categoryRepository}
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

func (s *categoryService) GetCategories() {}

func (s *categoryService) StoreCategory(ctx context.Context, category models.Category) (interface{}, error) {
	categoryID, err := s.categoryRepository.StoreCategory(ctx, category)
	if err != nil {
		return nil, err
	}

	return categoryID, nil
}

func (s *categoryService) UpdateCategory() {}
func (s *categoryService) DeleteCategory() {}
