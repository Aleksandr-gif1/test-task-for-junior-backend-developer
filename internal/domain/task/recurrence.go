package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// RecurrenceType определяет тип периодичности задачи.
type RecurrenceType string

const (
	RecurrenceDaily       RecurrenceType = "daily"
	RecurrenceMonthlyDate RecurrenceType = "monthly_date"
	RecurrenceSpecific    RecurrenceType = "specific_dates"
	RecurrenceEvenOdd     RecurrenceType = "even_odd"
)

func (r RecurrenceType) Valid() bool {
	switch r {
	case RecurrenceDaily, RecurrenceMonthlyDate, RecurrenceSpecific, RecurrenceEvenOdd:
		return true
	default:
		return false
	}
}

// RecurrenceConfig хранит настройки периодичности в зависимости от типа.
// daily:        {"every": 2}            — каждый 2-й день
// monthly_date: {"day_of_month": 15}    — 15-го числа каждого месяца
// specific_dates: {"dates": ["2026-04-20", "2026-05-10"]} — конкретные даты
// even_odd:     {"parity": "even"}      — чётные (или "odd" — нечётные)
type RecurrenceConfig struct {
	Every        int      `json:"every,omitempty"`
	DayOfMonth   int      `json:"day_of_month,omitempty"`
	Dates        []string `json:"dates,omitempty"`
	Parity       string   `json:"parity,omitempty"` // "even" или "odd"
}

// Valid проверяет корректность конфигурации периодичности.
func (c RecurrenceConfig) Valid(t RecurrenceType) error {
	if !t.Valid() {
		return errors.New("invalid recurrence type")
	}

	switch t {
	case RecurrenceDaily:
		if c.Every < 1 {
			return errors.New("daily recurrence requires 'every' >= 1")
		}
	case RecurrenceMonthlyDate:
		if c.DayOfMonth < 1 || c.DayOfMonth > 31 {
			return errors.New("monthly_date recurrence requires 'day_of_month' between 1 and 31")
		}
	case RecurrenceSpecific:
		if len(c.Dates) == 0 {
			return errors.New("specific_dates recurrence requires at least one date")
		}
		for _, d := range c.Dates {
			if _, err := time.Parse("2006-01-02", d); err != nil {
				return fmt.Errorf("invalid date %q in specific_dates: %w", d, err)
			}
		}
	case RecurrenceEvenOdd:
		if c.Parity != "even" && c.Parity != "odd" {
			return errors.New("even_odd recurrence requires 'parity' to be 'even' or 'odd'")
		}
	}

	return nil
}

// MarshalJSON сериализует RecurrenceConfig в JSON.
func (c *RecurrenceConfig) MarshalJSON() ([]byte, error) {
	type Alias RecurrenceConfig
	return json.Marshal((*Alias)(c))
}

// UnmarshalJSON десериализует JSON в RecurrenceConfig.
func (c *RecurrenceConfig) UnmarshalJSON(data []byte) error {
	type Alias RecurrenceConfig
	aux := (*Alias)(c)
	return json.Unmarshal(data, aux)
}

// TaskInstance представляет экземпляр периодической задачи.
type TaskInstance struct {
	ID          int64     `json:"id"`
	ParentTaskID int64    `json:"parent_task_id"`
	ScheduledDate time.Time `json:"scheduled_date"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

