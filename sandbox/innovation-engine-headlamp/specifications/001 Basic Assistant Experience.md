# Assistant Plugin: Basic User Experience Specification

**Status:** In Progress

## Introduction

**User Story:**
As a new Headlamp and Kubernetes user, I want to use a Headlamp plugin that presents a user interface within Headlamp to author Executable Documents with the help of GitHub Copilot, so that I can build and manage my application on a local Kubernetes cluster without leaving the Headlamp environment.

This specification defines the minimum user experience for the Assistant Headlamp plugin. The output of this work is a Headlamp plugin that enables new users to interact with GitHub Copilot for authoring Executable Documents and managing local Kubernetes clusters, all from within the Headlamp UI.

## Requirements

### Functional Requirements
- The output must be a Headlamp plugin that integrates directly into the Headlamp UI.
- The plugin must provide a user interface within Headlamp for interacting with GitHub Copilot.
- Users must be able to enter natural language queries or requests (e.g., "Create a deployment for my app").
- The plugin must enable users to author Executable Documents with Copilot's assistance, suitable for use in Innovation Engine.
- The plugin must display Copilot's responses, including code snippets, explanations, and step-by-step instructions.

### Non-Functional Requirements
- The UI must be simple, clear, and accessible to Kubernetes beginners.
- All interactions must occur within the Headlamp UI (no external browser windows or command shells required for core flows).
- The plugin must not require cluster admin privileges for basic usage.
- The plugin must not store user queries or responses beyond the current session.

## Design

### Architecture
- The plugin will be a frontend-only Headlamp plugin (initially), using the Headlamp plugin API and React.
- The UI will consist of:
  - An "Assistant" sidebar entry.
  - A main panel with:
    - A text input for user queries.
    - A display area for Copilot responses.
    - Action buttons for copying code, running suggested commands, or opening documentation.
    - Quick start links for common tasks.

### Components & Interfaces
- Sidebar entry: "Assistant"
- Main panel: React component(s) for chat, response display, and quick actions

### Implementation Details

#### Chat Interface
- **Layout**: Chat-based interface with clear distinction between user and assistant messages
  - User messages: Right-aligned with blue background
  - Assistant messages: Left-aligned with white background and subtle border
  - Auto-scrolling behavior to always show the latest messages

- **Input Area**:
  - Multi-line textarea for flexible query input
  - Supports both short questions and longer requests
  - Shift+Enter for line breaks, Enter to submit
  - Placeholder text provides example queries (e.g., "Create a deployment for my app")
  - Send button with visual feedback for processing state

- **Chat History Display**:
  - Scrollable container with fixed height (400px)
  - Messages styled as conversation bubbles
  - Clear visual distinction between user and assistant messages
  - Loading indicator when processing a request

- **Quick Start Suggestions**:
  - Predefined common queries for one-click selection
  - Initial set includes: 
    - "Create a simple deployment"
    - "Expose a service"
    - "Author an Executable Document"
    - "Learn about Kubernetes basics"
  - Clicking a suggestion populates the query input field

- **State Management**:
  - Maintains chat history during the session
  - Tracks processing state to prevent duplicate submissions
  - Provides visual feedback during response generation

#### Response Handling (Planned)
- Integration with GitHub Copilot API for intelligent responses
- Support for rendering formatted code snippets in responses
- Copy-to-clipboard functionality for code blocks
- Action buttons within responses for executing commands directly

## Testing

### Unit Tests
- UI components render and update correctly
- User input is accepted and displayed
- Copilot responses are rendered as expected

### Integration Tests
- End-to-end flow: user enters a query, receives a response, and can copy/use the result
- Quick start actions trigger the correct flows

### Acceptance Criteria
- A new user can open the Assistant panel, ask a question, and receive a helpful response from Copilot
- The user can copy code or commands from the response
- The user can access quick start actions for common Kubernetes tasks

## References
- [Headlamp Plugin Development Docs](https://headlamp.dev/docs/latest/development/plugins/building)
- [GitHub Copilot Documentation](https://docs.github.com/en/copilot)
- [Kubernetes Getting Started Guide](https://kubernetes.io/docs/tutorials/)
- [Minikube Documentation](https://minikube.sigs.k8s.io/docs/)
