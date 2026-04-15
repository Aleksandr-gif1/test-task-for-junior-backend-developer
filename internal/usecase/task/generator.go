package task

import (
	"context"
	"fmt"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

// GenerateTasks генерирует экземпляры периодических задач на daysAhead дней вперёд.
// Возвращает количество созданных экземпляров.
func (s *Service) GenerateTasks(ctx context.Context, daysAhead int) (int, error) {
	if daysAhead < 1 {
		return 0, fmt.Errorf("%w: daysAhead must be >= 1", ErrInvalidInput)
	}

	tasks, err := s.repo.List(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to list tasks: %w", err)
	}

	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	endDate := today.AddDate(0, 0, daysAhead)

	created := 0

	for i := range tasks {
		task := &tasks[i]
		if !task.IsRecurring() {
			continue
		}

		dates, err := s.calculateDates(task, today, endDate)
		if err != nil {
			return created, fmt.Errorf("failed to calculate dates for task %d: %w", task.ID, err)
		}

		for _, date := range dates {
			exists, err := s.instanceRepo.Exists(ctx, task.ID, date)
			if err != nil {
				return created, fmt.Errorf("failed to check existence for task %d on %s: %w", task.ID, date.Format("2006-01-02"), err)
			}

			if exists {
				continue
			}

			instance := &taskdomain.TaskInstance{
				ParentTaskID:  task.ID,
				ScheduledDate: date,
				Title:         task.Title,
				Description:   task.Description,
				Status:        taskdomain.StatusNew,
				CreatedAt:     s.now(),
				UpdatedAt:     s.now(),
			}

			if _, err := s.instanceRepo.Create(ctx, instance); err != nil {
				return created, fmt.Errorf("failed to create instance for task %d on %s: %w", task.ID, date.Format("2006-01-02"), err)
			}

			created++
		}
	}

	return created, nil
}

// calculateDates вычисляет даты, на которые нужно создать экземпляры для задачи.
func (s *Service) calculateDates(task *taskdomain.Task, from, to time.Time) ([]time.Time, error) {
	if task.RecurrenceConfig == nil {
		return nil, fmt.Errorf("recurrence config is nil for task %d", task.ID)
	}

	switch task.RecurrenceType {
	case taskdomain.RecurrenceDaily:
		return s.calculateDaily(task, from, to)
	case taskdomain.RecurrenceMonthlyDate:
		return s.calculateMonthlyDate(task, from, to)
	case taskdomain.RecurrenceSpecific:
		return s.calculateSpecificDates(task, from, to)
	case taskdomain.RecurrenceEvenOdd:
		return s.calculateEvenOdd(task, from, to)
	default:
		return nil, fmt.Errorf("unsupported recurrence type: %s", task.RecurrenceType)
	}
}

// calculateDaily — каждый N-й день, начиная с даты создания задачи.
func (s *Service) calculateDaily(task *taskdomain.Task, from, to time.Time) ([]time.Time, error) {
	every := task.RecurrenceConfig.Every
	if every < 1 {
		every = 1
	}

	var dates []time.Time
	start := task.CreatedAt
	if start.After(from) {
		start = from
	}

	for d := start; !d.After(to); d = d.AddDate(0, 0, 1) {
		if d.Before(from) {
			continue
		}

		daysDiff := int(d.Sub(start).Hours()/24) + 1
		if daysDiff%every == 0 {
			dates = append(dates, d)
		}
	}

	return dates, nil
}

// calculateMonthlyDate — каждое N-е число месяца.
func (s *Service) calculateMonthlyDate(task *taskdomain.Task, from, to time.Time) ([]time.Time, error) {
	dayOfMonth := task.RecurrenceConfig.DayOfMonth
	if dayOfMonth < 1 || dayOfMonth > 31 {
		return nil, fmt.Errorf("invalid day_of_month: %d", dayOfMonth)
	}

	var dates []time.Time

	for d := from; !d.After(to); d = d.AddDate(0, 1, 0) {
		// Определяем максимальный день в месяце
		lastDay := daysInMonth(d.Year(), d.Month())
		targetDay := dayOfMonth
		if targetDay > lastDay {
			targetDay = lastDay
		}

		candidate := time.Date(d.Year(), d.Month(), targetDay, 0, 0, 0, 0, time.UTC)
		if !candidate.Before(from) && !candidate.After(to) {
			dates = append(dates, candidate)
		}
	}

	return dates, nil
}

// calculateSpecificDates — только указанные даты.
func (s *Service) calculateSpecificDates(task *taskdomain.Task, from, to time.Time) ([]time.Time, error) {
	var dates []time.Time

	for _, dateStr := range task.RecurrenceConfig.Dates {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}

		if !date.Before(from) && !date.After(to) {
			dates = append(dates, date)
		}
	}

	return dates, nil
}

// calculateEvenOdd — чётные или нечётные дни месяца.
func (s *Service) calculateEvenOdd(task *taskdomain.Task, from, to time.Time) ([]time.Time, error) {
	parity := task.RecurrenceConfig.Parity
	if parity != "even" && parity != "odd" {
		return nil, fmt.Errorf("invalid parity: %s", parity)
	}

	var dates []time.Time

	for d := from; !d.After(to); d = d.AddDate(0, 0, 1) {
		day := d.Day()
		isEven := day%2 == 0

		if (parity == "even" && isEven) || (parity == "odd" && !isEven) {
			dates = append(dates, d)
		}
	}

	return dates, nil
}

// daysInMonth возвращает количество дней в месяце.
func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
