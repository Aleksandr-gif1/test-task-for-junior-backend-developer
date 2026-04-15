package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type InstanceRepository struct {
	pool *pgxpool.Pool
}

func NewInstanceRepository(pool *pgxpool.Pool) *InstanceRepository {
	return &InstanceRepository{pool: pool}
}

func (r *InstanceRepository) Create(ctx context.Context, instance *taskdomain.TaskInstance) (*taskdomain.TaskInstance, error) {
	const query = `
		INSERT INTO task_instances (parent_task_id, scheduled_date, title, description, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, parent_task_id, scheduled_date, title, description, status, created_at, updated_at
	`

	row := r.pool.QueryRow(ctx, query,
		instance.ParentTaskID,
		instance.ScheduledDate,
		instance.Title,
		instance.Description,
		instance.Status,
		instance.CreatedAt,
		instance.UpdatedAt,
	)

	if err := scanInstance(row, instance); err != nil {
		return nil, err
	}

	return instance, nil
}

func (r *InstanceRepository) GetByID(ctx context.Context, id int64) (*taskdomain.TaskInstance, error) {
	const query = `
		SELECT id, parent_task_id, scheduled_date, title, description, status, created_at, updated_at
		FROM task_instances
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
	return scanInstanceRow(row)
}

func (r *InstanceRepository) ListBetween(ctx context.Context, from, to time.Time) ([]taskdomain.TaskInstance, error) {
	const query = `
		SELECT id, parent_task_id, scheduled_date, title, description, status, created_at, updated_at
		FROM task_instances
		WHERE scheduled_date >= $1 AND scheduled_date <= $2
		ORDER BY scheduled_date ASC, id ASC
	`

	rows, err := r.pool.Query(ctx, query, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	instances := make([]taskdomain.TaskInstance, 0)
	for rows.Next() {
		instance, err := scanInstanceRow(rows)
		if err != nil {
			return nil, err
		}

		instances = append(instances, *instance)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return instances, nil
}

func (r *InstanceRepository) Exists(ctx context.Context, parentTaskID int64, scheduledDate time.Time) (bool, error) {
	const query = `
		SELECT EXISTS(SELECT 1 FROM task_instances WHERE parent_task_id = $1 AND scheduled_date = $2)
	`

	var exists bool
	if err := r.pool.QueryRow(ctx, query, parentTaskID, scheduledDate).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (r *InstanceRepository) UpdateStatus(ctx context.Context, id int64, status taskdomain.Status) (*taskdomain.TaskInstance, error) {
	const query = `
		UPDATE task_instances
		SET status = $1, updated_at = $2
		WHERE id = $3
		RETURNING id, parent_task_id, scheduled_date, title, description, status, created_at, updated_at
	`

	row := r.pool.QueryRow(ctx, query, status, time.Now().UTC(), id)
	return scanInstanceRow(row)
}

type instanceScanner interface {
	Scan(dest ...any) error
}

func scanInstance(scanner instanceScanner, instance *taskdomain.TaskInstance) error {
	var status string

	if err := scanner.Scan(
		&instance.ID,
		&instance.ParentTaskID,
		&instance.ScheduledDate,
		&instance.Title,
		&instance.Description,
		&status,
		&instance.CreatedAt,
		&instance.UpdatedAt,
	); err != nil {
		return err
	}

	instance.Status = taskdomain.Status(status)

	return nil
}

func scanInstanceRow(scanner instanceScanner) (*taskdomain.TaskInstance, error) {
	instance := &taskdomain.TaskInstance{}
	if err := scanInstance(scanner, instance); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, taskdomain.ErrNotFound
		}
		return nil, err
	}

	return instance, nil
}
