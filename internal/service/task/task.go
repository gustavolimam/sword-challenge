package task

import (
	"database/sql"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/sword-challenge/internal/constants"
	"github.com/sword-challenge/internal/model"
	"github.com/sword-challenge/internal/repository/taskRepository"
	"github.com/sword-challenge/internal/repository/userRepository"
	"github.com/sword-challenge/internal/service/notification"
	"github.com/sword-challenge/pkg/rabbitmq"
)

type Service interface {
	GetTasks(c echo.Context) error
	CreateTask(c echo.Context) error
	UpdateTask(c echo.Context) error
	DeleteTask(c echo.Context) error
}

type service struct {
	userQueries      userRepository.UserQueries
	taskNotification notification.Service
	taskQueries      taskRepository.TaskQueries
}

func New(db *sqlx.DB, publisher rabbitmq.Publisher) Service {
	return &service{userRepository.New(db), notification.New(publisher), taskRepository.New(db)}
}

func (s *service) GetTasks(c echo.Context) error {
	var err error
	user := &model.User{}
	tasks := []model.Task{}

	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if user.Role.Name == constants.AdminRole {
		tasks, err = s.taskQueries.GetAllTasksFromStore()
	} else {
		tasks, err = s.taskQueries.GetTasksFromStore(user.ID)
	}
	if err != nil && err != sql.ErrNoRows {
		log.Warn().Msgf("Failed to get task from storage - error: %v", err)

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, tasks)
}

func (s *service) CreateTask(c echo.Context) error {
	task := &model.Task{}

	if err := c.Bind(task); err != nil {
		log.Warn().Msgf("Failed to parse task request body - error: %v", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	authenticatedUser := c.Get("user").(*model.User)

	if task.User.ID != authenticatedUser.ID && authenticatedUser.Role.Name != constants.AdminRole {
		return c.JSON(http.StatusForbidden, "userNotAuthenticatedOrNotManager")
	}

	id, err := s.taskQueries.AddTaskToStore(task)
	if err != nil {
		log.Warn().Msgf("Failed to add task to storage - error: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	task.ID = id

	return c.JSON(http.StatusCreated, task)
}

func (s *service) UpdateTask(c echo.Context) error {
	id := c.Get("task_id").(int)

	task := &model.Task{}
	if err := c.Bind(task); err != nil {
		log.Warn().Msgf("Failed to parse task from body while updating", "error", err)
		return c.JSON(http.StatusBadRequest, err)
	}
	task.ID = id

	authenticatedUser := c.Get("user").(*model.User)
	if task.User.ID != authenticatedUser.ID && authenticatedUser.Role.Name != constants.AdminRole {
		return c.JSON(http.StatusForbidden, "userNotAuthenticatedOrNotManager")
	}

	taskToUpdate, err := s.taskQueries.GetTaskFromStore(task.ID)
	if err == sql.ErrNoRows {
		log.Warn().Msgf("Failed to find task while updating - taskId %v", id)
		return c.JSON(http.StatusNotFound, err)
	} else if err != nil {
		log.Warn().Msgf("Failed to get task while updating taskId %v error %v", id, err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	// Check whether the task belongs to the user making the change or if the user is manager, we can only do this after we fetch the task from the database
	if authenticatedUser.Role.Name != constants.AdminRole && taskToUpdate.User.ID != authenticatedUser.ID {
		c.JSON(http.StatusForbidden, "userNotAuthenticatedOrNotManager")
	}

	updatedTask, err := s.taskQueries.UpdateTaskInStore(task)
	if err != nil {
		log.Warn().Msgf("Failed to update task in database - error %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	if taskToUpdate.CompletedDate == nil && updatedTask.CompletedDate != nil && authenticatedUser.Role.Name != constants.AdminRole {
		go func(t model.Task) {
			users, err := s.userQueries.GetUsersByRole(constants.AdminRole)
			if err != nil {
				log.Warn().Msgf("Failed to get users by role when sending notification", "error", err)
				return
			}

			for _, u := range users {
				u := u
				_ = s.taskNotification.TaskNotification(model.Notification{ID: t.ID, Manager: u.Username, CompletedDate: t.CompletedDate, User: t.User})
			}

		}(*updatedTask)
	}

	return c.JSON(http.StatusOK, task)
}

func (s *service) DeleteTask(c echo.Context) error {
	id := c.Get("task_id").(int)

	authenticatedUser := c.Get("user").(*model.User)
	if authenticatedUser.Role.Name != constants.AdminRole {
		return c.JSON(http.StatusForbidden, "userNotAuthenticatedOrNotManager")
	}

	rowsAffected, err := s.taskQueries.DeleteTaskFromStore(id)
	if err != nil {
		log.Warn().Msgf("Failed to delete task taskId %v error %v", id, err)

		return c.JSON(http.StatusInternalServerError, err)
	} else if rowsAffected == 0 {
		log.Warn().Msgf("Failed to find task while deleting taskId %v", id)

		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, "success")
}
