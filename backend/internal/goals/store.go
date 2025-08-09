package goals

import (
	"sync"
	"time"
)

type Store struct {
	mu          sync.RWMutex
	goals       map[string]Goal
	completions []Completion
}

func NewStore() *Store {
	s := &Store{
		goals:       make(map[string]Goal),
		completions: make([]Completion, 0, 64),
	}
	// seed example data
	s.goals["g1"] = Goal{ID: "g1", Name: "Drink water"}
	s.goals["g2"] = Goal{ID: "g2", Name: "Read 20 min"}
	s.goals["g3"] = Goal{ID: "g3", Name: "Exercise"}
	s.completions = append(s.completions, Completion{GoalID: "g1", Date: time.Now().AddDate(0, 0, -1)})
	s.completions = append(s.completions, Completion{GoalID: "g2", Date: time.Now().AddDate(0, 0, -6)})
	s.completions = append(s.completions, Completion{GoalID: "zzz", Date: time.Now().AddDate(0, 0, -20)})
	return s
}

func (s *Store) ListGoals() []Goal {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Goal, 0, len(s.goals))
	for _, g := range s.goals {
		out = append(out, g)
	}
	return out
}

func (s *Store) IsCompletedToday(goalID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	today := time.Now()
	for _, c := range s.completions {
		if c.GoalID == goalID && sameDay(c.Date, today) {
			return true
		}
	}
	return false
}

func (s *Store) ToggleCompleteToday(goalID string, done bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	today := time.Now()
	// remove any existing today completion
	idx := -1
	for i, c := range s.completions {
		if c.GoalID == goalID && sameDay(c.Date, today) {
			idx = i
			break
		}
	}
	if idx >= 0 {
		// uncomplete first
		s.completions = append(s.completions[:idx], s.completions[idx+1:]...)
	}
	if done {
		s.completions = append(s.completions, Completion{GoalID: goalID, Date: today})
	}
}

func (s *Store) TodayState() []TodayState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]TodayState, 0, len(s.goals))
	for _, g := range s.goals {
		out = append(out, TodayState{Goal: g, Completed: s.isCompletedTodayLocked(g.ID)})
	}
	return out
}

func (s *Store) isCompletedTodayLocked(goalID string) bool {
	today := time.Now()
	for _, c := range s.completions {
		if c.GoalID == goalID && sameDay(c.Date, today) {
			return true
		}
	}
	return false
}

func (s *Store) Stats(window string) (completed int, total int) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	total = len(s.goals)
	now := time.Now()
	for _, c := range s.completions {
		switch window {
		case "day":
			if sameDay(c.Date, now) {
				completed++
			}
		case "week":
			if sameWeek(c.Date, now) {
				completed++
			}
		case "month":
			if sameMonth(c.Date, now) {
				completed++
			}
		}
	}
	if window == "week" {
		total = total * 7
	}
	if window == "month" {
		y, m, _ := now.Date()
		first := time.Date(y, m, 1, 0, 0, 0, 0, now.Location())
		totalDays := first.AddDate(0, 1, -1).Day()
		total = total * totalDays
	}
	return
}

func sameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

func sameWeek(a, b time.Time) bool {
	aa, aw := a.ISOWeek()
	ba, bw := b.ISOWeek()
	return aa == ba && aw == bw
}

func sameMonth(a, b time.Time) bool {
	ay, am, _ := a.Date()
	by, bm, _ := b.Date()
	return ay == by && am == bm
}
