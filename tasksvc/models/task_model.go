package models

import (
	"strings"
	"time"
)

// Task Priority => 1 - low, 2 - medium, 3 - high, 4 - highest

type Task struct {
	Title       string     `json:"title,omitempty" db:"title"`
	Description string     `json:"description,omitempty" db:"description"`
	Priority    int64      `json:"priority,omitempty" db:"priority"`
	DueDate     CustomTime `json:"dueDate,omitempty" db:"due_date"`
}

type TaskResponse struct {
	ID          *int64     `json:"id,omitempty" db:"id"`
	Title       *string    `json:"title,omitempty" db:"title"`
	Description *string    `json:"description,omitempty" db:"description"`
	Priority    *int64     `json:"priority,omitempty" db:"priority"`
	DueDate     *time.Time `json:"dueDate,omitempty" db:"due_date"`
	TaskID      *int64     `json:"taskId,omitempty" db:"task_id"`
	UserID      *int64     `json:"userId,omitempty" db:"user_id"`
	IsCompleted *bool      `json:"isCompleted,omitempty" db:"is_completed"`
	CreatedAt   *time.Time `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty" db:"updated_at"`
}

type TaskDetails struct {
	TaskID      int64 `json:"taskId,omitempty" db:"task_id"`
	UserID      int64 `json:"userId,omitempty" db:"user_id"`
	IsCompleted bool  `json:"isCompleted,omitempty" db:"is_completed"`
}

type CustomTime struct {
	time.Time
}

type AdminResponse struct {
	IsAdmin bool `json:"isAdmin"`
}

func (t *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")

	date, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return err
	}
	t.Time = date
	return
}
