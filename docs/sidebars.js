/**
 * Creating a sidebar enables you to:
 - create an ordered group of docs
 - render a sidebar for each doc of that group
 - provide next/previous navigation

 The sidebars can be generated from the filesystem, or explicitly defined here.

 Create as many sidebars as you want.
 */

// @ts-check

/** @type {import('@docusaurus/plugin-content-docs').SidebarsConfig} */
const sidebars = {
  // By default, Docusaurus generates a sidebar from the docs folder structure
  tutorialSidebar: [
    'introduction',
    'installation',
    {
      type: 'category',
      label: 'Core Concepts',
      items: [
        'core-concepts/tasks',
        'core-concepts/projects',
        'core-concepts/recurring-tasks',
      ],
    },
    {
      type: 'category',
      label: 'Features',
      items: [
        'features/dashboard',
        'features/focus-engine',
        'features/ai-assistant',
        'features/search-shortcuts',
        'features/undo-redo',
        'features/markdown-preview',
      ],
    },
    {
      type: 'category',
      label: 'Advanced Usage',
      items: [
        'advanced/configuration',
        'advanced/sync',
        'advanced/import-export',
        'advanced/tag-highlighting',
      ],
    },
    {
      type: 'category',
      label: 'Extensibility',
      items: [
        'extensibility/lua-plugins',
        'extensibility/cli-api',
        'extensibility/mcp-server',
      ],
    },
    'contributing',
  ],
};

module.exports = sidebars;
