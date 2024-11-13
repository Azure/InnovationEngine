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

1. **Trigger**: Runs on push or pull request events affecting Markdown files, or when manually triggered.

2. **Checkout Repository**: Uses `actions/checkout@v2` to clone the repository.

3. **Azure CLI Login**: Logs into Azure using the provided credentials with `azure/login@v2`.

4. **Set Up Python**: Sets up a Python 3.x environment using `actions/setup-python@v2`.

5. **Install Dependencies**: Installs required Python packages:
   ```bash
   pip install --upgrade pip
   pip install PyGithub pyyaml
   ```

6. **Run Tests and Handle PRs**: Executes a Python script that:
   - Installs Innovation Engine (`ie`) if not already installed.
   - Retrieves repository information.
   - If on a branch other than `main`, it creates or updates a pull request to merge changes into `main`.
   - Iterates over all Markdown files in the repository:
     - Runs `ie execute` on each file to test the executable documentation.
     - If tests fail:
       - Extracts error logs from the ie.log file, which contains logs of Innovation Engine execution.
       - Creates or updates a GitHub issue with the error details.
       - Records the failure in the test results.
     - If tests pass, records the success.
   - Posts the test results as a comment on the pull request or outputs them in the workflow logs.

## Customization

You can modify the GitHub Action to suit your project's needs:

- **Event Triggers**: Adjust the `on` section to change when the action runs (e.g., on different branches or file types).

- **Azure Login Step**: If your documentation doesn't require Azure resources, you can remove the Azure login step.

- **Dependencies**: Add or remove Python packages in the `Install dependencies` step as needed for your tests.

- **Testing Script**: Update the Python script to change how tests are executed, how errors are handled, or how results are reported.

- **Secrets Management**: Ensure all necessary secrets are correctly named and stored in your repository's settings.

## Conclusion

By integrating this GitHub Action into your repository, you automate the testing of your executable Markdown documents using Innovation Engine. This helps maintain documentation accuracy and, combine with branch protection rules, prevents broken documentation from being merged into your main branch.