package goals

import "time"

type Goal struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Completion struct {
	GoalID string    `json:"goal_id"`
	Date   time.Time `json:"date"`
}

type TodayState struct {
	Goal      Goal `json:"goal"`
	Completed bool `json:"completed"`
}

type Stats struct {
	Window    string `json:"window"`
	Completed int    `json:"completed"`
	Total     int    `json:"total"`
}
