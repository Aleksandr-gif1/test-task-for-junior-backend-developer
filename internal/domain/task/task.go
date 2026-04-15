package task

import "time"

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type Task struct {
	ID               int64             `json:"id"`
	Title            string            `json:"title"`
	Description      string            `json:"description"`
	Status           Status            `json:"status"`
	RecurrenceType   RecurrenceType    `json:"recurrence_type,omitempty"`
	RecurrenceConfig *RecurrenceConfig `json:"recurrence_config,omitempty"`
	RecurrenceActive bool              `json:"recurrence_active"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
}

// IsRecurring возвращает true, если задача является периодическим шаблоном.
func (t *Task) IsRecurring() bool {
	return t.RecurrenceType != "" && t.RecurrenceActive
}

func (s Status) Valid() bool {
	switch s {
	case StatusNew, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}
