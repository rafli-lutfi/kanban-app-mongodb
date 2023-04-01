package services

import (
	"context"

	"github.com/rafli-lutfi/kanban-app-mongodb/src/models"
	"github.com/rafli-lutfi/kanban-app-mongodb/src/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskService interface {
	GetTaskByID(ctx context.Context, taskID primitive.ObjectID) (models.Task, error)
	StoreTask(ctx context.Context, task models.Task) (interface{}, error)
	UpdateTask(ctx context.Context, task models.Task) error
	DeleteTask(ctx context.Context, taskID primitive.ObjectID) error
	UpdateTaskCategory(ctx context.Context, categoryID primitive.ObjectID, taskID primitive.ObjectID) error
}

type taskService struct {
	taskRepository     repository.TaskRepository
	categoryRepository repository.CategoryRepository
}

func NewTaskService(taskRepository repository.TaskRepository, categoryService repository.CategoryRepository) *taskService {
	return &taskService{taskRepository, categoryService}
}

func (s *taskService) GetTaskByID(ctx context.Context, taskID primitive.ObjectID) (models.Task, error) {
	task, err := s.taskRepository.GetTaskByID(ctx, taskID)
	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func (s *taskService) StoreTask(ctx context.Context, task models.Task) (interface{}, error) {
	_, err := s.categoryRepository.GetCategoryByID(ctx, task.CategoryId)
	if err == mongo.ErrNoDocuments {
		return nil, models.ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}

	taskID, err := s.taskRepository.StoreTask(ctx, task)
	if err != nil {
		return nil, err
	}

	return taskID, err
}

func (s *taskService) UpdateTask(ctx context.Context, task models.Task) error {
	_, err := s.taskRepository.GetTaskByID(ctx, task.Id)
	if err == mongo.ErrNoDocuments {
		return models.ErrRecordNotFound
	}
	if err != nil {
		return err
	}

	err = s.taskRepository.UpdateTask(ctx, task)
	if err != nil {
		return err
	}

	return nil
}

func (s *taskService) DeleteTask(ctx context.Context, taskID primitive.ObjectID) error {
	task, err := s.taskRepository.GetTaskByID(ctx, taskID)
	if err != nil {
		return err
	}

	err = s.categoryRepository.DeleteCategoryTask(ctx, task.CategoryId, taskID)
	if err != nil {
		return err
	}

	return s.taskRepository.DeleteTask(ctx, taskID)
}

func (s *taskService) UpdateTaskCategory(ctx context.Context, categoryID primitive.ObjectID, taskID primitive.ObjectID) error {
	taskDB, err := s.taskRepository.GetTaskByID(ctx, taskID)
	if err != nil {
		return err
	}

	updateTask := models.Task{
		Id:          taskID,
		CategoryId:  categoryID,
		Title:       taskDB.Title,
		Description: taskDB.Description,
		UserId:      taskDB.UserId,
	}

	err = s.taskRepository.UpdateTask(ctx, updateTask)
	if err != nil {
		return err
	}

	err = s.taskRepository.AppendTaskToCategoryTasks(ctx, updateTask)
	if err != nil {
		return err
	}

	err = s.categoryRepository.DeleteCategoryTask(ctx, taskDB.CategoryId, taskID)
	if err != nil {
		return err
	}

	return nil
}
