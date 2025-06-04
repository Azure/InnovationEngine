# 002-Assistant Plugin UI Design for Iterative Exec Doc Authoring

**Status:** IMPLEMENTED

[View related GitHub Issue](https://github.com/SorraTheOrc/InnovationEngine/issues/4)

> **Note:** This implementation is based on the FINAL specification for UI design and enhances the existing Basic Assistant Experience (001).

## Introduction

**User Story:**
As a user of the Assistant Headlamp plugin, I want an intuitive user interface that allows me to iteratively author, edit, and validate Executable Documents (Exec Docs) with the help of GitHub Copilot, so that I can efficiently create, test, and manage step-by-step automation for Kubernetes workflows.

## Requirements

### Functional Requirements
1. The UI must allow the user to author the overview of a new Exec Doc by submitting a short prompt to Copilot.
2. The user must be able to edit the overview either by direct text editing or by providing further instructions to Copilot for revision.
3. The user must be able to approve and submit the overview for conversion into a full Executable Document.
4. The UI must display the Exec Doc as a set of collapsible areas, each representing a step in the documented process.
5. The user must be able to edit the content of each step, either directly in the UI or by interacting with Copilot for suggestions or rewrites.
6. The user must be able to run each step of the Exec Doc individually from the UI to validate its correctness and behavior.
7. The user must be able to save Exec Docs to disk and load them from disk for further editing or execution.
8. The UI must allow users to select or switch Kubernetes contexts and namespaces for step execution.
9. The current context and namespace should be clearly displayed and easily changeable from within the plugin UI.
10. Step execution should respect the selected context/namespace, and users should receive feedback if a step is run in a different context than expected.
11. The plugin should provide warnings or require confirmation if a step attempts to perform operations outside the current context or namespace.

### Non-Functional Requirements
- The UI must be responsive and accessible, following Headlamp and general web accessibility guidelines.
- All editing and execution actions must be intuitive and require minimal user training.
- The UI must provide clear feedback for Copilot interactions, execution results, and file operations.
- The plugin must not require cluster admin privileges for authoring or editing Exec Docs.
- User data must not be stored beyond the current session unless explicitly saved to disk.

### Accessibility Requirements
- The UI must be fully navigable via keyboard, with logical tab order and visible focus indicators.
- All interactive elements (buttons, inputs, panels) must have accessible labels and ARIA attributes as appropriate.
- The plugin must support screen readers, providing meaningful descriptions for all controls and dynamic content.
- High-contrast mode and sufficient color contrast must be ensured for all UI elements.
- All feedback (success, error, Copilot responses, execution results) must be accessible to assistive technologies.
- Accessibility must be considered in all user flows, including editing, running, saving, and loading Exec Docs.
- Accessibility testing must be included in the acceptance criteria and test plans.

## Design

### Architecture
- The UI will be implemented as a set of React components within the Headlamp plugin framework.
- State management will be handled locally in the browser, with optional persistence to disk.
- Copilot integration will be via API calls, with results rendered in context.

### Components & Interfaces
- **Overview Authoring Panel:**
  - Prompt input for Copilot
  - Editable overview text area
  - Approve/submit button
- **Exec Doc Steps View:**
  - Collapsible panels for each step
  - Direct edit and Copilot-assist options for each step
  - Run button for each step with output/result display
- **File Operations:**
  - Save to disk
  - Load from disk
  - Feedback for success/failure
- **Kubernetes Context/Namespace Selector:**
  - Dropdowns or input fields for context and namespace
  - Display of current context/namespace
  - Feedback for context/namespace change and step execution results

### Implementation Details

#### Overview Authoring UI
- **Layout**: Two-panel approach with prompt panel and preview panel
  - Left panel: Prompt input and conversation with Copilot
  - Right panel: Live preview of the Executable Document overview
  - Toggle button to expand either panel to full width

- **Prompt Panel**:
  - Goal description field (e.g., "Create a deployment for a Node.js application")
  - Submit button with visual feedback for processing state
  - Conversation history with Copilot regarding the overview
  - Revision request field to refine generated overview

- **Preview Panel**:
  - Formatted Markdown preview of the overview
  - Edit button to switch to raw Markdown editing
  - Approve button to finalize overview and proceed to step generation

#### Exec Doc Editor UI
- **Document Structure**:
  - Hierarchical navigation sidebar showing document outline
  - Main content area with collapsible sections for each step
  - Step execution controls and status indicators

- **Step Editing**:
  - Rich text/Markdown editor for each step
  - Copilot assistance button to get help with specific steps
  - Syntax highlighting for code blocks
  - Preview toggle to see rendered content

- **Step Execution**:
  - Play button for each step
  - Live terminal output display
  - Success/failure indicators
  - Context/namespace warning badges

#### File Operations UI
- **Save/Load Controls**:
  - File browser dialog for selecting save location or file to load
  - Recent files list for quick access
  - Auto-save toggle with interval selection
  - Export options (Markdown, HTML, PDF)

- **Kubernetes Context Control**:
  - Context selector dropdown with current selection highlighted
  - Namespace selector dropdown with current selection highlighted
  - Visual indicator for active context throughout the UI
  - Permission level indicator (admin/non-admin)

## Testing

### Unit Tests
- Components render and update as expected
- Editing and Copilot interactions work for overview and steps
- File operations trigger correct UI feedback
- Context/namespace selection and display functions correctly

### Integration Tests
- End-to-end flow: author overview, edit, approve, generate steps, edit steps, run steps, save/load
- Error handling for Copilot/API failures and file I/O
- Context/namespace switching and its effect on step execution

### Acceptance Criteria
- A user can author, edit, and approve an overview, generate and edit steps, run steps, and save/load Exec Docs, all through the UI
- All actions provide clear feedback and are accessible
- Kubernetes context and namespace can be selected and displayed, affecting step execution as expected

## References
- [Headlamp Plugin Development Docs](https://headlamp.dev/docs/latest/development/plugins/building)
- [GitHub Copilot Documentation](https://docs.github.com/en/copilot)
- [Innovation Engine Executable Document Format](/docs/specs/test-reporting.md)
- [Web Content Accessibility Guidelines (WCAG)](https://www.w3.org/WAI/standards-guidelines/wcag/)
