package model

import "time"

type Task struct {
	ID            int        `json:"id,omitempty"`
	Summary       string     `json:"summary,omitempty" db:"summary" validate:"max=2500"`
	CompletedDate *time.Time `json:"completedDate" db:"completed_date"`
	User          *User      `json:"user,omitempty"`
}

type Notification struct {
	ID            int        `json:"id"`
	Manager       string     `json:"manager"`
	CompletedDate *time.Time `json:"completedDate"`
	User          *User      `json:"user"`
}
