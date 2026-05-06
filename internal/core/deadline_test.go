package core

import (
	"testing"
	"time"
)

func TestDeadlineRoundTrip(t *testing.T) {
	loc, _ := time.LoadLocation("UTC")

	cases := []struct {
		name string
		dt   time.Time
		want string
	}{
		{name: "case1", dt: time.Date(2026, 5, 5, 18, 0, 0, 0, loc), want: "2026-05-05 18:00"},
		{name: "case2", dt: time.Date(2026, 6, 10, 9, 30, 0, 0, loc), want: "2026-06-10 09:30"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			task := Task{Deadline: &tc.dt}

			serialized, err := task.MarshalJSON()
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			var unmarshaled Task
			err = unmarshaled.UnmarshalJSON(serialized)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			if unmarshaled.Deadline == nil {
				t.Fatal("Deadline should not be nil")
			}

			if got := unmarshaled.Deadline.Format("2006-01-02 15:04"); got != tc.want {
				t.Errorf("Expected %s, got %s", tc.want, got)
			}
		})
	}
}
