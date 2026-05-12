# Recurring Tasks

Kairo automates repetitive work through a robust recurring task system. When a recurring task is completed, Kairo automatically generates the next instance.

<img src={require('../../assets/recurring_tasks.png').default} alt="Recurring Tasks" />

## Recurrence Types

Kairo supports three types of recurrence:

| Type | Description |
|---|---|
| **None** | A standard, one-off task. |
| **Weekly** | Repeats on specific days of the week. |
| **Monthly** | Repeats on a specific day of the month. |

## Setting Up Recurrence

In the task editor (`e`), you can configure recurrence:

### Weekly Recurrence
1. Set **Recurrence** to `weekly`.
2. In the **Weekly Days** field, enter comma-separated days (e.g., `mon,wed,fri`).
3. Kairo will schedule the next instance on the earliest upcoming day in your list.

### Monthly Recurrence
1. Set **Recurrence** to `monthly`.
2. In the **Monthly Day** field, enter a number from `1` to `31`.
3. Kairo will schedule the task for that day every month.

## Smart Due Dates

When a recurring task is completed:
1. The current task is marked as `done`.
2. A new task is created with the same title, tags, priority, and project.
3. The **Due Date** of the new task is automatically calculated based on the recurrence schedule.

## Advanced Controls

### Wait Until
Use the `wait_until` field (in `config.toml` or via API) to hide a recurring task until it's actually relevant. New instances won't appear in your list until the `wait_until` time has passed.

### Until
The `until` field allows you to set an "expiration date" for the recurrence. After this date, Kairo will stop generating new instances of the task.

## Behavior Notes
- **Frictionless Completion**: Just press `z` on a recurring task. Kairo handles the logic instantly.
- **Manual Overrides**: You can always edit the due date of a generated instance without affecting the overall recurrence schedule.
