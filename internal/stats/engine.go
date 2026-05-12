package stats

import (
	"fmt"
	"sort"
	"time"

	"github.com/programmersd21/kairo/internal/core"
)

type DNAData struct {
	PeakHours   []float64 // 24 hours, intensity 0-1
	Velocity    float64   // completions per hour
	Consistency float64   // 0-1
}

type Insight struct {
	Title string
	Value string
	Type  string // positive, negative, neutral
}

type TimelinePoint struct {
	Date      time.Time
	Completed int
	Focus     int
}

type TagCluster struct {
	Tag   string
	Count int
	Score float64
}

type StreakData struct {
	Current int
	Longest int
	Warning bool
}

type DashboardData struct {
	DNA         DNAData
	Momentum    float64
	Insights    []Insight
	Timeline    []TimelinePoint
	TagClusters []TagCluster
	Streaks     StreakData
}

func ComputeDashboard(tasks []core.Task, sessions []core.Session, events []core.Event) DashboardData {
	return DashboardData{
		DNA:         computeDNA(tasks, sessions, events),
		Momentum:    computeMomentum(tasks, sessions),
		Insights:    generateInsights(tasks, sessions, events),
		Timeline:    computeTimeline(tasks, sessions),
		TagClusters: computeTagClusters(tasks),
		Streaks:     computeStreaks(tasks),
	}
}

func computeDNA(tasks []core.Task, sessions []core.Session, events []core.Event) DNAData {
	peaks := make([]float64, 24)
	counts := make([]int, 24)

	for _, t := range tasks {
		if t.CompletedAt != nil {
			hour := t.CompletedAt.Hour()
			counts[hour]++
		}
	}

	maxCount := 0
	for _, c := range counts {
		if c > maxCount {
			maxCount = c
		}
	}

	if maxCount > 0 {
		for i, c := range counts {
			peaks[i] = float64(c) / float64(maxCount)
		}
	}

	return DNAData{
		PeakHours:   peaks,
		Velocity:    float64(len(tasks)) / 168.0, // dummy velocity for now
		Consistency: 0.85,
	}
}

func computeMomentum(tasks []core.Task, sessions []core.Session) float64 {
	// Momentum = (Recent Completions * 20) + (Recent Active Sessions * 5)
	score := 0.0
	now := time.Now()

	// Completions in last 3 days
	recentCompletions := 0
	for _, t := range tasks {
		if t.CompletedAt != nil && now.Sub(*t.CompletedAt) < 72*time.Hour {
			recentCompletions++
		}
	}
	score += float64(recentCompletions) * 20.0

	// Sessions in last 3 days
	recentSessions := 0
	for _, s := range sessions {
		if now.Sub(s.StartTime) < 72*time.Hour {
			recentSessions++
		}
	}
	score += float64(recentSessions) * 5.0

	if score > 100 {
		score = 100
	}
	return score
}

func generateInsights(tasks []core.Task, sessions []core.Session, events []core.Event) []Insight {
	var insights []Insight

	// 1. Peak Productivity Insight
	dayTasks := 0
	nightTasks := 0
	for _, t := range tasks {
		if t.CompletedAt != nil {
			h := t.CompletedAt.Hour()
			if h >= 6 && h < 18 {
				dayTasks++
			} else {
				nightTasks++
			}
		}
	}
	if nightTasks > dayTasks && nightTasks > 0 {
		diff := float64(nightTasks-dayTasks) / float64(dayTasks+1) * 100
		insights = append(insights, Insight{
			Title: "Peak Productivity",
			Value: fmt.Sprintf("You complete %.0f%% more tasks at night", diff),
			Type:  "positive",
		})
	} else if dayTasks > 0 {
		insights = append(insights, Insight{
			Title: "Peak Productivity",
			Value: "You are most productive during the day",
			Type:  "positive",
		})
	}

	// 2. Focus Insight
	if len(sessions) > 5 {
		lowFocus := 0
		for _, s := range sessions {
			if s.FocusScore < 50 {
				lowFocus++
			}
		}
		if lowFocus > len(sessions)/2 {
			insights = append(insights, Insight{
				Title: "Focus Trend",
				Value: "Your average focus score is on the lower side",
				Type:  "negative",
			})
		} else {
			insights = append(insights, Insight{
				Title: "Focus Trend",
				Value: "Your focus sessions are performing well",
				Type:  "positive",
			})
		}
	}

	// 3. Tag Mastery
	tagCounts := make(map[string]int)
	for _, t := range tasks {
		for _, tag := range t.Tags {
			tagCounts[tag]++
		}
	}
	bestTag := ""
	maxC := 0
	for tag, count := range tagCounts {
		if count > maxC {
			maxC = count
			bestTag = tag
		}
	}
	if bestTag != "" {
		insights = append(insights, Insight{
			Title: "Tag Mastery",
			Value: fmt.Sprintf("Tag '%s' is your most frequent focus area", bestTag),
			Type:  "neutral",
		})
	}

	if len(insights) == 0 {
		insights = append(insights, Insight{Title: "Keep Going", Value: "Start completing tasks to see insights!", Type: "neutral"})
	}

	return insights
}

func computeTimeline(tasks []core.Task, sessions []core.Session) []TimelinePoint {
	// Last 7 days
	points := make([]TimelinePoint, 7)
	now := time.Now()
	for i := 0; i < 7; i++ {
		date := now.AddDate(0, 0, -i)
		points[6-i].Date = date

		count := 0
		for _, t := range tasks {
			if t.CompletedAt != nil && t.CompletedAt.YearDay() == date.YearDay() && t.CompletedAt.Year() == date.Year() {
				count++
			}
		}
		points[6-i].Completed = count
	}
	return points
}

func computeTagClusters(tasks []core.Task) []TagCluster {
	m := make(map[string]int)
	for _, t := range tasks {
		for _, tag := range t.Tags {
			m[tag]++
		}
	}
	var out []TagCluster
	for tag, count := range m {
		out = append(out, TagCluster{Tag: tag, Count: count, Score: float64(count)})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Count > out[j].Count })
	if len(out) > 8 {
		out = out[:8]
	}
	return out
}

func computeStreaks(tasks []core.Task) StreakData {
	completedDays := make(map[time.Time]struct{})
	for _, task := range tasks {
		if task.CompletedAt == nil {
			continue
		}
		completedAt := task.CompletedAt.Local()
		day := time.Date(completedAt.Year(), completedAt.Month(), completedAt.Day(), 0, 0, 0, 0, time.Local)
		completedDays[day] = struct{}{}
	}

	if len(completedDays) == 0 {
		return StreakData{Warning: false}
	}

	days := make([]time.Time, 0, len(completedDays))
	for day := range completedDays {
		days = append(days, day)
	}
	sort.Slice(days, func(i, j int) bool { return days[i].Before(days[j]) })

	longest := 1
	run := 1
	for i := 1; i < len(days); i++ {
		if days[i].Equal(days[i-1].AddDate(0, 0, 1)) {
			run++
		} else {
			run = 1
		}
		if run > longest {
			longest = run
		}
	}

	now := time.Now().Local()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	current := streakEndingAt(completedDays, today)
	if current == 0 {
		current = streakEndingAt(completedDays, today.AddDate(0, 0, -1))
	}

	return StreakData{
		Current: current,
		Longest: longest,
		Warning: false,
	}
}

func streakEndingAt(completedDays map[time.Time]struct{}, end time.Time) int {
	streak := 0
	for day := end; ; day = day.AddDate(0, 0, -1) {
		if _, ok := completedDays[day]; !ok {
			break
		}
		streak++
	}
	return streak
}
