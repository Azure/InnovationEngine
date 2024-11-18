# Automating Documentation Testing with Innovation Engine

## Overview

This guide explains how to set up and use a GitHub Action that automates testing of your Markdown documentation using Innovation Engine. The action ensures your executable documents are accurate and robust by running tests on them. It will create comments on pull requests for any errors found and, when combined with branch protection rules, can prevent pull requests from merging until those errors are resolved.

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

- **GitHub Personal Access Token (PAT)**: The action interacts with the GitHub API to create comments on pull requests. Store your PAT in a repository secret named `PAT`.

## GitHub Action Workflow

The GitHub Action performs the following steps:

1. **Trigger**: Runs on pull request events when a pull request is opened or updated.

2. **Checkout Repository**: Uses `actions/checkout@v3` to clone the repository with full history for accurate diffs.

3. **Determine Commit SHAs**:
   - Identifies the base and head commits of the pull request.
   - Calculates the number of commits between the base and head.
   - Sets output variables for use in subsequent steps.

4. **Check for Changed Markdown Files**:

   - Compares changes between the base and head commits.
   - Generates a list of changed Markdown files to be tested.
   - If no Markdown files have changed, the workflow exits early.

5. **Azure CLI Login**: Logs into Azure using the provided credentials with `azure/login@v2` (only if Markdown files have changed).

6. **Set Up Python**: Sets up a Python 3.12 environment using `actions/setup-python@v4` (only if Markdown files have changed).

7. **Install Dependencies**: Installs required Python packages (only if Markdown files have changed):

   ```bash
   pip install --upgrade pip==24.3.1 
   pip install setuptools wheel
   pip install PyGithub==2.1.1 pyyaml==6.0.2 
   ```

8. **Run Tests and Handle Pull Request**: Executes a Python script that (only if Markdown files have changed):

   - **Innovation Engine Installation**: Installs `ie` (Innovation Engine) if not already installed.
   - **Repository Information**: Retrieves repository details such as owner, name, and pull request number.
   - **Testing Markdown Files**:
     - Iterates over the list of changed Markdown files.
     - Runs `ie execute` on each file to test the executable documentation.
     - **On Test Failure**:
       - Extracts error logs from the `ie.log` file.
       - Creates a comment on the pull request with the error details.
       - Records the failure in the test results.
     - **On Test Success**:
       - Records the success in the test results.
   - **Reporting Results**:
     - Posts the test results as a comment on the pull request.
     - If any failures occurred, exits with a non-zero status to fail the workflow.

## Customization

You can modify the GitHub Action to suit your project's needs:

- **Event Triggers**: Adjust the `on` section to change when the action runs. Currently, it triggers on pull request events (`opened` and `synchronize`).

- **Azure Login Step**: If your documentation doesn't require Azure resources, you can remove the Azure login step.

- **Python Version and Dependencies**:
  - Update the `python-version` in the **Set Up Python** step to match your project's requirements.
  - In the **Install Dependencies** step, specify the versions of `pip`, `PyGithub`, and `pyyaml` as needed.

- **Testing Script**: Update the Python script to change how tests are executed, how errors are handled, or how results are reported. For example, you might customize the error extraction logic or modify the comment creation process.

- **Secrets Management**: Ensure all necessary secrets (`AZURE_CLIENT_ID`, `AZURE_TENANT_ID`, `AZURE_SUBSCRIPTION_ID`, `PAT`) are correctly named and stored in your repository's settings under **Settings > Secrets and variables > Actions**.

## Conclusion

By integrating this GitHub Action into your repository, you automate the testing of your executable Markdown documents using Innovation Engine. This helps maintain documentation accuracy and, combined with branch protection rules, prevents broken documentation from being merged into your main branch.