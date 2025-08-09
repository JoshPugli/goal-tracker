package goals

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateGoal(userID string, req CreateGoalRequest) (*Goal, error) {
	goal := &Goal{
		ID:          uuid.New().String(),
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		GoalType:    req.GoalType,
		TargetValue: req.TargetValue,
		Unit:        req.Unit,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	query := `
		INSERT INTO goals (id, user_id, title, description, goal_type, target_value, unit, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.Exec(query, goal.ID, goal.UserID, goal.Title, goal.Description, goal.GoalType, goal.TargetValue, goal.Unit, goal.IsActive, goal.CreatedAt, goal.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create goal: %w", err)
	}

	return goal, nil
}

func (r *Repository) GetGoalsByUserID(userID string) ([]Goal, error) {
	query := `
		SELECT id, user_id, title, description, goal_type, target_value, unit, is_active, created_at, updated_at
		FROM goals
		WHERE user_id = $1 AND is_active = true
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get goals: %w", err)
	}
	defer rows.Close()

	var goals []Goal
	for rows.Next() {
		var goal Goal
		err := rows.Scan(&goal.ID, &goal.UserID, &goal.Title, &goal.Description, &goal.GoalType, &goal.TargetValue, &goal.Unit, &goal.IsActive, &goal.CreatedAt, &goal.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan goal: %w", err)
		}
		goals = append(goals, goal)
	}

	return goals, nil
}

func (r *Repository) GetGoalByID(goalID, userID string) (*Goal, error) {
	query := `
		SELECT id, user_id, title, description, goal_type, target_value, unit, is_active, created_at, updated_at
		FROM goals
		WHERE id = $1 AND user_id = $2
	`
	var goal Goal
	err := r.db.QueryRow(query, goalID, userID).Scan(&goal.ID, &goal.UserID, &goal.Title, &goal.Description, &goal.GoalType, &goal.TargetValue, &goal.Unit, &goal.IsActive, &goal.CreatedAt, &goal.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("goal not found")
		}
		return nil, fmt.Errorf("failed to get goal: %w", err)
	}

	return &goal, nil
}

func (r *Repository) UpdateGoal(goalID, userID string, req UpdateGoalRequest) (*Goal, error) {
	goal, err := r.GetGoalByID(goalID, userID)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		goal.Title = *req.Title
	}
	if req.Description != nil {
		goal.Description = req.Description
	}
	if req.TargetValue != nil {
		goal.TargetValue = req.TargetValue
	}
	if req.Unit != nil {
		goal.Unit = req.Unit
	}
	if req.IsActive != nil {
		goal.IsActive = *req.IsActive
	}
	goal.UpdatedAt = time.Now()

	query := `
		UPDATE goals 
		SET title = $1, description = $2, target_value = $3, unit = $4, is_active = $5, updated_at = $6
		WHERE id = $7 AND user_id = $8
	`
	_, err = r.db.Exec(query, goal.Title, goal.Description, goal.TargetValue, goal.Unit, goal.IsActive, goal.UpdatedAt, goalID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to update goal: %w", err)
	}

	return goal, nil
}

func (r *Repository) DeleteGoal(goalID, userID string) error {
	query := `UPDATE goals SET is_active = false WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(query, goalID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete goal: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("goal not found")
	}

	return nil
}

func (r *Repository) GetOrCreateDailyInstance(goalID, userID string, date time.Time) (*DailyGoalInstance, error) {
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	
	query := `
		SELECT id, goal_id, user_id, date, target_value, completed_value, is_completed, completed_at, created_at
		FROM daily_goal_instances
		WHERE goal_id = $1 AND user_id = $2 AND date = $3
	`
	var instance DailyGoalInstance
	err := r.db.QueryRow(query, goalID, userID, dateOnly).Scan(
		&instance.ID, &instance.GoalID, &instance.UserID, &instance.Date,
		&instance.TargetValue, &instance.CompletedValue, &instance.IsCompleted,
		&instance.CompletedAt, &instance.CreatedAt)
	
	if err == nil {
		return &instance, nil
	}
	
	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get daily instance: %w", err)
	}

	goal, err := r.GetGoalByID(goalID, userID)
	if err != nil {
		return nil, err
	}

	instance = DailyGoalInstance{
		ID:             uuid.New().String(),
		GoalID:         goalID,
		UserID:         userID,
		Date:           dateOnly,
		TargetValue:    goal.TargetValue,
		CompletedValue: nil,
		IsCompleted:    false,
		CompletedAt:    nil,
		CreatedAt:      time.Now(),
	}

	insertQuery := `
		INSERT INTO daily_goal_instances (id, goal_id, user_id, date, target_value, completed_value, is_completed, completed_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err = r.db.Exec(insertQuery, instance.ID, instance.GoalID, instance.UserID, instance.Date, 
		instance.TargetValue, instance.CompletedValue, instance.IsCompleted, instance.CompletedAt, instance.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create daily instance: %w", err)
	}

	return &instance, nil
}

func (r *Repository) UpdateDailyInstance(instanceID, userID string, req UpdateDailyInstanceRequest) (*DailyGoalInstance, error) {
	query := `
		SELECT id, goal_id, user_id, date, target_value, completed_value, is_completed, completed_at, created_at
		FROM daily_goal_instances
		WHERE id = $1 AND user_id = $2
	`
	var instance DailyGoalInstance
	err := r.db.QueryRow(query, instanceID, userID).Scan(
		&instance.ID, &instance.GoalID, &instance.UserID, &instance.Date,
		&instance.TargetValue, &instance.CompletedValue, &instance.IsCompleted,
		&instance.CompletedAt, &instance.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("daily instance not found")
		}
		return nil, fmt.Errorf("failed to get daily instance: %w", err)
	}

	if req.CompletedValue != nil {
		instance.CompletedValue = req.CompletedValue
	}
	if req.IsCompleted != nil {
		instance.IsCompleted = *req.IsCompleted
		if *req.IsCompleted && instance.CompletedAt == nil {
			now := time.Now()
			instance.CompletedAt = &now
		} else if !*req.IsCompleted {
			instance.CompletedAt = nil
		}
	}

	updateQuery := `
		UPDATE daily_goal_instances 
		SET completed_value = $1, is_completed = $2, completed_at = $3
		WHERE id = $4 AND user_id = $5
	`
	_, err = r.db.Exec(updateQuery, instance.CompletedValue, instance.IsCompleted, instance.CompletedAt, instanceID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to update daily instance: %w", err)
	}

	return &instance, nil
}

func (r *Repository) GetGoalsWithTodayInstances(userID string) ([]GoalWithTodayInstance, error) {
	today := time.Now().UTC()
	dateOnly := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)

	query := `
		SELECT 
			g.id, g.user_id, g.title, g.description, g.goal_type, g.target_value, g.unit, g.is_active, g.created_at, g.updated_at,
			dgi.id, dgi.goal_id, dgi.user_id, dgi.date, dgi.target_value, dgi.completed_value, dgi.is_completed, dgi.completed_at, dgi.created_at
		FROM goals g
		LEFT JOIN daily_goal_instances dgi ON g.id = dgi.goal_id AND dgi.date = $2
		WHERE g.user_id = $1 AND g.is_active = true
		ORDER BY g.created_at DESC
	`

	rows, err := r.db.Query(query, userID, dateOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to get goals with today instances: %w", err)
	}
	defer rows.Close()

	var results []GoalWithTodayInstance
	for rows.Next() {
		var goal Goal
		var instance DailyGoalInstance
		var instanceID, instanceGoalID, instanceUserID, instanceCreatedAt sql.NullString
		var instanceDate, instanceCompletedAt sql.NullTime
		var instanceTargetValue, instanceCompletedValue sql.NullFloat64
		var instanceIsCompleted sql.NullBool

		err := rows.Scan(
			&goal.ID, &goal.UserID, &goal.Title, &goal.Description, &goal.GoalType, &goal.TargetValue, &goal.Unit, &goal.IsActive, &goal.CreatedAt, &goal.UpdatedAt,
			&instanceID, &instanceGoalID, &instanceUserID, &instanceDate, &instanceTargetValue, &instanceCompletedValue, &instanceIsCompleted, &instanceCompletedAt, &instanceCreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan goal with instance: %w", err)
		}

		result := GoalWithTodayInstance{
			Goal: goal,
		}

		if instanceID.Valid {
			instance.ID = instanceID.String
			instance.GoalID = instanceGoalID.String
			instance.UserID = instanceUserID.String
			instance.Date = instanceDate.Time
			if instanceTargetValue.Valid {
				instance.TargetValue = &instanceTargetValue.Float64
			}
			if instanceCompletedValue.Valid {
				instance.CompletedValue = &instanceCompletedValue.Float64
			}
			instance.IsCompleted = instanceIsCompleted.Bool
			if instanceCompletedAt.Valid {
				instance.CompletedAt = &instanceCompletedAt.Time
			}
			instanceCreatedTime, _ := time.Parse(time.RFC3339, instanceCreatedAt.String)
			instance.CreatedAt = instanceCreatedTime

			result.TodayInstance = &instance
		}

		results = append(results, result)
	}

	return results, nil
}

func (r *Repository) GetDailyInstancesByGoal(goalID, userID string, startDate, endDate time.Time) ([]DailyGoalInstance, error) {
	query := `
		SELECT id, goal_id, user_id, date, target_value, completed_value, is_completed, completed_at, created_at
		FROM daily_goal_instances
		WHERE goal_id = $1 AND user_id = $2 AND date >= $3 AND date <= $4
		ORDER BY date DESC
	`

	rows, err := r.db.Query(query, goalID, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily instances: %w", err)
	}
	defer rows.Close()

	var instances []DailyGoalInstance
	for rows.Next() {
		var instance DailyGoalInstance
		err := rows.Scan(&instance.ID, &instance.GoalID, &instance.UserID, &instance.Date,
			&instance.TargetValue, &instance.CompletedValue, &instance.IsCompleted,
			&instance.CompletedAt, &instance.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan daily instance: %w", err)
		}
		instances = append(instances, instance)
	}

	return instances, nil
}