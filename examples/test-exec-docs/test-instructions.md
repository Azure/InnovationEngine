# Automating Documentation Testing with Innovation Engine

## Overview

This guide explains how to set up and use a GitHub Action that automates testing of your Markdown documentation using Innovation Engine. The action ensures your executable documents are accurate and robust by running tests on them. It will create issues for any errors found and, when combined with branch protection rules, can prevent pull requests from merging until those errors are resolved.

## Prerequisites

Before setting up the GitHub Action, you need the following:

- **Azure Credentials**: The action logs into Azure using the Azure CLI, so you need to provide the following repository secrets:
  - `AZURE_CLIENT_ID`
  - `AZURE_TENANT_ID`
  - `AZURE_SUBSCRIPTION_ID`

  To obtain these values:
  1. Register an application in Azure Active Directory.
  2. Assign the necessary permissions to the application.
  3. Note the application (client) ID, tenant ID, and subscription ID.

- **GitHub Personal Access Token (PAT)**: The action interacts with the GitHub API to create issues and pull requests. Store your PAT in a repository secret named `PAT`.

## GitHub Action Workflow

The GitHub Action performs the following steps:

1. **Trigger**: Runs on push events affecting Markdown files or when manually triggered.

2. **Checkout Repository**: Uses `actions/checkout@v2` to clone the repository with full history for accurate diffs.

3. **Check for Changed Markdown Files**:

   - **Initial Commit**:
     - If this is the initial commit and the default branch is `main` or `master`, it compares changes with the default branch.
     - If the default branch is not `main` or `master`, it considers all Markdown files in the repository.
   - **Subsequent Commits**:
     - Compares changes between the latest commit and the previous one.
   - **Outputs**:
     - Sets an output variable indicating whether Markdown files have changed.
     - Lists the changed Markdown files to be tested.

4. **Azure CLI Login**: Logs into Azure using the provided credentials with `azure/login@v2` (only if Markdown files have changed).

5. **Set Up Python**: Sets up a Python 3.12 environment using `actions/setup-python@v2` (only if Markdown files have changed).

6. **Install Dependencies**: Installs required Python packages (only if Markdown files have changed):

   ```bash
   pip install --upgrade pip==24.3.1 || pip install --upgrade pip  # Try specific pip version, fallback to latest
   pip install PyGithub==2.1.1 pyyaml==5.4.1 || pip install PyGithub pyyaml  # Try specific package versions, fallback to latest
   ```

7. **Run Tests and Handle PRs**: Executes a Python script that (only if Markdown files have changed):

   - **Innovation Engine Installation**: Installs `ie` (Innovation Engine) if not already installed.
   - **Repository Information**: Retrieves repository details such as owner, name, and branch.
   - **Pull Request Handling**:
     - If on a branch other than `main`, it creates or updates a pull request to merge changes into `main`.
     - Checks for existing pull requests with the same title to avoid duplicates.
     - **CREATE_PR Environment Variable**:
       - The `CREATE_PR` variable is set to `true` if the event is a push, otherwise it uses the value from `workflow_dispatch` input or defaults to `false`.
   - **Testing Markdown Files**:
     - Iterates over the list of changed Markdown files.
     - Runs `ie execute` on each file to test the executable documentation.
     - **On Test Failure**:
       - Extracts error logs from the ie.log file.
       - Creates or updates a GitHub issue with the error details.
       - Records the failure in the test results.
       - Does not assign the issue to any user to avoid permission issues.
     - **On Test Success**:
       - Records the success in the test results.
   - **Reporting Results**:
     - Posts the test results as a comment on the pull request.
     - If no pull request exists (i.e., on the `main` branch), outputs the results in the workflow logs.

## Customization

You can modify the GitHub Action to suit your project's needs:

- **Event Triggers**: Adjust the `on` section to change when the action runs (e.g., on different branches or file types).

- **Azure Login Step**: If your documentation doesn't require Azure resources, you can remove the Azure login step.

- **Python Version and Dependencies**:
  - Update the `python-version` in the `Set Up Python` step to match your project's requirements.
  - In the `Install Dependencies` step, specify the versions of `pip`, `PyGithub`, and `pyyaml` as needed. The action tries to install specific versions first and falls back to the latest versions if the specified versions are not available.

- **Testing Script**: Update the Python script to change how tests are executed, how errors are handled, or how results are reported. For example, you might customize the error extraction logic or modify the issue creation process.

- **Secrets Management**: Ensure all necessary secrets (`AZURE_CLIENT_ID`, `AZURE_TENANT_ID`, `AZURE_SUBSCRIPTION_ID`, `PAT`) are correctly named and stored in your repository's settings under **Settings > Secrets and variables > Actions**.

## Conclusion

By integrating this GitHub Action into your repository, you automate the testing of your executable Markdown documents using Innovation Engine. This helps maintain documentation accuracy and, combined with branch protection rules, prevents broken documentation from being merged into your main branch.