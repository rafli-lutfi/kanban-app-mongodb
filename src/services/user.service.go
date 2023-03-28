package services

import (
	"context"

	"github.com/rafli-lutfi/kanban-app-mongodb/src/models"
	"github.com/rafli-lutfi/kanban-app-mongodb/src/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, user models.UserRegister) (interface{}, error)
	Login(ctx context.Context, user models.UserLogin) (models.User, error)
}

type userService struct {
	userRepository  repository.UserRepository
	categoryService repository.CategoryRepository
}

func NewUserService(userRepository repository.UserRepository, categoryService repository.CategoryRepository) *userService {
	return &userService{userRepository, categoryService}
}

func (s *userService) Register(ctx context.Context, user models.UserRegister) (interface{}, error) {
	hashPassword(&user.Password)

	_, err := s.userRepository.FindUserByEmail(ctx, user.Email)
	if err == nil {
		return nil, models.ErrEmailAlredyExist
	}

	userID, err := s.userRepository.Register(ctx, user)
	if err != nil {
		return nil, err
	}

	objectID := userID.(primitive.ObjectID)

	categories := []interface{}{
		models.Category{Id: primitive.NewObjectID(), Type: "Todo", UserId: objectID},
		models.Category{Id: primitive.NewObjectID(), Type: "In Progress", UserId: objectID},
		models.Category{Id: primitive.NewObjectID(), Type: "Done", UserId: objectID},
		models.Category{Id: primitive.NewObjectID(), Type: "Backlog", UserId: objectID},
	}

	_, err = s.categoryService.StoreManyCategory(ctx, categories)
	if err != nil {
		return nil, err
	}

	return userID, nil
}

func (s *userService) Login(ctx context.Context, user models.UserLogin) (models.User, error) {
	userDB, err := s.userRepository.FindUserByEmail(ctx, user.Email)
	if err == mongo.ErrNoDocuments {
		return models.User{}, models.ErrEmailPasswordNotMatched
	}

	if err != nil {
		return models.User{}, err
	}

	valid := validatePassword(userDB.Password, user.Password)
	if !valid {
		return models.User{}, models.ErrEmailPasswordNotMatched
	}

	return userDB, nil
}

func hashPassword(pass *string) {
	bytePass := []byte(*pass)
	hash, _ := bcrypt.GenerateFromPassword(bytePass, 14)

	*pass = string(hash)
}

func validatePassword(hash string, pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	return err == nil
}
