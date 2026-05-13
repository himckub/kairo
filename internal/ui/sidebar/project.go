package sidebar

import (
	"fmt"
	"io"
	"sort"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/programmersd21/kairo/internal/ui/styles"
)

type Item struct {
	ID    string
	Title string
}

func (i Item) FilterValue() string { return i.Title }

type itemDelegate struct {
	activeProject string
	styles        styles.Styles
	isFocused     bool
}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}

	str := i.Title
	if i.ID == "" {
		str = "All Projects"
	}

	isActive := i.ID == d.activeProject
	isSelected := index == m.Index()

	textStyle := d.styles.Text
	if isSelected && d.isFocused {
		textStyle = textStyle.
			Foreground(d.styles.Theme.Bg).
			Background(d.styles.Theme.Accent).
			Padding(0, 1)
	} else if isSelected {
		textStyle = textStyle.
			Foreground(d.styles.Theme.Accent)
	} else if isActive {
		textStyle = textStyle.
			Foreground(d.styles.Theme.Accent).
			Bold(true)
	} else {
		textStyle = textStyle.
			Foreground(d.styles.Theme.Muted)
	}

	prefix := "  "
	if isActive && isSelected && d.isFocused {
		prefix = "• "
	} else if isActive {
		prefix = "• "
	} else if isSelected && d.isFocused {
		prefix = "> "
	}

	_, _ = fmt.Fprintf(w, "%s", textStyle.Render(prefix+str))
}

type Model struct {
	list          list.Model
	delegate      itemDelegate
	activeProject string
	styles        styles.Styles
	active        bool
	width         int
	height        int
}

func New(s styles.Styles) Model {
	delegate := itemDelegate{styles: s}
	l := list.New([]list.Item{}, delegate, 0, 0)
	l.SetShowTitle(true)
	l.Title = " PROJECTS"
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(s.Theme.Muted).
		Bold(true).
		Padding(0, 0, 1, 0)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)
	l.Styles.NoItems = lipgloss.NewStyle().MarginLeft(2)

	return Model{
		list:     l,
		delegate: delegate,
		styles:   s,
	}
}

func (m *Model) SetSize(w, h int) {
	m.width, m.height = w, h
	m.list.SetSize(w, h)
}

func (m *Model) SetActive(project string) {
	m.activeProject = project
	m.delegate.activeProject = project
	m.list.SetDelegate(m.delegate)

	// Try to select the item in the list too
	for i, item := range m.list.Items() {
		if it, ok := item.(Item); ok && it.ID == project {
			m.list.Select(i)
			break
		}
	}
}

func (m *Model) SetProjects(projects []string, order string, recent []string) {
	// Add "All Projects"
	sortable := make([]string, len(projects))
	copy(sortable, projects)

	if order == "recent" {
		recentMap := make(map[string]int)
		for i, p := range recent {
			recentMap[p] = i
		}
		sort.Slice(sortable, func(i, j int) bool {
			idxI, okI := recentMap[sortable[i]]
			idxJ, okJ := recentMap[sortable[j]]
			if okI && okJ {
				return idxI < idxJ
			}
			if okI {
				return true
			}
			if okJ {
				return false
			}
			return sortable[i] < sortable[j]
		})
	} else {
		sort.Strings(sortable)
	}

	finalProjects := append([]string{""}, sortable...)
	items := make([]list.Item, len(finalProjects))
	for i, p := range finalProjects {
		items[i] = Item{ID: p, Title: p}
	}
	m.list.SetItems(items)
	m.SetActive(m.activeProject) // Restore selection after items change
}

func (m *Model) Focus(active bool) {
	m.active = active
}

type SelectMsg struct {
	Project string
}

func (m *Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.delegate.isFocused != m.active {
		m.delegate.isFocused = m.active
		m.list.SetDelegate(m.delegate)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}

		if m.active {
			switch msg.String() {
			case "enter":
				if i, ok := m.list.SelectedItem().(Item); ok {
					return *m, func() tea.Msg {
						return SelectMsg{Project: i.ID}
					}
				}
			}
		}
	}

	m.list, cmd = m.list.Update(msg)
	return *m, cmd
}

func (m *Model) View() string {
	return m.list.View()
}
