package stats

import (
	"testing"
	"time"

	"github.com/programmersd21/kairo/internal/core"
)

func completedTaskAt(t time.Time) core.Task {
	return core.Task{CompletedAt: &t}
}

func localDay(t time.Time) time.Time {
	lt := t.Local()
	return time.Date(lt.Year(), lt.Month(), lt.Day(), 12, 0, 0, 0, time.Local)
}

func TestComputeStreaksEmpty(t *testing.T) {
	got := computeStreaks(nil)
	if got.Current != 0 || got.Longest != 0 || got.Warning {
		t.Fatalf("computeStreaks(nil) = %+v, want zero streaks with no warning", got)
	}
}

func TestComputeStreaksCountsMultipleTasksOnSameDayOnce(t *testing.T) {
	now := time.Now()
	today := localDay(now)
	yesterday := today.AddDate(0, 0, -1)

	got := computeStreaks([]core.Task{
		completedTaskAt(today),
		completedTaskAt(today.Add(2 * time.Hour)),
		completedTaskAt(yesterday),
	})

	if got.Current != 2 || got.Longest != 2 {
		t.Fatalf("computeStreaks duplicate day = %+v, want current=2 longest=2", got)
	}
}

func TestComputeStreaksAllowsCurrentStreakThroughYesterday(t *testing.T) {
	now := time.Now()
	yesterday := localDay(now).AddDate(0, 0, -1)
	twoDaysAgo := yesterday.AddDate(0, 0, -1)

	got := computeStreaks([]core.Task{
		completedTaskAt(twoDaysAgo),
		completedTaskAt(yesterday),
	})

	if got.Current != 2 || got.Longest != 2 {
		t.Fatalf("computeStreaks through yesterday = %+v, want current=2 longest=2", got)
	}
}

func TestComputeStreaksCurrentIsZeroWhenLatestCompletionIsOlder(t *testing.T) {
	old := localDay(time.Now()).AddDate(0, 0, -3)

	got := computeStreaks([]core.Task{completedTaskAt(old)})

	if got.Current != 0 || got.Longest != 1 {
		t.Fatalf("computeStreaks stale completion = %+v, want current=0 longest=1", got)
	}
}

func TestComputeStreaksLongestCanExceedCurrent(t *testing.T) {
	today := localDay(time.Now())
	yesterday := today.AddDate(0, 0, -1)
	lastWeek := today.AddDate(0, 0, -7)

	got := computeStreaks([]core.Task{
		completedTaskAt(lastWeek),
		completedTaskAt(lastWeek.AddDate(0, 0, 1)),
		completedTaskAt(lastWeek.AddDate(0, 0, 2)),
		completedTaskAt(yesterday),
	})

	if got.Current != 1 || got.Longest != 3 {
		t.Fatalf("computeStreaks longest vs current = %+v, want current=1 longest=3", got)
	}
}
