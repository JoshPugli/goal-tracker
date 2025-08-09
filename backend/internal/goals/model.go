package goals

import "time"

type GoalType string

const (
	GoalTypeBoolean  GoalType = "boolean"
	GoalTypeNumeric  GoalType = "numeric"
	GoalTypeDuration GoalType = "duration"
)

type Goal struct {
	ID           string     `json:"id" db:"id"`
	UserID       string     `json:"user_id" db:"user_id"`
	Title        string     `json:"title" db:"title"`
	Description  *string    `json:"description" db:"description"`
	GoalType     GoalType   `json:"goal_type" db:"goal_type"`
	TargetValue  *float64   `json:"target_value" db:"target_value"`
	Unit         *string    `json:"unit" db:"unit"`
	IsActive     bool       `json:"is_active" db:"is_active"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

type DailyGoalInstance struct {
	ID             string     `json:"id" db:"id"`
	GoalID         string     `json:"goal_id" db:"goal_id"`
	UserID         string     `json:"user_id" db:"user_id"`
	Date           time.Time  `json:"date" db:"date"`
	TargetValue    *float64   `json:"target_value" db:"target_value"`
	CompletedValue *float64   `json:"completed_value" db:"completed_value"`
	IsCompleted    bool       `json:"is_completed" db:"is_completed"`
	CompletedAt    *time.Time `json:"completed_at" db:"completed_at"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
}

type CreateGoalRequest struct {
	Title       string   `json:"title"`
	Description *string  `json:"description"`
	GoalType    GoalType `json:"goal_type"`
	TargetValue *float64 `json:"target_value"`
	Unit        *string  `json:"unit"`
}

type UpdateGoalRequest struct {
	Title       *string  `json:"title"`
	Description *string  `json:"description"`
	TargetValue *float64 `json:"target_value"`
	Unit        *string  `json:"unit"`
	IsActive    *bool    `json:"is_active"`
}

type UpdateDailyInstanceRequest struct {
	CompletedValue *float64 `json:"completed_value"`
	IsCompleted    *bool    `json:"is_completed"`
}

type GoalWithTodayInstance struct {
	Goal          Goal               `json:"goal"`
	TodayInstance *DailyGoalInstance `json:"today_instance"`
}