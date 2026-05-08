package history

import (
	"fmt"
	"sync"
	"time"

	"github.com/programmersd21/kairo/internal/core"
)

// OperationType defines the type of operation recorded in history
type OperationType int

const (
	OpCreate OperationType = iota
	OpUpdate
	OpDelete
	OpBulkDelete
	OpToggleStatus
	OpBulkToggleStatus
	OpChangePriority
	OpChangeDeadline
	OpChangeCollapsed
	OpBulkOperation
)

func (op OperationType) String() string {
	switch op {
	case OpCreate:
		return "Create"
	case OpUpdate:
		return "Update"
	case OpDelete:
		return "Delete"
	case OpBulkDelete:
		return "Delete Multiple"
	case OpToggleStatus:
		return "Toggle Status"
	case OpBulkToggleStatus:
		return "Toggle Status (Multiple)"
	case OpChangePriority:
		return "Change Priority"
	case OpChangeDeadline:
		return "Change Deadline"
	case OpChangeCollapsed:
		return "Toggle Collapse"
	case OpBulkOperation:
		return "Bulk Operation"
	default:
		return "Unknown"
	}
}

// Operation represents a single undoable operation
type Operation struct {
	Type        OperationType `json:"type"`
	Description string        `json:"description"`
	Timestamp   time.Time     `json:"timestamp"`

	// Affected task IDs
	TaskIDs []string `json:"task_ids"`

	// Snapshots of tasks before and after the operation
	// Before: original state before operation
	// After: state after operation (for redo)
	Before []core.Task `json:"before"`
	After  []core.Task `json:"after"`
}

// History manages undo/redo stacks
type History struct {
	mu        sync.Mutex
	undoStack []*Operation
	redoStack []*Operation
	maxSize   int // Maximum operations to keep in history
}

// New creates a new history manager
func New(maxSize int) *History {
	if maxSize <= 0 {
		maxSize = 100 // Default to 100 operations
	}
	return &History{
		undoStack: make([]*Operation, 0, maxSize),
		redoStack: make([]*Operation, 0, maxSize),
		maxSize:   maxSize,
	}
}

// Record adds an operation to the undo stack and clears the redo stack
func (h *History) Record(op *Operation) {
	if op == nil {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	// Trim undo stack if it exceeds max size
	if len(h.undoStack) >= h.maxSize {
		// Remove oldest item
		h.undoStack = h.undoStack[1:]
	}

	h.undoStack = append(h.undoStack, op)

	// Clear redo stack when a new operation is recorded
	h.redoStack = h.redoStack[:0]
}

// Undo returns the operation to undo and moves it to the redo stack
func (h *History) Undo() *Operation {
	h.mu.Lock()
	defer h.mu.Unlock()

	if len(h.undoStack) == 0 {
		return nil
	}

	// Pop from undo stack
	op := h.undoStack[len(h.undoStack)-1]
	h.undoStack = h.undoStack[:len(h.undoStack)-1]

	// Push to redo stack
	h.redoStack = append(h.redoStack, op)

	return op
}

// Redo returns the operation to redo and moves it back to the undo stack
func (h *History) Redo() *Operation {
	h.mu.Lock()
	defer h.mu.Unlock()

	if len(h.redoStack) == 0 {
		return nil
	}

	// Pop from redo stack
	op := h.redoStack[len(h.redoStack)-1]
	h.redoStack = h.redoStack[:len(h.redoStack)-1]

	// Push to undo stack
	h.undoStack = append(h.undoStack, op)

	return op
}

// CanUndo returns true if undo is available
func (h *History) CanUndo() bool {
	return len(h.undoStack) > 0
}

// CanRedo returns true if redo is available
func (h *History) CanRedo() bool {
	return len(h.redoStack) > 0
}

// GetUndoStack returns a copy of the undo stack for display (newest first)
func (h *History) GetUndoStack() []*Operation {
	result := make([]*Operation, len(h.undoStack))
	copy(result, h.undoStack)
	// Reverse so newest is first
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return result
}

// Clear resets all history
func (h *History) Clear() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.undoStack = h.undoStack[:0]
	h.redoStack = h.redoStack[:0]
}

// Len returns the current sizes of undo and redo stacks
func (h *History) Len() (int, int) {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.undoStack), len(h.redoStack)
}

// UndoSize returns the number of undo operations available
func (h *History) UndoSize() int {
	return len(h.undoStack)
}

// RedoSize returns the number of redo operations available
func (h *History) RedoSize() int {
	return len(h.redoStack)
}

// CreateOperation is a helper to create a new operation record
func CreateOperation(opType OperationType, description string, taskIDs []string, before, after []core.Task) *Operation {
	return &Operation{
		Type:        opType,
		Description: description,
		Timestamp:   time.Now(),
		TaskIDs:     taskIDs,
		Before:      before,
		After:       after,
	}
}

// GetOperationDescription returns a human-readable description of an operation
func GetOperationDescription(op *Operation) string {
	if op == nil {
		return ""
	}

	if op.Description != "" {
		return op.Description
	}

	// Fallback to generated description
	switch op.Type {
	case OpCreate:
		if len(op.After) == 1 {
			return fmt.Sprintf("Created task: \"%s\"", truncateTitle(op.After[0].Title, 50))
		}
		return fmt.Sprintf("Created %d tasks", len(op.After))

	case OpDelete:
		if len(op.Before) == 1 {
			return fmt.Sprintf("Deleted task: \"%s\"", truncateTitle(op.Before[0].Title, 50))
		}
		return fmt.Sprintf("Deleted %d tasks", len(op.Before))

	case OpBulkDelete:
		return fmt.Sprintf("Deleted %d tasks", len(op.Before))

	case OpUpdate:
		if len(op.After) == 1 {
			return fmt.Sprintf("Edited task: \"%s\"", truncateTitle(op.After[0].Title, 50))
		}
		return fmt.Sprintf("Updated %d tasks", len(op.After))

	case OpToggleStatus:
		if len(op.After) == 1 {
			status := "completed"
			if op.After[0].Status != core.StatusDone {
				status = "reopened"
			}
			return fmt.Sprintf("%s task: \"%s\"", capitalizeFirst(status), truncateTitle(op.After[0].Title, 40))
		}
		return fmt.Sprintf("Toggled status for %d tasks", len(op.After))

	case OpBulkToggleStatus:
		return fmt.Sprintf("Toggled status for %d tasks", len(op.After))

	case OpChangePriority:
		if len(op.After) == 1 && len(op.Before) == 1 {
			return fmt.Sprintf("Changed priority: P%d → P%d", int(op.Before[0].Priority), int(op.After[0].Priority))
		}
		return fmt.Sprintf("Changed priority for %d tasks", len(op.After))

	case OpChangeDeadline:
		if len(op.After) == 1 {
			if op.After[0].Deadline.IsZero() {
				return fmt.Sprintf("Removed deadline from: \"%s\"", truncateTitle(op.After[0].Title, 40))
			}
			return fmt.Sprintf("Set deadline to %s", op.After[0].Deadline.Format("Jan 02 2006"))
		}
		return fmt.Sprintf("Changed deadline for %d tasks", len(op.After))

	case OpChangeCollapsed:
		if len(op.After) == 1 {
			if op.After[0].Collapsed {
				return fmt.Sprintf("Collapsed: \"%s\"", truncateTitle(op.After[0].Title, 40))
			}
			return fmt.Sprintf("Expanded: \"%s\"", truncateTitle(op.After[0].Title, 40))
		}

	default:
		return fmt.Sprintf("Modified %d task(s)", len(op.TaskIDs))
	}

	return op.Description
}

func truncateTitle(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]-'a'+'A') + s[1:]
}
