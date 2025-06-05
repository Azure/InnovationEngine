# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.8] - 2025-06-05 - Commit: (Bug fix for overview suggestions)

### Fixed
- Fixed issue with "Apply Suggestions" button in overview help panel deleting existing content
- Improved suggestion handling to always append rather than replace content
- Enhanced user experience by preserving existing overview text when applying Copilot suggestions

## [0.2.7] - 2025-06-05 - Commit: (Help feature for overview editing)

### Added
- Added "Get Help" button to the overview editing interface
- Implemented a Copilot assistance panel for discussing and improving the overview
- Added functionality to apply Copilot suggestions directly to the overview content
- Enhanced the overview editing experience with AI assistance similar to step editing

## [0.2.6] - 2025-06-05 - Commit: (Streamlined UI workflow)

### Changed
- Streamlined the authoring process by eliminating the intermediate "Approve & Create Overview" step
- Renamed button to "Generate Steps" for direct step creation after overview editing
- Updated documentation and specifications to reflect the simplified three-phase workflow
- Improved user flow by merging overview creation and refinement phases

## [0.2.5] - 2025-06-04 - Commit: (Keyboard shortcut enhancements)

### Added
- Added CTRL+ENTER keyboard shortcut functionality to all text input areas:
  - Description textarea in ExecDocStepEditor now supports CTRL+ENTER to save changes
  - Code textarea in ExecDocStepEditor now supports CTRL+ENTER to save changes
  - Assistance prompt textarea in ExecDocStepEditor supports CTRL+ENTER to submit
  - Overview textarea and prompt input in OverviewAuthoring support CTRL+ENTER
- Improved user experience with keyboard shortcuts and better placeholders

## [0.2.4] - 2025-06-04 - Commit: (UI cleanup)

### Changed
- Removed "Step 4: Refine..." banner from step editor
- Removed "Phase Guidance" section from step content
- Simplified UI by reducing redundant phase information
- Improved focus on content editing rather than process guidance

## [0.2.3] - 2025-06-04 - Commit: (Complete step editing capabilities)

### Added
- Enhanced step editor to enable editing of all step properties, not just description
- Added ability to toggle between code block and non-code block content for any step
- Added controls to expand/collapse steps and reset execution status
- Improved step header with better visual indicators for executable steps
- Smart handling of Copilot assistance with context-aware suggestion application
- Option to apply suggestions to either description or code sections

## [0.2.2] - 2025-06-04 - Commit: (UI improvements)

### Changed
- Increased height of text editing areas throughout the application for better usability:
  - Overview editor now has a minimum height of 400px
  - Step description editor increased to 200px minimum height
  - Code editor increased to 250px minimum height
  - Assistance prompt textarea increased to 150px
- Simplified export options to only support Markdown format, removed HTML and PDF export options

## [0.2.1] - 2025-06-04 - Commit: (structured authoring process)

### Added
- Implemented structured four-phase authoring process:
  1. Create overview
  2. Manually edit or ask Copilot to make changes to the overview
  3. Write executable content to implement the described overview
  4. Manually edit or ask Copilot to make changes to the executable content
- Phase-specific UI guidance throughout the authoring process
- Visual progress indicators for the authoring workflow
- Phase-appropriate button labels and instructions
- Updated specification to reflect the structured authoring approach

## [0.2.0] - 2025-06-04 - Commits: e433589 (002 spec implementation), 3b17ce1 (Architect sidebar)

### Added
- Implementation of UI for iterative Executable Document authoring
- Two-panel interface for overview authoring with Copilot assistance
- Hierarchical document structure for Executive Document steps
- Collapsible step editor with rich text/Markdown support
- Step execution controls with live terminal output
- Kubernetes context and namespace selector
- File operations for saving and loading Executable Documents
- Accessibility improvements for keyboard navigation and screen readers
- Added "Architect" sidebar entry for quick access to the Exec Doc Editor

## [0.1.0] - 2025-06-04 - Commits: 903def7 (chat interface), d2c9223 (changelog)

### Added
- Initial Headlamp plugin setup
- "Getting Started" sidebar entry and page
- "Innovation Engine" sidebar entry and shell execution page
- "Assistant" sidebar entry and chat interface (903def7)
  - Multi-line text input for user queries
  - Chat history display with message styling
  - Processing state indicators
  - Quick start suggestion buttons
