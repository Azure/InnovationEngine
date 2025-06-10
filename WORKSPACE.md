# Innovation Engine VS Code Workspace

The workspace configuration file (innovation-engine.code-workspace) helps with developing different components of the Innovation Engine project, including the Headlamp plugin.

## How to Use This Workspace

1. Open this file in VS Code:
   ```
   code innovation-engine.code-workspace
   ```

2. This will open VS Code with two workspace folders:
   - **Innovation Engine**: The main project
   - **Headlamp Plugin**: The Headlamp plugin in the sandbox directory

## Available Debug Configurations

- **Debug Headlamp Server**: Debug the Node.js backend of the Headlamp plugin
- **Debug Headlamp Client**: Debug the React frontend of the Headlamp plugin
- **Debug Headlamp Tests**: Run and debug tests for the Headlamp plugin
- **Debug Headlamp Full App**: Debug both the server and client simultaneously

For more details on debugging the Headlamp plugin, see the [DEBUG.md](./sandbox/innovation-engine-headlamp/DEBUG.md) file in the plugin directory.
