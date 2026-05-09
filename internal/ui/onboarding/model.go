package onboarding

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/programmersd21/kairo/internal/ui/keymap"
	"github.com/programmersd21/kairo/internal/ui/styles"
)

type Step int

const (
	StepWelcome Step = iota
	StepNavigation
	StepCreation
	StepCompletion
	StepFilter
	StepDashboard
	StepAI
	StepRecurring
	StepAdvanced
	StepFocus
	StepUndoRedo
	StepFinish
)

type CloseMsg struct {
	Skipped bool
}

type tickMsg time.Time

type Model struct {
	styles styles.Styles
	km     keymap.Keymap
	step   Step
	width  int
	height int
	frame  int
}

const (
	// Ensure there are NO leading spaces/tabs before the slashes on each line
	logo = `
    __         _            
   / /______ _(_)________   
  / //_/ __ ` + "`" + ` / / ___/ __\  
 / ,< / /_/ / / /  / /_/ /  
/_/|_|\__,_/_/_/   \____/   
`
)

func New(s styles.Styles, km keymap.Keymap) Model {
	return Model{
		styles: s,
		km:     km,
		step:   StepWelcome,
	}
}

func (m *Model) SetSize(w, h int) {
	m.width, m.height = w, h
}

func (m Model) Init() tea.Cmd {
	return m.tick()
}

func (m *Model) tick() tea.Cmd {
	return tea.Tick(time.Millisecond*150, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch x := msg.(type) {
	case tickMsg:
		m.frame++
		return m, m.tick()
	case tea.KeyMsg:
		if x.String() == "esc" {
			return m, func() tea.Msg { return CloseMsg{Skipped: true} }
		}

		switch m.step {
		case StepWelcome:
			if x.String() == "enter" {
				m.step = StepNavigation
			}
		case StepNavigation:
			if x.String() == "tab" || x.String() == "shift+tab" {
				m.step = StepCreation
			}
		case StepCreation:
			// In the real app, we'd wait for the task creation event,
			// but for the tutorial component itself, we just listen for the key 'n'.
			if x.String() == "n" {
				m.step = StepCompletion
			}
		case StepCompletion:
			if x.String() == "z" {
				m.step = StepFilter
			}
		case StepFilter:
			if x.Type == tea.KeyCtrlF || strings.EqualFold(x.String(), "ctrl+f") {
				m.step = StepDashboard
			}
		case StepDashboard:
			if x.String() == "s" {
				m.step = StepAI
			}
		case StepAI:
			// Some terminals/OSes don't reliably set KeyMsg.Type for Ctrl combos.
			// Accept both the normalized string and the control-key type.
			if x.Type == tea.KeyCtrlA || strings.EqualFold(x.String(), "ctrl+a") {
				m.step = StepRecurring
			}
		case StepRecurring:
			if x.String() == "enter" {
				m.step = StepAdvanced
			}
		case StepAdvanced:
			if x.String() == "ctrl+s" {
				m.step = StepFocus
			}
		case StepFocus:
			if x.String() == "f" {
				m.step = StepUndoRedo
			}
		case StepUndoRedo:
			if x.Type == tea.KeyCtrlZ || strings.EqualFold(x.String(), "ctrl+z") {
				m.step = StepFinish
			}
		case StepFinish:
			if x.String() == "enter" {
				return m, func() tea.Msg { return CloseMsg{Skipped: false} }
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.width <= 0 || m.height <= 0 {
		return ""
	}

	cardW := min(60, m.width-4)
	if cardW < 40 {
		cardW = m.width - 2
	}

	var title, body, action string

	colors := []lipgloss.Color{
		m.styles.Theme.Accent,
		m.styles.Theme.Good,
		m.styles.Theme.Warn,
		m.styles.Theme.Bad,
	}
	logoColor := colors[m.frame%len(colors)]
	animatedLogo := lipgloss.NewStyle().Foreground(logoColor).Bold(true).Render(logo)

	switch m.step {
	case StepWelcome:
		title = "WELCOME TO KAIRO"
		body = "A minimal, keyboard-first task manager designed for speed and focus."
		action = "Press [ENTER] to start the tour"
	case StepNavigation:
		title = "FAST NAVIGATION"
		body = "Kairo uses tabs for different views. It's the core of the workflow."
		action = "Press [TAB] to cycle through views"
	case StepCreation:
		title = "CREATE YOUR FIRST TASK"
		body = "Tasks are the heart of Kairo. Keep them crisp and actionable."
		action = "Press [N] to create a new task"
	case StepCompletion:
		title = "MARK AS DONE"
		body = "The best part of productivity is checking things off."
		action = "Press [Z] to complete a task"
	case StepFilter:
		title = "ORGANIZE BY TAG"
		body = "Keep your tasks categorized and find what you need instantly."
		action = "Press [CTRL+F] to filter by tag"
	case StepDashboard:
		title = "COMMAND CENTER"
		body = "Analyze your productivity, momentum, and focus trends in the Command Center."
		action = "Press [S] to open Dashboard"
	case StepAI:
		title = "AI ASSISTANT"
		body = "Leverage Gemini to manage tasks, get insights, or control the app via natural language."
		action = "Press [CTRL+A] to open AI panel"
	case StepRecurring:
		title = "RECURRING TASKS"
		body = "Automatically handle repeating work. Kairo handles the scheduling for you."
		action = "Press [ENTER] to continue"
	case StepAdvanced:
		title = "ADVANCED CONFIG"
		body = "Fine-tune Kairo's behavior, themes, and keybindings in the Settings menu."
		action = "Press [CTRL+S] for Settings"
	case StepFocus:
		title = "FOCUS ENGINE"
		body = "Bridge planning and execution with the integrated Focus Engine. Track deep work sessions directly against your tasks."
		action = "Press [F] to open Focus"
	case StepUndoRedo:
		title = "UNDO & REDO"
		body = "Made a mistake? Instantly reverse any action with Undo [CTRL+Z] or re-apply with Redo [CTRL+Y]."
		action = "Press [CTRL+Z] to continue"
	case StepFinish:
		title = "YOU'RE ALL SET"
		body = "You're now proficient with Kairo! Explore more via the palette [CTRL+P].\n\n(Tip: Relaunch tour anytime with [CTRL+W])"
		action = "Press [ENTER] to finish"
	}

	progress := m.renderProgress()

	content := lipgloss.JoinVertical(lipgloss.Center,
		animatedLogo,
		"",
		m.styles.Title.Render(" "+title+" "),
		"",
		lipgloss.NewStyle().Width(cardW-4).Align(lipgloss.Center).Render(body),
		"",
		progress,
		"",
		m.styles.Muted.Render(action),
		"",
		m.styles.Muted.Render("press [esc] to skip"),
	)

	card := lipgloss.NewStyle().
		Width(cardW).
		Border(lipgloss.ThickBorder()).
		BorderForeground(m.styles.Theme.Accent).
		Background(m.styles.Theme.Bg).
		Padding(2, 2).
		Align(lipgloss.Center).
		Render(content)

	return card
}

func (m Model) renderProgress() string {
	steps := 12
	var b strings.Builder
	for i := 0; i < steps; i++ {
		if i == int(m.step) {
			b.WriteString(lipgloss.NewStyle().Foreground(m.styles.Theme.Accent).Render("● "))
		} else if i < int(m.step) {
			b.WriteString(lipgloss.NewStyle().Foreground(m.styles.Theme.Good).Render("○ "))
		} else {
			b.WriteString(m.styles.Muted.Render("○ "))
		}
	}
	return b.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
