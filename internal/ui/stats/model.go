package stats

import (
	"fmt"
	"math"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/programmersd21/kairo/internal/stats"
	"github.com/programmersd21/kairo/internal/ui/styles"
)

type Model struct {
	styles styles.Styles
	width  int
	height int

	data stats.DashboardData

	activeSection int // 0: DNA, 1: Timeline, 2: Momentum, 3: Insights, 4: Tags

	animationOffset int
}

func New(s styles.Styles) Model {
	return Model{
		styles: s,
	}
}

func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
}

func (m *Model) SetData(data stats.DashboardData) {
	m.data = data
}

func (m Model) Init() tea.Cmd {
	return m.tickCmd()
}

func (m Model) tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

type tickMsg struct{}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.animationOffset++
		return m, m.tickCmd()
	case tea.KeyMsg:
		switch msg.String() {
		case "h", "left":
			m.activeSection--
			if m.activeSection < 0 {
				m.activeSection = 4
			}
		case "l", "right":
			m.activeSection = (m.activeSection + 1) % 5
		case "1", "2", "3", "4", "5":
			m.activeSection = int(msg.String()[0] - '1')
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	// Calculate panel widths
	innerW := m.width - 12
	if innerW < 40 {
		innerW = 40
	}
	if innerW > 100 {
		innerW = 100
	}
	halfW := (innerW / 2) - 2

	// DNA & Momentum row
	col1 := m.renderDNA(halfW)
	col2 := m.renderMomentum(halfW)
	topRow := lipgloss.JoinHorizontal(lipgloss.Top, col1, "  ", col2)

	// Timeline
	timeline := m.renderTimeline(innerW)

	// Insights & Tags row
	insights := m.renderInsights(halfW)
	tags := m.renderTags(halfW)
	bottomRow := lipgloss.JoinHorizontal(lipgloss.Top, insights, "  ", tags)

	// Build main container
	dashboard := lipgloss.JoinVertical(lipgloss.Left,
		topRow,
		m.styles.Divider.Width(innerW).MarginTop(1).MarginBottom(1).Render(""),
		timeline,
		m.styles.Divider.Width(innerW).MarginTop(1).MarginBottom(1).Render(""),
		bottomRow,
	)

	// Wrap in a premium border
	containerStyle := m.styles.Overlay.
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(m.styles.Theme.Accent).
		Padding(1, 2).
		Width(innerW + 4)

	title := m.styles.Title.
		Background(m.styles.Theme.Accent).
		Foreground(m.styles.Theme.Bg).
		Padding(0, 2).
		Render(" COMMAND CENTER ")

	// Position the title on the top border - JoinVertical adds no gap
	// Adjust container to pull up into title area by 1 row
	finalView := lipgloss.JoinVertical(lipgloss.Center,
		title,
		containerStyle.MarginTop(-1).Render(dashboard),
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, finalView)
}

func (m Model) renderDNA(w int) string {
	var b strings.Builder
	b.WriteString(m.styles.Accent.Bold(true).Render("🧬 PRODUCTIVITY DNA") + "\n\n")

	// Visualizing Peak Hours
	for i, intensity := range m.data.DNA.PeakHours {
		char := "░"
		color := m.styles.Theme.Muted

		if intensity > 0.8 {
			char = "█"
			color = m.styles.Theme.Accent
		} else if intensity > 0.5 {
			char = "▓"
			color = m.styles.Theme.Good
		} else if intensity > 0.2 {
			char = "▒"
			color = m.styles.Theme.Warn
		}

		if i == time.Now().Hour() {
			char = "┃"
			color = m.styles.Theme.Bad
		}

		b.WriteString(lipgloss.NewStyle().Foreground(color).Render(char))
	}

	b.WriteString("\n" + m.styles.Muted.Render("00h . . . 06h . . . 12h . . . 18h . . . 23h") + "\n\n")

	fmt.Fprintf(&b, "%-12s %s\n", m.styles.Muted.Render("VELOCITY:"), m.styles.Bold.Foreground(m.styles.Theme.Good).Render(fmt.Sprintf("%.2f tasks/hr", m.data.DNA.Velocity)))
	fmt.Fprintf(&b, "%-12s %s\n", m.styles.Muted.Render("CONSISTENCY:"), m.styles.Bold.Foreground(m.styles.Theme.Accent).Render(fmt.Sprintf("%.0f%%", m.data.DNA.Consistency*100)))

	return lipgloss.NewStyle().Width(w).Foreground(m.styles.Theme.Fg).Render(b.String())
}

func (m Model) renderMomentum(w int) string {
	var b strings.Builder
	b.WriteString(m.styles.Accent.Bold(true).Render("⚡ MOMENTUM ENGINE") + "\n\n")

	// Momentum gauge
	gaugeW := w - 8
	if gaugeW < 10 {
		gaugeW = 10
	}
	filled := int(float64(gaugeW) * (m.data.Momentum / 100.0))

	// Animated breathing pulse
	offset := float64(m.animationOffset) * 0.15
	pulse := (math.Sin(offset) + 1) / 2

	style := lipgloss.NewStyle().Foreground(m.styles.Theme.Good)
	if m.data.Momentum < 30 {
		style = style.Foreground(m.styles.Theme.Bad)
	} else if m.data.Momentum < 70 {
		style = style.Foreground(m.styles.Theme.Warn)
	}

	if pulse > 0.6 {
		style = style.Bold(true).Foreground(m.styles.Theme.Accent)
	}

	bar := style.Render(strings.Repeat("━", filled)) + m.styles.Muted.Render(strings.Repeat("─", gaugeW-filled))
	b.WriteString(bar + fmt.Sprintf(" %s\n\n", m.styles.Bold.Render(fmt.Sprintf("%.0f%%", m.data.Momentum))))

	fmt.Fprintf(&b, "%-12s %s\n", m.styles.Muted.Render("STREAK:"), m.styles.Bold.Foreground(m.styles.Theme.Good).Render(fmt.Sprintf("%d days", m.data.Streaks.Current)))
	fmt.Fprintf(&b, "%-12s %s\n", m.styles.Muted.Render("BEST:"), m.styles.Muted.Render(fmt.Sprintf("%d days", m.data.Streaks.Longest)))

	return lipgloss.NewStyle().Width(w).Foreground(m.styles.Theme.Fg).Render(b.String())
}

func (m Model) renderTimeline(w int) string {
	var b strings.Builder
	b.WriteString(m.styles.Accent.Bold(true).Render("📈 ACTIVITY TIMELINE") + "\n\n")

	max := 0
	for _, p := range m.data.Timeline {
		if p.Completed > max {
			max = p.Completed
		}
	}

	height := 6
	for h := height; h > 0; h-- {
		// Only show y-axis if we have data
		label := ""
		if max > 0 {
			label = fmt.Sprintf("%2d ", h*max/height)
		} else {
			label = "   "
		}
		b.WriteString(m.styles.Muted.Render(label))

		for _, p := range m.data.Timeline {
			level := 0
			if max > 0 {
				level = (p.Completed * height) / max
			}

			if level >= h {
				char := "▄"
				color := m.styles.Theme.Accent
				if h == level {
					char = "█"
				}
				b.WriteString(lipgloss.NewStyle().Foreground(color).Render(char + " "))
			} else {
				b.WriteString(m.styles.Muted.Render("░ "))
			}
		}
		b.WriteString("\n")
	}

	b.WriteString("   ")
	for _, p := range m.data.Timeline {
		style := m.styles.Muted
		if p.Date.Day() == time.Now().Day() {
			style = m.styles.Bold.Foreground(m.styles.Theme.Accent)
		}
		b.WriteString(style.Render(p.Date.Format("02") + " "))
	}

	return lipgloss.NewStyle().Width(w).Foreground(m.styles.Theme.Fg).Render(b.String())
}

func (m Model) renderInsights(w int) string {
	var b strings.Builder
	b.WriteString(m.styles.Accent.Bold(true).Render("🧠 BEHAVIORAL INSIGHTS") + "\n\n")

	for _, in := range m.data.Insights {
		icon := "•"
		color := m.styles.Theme.Muted
		switch in.Type {
		case "positive":
			icon = "󰄲"
			color = m.styles.Theme.Good
		case "negative":
			icon = "󰅚"
			color = m.styles.Theme.Bad
		}

		b.WriteString(lipgloss.NewStyle().Foreground(color).Render(icon) + " ")
		b.WriteString(m.styles.Bold.Render(in.Title) + "\n")
		b.WriteString("  " + m.styles.Muted.Render(in.Value) + "\n\n")
	}

	return lipgloss.NewStyle().Width(w).Foreground(m.styles.Theme.Fg).Render(b.String())
}

func (m Model) renderTags(w int) string {
	var b strings.Builder
	b.WriteString(m.styles.Accent.Bold(true).Render("🏷  TAG INTELLIGENCE") + "\n\n")

	for _, tc := range m.data.TagClusters {
		barW := int(tc.Score)
		if barW > 12 {
			barW = 12
		}
		label := m.styles.Text.Width(12).Render(tc.Tag)
		bar := m.styles.Accent.Render(strings.Repeat("■", barW))
		fmt.Fprintf(&b, "%s %s %s\n", label, bar, m.styles.Muted.Render(fmt.Sprintf("(%d)", tc.Count)))
	}

	return lipgloss.NewStyle().Width(w).Foreground(m.styles.Theme.Fg).Render(b.String())
}
