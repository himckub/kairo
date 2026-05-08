package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/programmersd21/kairo/internal/core"
	"github.com/programmersd21/kairo/internal/history"
)

func (m *Model) deleteTasksCmd(ids []string) tea.Cmd {
	return func() tea.Msg {
		// Get the tasks before deleting them for history
		var before []core.Task
		for _, id := range ids {
			if task, err := m.svc.GetByID(m.ctx, id); err == nil {
				before = append(before, task)
			}
		}

		if err := m.svc.DeleteTasks(m.ctx, ids); err != nil {
			return errMsg{Err: err}
		}

		// Record deletion in history
		op := history.CreateOperation(history.OpBulkDelete, "", ids, before, nil)
		m.hist.Record(op)

		return taskUpdatedMsg{}
	}
}
