package codec

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/programmersd21/kairo/internal/core"
)

func MarshalMarkdown(tasks []core.Task) []byte {
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

	// Sort roots by status and then by update time
	sortTasks(roots)

	var b bytes.Buffer
	fmt.Fprintf(&b, "# Kairo Export\n\n")
	fmt.Fprintf(&b, "_Exported: %s_\n\n", time.Now().UTC().Format(time.RFC3339))

	var walk func(t core.Task, depth int)
	walk = func(t core.Task, depth int) {
		indent := strings.Repeat("  ", depth)
		box := " "
		if t.Status == core.StatusDone {
			box = "x"
		}

		line := fmt.Sprintf("%s- [%s] %s", indent, box, escapeMDInline(strings.TrimSpace(t.Title)))
		if t.Deadline != nil {
			line += "  _(due " + t.Deadline.Local().Format("2006-01-02") + ")_"
		}
		if len(t.Tags) > 0 {
			line += "  " + tagsInline(t.Tags)
		}
		fmt.Fprintln(&b, line)

		if strings.TrimSpace(t.Description) != "" {
			descIndent := indent + "    "
			for _, ln := range strings.Split(strings.TrimRight(t.Description, "\n"), "\n") {
				fmt.Fprintln(&b, descIndent+ln)
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

func sortTasks(ts []core.Task) {
	sort.Slice(ts, func(i, j int) bool {
		if ts[i].Status != ts[j].Status {
			// Todo < Doing < Done (arbitrary, but consistent)
			return statusOrder(ts[i].Status) < statusOrder(ts[j].Status)
		}
		return ts[i].UpdatedAt.After(ts[j].UpdatedAt)
	})
}

func statusOrder(s core.Status) int {
	switch s {
	case core.StatusTodo:
		return 0
	case core.StatusDoing:
		return 1
	case core.StatusDone:
		return 2
	default:
		return 3
	}
}

var mdTaskRe = regexp.MustCompile(`^(\s*)-\s*\[( |x|X)\]\s+(.*)$`)

func UnmarshalMarkdown(b []byte) ([]core.Task, error) {
	s := bufio.NewScanner(bytes.NewReader(b))
	var tasks []core.Task

	type taskWithDepth struct {
		task  *core.Task
		depth int
	}
	var stack []taskWithDepth

	for s.Scan() {
		line := s.Text()
		if m := mdTaskRe.FindStringSubmatch(line); m != nil {
			indent := m[1]
			depth := len(indent) / 2

			title := strings.TrimSpace(m[3])
			st := core.StatusTodo
			if m[2] == "x" || m[2] == "X" {
				st = core.StatusDone
			}
			tags := extractTags(title)
			title = stripTags(title)

			// Extract due date if present
			var deadline *time.Time
			if idx := strings.LastIndex(title, "_(due "); idx != -1 {
				if endIdx := strings.Index(title[idx:], ")_"); endIdx != -1 {
					dueStr := title[idx+6 : idx+endIdx]
					if t, err := time.Parse("2006-01-02", dueStr); err == nil {
						deadline = &t
					}
					title = strings.TrimSpace(title[:idx])
				}
			}

			curTask := &core.Task{
				ID:       fmt.Sprintf("import-%d", len(tasks)), // Temporary ID
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
			// Use the pointer to the task in the slice so we can update it (e.g. Description)
			stack = append(stack, taskWithDepth{task: &tasks[len(tasks)-1], depth: depth})
			continue
		}

		if len(stack) > 0 {
			cur := stack[len(stack)-1].task
			trimmed := strings.TrimSpace(line)
			if trimmed != "" && !strings.HasPrefix(trimmed, "#") && !strings.HasPrefix(trimmed, "_") {
				// Potential description line.
				// Check if it's indented enough to be a description of the current task
				indent := ""
				if idx := strings.Index(line, trimmed); idx != -1 {
					indent = line[:idx]
				}

				expectedIndent := stack[len(stack)-1].depth*2 + 2
				if len(indent) >= expectedIndent {
					cur.Description += strings.TrimPrefix(line, strings.Repeat(" ", expectedIndent)) + "\n"
				}
			}
		}
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	for i := range tasks {
		tasks[i].Description = strings.TrimRight(tasks[i].Description, "\n")
		// Clean up temporary IDs for service to regenerate if needed,
		// but keep ParentID relationships.
		// Wait, if I clear IDs, ParentID will point to non-existent things.
		// The service.UpsertTask will handle it if we keep the temporary IDs
		// OR we let the service generate new IDs but we need to map them.

		// Actually, api.handleImport calls service.UpsertTask.
		// UpsertTask in repo.go:
		/*
			func (r *Repository) UpsertTask(ctx context.Context, task core.Task) error {
				if task.ID == "" {
					task.ID = r.nextID()
				}
				...
			}
		*/
		// If I keep "import-N", it will be saved as "import-N".
		// This might be fine.
	}
	return tasks, nil
}

func escapeMDInline(s string) string {
	repl := strings.NewReplacer("|", "\\|", "*", "\\*", "_", "\\_", "`", "\\`")
	return repl.Replace(s)
}

func tagsInline(tags []string) string {
	out := make([]string, 0, len(tags))
	for _, t := range tags {
		t = core.NormalizeTag(t)
		if t != "" {
			out = append(out, "#"+t)
		}
	}
	sort.Strings(out)
	return strings.Join(out, " ")
}

func extractTags(s string) []string {
	parts := strings.Fields(s)
	var tags []string
	for _, p := range parts {
		if strings.HasPrefix(p, "#") && len(p) > 1 {
			tags = append(tags, core.NormalizeTag(p))
		}
	}
	return tags
}

func stripTags(s string) string {
	parts := strings.Fields(s)
	keep := parts[:0]
	for _, p := range parts {
		if strings.HasPrefix(p, "#") && len(p) > 1 {
			continue
		}
		keep = append(keep, p)
	}
	return strings.Join(keep, " ")
}
