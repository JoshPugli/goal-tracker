package goals

import (
	"sort"
	"sync"
	"time"
)

type Store struct {
	mu                sync.RWMutex
	goalsByUser       map[string]map[string]Goal
	completionsByUser map[string][]Completion
}

func NewStore() *Store {
	s := &Store{
		goalsByUser:       make(map[string]map[string]Goal),
		completionsByUser: make(map[string][]Completion),
	}
	// seed example data for demo user
	demo := "demo"
	s.goalsByUser[demo] = map[string]Goal{
		"g1": {ID: "g1", Name: "Drink water"},
		"g2": {ID: "g2", Name: "Read 20 min"},
		"g3": {ID: "g3", Name: "Exercise"},
	}
	now := time.Now()
	s.completionsByUser[demo] = []Completion{
		{GoalID: "g1", Date: now},                    // today
		{GoalID: "g2", Date: now.AddDate(0, 0, -1)},  // yesterday
		{GoalID: "g2", Date: now.AddDate(0, 0, -6)},  // last week
		{GoalID: "g3", Date: now.AddDate(0, 0, -20)}, // last month
		{GoalID: "g1", Date: now.AddDate(0, 0, -8)},  // ~ last week
	}
	return s
}

func (s *Store) ListGoals(userID string) []Goal {
	s.mu.RLock()
	defer s.mu.RUnlock()
	userGoals := s.goalsByUser[userID]
	out := make([]Goal, 0, len(userGoals))
	for _, g := range userGoals {
		out = append(out, g)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Name == out[j].Name {
			return out[i].ID < out[j].ID
		}
		return out[i].Name < out[j].Name
	})
	return out
}

func (s *Store) IsCompletedToday(userID, goalID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	today := time.Now()
	for _, c := range s.completionsByUser[userID] {
		if c.GoalID == goalID && sameDay(c.Date, today) {
			return true
		}
	}
	return false
}

func (s *Store) ToggleCompleteToday(userID, goalID string, done bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	today := time.Now()
	comps := s.completionsByUser[userID]
	// remove any existing today completion
	idx := -1
	for i, c := range comps {
		if c.GoalID == goalID && sameDay(c.Date, today) {
			idx = i
			break
		}
	}
	if idx >= 0 {
		// uncomplete first
		comps = append(comps[:idx], comps[idx+1:]...)
	}
	if done {
		comps = append(comps, Completion{GoalID: goalID, Date: today})
	}
	s.completionsByUser[userID] = comps
}

func (s *Store) TodayState(userID string) []TodayState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	userGoals := s.goalsByUser[userID]
	out := make([]TodayState, 0, len(userGoals))
	for _, g := range userGoals {
		out = append(out, TodayState{Goal: g, Completed: s.isCompletedTodayLocked(userID, g.ID)})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Goal.Name == out[j].Goal.Name {
			return out[i].Goal.ID < out[j].Goal.ID
		}
		return out[i].Goal.Name < out[j].Goal.Name
	})
	return out
}

func (s *Store) isCompletedTodayLocked(userID, goalID string) bool {
	today := time.Now()
	for _, c := range s.completionsByUser[userID] {
		if c.GoalID == goalID && sameDay(c.Date, today) {
			return true
		}
	}
	return false
}

func (s *Store) Stats(userID string, window string) (completed int, total int) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	total = len(s.goalsByUser[userID])
	now := time.Now()
	for _, c := range s.completionsByUser[userID] {
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
