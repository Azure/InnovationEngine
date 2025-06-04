# 002-Assistant Plugin UI Design for Iterative Exec Doc Authoring

**Status:** IMPLEMENTED

[View related GitHub Issue](https://github.com/SorraTheOrc/InnovationEngine/issues/4)

> **Note:** This implementation is based on the FINAL specification for UI design and enhances the existing Basic Assistant Experience (001).

## Introduction

**User Story:**
As a user of the Assistant Headlamp plugin, I want an intuitive user interface that allows me to iteratively author, edit, and validate Executable Documents (Exec Docs) with the help of GitHub Copilot, so that I can efficiently create, test, and manage step-by-step automation for Kubernetes workflows.

**Authoring Process:**
The Executable Document authoring process follows a structured four-step approach:
1. **Create overview** - Generate a high-level description of what the document will accomplish
2. **Manually edit or ask Copilot to make changes** - Refine the overview to ensure clarity and completeness
3. **Write executable content to implement the described overview** - Generate step-by-step instructions with executable code
4. **Manually edit or ask Copilot to make changes** - Refine the executable content for accuracy and usability

## Requirements

### Functional Requirements
1. The UI must support a structured process for authoring an Executable Document:
   - Create overview
   - Manually edit or ask Copilot to make changes to the overview
   - Write executable content to implement the described overview
   - Manually edit or ask Copilot to make changes to the executable content
2. The UI must allow the user to author the overview of a new Exec Doc by submitting a short prompt to Copilot.
3. The user must be able to edit the overview either by direct text editing or by providing further instructions to Copilot for revision.
4. The user must be able to approve and submit the overview for conversion into a full Executable Document.
5. The UI must display the Exec Doc as a set of collapsible areas, each representing a step in the documented process.
6. The user must be able to edit the content of each step, either directly in the UI or by interacting with Copilot for suggestions or rewrites.
7. The user must be able to run each step of the Exec Doc individually from the UI to validate its correctness and behavior.
8. The user must be able to save Exec Docs to disk and load them from disk for further editing or execution.
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

#### Doc Authoring Process UI
- **Structured Authoring Flow**:
  - Step 1: Create Overview
    - Initial prompt input for generating document outline
    - Overview preview and approval
  - Step 2: Overview Refinement
    - Direct editing interface or Copilot-assisted refinement
    - Version tracking of major edits
  - Step 3: Executable Content Creation
    - Automatic step generation based on approved overview
    - Step-by-step implementation guidance
  - Step 4: Content Refinement
    - Direct editing interface or Copilot-assisted refinement for each step
    - Validation through step execution

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

### Implementation Notes

#### Structured Authoring Process Implementation
- **Phase Tracking**:
  - Implemented authoring phase state management in `ExecDocEditor.tsx` with the states: 'create-overview', 'refine-overview', 'implement-content', and 'refine-content'
  - Added phase transition functions: `moveToRefineOverview()`, `moveToImplementContent()`, and `moveToRefineContent()`
  - Connected phase transitions to appropriate user actions (overview creation, step generation, step edits)

- **Phase-Aware UI**:
  - Added a visual phase indicator with progress tracking at the top of the editor
  - Implemented phase-specific guidance messages for each step of the process
  - Updated component props to pass current phase information between components
  - Modified button text and UI elements to reflect the current authoring phase

- **Component Updates**:
  - Enhanced `OverviewAuthoring.tsx` to display different UI elements and instructions based on current phase
  - Extended `ExecDocStepEditor.tsx` to show phase-specific guidance for implementing or refining content
  - Updated existing UI components to ensure phase transitions maintain proper state

- **Version Management**:
  - Updated version in `package.json` from 0.2.0 to 0.2.1
  - Added new version entry in `CHANGELOG.md` detailing the structured authoring process implementation

## Testing

### Unit Tests
- Components render and update as expected
- Editing and Copilot interactions work for overview and steps
- File operations trigger correct UI feedback
- Context/namespace selection and display functions correctly

### Integration Tests
- End-to-end flow following the structured process: create overview, refine overview, implement executable content, refine content, run steps, save/load
- Error handling for Copilot/API failures and file I/O
- Context/namespace switching and its effect on step execution
- Transition between authoring process phases with proper state preservation

### Acceptance Criteria
- A user can complete all four steps of the authoring process (create overview, refine overview, implement executable content, refine content) within the UI
- The structured authoring process provides clear guidance and transitions between phases
- All actions provide clear feedback and are accessible
- Kubernetes context and namespace can be selected and displayed, affecting step execution as expected

## References
- [Headlamp Plugin Development Docs](https://headlamp.dev/docs/latest/development/plugins/building)
- [GitHub Copilot Documentation](https://docs.github.com/en/copilot)
- [Innovation Engine Executable Document Format](/docs/specs/test-reporting.md)
- [Web Content Accessibility Guidelines (WCAG)](https://www.w3.org/WAI/standards-guidelines/wcag/)
