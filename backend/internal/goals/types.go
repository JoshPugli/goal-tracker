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

// DashboardResponse aggregates common data needed by the client in one request
type DashboardResponse struct {
    StatsDay   Stats        `json:"stats_day"`
    StatsWeek  Stats        `json:"stats_week"`
    StatsMonth Stats        `json:"stats_month"`
    Today      []TodayState `json:"today"`
}
