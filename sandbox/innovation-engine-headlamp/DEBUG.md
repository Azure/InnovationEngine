# Debug Configuration for Innovation Engine Headlamp

This document explains how to use the debug configurations for the Innovation Engine Headlamp plugin.

## Available Debug Configurations

### Using VS Code

1. **Debug Server** - Debug only the backend server (Node.js/Express)
2. **Debug Client** - Debug only the frontend application (React)
3. **Debug Full App** - Debug both the server and client simultaneously
4. **Debug Tests** - Debug test cases

## Opening the Correct Workspace

Since this plugin is a subfolder of the main Innovation Engine project, you have two options for opening it in VS Code:

### Option 1: Open as part of the parent project
1. Open the parent project: `code /home/rgardler/projects/InnovationEngine`
2. Use the workspace file: Open `innovation-engine.code-workspace` in the root directory
3. Debug configurations will be available under the "Debug Headlamp..." names

### Option 2: Open as a standalone project
1. Open the plugin directory directly: `code /home/rgardler/projects/InnovationEngine/sandbox/innovation-engine-headlamp`
2. For better configuration, open the workspace file: `headlamp-plugin.code-workspace`
3. Debug configurations will be available with normal names like "Debug Server"

## How to Debug

### Method 1: Using VS Code Debug Panel

1. Open the VS Code Debug panel (Ctrl+Shift+D or click on the Debug icon in the Activity Bar)
2. Select the configuration you want to use from the dropdown menu
3. Click the green Play button to start debugging

### Method 2: Using the Debug Script

Run the debug script from your terminal:

```bash
./start_debug.sh
```

This will:
1. Start the server with debugging enabled
2. Start the frontend application
3. Allow you to connect to the debug session from VS Code or Chrome DevTools

### Method 3: Using npm Scripts

```bash
# Debug mode (server with --inspect-brk flag)
npm run debug

# Just the server in debug mode
npm run server:debug
```

## Debugging Tips

1. **Breakpoints**: Click in the gutter next to line numbers to set breakpoints
2. **Watch Expressions**: Add variables to the Watch panel to monitor their values
3. **Call Stack**: View the call stack to understand the execution flow
4. **Variables**: Inspect local and global variables during debugging
5. **Console**: Use `console.log()` statements to output values (visible in the Debug Console)

## Environment Variables

When debugging, the following environment variables are set:
- `NODE_ENV=development`
- `PORT=4001`

Make sure your `.env` file contains any other required variables.
