package taskRepository

import (
	"github.com/jmoiron/sqlx"
	"github.com/sword-challenge/internal/model"
)

type TaskQueries interface {
	DeleteTaskFromStore(id int) (int, error)
	GetTaskFromStore(id int) (*model.Task, error)
	GetTasksFromStore(id int) ([]model.Task, error)
	GetAllTasksFromStore() ([]model.Task, error)
	AddTaskToStore(task *model.Task) (int, error)
	UpdateTaskInStore(task *model.Task) (*model.Task, error)
}

type Task struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) TaskQueries {
	return &Task{db}
}

func (t *Task) DeleteTaskFromStore(id int) (int, error) {
	res, err := t.db.Exec("DELETE FROM tasks t WHERE t.id = ?;", id)
	if err != nil {
		return 0, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affected), nil
}

func (s *Task) GetTaskFromStore(id int) (*model.Task, error) {
	task := &model.Task{}
	err := s.db.Get(task, "SELECT t.id, t.summary, t.completed_date, u.id as 'user.id', u.username as 'user.username' FROM tasks t INNER JOIN users u on t.user_id = u.id WHERE t.id = ?;", id)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (t *Task) GetTasksFromStore(id int) ([]model.Task, error) {
	task := []model.Task{}
	err := t.db.Select(&task, "SELECT t.id, t.summary, t.completed_date, u.id as 'user.id', u.username as 'user.username' FROM tasks t INNER JOIN users u on t.user_id = u.id WHERE t.user_id = ?;", id)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (t *Task) GetAllTasksFromStore() ([]model.Task, error) {
	task := []model.Task{}
	err := t.db.Select(&task, "SELECT t.id, t.summary, t.completed_date, u.id as 'user.id', u.username as 'user.username' FROM tasks t INNER JOIN users u on t.user_id = u.id;")
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (t *Task) AddTaskToStore(task *model.Task) (int, error) {
	result, err := t.db.Exec("INSERT INTO tasks (user_id, summary) VALUES (?, ?);", task.User.ID, task.Summary)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (t *Task) UpdateTaskInStore(task *model.Task) (*model.Task, error) {
	_, err := t.db.Exec(
		// Coalesce the fields so we only update the ones that were not sent as empty to the API
		"UPDATE tasks SET user_id = COALESCE(?, user_id), summary = COALESCE(?, summary), completed_date = ? WHERE id = ?;",
		task.User.ID, task.Summary, task.CompletedDate, task.ID)
	if err != nil {
		return nil, err
	}

	return t.GetTaskFromStore(task.ID)
}
