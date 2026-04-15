package task

import (
	"context"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Repository interface {
	Create(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error)
	GetByID(ctx context.Context, id int64) (*taskdomain.Task, error)
	Update(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]taskdomain.Task, error)
}

type InstanceRepository interface {
	Create(ctx context.Context, instance *taskdomain.TaskInstance) (*taskdomain.TaskInstance, error)
	GetByID(ctx context.Context, id int64) (*taskdomain.TaskInstance, error)
	ListBetween(ctx context.Context, from, to time.Time) ([]taskdomain.TaskInstance, error)
	Exists(ctx context.Context, parentTaskID int64, scheduledDate time.Time) (bool, error)
	UpdateStatus(ctx context.Context, id int64, status taskdomain.Status) (*taskdomain.TaskInstance, error)
}

type Usecase interface {
	Create(ctx context.Context, input CreateInput) (*taskdomain.Task, error)
	GetByID(ctx context.Context, id int64) (*taskdomain.Task, error)
	Update(ctx context.Context, id int64, input UpdateInput) (*taskdomain.Task, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]taskdomain.Task, error)
	GenerateTasks(ctx context.Context, daysAhead int) (int, error)
	GetInstanceByID(ctx context.Context, id int64) (*taskdomain.TaskInstance, error)
	ListInstancesBetween(ctx context.Context, from, to time.Time) ([]taskdomain.TaskInstance, error)
	UpdateInstanceStatus(ctx context.Context, id int64, status taskdomain.Status) (*taskdomain.TaskInstance, error)
}

type CreateInput struct {
	Title            string
	Description      string
	Status           taskdomain.Status
	RecurrenceType   taskdomain.RecurrenceType
	RecurrenceConfig *taskdomain.RecurrenceConfig
	RecurrenceActive bool
}

type UpdateInput struct {
	Title            string
	Description      string
	Status           taskdomain.Status
	RecurrenceType   taskdomain.RecurrenceType
	RecurrenceConfig *taskdomain.RecurrenceConfig
	RecurrenceActive bool
}
