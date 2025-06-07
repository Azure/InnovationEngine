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
- Azure OpenAI API credentials (for AI Assistant functionality)

### Azure AI Integration

The plugin now includes three AI-powered features that integrate with Azure OpenAI:

1. **Assistant** - Provides intelligent responses to questions about Kubernetes and Innovation Engine.
2. **Azure Architecture Overview Generator** - Creates comprehensive architectural overviews for Azure workloads and solutions.
3. **Executable Document Overview Generator** - Automatically creates detailed overviews for Kubernetes-focused executable documents.

To configure the Azure AI integration:

1. Copy the `env.example` file to `.env` in the project root:
   ```bash
   cp env.example .env
   ```

2. Edit `.env` and add your Azure OpenAI credentials:
   ```
   AZURE_OPENAI_API_KEY=your_api_key
   AZURE_OPENAI_ENDPOINT=https://your-resource.openai.azure.com
   AZURE_OPENAI_DEPLOYMENT_ID=your_deployment_name
   ```

3. Start the development server with both frontend and backend (includes environment validation):
   ```bash
   npm run dev
   ```

### Executable Document Overview Generator

The Executable Document Editor now includes an AI-powered overview generator that can automatically create detailed overviews for Kubernetes-focused operations. To use this feature:

1. Open the ExecDoc Editor component
2. In the "Create & Edit Document Overview" section, find the "Generate with Azure AI" panel
3. Enter a Kubernetes-related topic (e.g., "Deploy a stateful application with persistent storage")
4. Click "Generate Overview"
5. Review and edit the generated overview as needed
6. Click "Generate Steps" to continue building your executable document

The generated overview includes:
- A clear title for your document
- A concise overview of the process
- Prerequisites needed to follow the document
- Expected outcome after completing the steps
- Major steps involved in the process
- Additional considerations or warnings

#### Rate Limiting Considerations

The application includes built-in handling for Azure OpenAI API rate limits:

- Automatic retries with exponential backoff when rate limits are hit
- Respects the `Retry-After` header when provided by the Azure API
- Configurable maximum retry attempts to prevent excessive retries
- Detailed logging of rate limiting events for debugging
- Clear failure after maximum retries with descriptive error messages
- Tests that validate the rate limiting behavior

If you encounter rate limiting errors, you can modify the retry behavior in `/src/services/azureAI.ts` by adjusting:

- `maxRetries`: Maximum number of retry attempts for any error (default: 3)
- `retryDelay`: Base delay in milliseconds before retrying (default: 1000)
- `maxRateLimitRetries`: Maximum retries specifically for rate limiting errors (default: 2)

When tests run in CI mode (with `CI=true` environment variable), certain intensive rate limiting tests are automatically skipped to prevent unnecessary API calls.

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

### Running Tests

Before deploying or using the plugin, you should run the tests to ensure everything is configured correctly:

```bash
# Run environment validation tests first to check your setup
npx vitest run src/__tests__/a_environment.test.ts
```

The environment validation test will check if all required Azure OpenAI environment variables are properly set. If they are missing, you'll see a warning message with instructions on how to set them up.

**Important:** The validation test requires the `.env` file to be present and properly formatted. If you experience issues with environment validation failing despite having a `.env` file, ensure that:
1. The file is properly formatted with no extra spaces or special characters
2. The file is in the root directory of the project
3. The variables are named exactly as expected: `AZURE_OPENAI_API_KEY`, `AZURE_OPENAI_ENDPOINT`, and `AZURE_OPENAI_DEPLOYMENT_ID`

After confirming your environment is properly configured:

```bash
# Run all tests
npm test

# Run specific test suites
npm run test:rate-limits   # Run only the rate limiting tests
```

#### Rate Limiting Tests

The rate limiting tests have been specifically designed to:

1. Respect Azure OpenAI's rate limitations
2. Fail when there's a persistent rate limiting problem
3. Skip intensive tests in CI environments

To run only the rate limit tests:

```bash
npm run test:rate-limits
```

These tests verify that:
- The application correctly backs off when receiving 429 responses
- It respects the Retry-After header from Azure
- It fails with an appropriate error after exceeding maximum retry attempts
- It handles exponential backoff correctly

When running in CI environments, you can use these commands:

```bash
# Run all tests with CI flag to skip intensive tests
npm run test:ci

# Run only rate limit tests with CI flag to skip intensive ones
npm run test:ci-rate-limits
```

The `start` and `dev` commands now run environment validation by default with these behaviors:

- `npm run start`: Runs in **production** mode - validation will **fail and prevent execution** if Azure OpenAI credentials are missing
- `npm run dev`: Runs in **development** mode - validation will show warnings but continue even if credentials are missing

**Troubleshooting:**
If you're seeing validation failures with `npm run start` even though you have credentials in your `.env` file:

1. Make sure your `.env` file is properly formatted and contains the correct environment variables
2. Try running the validation test directly with: `NODE_ENV=production npx vitest run src/__tests__/a_environment.test.ts --reporter=verbose`
3. If the test passes directly but fails with npm scripts, it might be a shell environment issue - try modifying the scripts in package.json to use a different method of setting NODE_ENV

In both modes, the server itself will also check for credentials:
- In production mode: Returns 500 errors if credentials are missing
- In development mode: Returns fallback responses if credentials are missing

If you want to skip validation entirely (not recommended for production), use:

```bash
# Start the frontend without any environment validation
npm run start-without-validation

# Start both frontend and backend without any environment validation
npm run dev-without-validation
```

You can also run the validation separately:

```bash
# Check environment with production requirements (fails if not configured)
npm run validate-env:prod

# Check environment with development requirements (warns only)
npm run validate-env:dev
```

### Running the Shell Exec Backend and the Headlamp Plugin

The backend service provides a secure API for executing allowlisted shell commands—by default, only the Innovation Engine CLI (`ie`)—from the Headlamp UI. This allows users to interact with Innovation Engine features directly in the browser, while ensuring that only approved commands can be run for security. By default, the backend listens on port 4000.

In a terminal, from anywhere in the repository, start the backend service:

```bash
# Development mode - environment validation warnings only
pushd sandbox/innovation-engine-headlamp/
npm run dev
popd

# Production mode - will fail if environment not properly configured
pushd sandbox/innovation-engine-headlamp/
npm run start
popd

# Or start them separately:
pushd sandbox/innovation-engine-headlamp/
# Start the backend API server in development mode
npm run server:dev &
# Start the frontend development server (with validation)
npm run start &
popd

# Start the backend API server in production mode
npm run server:prod

# Skip validation entirely (not recommended for production):
npm run start-without-validation
npm run dev-without-validation
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
