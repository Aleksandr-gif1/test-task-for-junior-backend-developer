package task

import (
	"fmt"
	"strings"

	taskdomain "example.com/taskservice/internal/domain/task"
)

// validateCreateInput валидирует и нормализует входные данные для создания задачи.
func validateCreateInput(input CreateInput) (CreateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return CreateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	// Если не указан статус — ставим new
	if input.Status == "" {
		input.Status = taskdomain.StatusNew
	}

	if !input.Status.Valid() {
		return CreateInput{}, fmt.Errorf("%w: invalid status %q", ErrInvalidInput, input.Status)
	}

	// Валидация периодичности
	if input.RecurrenceType != "" {
		if !input.RecurrenceType.Valid() {
			return CreateInput{}, fmt.Errorf("%w: invalid recurrence_type %q", ErrInvalidInput, input.RecurrenceType)
		}

		if input.RecurrenceConfig == nil {
			return CreateInput{}, fmt.Errorf("%w: recurrence_config is required when recurrence_type is set", ErrInvalidInput)
		}

		if err := input.RecurrenceConfig.Valid(input.RecurrenceType); err != nil {
			return CreateInput{}, fmt.Errorf("%w: %w", ErrInvalidInput, err)
		}

		input.RecurrenceActive = true
	}

	return input, nil
}

// validateUpdateInput валидирует входные данные для обновления задачи.
func validateUpdateInput(input UpdateInput) (UpdateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return UpdateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if !input.Status.Valid() {
		return UpdateInput{}, fmt.Errorf("%w: invalid status %q", ErrInvalidInput, input.Status)
	}

	// Валидация периодичности
	if input.RecurrenceType != "" {
		if !input.RecurrenceType.Valid() {
			return UpdateInput{}, fmt.Errorf("%w: invalid recurrence_type %q", ErrInvalidInput, input.RecurrenceType)
		}

		if input.RecurrenceConfig == nil {
			return UpdateInput{}, fmt.Errorf("%w: recurrence_config is required when recurrence_type is set", ErrInvalidInput)
		}

		if err := input.RecurrenceConfig.Valid(input.RecurrenceType); err != nil {
			return UpdateInput{}, fmt.Errorf("%w: %w", ErrInvalidInput, err)
		}

		input.RecurrenceActive = true
	}

	return input, nil
}
