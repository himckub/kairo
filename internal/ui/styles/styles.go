package styles

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/programmersd21/kairo/internal/core"
	"github.com/programmersd21/kairo/internal/ui/theme"
)

// Icons used throughout the app.
// Designed with a "Premium & Sentimental" aesthetic for modern terminals.
const (
	IconTodo      = "󰄱 "
	IconDoing     = "󰔟 "
	IconDone      = "󰄲 "
	IconPriority0 = "󰼎 "
	IconPriority1 = "󰼏 "
	IconPriority2 = "󰼐 "
	IconPriority3 = "󰼑 "
	IconDeadline  = "󰃰 "
	IconWaitUntil = "󰥔 "
	IconUntil     = "󰦞 "
	IconTag       = "󰓹 "
	IconSync      = "󰑓 "
	IconError     = "󰅚 "
	IconSuccess   = "󰄲 "
	IconInfo      = "󰋽 "
	IconHelp      = "󰋗 "
	IconTask      = "󰈈 "
	IconPlugin    = "󰡀 "
	// UI Affordances (Safe Unicode for cross-platform terminal compatibility)
	IconPalette   = "⌘ "
	IconNew       = "+ "
	IconDelete    = "× "
	IconView      = "◎ "
	IconStrike    = "✓ "
	IconIssues    = "! "
	IconChangelog = "≡ "
	IconBack      = "« "
	IconEdit      = "✎ "
	IconClose     = "× "
	IconUp        = "↑ "
	IconDown      = "↓ "
	IconEnter     = "↵ "
	IconDiscuss   = "󰭹 "
)

// Design System Constants
const (
	// Spacing (refined, airy grid)
	Spacing0 = 0
	Spacing1 = 1
	Spacing2 = 2
	Spacing4 = 4
)

type Styles struct {
	Theme theme.Theme

	// Base
	App    lipgloss.Style
	Header lipgloss.Style
	Footer lipgloss.Style
	Panel  lipgloss.Style

	// Typography
	Title    lipgloss.Style
	Subtitle lipgloss.Style
	Bold     lipgloss.Style
	Muted    lipgloss.Style
	Text     lipgloss.Style
	Accent   lipgloss.Style

	// Tabs & Navigation
	TabActive   lipgloss.Style
	TabInactive lipgloss.Style
	Separator   lipgloss.Style

	// Rows & List Items
	RowSelected lipgloss.Style
	RowNormal   lipgloss.Style
	RowHovered  lipgloss.Style
	RowDimmed   lipgloss.Style

	// Badges & Status
	Badge       lipgloss.Style
	BadgeGood   lipgloss.Style
	BadgeDoing  lipgloss.Style
	BadgeWarn   lipgloss.Style
	BadgeBad    lipgloss.Style
	BadgeMuted  lipgloss.Style
	BadgeDelete lipgloss.Style
	BadgeQuit   lipgloss.Style

	// Detail & Form
	DetailKey   lipgloss.Style
	DetailValue lipgloss.Style
	DetailLabel lipgloss.Style
	Tag         lipgloss.Style
	TagLeft     lipgloss.Style
	TagRight    lipgloss.Style
	FormLabel   lipgloss.Style

	// Components
	Card        lipgloss.Style
	Input       lipgloss.Style
	InputActive lipgloss.Style
	Overlay     lipgloss.Style
	Divider     lipgloss.Style

	// States
	Empty   lipgloss.Style
	Loading lipgloss.Style
	Error   lipgloss.Style
	Success lipgloss.Style
}

func New(t theme.Theme) Styles {
	base := lipgloss.NewStyle().Foreground(t.Fg).Background(t.Bg)
	accent := lipgloss.NewStyle().Foreground(t.Accent).Background(t.Bg)
	muted := lipgloss.NewStyle().Foreground(t.Muted).Background(t.Bg)

	// Minimalist accent bar for selection
	selected := lipgloss.NewStyle().
		Background(t.Overlay).
		Foreground(t.Accent).
		Bold(true).
		PaddingLeft(Spacing2)

	return Styles{
		Theme: t,

		App:    base,
		Header: base.Padding(0, Spacing2),
		Footer: base.Padding(0, Spacing2),
		Panel:  base,

		Title:    base.Bold(true).Foreground(t.Accent),
		Subtitle: base.Bold(true).Foreground(t.Muted),
		Bold:     base.Bold(true),
		Muted:    muted,
		Text:     base,
		Accent:   accent,

		TabActive:   accent.Bold(true).Padding(0, Spacing1),
		TabInactive: muted.Padding(0, Spacing1),
		Separator:   muted.SetString("│"),

		RowSelected: selected,
		RowNormal:   base.PaddingLeft(Spacing2),
		RowHovered:  base.PaddingLeft(Spacing2).Foreground(t.Accent),
		RowDimmed:   muted.PaddingLeft(Spacing2),

		Badge:       muted,
		BadgeGood:   lipgloss.NewStyle().Foreground(t.Bg).Background(t.Good).Padding(0, Spacing1),
		BadgeDoing:  lipgloss.NewStyle().Foreground(t.Bg).Background(t.Accent).Padding(0, Spacing1),
		BadgeWarn:   lipgloss.NewStyle().Foreground(t.Bg).Background(t.Warn).Padding(0, Spacing1),
		BadgeBad:    lipgloss.NewStyle().Foreground(t.Bg).Background(t.Bad).Padding(0, Spacing1),
		BadgeMuted:  muted,
		BadgeDelete: lipgloss.NewStyle().Foreground(t.Bg).Background(t.Bad).Padding(0, Spacing1),
		BadgeQuit:   lipgloss.NewStyle().Foreground(t.Bg).Background(t.Warn).Padding(0, Spacing1),

		DetailKey:   muted.Width(12).MarginRight(Spacing1),
		DetailValue: base,
		DetailLabel: muted.Bold(true),
		Tag:         accent.Padding(0, Spacing1),
		TagLeft:     accent.SetString(""),
		TagRight:    accent.SetString(""),
		FormLabel:   muted.Bold(true),

		Card:        base.Padding(Spacing1, Spacing2),
		Input:       base.Padding(0, Spacing1).Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(t.Muted),
		InputActive: base.Padding(0, Spacing1).Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(t.Accent),
		Overlay:     base.BorderStyle(lipgloss.RoundedBorder()).BorderForeground(t.Accent).Padding(Spacing1, Spacing2),
		Divider:     muted.SetString("─"),

		Empty:   muted.Italic(true),
		Loading: accent,
		Error:   lipgloss.NewStyle().Foreground(t.Bad),
		Success: lipgloss.NewStyle().Foreground(t.Good),
	}
}

func (s Styles) StatusBadge(st core.Status) string {
	var style lipgloss.Style
	var text string

	switch st {
	case core.StatusTodo:
		style = s.BadgeMuted.Background(s.Theme.Muted).Foreground(s.Theme.Bg)
		text = IconTodo + "TODO"
	case core.StatusDoing:
		style = s.BadgeDoing
		text = IconDoing + "DOING"
	case core.StatusDone:
		style = s.BadgeGood
		text = IconDone + "DONE"
	default:
		style = s.BadgeMuted
		text = string(st)
	}

	return style.Render(text)
}

func (s Styles) PriorityBadge(p core.Priority) string {
	var style lipgloss.Style
	var text string

	switch p.Clamp() {
	case core.P0:
		style = s.BadgeMuted.Background(s.Theme.Muted).Foreground(s.Theme.Bg)
		text = IconPriority0 + "P0"
	case core.P1:
		style = s.BadgeMuted.Background(s.Theme.Muted).Foreground(s.Theme.Bg)
		text = IconPriority1 + "P1"
	case core.P2:
		style = s.BadgeWarn
		text = IconPriority2 + "P2"
	case core.P3:
		style = s.BadgeBad
		text = IconPriority3 + "P3"
	default:
		style = s.BadgeMuted
		text = fmt.Sprintf("P%d", int(p))
	}

	pill := lipgloss.JoinHorizontal(lipgloss.Left,
		s.TagLeft.Foreground(style.GetBackground()).Render(),
		style.Padding(0, 0).Render(text),
		s.TagRight.Foreground(style.GetBackground()).Render(),
	)

	return pill
}
