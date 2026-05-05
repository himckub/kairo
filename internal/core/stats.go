package core

import "time"

type Session struct {
	ID         string
	StartTime  time.Time
	EndTime    *time.Time
	FocusScore int
}

type Event struct {
	ID        int64
	Type      string
	TaskID    string
	Timestamp time.Time
	Metadata  string
}

const (
	EventTypeTaskCreated   = "task_created"
	EventTypeTaskCompleted = "task_completed"
	EventTypeTaskDeleted   = "task_deleted"
	EventTypeAppOpened     = "app_opened"
	EventTypeStatsOpened   = "stats_opened"
)
