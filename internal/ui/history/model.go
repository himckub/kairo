package history

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/programmersd21/kairo/internal/history"
	"github.com/programmersd21/kairo/internal/ui/keymap"
	"github.com/programmersd21/kairo/internal/ui/styles"
)

type item struct {
	op *history.Operation
}

func (i item) FilterValue() string {
	return i.op.Description
}

func (i item) Title() string {
	return i.op.Description
}

func (i item) Description() string {
	return fmt.Sprintf("[%s] %s", i.op.Timestamp.Format("Jan 02 15:04"), strings.Join(i.op.TaskIDs, ", "))
}

type ItemDelegate struct {
	styles styles.Styles
}

func (d ItemDelegate) Height() int                               { return 2 }
func (d ItemDelegate) Spacing() int                              { return 1 }
func (d ItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, li list.Item) {
	i, ok := li.(item)
	if !ok {
		return
	}

	var title, description string
	if index == m.Index() {
		title = d.styles.RowSelected.Render(i.Title())
		description = d.styles.RowSelected.Render(i.Description())
	} else {
		title = d.styles.RowNormal.Render(i.Title())
		description = d.styles.RowNormal.Render(i.Description())
	}

	_, _ = fmt.Fprintf(w, "%s\n%s", title, description)
}

type Model struct {
	list    list.Model
	km      keymap.Keymap
	styles  styles.Styles
	visible bool
}

func New(s styles.Styles, km keymap.Keymap) Model {
	l := list.New([]list.Item{}, ItemDelegate{styles: s}, 0, 0)
	l.Title = "History"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)

	k := keymap.GetHistoryListKeyMap()
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{k.Back, k.Confirm}
	}

	return Model{
		list:   l,
		km:     km,
		styles: s,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
func (m *Model) SetSize(width, height int) {
	m.list.SetSize(width, height)
}

func (m *Model) SetHistory(history []*history.Operation) {
	items := make([]list.Item, len(history))
	for i, op := range history {
		items[i] = item{op: op}
	}
	m.list.SetItems(items)
}

func (m *Model) View() string {
	if !m.visible {
		return ""
	}
	return m.styles.Overlay.Render(m.list.View())
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.visible {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.km.Back):
			m.visible = false
			return m, nil
		case key.Matches(msg, m.km.OpenTask):
			// TODO: Implement jump to specific undo point
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *Model) Show() {
	m.visible = true
}

func (m *Model) Hide() {
	m.visible = false
}

func (m *Model) Visible() bool {
	return m.visible
}
