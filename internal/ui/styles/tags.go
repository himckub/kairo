package styles

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/programmersd21/kairo/internal/ui/theme"
)

type TagHighlight struct {
	Fg        string `toml:"fg"`
	Bg        string `toml:"bg"`
	Bold      bool   `toml:"bold"`
	Italic    bool   `toml:"italic"`
	Underline bool   `toml:"underline"`
}

func ResolveColor(s string, t theme.Theme) lipgloss.Color {
	switch s {
	case "fg":
		return t.Fg
	case "bg":
		return t.Bg
	case "muted":
		return t.Muted
	case "border":
		return t.Border
	case "accent":
		return t.Accent
	case "good":
		return t.Good
	case "warn":
		return t.Warn
	case "bad":
		return t.Bad
	case "overlay":
		return t.Overlay
	case "black":
		return lipgloss.Color("0")
	case "red":
		return lipgloss.Color("1")
	case "green":
		return lipgloss.Color("2")
	case "yellow":
		return lipgloss.Color("3")
	case "blue":
		return lipgloss.Color("4")
	case "magenta":
		return lipgloss.Color("5")
	case "cyan":
		return lipgloss.Color("6")
	case "white":
		return lipgloss.Color("7")
	case "bright_black":
		return lipgloss.Color("8")
	case "bright_red":
		return lipgloss.Color("9")
	case "bright_green":
		return lipgloss.Color("10")
	case "bright_yellow":
		return lipgloss.Color("11")
	case "bright_blue":
		return lipgloss.Color("12")
	case "bright_magenta":
		return lipgloss.Color("13")
	case "bright_cyan":
		return lipgloss.Color("14")
	case "bright_white":
		return lipgloss.Color("15")
	}
	return lipgloss.Color(s)
}

func ApplyTagHighlight(base lipgloss.Style, highlight TagHighlight, t theme.Theme) lipgloss.Style {
	style := base
	if highlight.Fg != "" {
		style = style.Foreground(ResolveColor(highlight.Fg, t))
	}
	if highlight.Bg != "" {
		style = style.Background(ResolveColor(highlight.Bg, t))
	}
	if highlight.Bold {
		style = style.Bold(true)
	}
	if highlight.Italic {
		style = style.Italic(true)
	}
	if highlight.Underline {
		style = style.Underline(true)
	}
	return style
}

func ParseTagHighlightString(s string) TagHighlight {
	highlight := TagHighlight{}
	parts := strings.Split(s, ",")
	for _, part := range parts {
		kv := strings.Split(part, "=")
		if len(kv) == 2 {
			switch kv[0] {
			case "fg":
				highlight.Fg = kv[1]
			case "bg":
				highlight.Bg = kv[1]
			}
		} else {
			switch part {
			case "bold":
				highlight.Bold = true
			case "italic":
				highlight.Italic = true
			case "underline":
				highlight.Underline = true
			}
		}
	}
	return highlight
}
