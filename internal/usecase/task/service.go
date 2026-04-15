package task

import (
	"context"
	"fmt"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Service struct {
	repo             Repository
	instanceRepo     InstanceRepository
	now              func() time.Time
}

func NewService(repo Repository, instanceRepo InstanceRepository) *Service {
	return &Service{
		repo:         repo,
		instanceRepo: instanceRepo,
		now:          func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*taskdomain.Task, error) {
	normalized, err := validateCreateInput(input)
	if err != nil {
		return nil, err
	}

	model := &taskdomain.Task{
		Title:            normalized.Title,
		Description:      normalized.Description,
		Status:           normalized.Status,
		RecurrenceType:   normalized.RecurrenceType,
		RecurrenceConfig: normalized.RecurrenceConfig,
		RecurrenceActive: normalized.RecurrenceActive,
	}
	now := s.now()
	model.CreatedAt = now
	model.UpdatedAt = now

	created, err := s.repo.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	normalized, err := validateUpdateInput(input)
	if err != nil {
		return nil, err
	}

	model := &taskdomain.Task{
		ID:               id,
		Title:            normalized.Title,
		Description:      normalized.Description,
		Status:           normalized.Status,
		RecurrenceType:   normalized.RecurrenceType,
		RecurrenceConfig: normalized.RecurrenceConfig,
		RecurrenceActive: normalized.RecurrenceActive,
		UpdatedAt:        s.now(),
	}

	updated, err := s.repo.Update(ctx, model)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]taskdomain.Task, error) {
	return s.repo.List(ctx)
}

func (s *Service) GetInstanceByID(ctx context.Context, id int64) (*taskdomain.TaskInstance, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.instanceRepo.GetByID(ctx, id)
}

func (s *Service) ListInstancesBetween(ctx context.Context, from, to time.Time) ([]taskdomain.TaskInstance, error) {
	return s.instanceRepo.ListBetween(ctx, from, to)
}

func (s *Service) UpdateInstanceStatus(ctx context.Context, id int64, status taskdomain.Status) (*taskdomain.TaskInstance, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	if !status.Valid() {
		return nil, fmt.Errorf("%w: invalid status %q", ErrInvalidInput, status)
	}

	return s.instanceRepo.UpdateStatus(ctx, id, status)
}
