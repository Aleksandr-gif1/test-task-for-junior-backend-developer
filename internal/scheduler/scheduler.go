package scheduler

import (
	"context"
	"log/slog"
	"os"
	"time"

	taskusecase "example.com/taskservice/internal/usecase/task"
)

// Scheduler — фоновый планировщик, который периодически запускает генерацию
// экземпляров периодических задач.
type Scheduler struct {
	usecase    taskusecase.Usecase
	logger     *slog.Logger
	interval   time.Duration
	daysAhead  int
	enabled    bool
}

// New создаёт новый Scheduler.
// interval — как часто запускать генерацию (например, 24h).
// daysAhead — на сколько дней вперёд генерировать задачи.
func New(usecase taskusecase.Usecase, logger *slog.Logger, interval time.Duration, daysAhead int) *Scheduler {
	return &Scheduler{
		usecase:   usecase,
		logger:    logger,
		interval:  interval,
		daysAhead: daysAhead,
		enabled:   os.Getenv("SCHEDULER_ENABLED") != "false",
	}
}

// Run запускает фоновую генерацию задач. Блокирует goroutine до остановки контекста.
func (s *Scheduler) Run(ctx context.Context) {
	if !s.enabled {
		s.logger.Info("scheduler is disabled (SCHEDULER_ENABLED=false)")
		return
	}

	s.logger.Info("scheduler started",
		"interval", s.interval.String(),
		"days_ahead", s.daysAhead,
	)

	// Первый запуск сразу
	s.generate(ctx)

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("scheduler stopped")
			return
		case <-ticker.C:
			s.generate(ctx)
		}
	}
}

func (s *Scheduler) generate(ctx context.Context) {
	created, err := s.usecase.GenerateTasks(ctx, s.daysAhead)
	if err != nil {
		s.logger.Error("scheduler: failed to generate tasks", "error", err)
		return
	}

	if created > 0 {
		s.logger.Info("scheduler: generated task instances", "count", created)
	}
}
