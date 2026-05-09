package focus

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/programmersd21/kairo/internal/core"
	"github.com/programmersd21/kairo/internal/ui/styles"
)

type State int

const (
	StateIdle State = iota
	StateFocus
	StateShortBreak
	StateLongBreak
)

type TickMsg time.Time

type Model struct {
	styles styles.Styles
	State  State
	Task   *core.Task

	Timer        time.Duration
	TotalFocused time.Duration

	lastTick time.Time
	Active   bool

	width  int
	height int

	FocusDuration      time.Duration
	ShortBreakDuration time.Duration
	LongBreakDuration  time.Duration

	SessionID string
}

type SessionDoneMsg struct {
	Session core.FocusSession
}

func New(s styles.Styles) Model {
	return Model{
		styles:             s,
		State:              StateIdle,
		FocusDuration:      25 * time.Minute,
		ShortBreakDuration: 5 * time.Minute,
		LongBreakDuration:  15 * time.Minute,
		Timer:              25 * time.Minute,
	}
}

func (m *Model) SetSize(w, h int) {
	m.width, m.height = w, h
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch x := msg.(type) {
	case TickMsg:
		if !m.Active {
			return m, nil
		}

		now := time.Time(x)
		if m.lastTick.IsZero() {
			m.lastTick = now
		}

		elapsed := now.Sub(m.lastTick)
		m.lastTick = now
		m.Timer -= elapsed

		if m.State == StateFocus {
			m.TotalFocused += elapsed
		}

		if m.Timer <= 0 {
			m.Timer = 0
			m.Active = false

			if m.State == StateFocus {
				session := core.FocusSession{
					ID:        m.SessionID,
					StartTime: now.Add(-m.FocusDuration), // Approximate
					EndTime:   &now,
					Duration:  m.FocusDuration,
				}
				if m.Task != nil {
					session.TaskID = m.Task.ID
				}

				return m, func() tea.Msg {
					return SessionDoneMsg{Session: session}
				}
			}
			return m, nil
		}

		return m, m.tick()

	case tea.KeyMsg:
		switch x.String() {
		case "enter":
			if m.State == StateIdle {
				m.startFocus()
				return m, m.tick()
			}
			if !m.Active {
				m.Active = true
				m.lastTick = time.Now()
				return m, m.tick()
			} else {
				m.Active = false
			}
		case "r": // Reset
			m.Active = false
			m.State = StateIdle
			m.Timer = m.FocusDuration
		case "s": // Skip to break/focus
			m.Active = false
			if m.State == StateFocus {
				m.startShortBreak()
			} else {
				m.startFocus()
			}
		}
	}
	return m, nil
}

func (m *Model) startFocus() {
	m.State = StateFocus
	m.Timer = m.FocusDuration
	m.Active = true
	m.lastTick = time.Now()
	m.SessionID = fmt.Sprintf("fs_%d", time.Now().UnixNano())
}

func (m *Model) startShortBreak() {
	m.State = StateShortBreak
	m.Timer = m.ShortBreakDuration
	m.Active = true
	m.lastTick = time.Now()
}

func (m Model) View() string {
	if m.width <= 0 || m.height <= 0 {
		return ""
	}

	accent := m.styles.Theme.Accent
	switch m.State {
	case StateShortBreak, StateLongBreak:
		accent = m.styles.Theme.Good
	}

	// Calculate Progress for "Warp Track"
	total := m.FocusDuration
	switch m.State {
	case StateShortBreak:
		total = m.ShortBreakDuration
	case StateLongBreak:
		total = m.LongBreakDuration
	}

	progress := 1.0
	if total > 0 {
		progress = 1.0 - (float64(m.Timer) / float64(total))
	}

	// Status text with "Professional Sci-Fi" feel
	status := "SYSTEM STANDBY • READY TO FOCUS"
	if m.Active {
		if m.State == StateFocus {
			status = "WARP DRIVE ACTIVE • DEEP WORK"
		} else {
			status = "RECHARGING SYSTEMS • TAKE A BREAK"
		}
	} else if m.Timer < total && m.Timer > 0 {
		status = "STASIS MODE • TIMER PAUSED"
	}

	timerStr := fmt.Sprintf("%02d:%02d", int(m.Timer.Minutes()), int(m.Timer.Seconds())%60)

	taskTitle := "NO MISSION TARGETED"
	if m.Task != nil {
		taskTitle = strings.ToUpper(m.Task.Title)
	}

	// Centered UI dimensions
	cardW := min(60, m.width-4)

	// Build the Warp Track (linear progress)
	trackW := cardW - 10
	track := m.renderWarpTrack(trackW, progress, accent)

	// Styles
	titleStyle := lipgloss.NewStyle().
		Foreground(accent).
		Bold(true)

	timerStyle := lipgloss.NewStyle().
		Foreground(accent).
		Bold(true).
		Padding(0, 2).
		BorderStyle(lipgloss.DoubleBorder()).
		BorderForeground(accent)

	content := lipgloss.JoinVertical(lipgloss.Center,
		titleStyle.Render(status),
		"",
		track,
		"",
		timerStyle.Render(timerStr),
		"",
		m.styles.Muted.Render("CURRENT MISSION"),
		lipgloss.NewStyle().Foreground(m.styles.Theme.Fg).Bold(true).Render(utilTruncate(taskTitle, cardW-4)),
		"",
		m.styles.Muted.Render("enter: start/pause • r: abort • s: warp"),
	)

	// The "Liquid Glass" HUD
	hud := m.styles.Overlay.
		Width(cardW).
		Padding(2, 4).
		Align(lipgloss.Center, lipgloss.Center).
		Render(content)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, hud)
}

func (m Model) renderWarpTrack(w int, progress float64, color lipgloss.Color) string {
	if w < 10 {
		w = 10
	}

	// Rocket position
	pos := int(progress * float64(w-1))
	if pos < 0 {
		pos = 0
	}
	if pos >= w {
		pos = w - 1
	}

	// Build the string with the rocket
	var b strings.Builder
	mutedStyle := m.styles.Muted
	accentStyle := lipgloss.NewStyle().Foreground(color)

	// Left part (completed)
	b.WriteString(accentStyle.Render(strings.Repeat("━", pos)))
	// The Rocket
	b.WriteString(accentStyle.Bold(true).Render("🚀"))
	// Right part (remaining)
	if pos < w-1 {
		b.WriteString(mutedStyle.Render(strings.Repeat("─", w-1-pos)))
	}

	return b.String()
}

func utilTruncate(s string, w int) string {
	if w <= 0 {
		return ""
	}
	if lipgloss.Width(s) <= w {
		return s
	}
	if w <= 1 {
		return "…"
	}
	r := []rune(s)
	if len(r) <= w-1 {
		return string(r)
	}
	return string(r[:w-1]) + "…"
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
