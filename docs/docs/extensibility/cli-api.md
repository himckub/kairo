# CLI API

Kairo exposes a headless API via its CLI, making it perfect for shell scripts, cron jobs, and integration with other terminal tools.

<img src={require('../../assets/cli_api.png').default} alt="Kairo CLI API" />

## Base Command
All API calls are performed through the `kairo api` command.

```bash
kairo api [action] --payload 'json_payload'
```

## Actions & Payloads

### `create`
Create a new task.
```bash
kairo api create --payload '{"title": "Automated task", "priority": 1}'
```

### `list`
List tasks with filters.
```bash
kairo api list --payload '{"statuses": ["todo"], "tags": ["work"]}'
```

### `update`
Modify an existing task.
```bash
kairo api update --payload '{"id": "TASK_ID", "status": "done"}'
```

### `delete`
Remove a task.
```bash
kairo api delete --payload '{"id": "TASK_ID"}'
```

### `export`
Export tasks to various formats.
```bash
kairo api export --payload '{"format": "json"}'
```

## JSON Schema (TaskDTO)

When interacting with the API, tasks are represented by the following JSON structure:

```json
{
  "id": "string",
  "title": "string",
  "description": "string (optional)",
  "tags": ["string"],
  "priority": number (0-3),
  "status": "todo | doing | done",
  "deadline": "RFC3339 string (optional)",
  "project": "string",
  "parent_id": "string (optional)",
  "created_at": "RFC3339 string",
  "updated_at": "RFC3339 string"
}
```

## Real-world Examples

### Weekly Review Script
Automatically export your completed tasks for the week to a Markdown file:
```bash
kairo api export --payload '{"format": "markdown"}' > review.md
```

### Git Hook
Create a task in Kairo every time you push to a specific branch:
```bash
# In .git/hooks/pre-push
kairo api create --payload '{"title": "Review push to main", "tags": ["git"]}'
```
