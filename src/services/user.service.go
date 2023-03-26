package services

import (
	"context"

	"github.com/rafli-lutfi/kanban-app-mongodb/src/models"
	"github.com/rafli-lutfi/kanban-app-mongodb/src/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, user models.UserRegister) (interface{}, error)
	Login(ctx context.Context, user models.UserLogin) (models.User, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) *userService {
	return &userService{userRepository}
}

func (s *userService) Register(ctx context.Context, user models.UserRegister) (interface{}, error) {
	hashPassword(&user.Password)

	userID, err := s.userRepository.Register(ctx, user)
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
