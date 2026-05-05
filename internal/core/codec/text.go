package codec

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/programmersd21/kairo/internal/core"
)

// MarshalText encodes tasks into a simple plain text format.
func MarshalText(tasks []core.Task) []byte {
	taskMap := make(map[string]core.Task)
	for _, t := range tasks {
		taskMap[t.ID] = t
	}

	children := make(map[string][]core.Task)
	for _, t := range tasks {
		if t.ParentID != "" {
			children[t.ParentID] = append(children[t.ParentID], t)
		}
	}

	var roots []core.Task
	for _, t := range tasks {
		if t.ParentID == "" {
			roots = append(roots, t)
		} else if _, ok := taskMap[t.ParentID]; !ok {
			roots = append(roots, t)
		}
	}

	// Sort roots
	sortTasks(roots)

	var b bytes.Buffer
	var walk func(t core.Task, depth int)
	walk = func(t core.Task, depth int) {
		indent := strings.Repeat("  ", depth)
		status := " "
		if t.Status == core.StatusDone {
			status = "x"
		}
		tags := ""
		if len(t.Tags) > 0 {
			tags = " " + tagsInline(t.Tags)
		}
		deadline := ""
		if t.Deadline != nil {
			deadline = " (due " + t.Deadline.Local().Format("2006-01-02") + ")"
		}
		fmt.Fprintf(&b, "%s[%s] %s%s%s\n", indent, status, t.Title, tags, deadline)
		if t.Description != "" {
			descIndent := indent + "    "
			for _, line := range strings.Split(strings.TrimRight(t.Description, "\n"), "\n") {
				fmt.Fprintf(&b, "%s%s\n", descIndent, line)
			}
		}

		kids := children[t.ID]
		sortTasks(kids)
		for _, kid := range kids {
			walk(kid, depth+1)
		}
	}

	for _, r := range roots {
		walk(r, 0)
	}
	return b.Bytes()
}

// Helper for sorting tasks (same as in markdown.go but internal to avoid cross-package issues if needed,
// but here it's the same package 'codec' so I should probably put it in a shared file in codec,
// but for simplicity I will just redefine or use it if available.)
// Since markdown.go and text.go are in same package, they can share sortTasks.
// Wait, I didn't export sortTasks in markdown.go. I should probably make it exported or put it in a shared file.
// I'll check if I can just use it. Yes, same package.

// UnmarshalText decodes tasks from a simple plain text format.
func UnmarshalText(b []byte) ([]core.Task, error) {
	s := bufio.NewScanner(bytes.NewReader(b))
	var tasks []core.Task

	type taskWithDepth struct {
		task  *core.Task
		depth int
	}
	var stack []taskWithDepth

	for s.Scan() {
		line := s.Text()
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "[") && len(trimmed) >= 4 && trimmed[2] == ']' {
			// Handle "[ ] Title" or "[x] Title"
			indent := ""
			if idx := strings.Index(line, "["); idx != -1 {
				indent = line[:idx]
			}
			depth := len(indent) / 2

			st := core.StatusTodo
			if trimmed[1] == 'x' || trimmed[1] == 'X' {
				st = core.StatusDone
			}
			titleAndTags := strings.TrimSpace(trimmed[4:])
			tags := extractTags(titleAndTags)
			title := stripTags(titleAndTags)

			// Extract due date if present
			var deadline *time.Time
			if idx := strings.LastIndex(title, " (due "); idx != -1 {
				if endIdx := strings.Index(title[idx:], ")"); endIdx != -1 {
					dueStr := title[idx+6 : idx+endIdx]
					if t, err := time.Parse("2006-01-02", dueStr); err == nil {
						deadline = &t
					}
					title = strings.TrimSpace(title[:idx])
				}
			}

			curTask := &core.Task{
				ID:       fmt.Sprintf("import-%d", len(tasks)),
				Title:    title,
				Tags:     tags,
				Status:   st,
				Deadline: deadline,
				Priority: core.P1,
			}

			// Find parent based on depth
			for len(stack) > 0 && stack[len(stack)-1].depth >= depth {
				stack = stack[:len(stack)-1]
			}
			if len(stack) > 0 {
				curTask.ParentID = stack[len(stack)-1].task.ID
			}

			tasks = append(tasks, *curTask)
			stack = append(stack, taskWithDepth{task: &tasks[len(tasks)-1], depth: depth})
		} else if len(stack) > 0 {
			cur := stack[len(stack)-1].task
			if trimmed != "" {
				indent := ""
				if idx := strings.Index(line, trimmed); idx != -1 {
					indent = line[:idx]
				}
				expectedIndent := stack[len(stack)-1].depth*2 + 4
				if len(indent) >= expectedIndent {
					cur.Description += strings.TrimPrefix(line, strings.Repeat(" ", expectedIndent)) + "\n"
				}
			}
		}
	}

	for i := range tasks {
		tasks[i].Description = strings.TrimRight(tasks[i].Description, "\n")
	}
	return tasks, nil
}
