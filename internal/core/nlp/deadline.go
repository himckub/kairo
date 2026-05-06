package nlp

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/olebedev/when"
	en "github.com/olebedev/when/rules/en"
)

var parser = func() *when.Parser {
	p := when.New(nil)
	p.Add(en.All...)
	return p
}()

var (
	// Strict ISO-ish formats we support for deterministic parsing.
	// We intentionally parse these before the NLP parser because the NLP parser
	// may interpret the MM-dd portion as HH:mm (e.g. "2026-05-05" -> "05:05").
	reISODateTime = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}$`)
	reISODateOnly = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
)

func ParseDeadline(input string, now time.Time) (*time.Time, error) {
	s := strings.TrimSpace(input)
	if s == "" {
		return nil, nil
	}

	// Deterministic parsing for exact user-entered timestamps.
	loc := now.Location()
	switch {
	case reISODateTime.MatchString(s):
		if t, err := time.ParseInLocation("2006-01-02 15:04", s, loc); err == nil {
			return &t, nil
		}
	case reISODateOnly.MatchString(s):
		if t, err := time.ParseInLocation("2006-01-02", s, loc); err == nil {
			// Default to midnight local time when only a date is provided.
			t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
			return &t, nil
		}
	}

	res, err := parser.Parse(s, now)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, errors.New("could not parse deadline")
	}
	t := res.Time
	return &t, nil
}
