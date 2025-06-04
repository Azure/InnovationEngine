# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
