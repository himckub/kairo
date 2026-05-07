package styles

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/programmersd21/kairo/internal/ui/theme"
)

func ParseTagStyle(s string, t theme.Theme) lipgloss.Style {
	style := lipgloss.NewStyle()
	parts := strings.Split(s, ",")
	for _, part := range parts {
		kv := strings.Split(part, "=")
		if len(kv) != 2 {
			if part == "bold" {
				style = style.Bold(true)
			}
			continue
		}
		key, val := kv[0], kv[1]
		var col lipgloss.Color
		if strings.HasPrefix(val, "ui-") {
			if c, ok := t.GetColor(val); ok {
				col = c
			} else {
				continue // Theme alias not found, fall back
			}
		} else {
			col = lipgloss.Color(val)
		}

		switch key {
		case "fg":
			style = style.Foreground(col)
		case "bg":
			style = style.Background(col)
		}
	}
	return style
}
