package handlers

import (
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type taskMutationDTO struct {
	Title            string                      `json:"title"`
	Description      string                      `json:"description"`
	Status           taskdomain.Status           `json:"status"`
	RecurrenceType   taskdomain.RecurrenceType   `json:"recurrence_type,omitempty"`
	RecurrenceConfig *taskdomain.RecurrenceConfig `json:"recurrence_config,omitempty"`
	RecurrenceActive bool                        `json:"recurrence_active"`
}

type taskDTO struct {
	ID               int64                       `json:"id"`
	Title            string                      `json:"title"`
	Description      string                      `json:"description"`
	Status           taskdomain.Status           `json:"status"`
	RecurrenceType   taskdomain.RecurrenceType   `json:"recurrence_type,omitempty"`
	RecurrenceConfig *taskdomain.RecurrenceConfig `json:"recurrence_config,omitempty"`
	RecurrenceActive bool                        `json:"recurrence_active"`
	IsRecurring      bool                        `json:"is_recurring"`
	CreatedAt        time.Time                   `json:"created_at"`
	UpdatedAt        time.Time                   `json:"updated_at"`
}

func newTaskDTO(task *taskdomain.Task) taskDTO {
	return taskDTO{
		ID:               task.ID,
		Title:            task.Title,
		Description:      task.Description,
		Status:           task.Status,
		RecurrenceType:   task.RecurrenceType,
		RecurrenceConfig: task.RecurrenceConfig,
		RecurrenceActive: task.RecurrenceActive,
		IsRecurring:      task.IsRecurring(),
		CreatedAt:        task.CreatedAt,
		UpdatedAt:        task.UpdatedAt,
	}
}

type taskInstanceDTO struct {
	ID            int64             `json:"id"`
	ParentTaskID  int64             `json:"parent_task_id"`
	ScheduledDate time.Time         `json:"scheduled_date"`
	Title         string            `json:"title"`
	Description   string            `json:"description"`
	Status        taskdomain.Status `json:"status"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

func newInstanceDTO(instance *taskdomain.TaskInstance) taskInstanceDTO {
	return taskInstanceDTO{
		ID:            instance.ID,
		ParentTaskID:  instance.ParentTaskID,
		ScheduledDate: instance.ScheduledDate,
		Title:         instance.Title,
		Description:   instance.Description,
		Status:        instance.Status,
		CreatedAt:     instance.CreatedAt,
		UpdatedAt:     instance.UpdatedAt,
	}
}
