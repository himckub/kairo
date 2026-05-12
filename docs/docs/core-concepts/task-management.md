# Tasks

The fundamental unit of organization in Kairo is the **Task**. Every task is designed to capture the essence of what needs to be done with minimal friction.

<img src={require('../../assets/tasks_list.png').default} alt="Managing Tasks" />

## Task Fields

A task in Kairo contains the following fields:

| Field | Description |
|---|---|
| **Title** | A short, descriptive name for the task. |
| **Description** | Detailed notes, supporting Markdown. |
| **Status** | `todo`, `doing`, or `done`. |
| **Priority** | `0` (Critical), `1` (High), `2` (Medium), `3` (Low). |
| **Tags** | Comma-separated labels for categorization. |
| **Due Date** | An optional deadline (supports natural language). |
| **Project** | The project this task belongs to. |
| **Parent** | For nested tasks, the ID of the parent task. |

## Creating Tasks

Press `n` from the main list view to open the task editor. 

- Use `Tab` / `Shift+Tab` to navigate between fields.
- Press `ctrl+s` to save.
- Press `Esc` to cancel.

### Natural Language Deadlines

Kairo features a powerful NLP engine for deadlines. You can type:
- `tomorrow 10am`
- `next friday`
- `in 2 days`
- `end of month`
- `mon 3pm`

## Status Management

You can quickly change task status from the list view:
- Press `z` to mark a task as **Done**.
- Tasks in `doing` state are highlighted to indicate active focus.

## Bulk Actions

Select multiple tasks using `Space` and perform actions on all of them at once (e.g., complete, delete).
