# Import & Export

Kairo believes in data portability. You can easily bring your tasks into Kairo or export them for use in other applications.

<img src={require('../../assets/export_menu.png').default} alt="Export Menu" />

## Supported Formats

| Format | Extension | Description |
|---|---|---|
| **JSON** | `.json` | Full fidelity. Preserves IDs, hierarchy, and metadata. |
| **CSV** | `.csv` | Flat list. Good for spreadsheets. |
| **Markdown** | `.md` | Human-readable. Great for meeting notes or project docs. |
| **Plain Text** | `.txt` | Minimalist list of titles. |

## Exporting Data

Press `x` from the main list view to open the **Import/Export Menu**.
1. Select **Export**.
2. Choose your desired format.
3. Kairo will generate the file and show you the save path.

## Importing Data

1. Press `x` and select **Import**.
2. Point Kairo to your file (JSON, CSV, MD, or TXT).
3. Kairo will parse the file and upsert the tasks into your database.

### Smart Hierarchy Import
When importing JSON or Markdown with nested structures, Kairo automatically preserves the parent-child relationships and restores the visual hierarchy in your task list.

## Automation via CLI

You can also export data directly from your terminal:
```bash
kairo export --format markdown > my_tasks.md
```
See the [CLI API Reference](../extensibility/cli-api.md) for more details.
