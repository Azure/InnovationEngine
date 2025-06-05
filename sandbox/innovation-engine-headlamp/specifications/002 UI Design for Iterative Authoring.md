# 002-Assistant Plugin UI Design for Iterative Exec Doc Authoring

**Status:** IMPLEMENTED

[View related GitHub Issue](https://github.com/SorraTheOrc/InnovationEngine/issues/4)

> **Note:** This implementation is based on the FINAL specification for UI design and enhances the existing Basic Assistant Experience (001).

## Introduction

**User Story:**
As a user of the Assistant Headlamp plugin, I want an intuitive user interface that allows me to iteratively author, edit, and validate Executable Documents (Exec Docs) with the help of GitHub Copilot, so that I can efficiently create, test, and manage step-by-step automation for Kubernetes workflows.

**Authoring Process:**
The Executable Document authoring process follows a streamlined three-step approach:
1. **Create and refine overview** - Generate and edit a high-level description of what the document will accomplish
2. **Write executable content to implement the described overview** - Generate step-by-step instructions with executable code
3. **Manually edit or ask Copilot to make changes** - Refine the executable content for accuracy and usability

## Requirements

### Functional Requirements
1. The UI must support a streamlined process for authoring an Executable Document:
   - Create and refine overview
   - Write executable content to implement the described overview
   - Manually edit or ask Copilot to make changes to the executable content
2. The UI must allow the user to author the overview of a new Exec Doc by submitting a short prompt to Copilot.
3. The user must be able to edit the overview either by direct text editing or by providing further instructions to Copilot for revision.
4. The user must be able to directly submit the overview for conversion into a full Executable Document.
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
- Keyboard shortcuts (CTRL+ENTER) must be implemented for common actions like submitting prompts or saving edits.
- Keyboard shortcuts must be clearly indicated in placeholder text and tooltips for better discoverability.
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
    - Overview preview and direct editing
  - Step 2: Executable Content Creation
    - Automatic step generation based on approved overview
    - Step-by-step implementation guidance
  - Step 3: Content Refinement
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
  - "Generate Steps" button to directly create steps from the overview

#### Exec Doc Editor UI
- **Document Structure**:
  - Hierarchical navigation sidebar showing document outline
  - Main content area with collapsible sections for each step
  - Step execution controls and status indicators

- **Step Editing**:
  - Full-height rich text/Markdown editor for each step (minimum 10-12 lines visible)
  - Complete editing capabilities for all step properties (not just description)
  - Toggle controls for code/non-code blocks and expanded/collapsed state
  - Ability to reset execution status for steps
  - Copilot assistance button to get help with specific steps
  - Syntax highlighting for code blocks
  - Context-aware suggestion application to either description or code
  - Preview toggle to see rendered content
  - Keyboard shortcuts (CTRL+ENTER) to save changes or submit prompts

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
  - Save as Markdown format only

- **Kubernetes Context Control**:
  - Context selector dropdown with current selection highlighted
  - Namespace selector dropdown with current selection highlighted
  - Visual indicator for active context throughout the UI
  - Permission level indicator (admin/non-admin)

### Implementation Notes

#### Structured Authoring Process Implementation
- **Phase Tracking**:
  - Implemented authoring phase state management in `ExecDocEditor.tsx` with the states: 'create-overview', 'implement-content', and 'refine-content'
  - Added phase transition functions: `moveToImplementContent()` and `moveToRefineContent()`
  - Connected phase transitions to appropriate user actions (overview creation and editing, step generation, step edits)

- **Phase-Aware UI**:
  - Added a visual phase indicator with progress tracking at the top of the editor
  - Implemented phase-specific guidance messages for each step of the process
  - Updated component props to pass current phase information between components
  - Modified button text and UI elements to reflect the current authoring phase, with direct "Generate Steps" option

- **Component Updates**:
  - Enhanced `OverviewAuthoring.tsx` to display different UI elements and instructions based on current phase
  - Simplified the overview creation process by merging overview creation and refinement into a single step
  - Renamed the main action button from "Approve & Create Overview" to "Generate Steps" to streamline the workflow
  - Extended `ExecDocStepEditor.tsx` to show phase-specific guidance for implementing or refining content
  - Implemented comprehensive step property editing in `ExecDocStepEditor.tsx` for all step attributes
  - Added intelligent suggestion handling with contextual application to appropriate step sections
  - Improved step execution UI with better visual indicators and controls
  - Updated existing UI components to ensure simplified phase transitions maintain proper state
  - Implemented keyboard shortcuts (CTRL+ENTER) across all text input areas:
    - Description and code textareas in `ExecDocStepEditor.tsx` support CTRL+ENTER to save changes
    - Assistance prompt textarea in `ExecDocStepEditor.tsx` uses CTRL+ENTER to submit requests
    - Overview textarea and prompt inputs in `OverviewAuthoring.tsx` handle CTRL+ENTER for submitting

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
- End-to-end flow following the streamlined process: create/refine overview, implement executable content, refine content, run steps, save/load
- Error handling for Copilot/API failures and file I/O
- Context/namespace switching and its effect on step execution
- Transition between authoring process phases with proper state preservation

### Acceptance Criteria
- A user can complete all three steps of the streamlined authoring process (create/refine overview, implement executable content, refine content) within the UI
- The structured authoring process provides clear guidance and transitions between phases
- All actions provide clear feedback and are accessible
- Kubernetes context and namespace can be selected and displayed, affecting step execution as expected

## References
- [Headlamp Plugin Development Docs](https://headlamp.dev/docs/latest/development/plugins/building)
- [GitHub Copilot Documentation](https://docs.github.com/en/copilot)
- [Innovation Engine Executable Document Format](/docs/specs/test-reporting.md)
- [Web Content Accessibility Guidelines (WCAG)](https://www.w3.org/WAI/standards-guidelines/wcag/)
