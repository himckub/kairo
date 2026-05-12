# Kairo Documentation

This folder contains the source code for the Kairo documentation website, built with [Docusaurus](https://docusaurus.io/).

## Structure

- `docs/`: The actual documentation Markdown files.
- `src/`: Custom React pages and global CSS.
- `static/`: Static assets (images, etc.).
- `docusaurus.config.js`: Site configuration.
- `sidebars.js`: Sidebar organization.

## How to use

### Local Development

1. Navigate to the `docs/` directory.
2. Install dependencies: `npm install`.
3. Start the development server: `npm start`.
4. Open `http://localhost:3000` in your browser.

### Visual Assets

The documentation uses PNG images stored in `docs/assets/`. These are referenced in the Markdown files using the `require()` syntax for optimal Docusaurus processing.

To update an image:
1. Replace the file in `docs/assets/` with a new version (keep the same filename).
2. Or, add a new image to `docs/assets/` and update the corresponding Markdown file.

### Deployment

A GitHub Action is already configured in `.github/workflows/deploy-docs.yml`. It will automatically:
1. Build the documentation.
2. Deploy it to your GitHub Pages site (typically `https://<username>.github.io/kairo/`) whenever you push changes to the `main` branch.

## Contributing

If you find errors or want to add new sections, please edit the files in this directory and submit a Pull Request.
