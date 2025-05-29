# innovation-engine

[Headlamp](https://headlamp.dev/) is an extensible Kubernetes UI that enables users to manage and visualize their clusters through a modern, customizable interface. It supports plugins to extend its capabilities with custom UIs and backend integrations.

The **innovation-engine** plugin integrates the Innovation Engine (IE) CLI directly into Headlamp. This allows users to run IE commands from within the Headlamp dashboard, providing a seamless workflow for interacting with Innovation Engine features alongside cluster management tasks. The plugin includes a secure backend API that only allows allowlisted commands (by default, the IE CLI), ensuring safe execution of shell commands from the browser UI.

----------------------------------------------------------------------------
----------------------------------------------------------------------------
DO NOT USE IN PRODUCTION - see note on security in the Developer Notes below
----------------------------------------------------------------------------
----------------------------------------------------------------------------

## Use

### Prerequisites

- Node.js (for the backend service)
- Go (for building the Innovation Engine CLI)
- Headlamp with this plugin installed

### Installing Backend Dependencies

Before starting the plugin, you must install the required dependencies:

```bash
pushd sandbox/innovation-engine-headlamp/
npm install
popd
```

If you see an error about a missing module such as 'express', ensure you have run the above command in the `sandbox/innovation-engine-headlamp` directory before starting the backend service.

### Building the Innovation Engine CLI

From the root of the Innovation Engine repository:

1. Build the CLI:
   ```bash
   make build-ie
   ```
   The resulting binary will be at `bin/ie` in the repository root.

### Running the Shell Exec Backend and the Headlamp Plugin

The backend service provides a secure API for executing allowlisted shell commands—by default, only the Innovation Engine CLI (`ie`)—from the Headlamp UI. This allows users to interact with Innovation Engine features directly in the browser, while ensuring that only approved commands can be run for security. By default, the backend listens on port 4000.

In a terminal, from anywhere in the repository, start the backend service:

```bash
pushd sandbox/innovation-engine-headlamp/src
node shell-exec-backend.js &
cd ..
npm run start &
popd
```

The backend will automatically discover the Innovation Engine binary at `../../bin/ie` (relative to the plugin directory) and allow it as a allowlisted command if present.

### Using the Plugin in Headlamp

1. Open Headlamp and ensure the plugin is enabled. If you see a "Getting Started" option in the left menu then the plugin is installed.
2. Click `Getting Started`
3. On the resulting page there are some quicklinks to test IE in headlamp and a link to a freeform command entry form.
4. The output, error (if any), and exit code will be displayed.

## Developer Notes

- The path to the test exec doc is currently hard coded (L36 of src/index.tsx), arbitrary paths can by typed in on the IE execution page
- The backend service is implemented in `shell-exec-backend.js` and exposes a simple REST API at `POST /api/exec` on port 4000 (by default).
- Only the Innovation Engine CLI (`ie`) can be executed in the shell, though any command within an Exec Docs will run (presenting a significant security risk since, at present, arbitrary documents can be run). 
- To allow more commands in the shell execution, add their full paths to the `ALLOWED_COMMANDS` array in `shell-exec-backend.js`.
- The plugin frontend (in `src/index.tsx`) uses React and Headlamp's plugin API to provide a UI for command execution and output display.
- See the [`specifications/`](./specifications/) folder for design specifications.

## Developing Headlamp plugins

For more information on developing Headlamp plugins, please refer to:

- [Getting Started](https://headlamp.dev/docs/latest/development/plugins/), How to create a new Headlamp plugin.
- [API Reference](https://headlamp.dev/docs/latest/development/api/), API documentation for what you can do
- [UI Component Storybook](https://headlamp.dev/docs/latest/development/frontend/#storybook), pre-existing components you can use when creating your plugin.
- [Plugin Examples](https://github.com/kubernetes-sigs/headlamp/tree/main/plugins/examples), Example plugins you can look at to see how it's done.
