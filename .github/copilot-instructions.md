## Specifications

- When asked to write a specification for a feature, provide a detailed description of the feature, including its purpose, functionality, and any relevant technical details.
- Include examples or use cases to illustrate how the feature should work in practice.
- Ensure that the specification is clear and unambiguous, allowing developers to implement the feature without needing further clarification.
- Use consistent terminology and formatting throughout all specifications to enhance readability.
- If applicable, reference any existing specifications, documentation or standards that the feature should adhere to.
- Consider edge cases and potential limitations of the feature, and document how these should be handled.
- Avoid assumptions about the implementation details unless explicitly stated, focusing instead on the desired outcomes and behaviors of the feature.
- If the feature involves user interaction, describe the user interface elements and their expected behavior.
- Save the specification in a file named `[SpecificationID]-feature-specification.md` in the `specifications` directory of the repository. SpecificationID is a three digit number, with leading zero's that increments by one for each specification.
- Specifications should be written in Markdown format and include the following sections (each marked as a heading):
  - **Title**: A concise title for the feature.
  - **Status**: Lines recording the status of the feature (e.g., Proposed, In Progress, Completed) after each edit, with a single sentence description of the intent of the edit, long with the last update date and time.
  - **Overview**: A brief description of why this feature is needed.
  - **Functionality**: Detailed description of how the feature should work.
  - **Use Cases**: Examples or scenarios where the feature will be used. Used cases are written in the format of "As a [user type], I want to [action] so that [benefit]."
  - **Technical Details**: Any relevant technical information, including APIs, data structures, or algorithms involved.
  - **Edge Cases**: Considerations for unusual or unexpected situations.
  - **User Interface**: Description of any user interface elements related to the feature.
  - **References**: Links to any related documentation or standards.
- Ensure that the specification is self-contained and does not rely on external documents for understanding.
- Whenever a specification is edited, update the status section to reflect the current state of the feature and the date of the last update.
- Review the specification for clarity, completeness, and adherence to the above guidelines before saving it.

## Implementation Plan

- When asked to create an implementation plan for a feature, provide a step-by-step guide on how to implement the feature.
- Follow the Headlamp plugin architecture and React best practices.
- Include tasks, subtasks, and any dependencies that need to be addressed.
- Break down the implementation into manageable tasks that can be assigned to developers and, where possible, worked on in parallel.
- Each task should be specific, actionable, and clearly defined to avoid ambiguity. There should be no doubt about which libraries and frameworks are to be used. If you need more information before making a decision then ask for it.
- Each task will be assigned a unique identifier that will be used to reference it in the implementation plan. The identifier will be [SpecificationID-TaskID], where SpecificationID is the number used in the specification filename and TaskID is a three digit number, with leading zero's that increments by one for each task with that specificationID.
- Specify the order in which tasks should be completed and indicated dependencies and any prerequisites for each task.
- Each task will define acceptance criteria to determine when it is considered complete. These criteria should be written to enable tests to be written to validate them and should be comprehensive enough to validate progression to taks that depend on this one.
- Acceptance criteria will be formatted as checkbox list items (`- [ ]`) and written to be measurable and testable.
- Use clear and concise language to describe each task, ensuring that it is actionable and understandable.
- Ensure that the implementation plan aligns with the specifications provided for the feature.
- Save the implementation plan in a file named `[SpecificationID]-implementation-plan.md` in the `specifications` directory of the repository. SpecificationID is the same number that is used for the specification this plan implements.
- Implementation plans should be written in Markdown format and include the following sections (each marked as a heading):
  - **Title**: A concise title for the implementation plan.
  - **Status**: Lines recording the status of the feature (e.g., Proposed, In Progress, Completed) after each edit, with a single sentence description of the intent of the edit, long with the last update date and time.
  - **Overview**: A brief description of the feature being implemented and its purpose, followed by a table listing the major tasks to be completed. The table will have TaskID | Title | Dependencies | Status | GitHub Issye
  - **Tasks**: A detailed list of tasks required to implement the feature, including subtasks and dependencies.
  - **Order of Tasks**: The sequence in which tasks should be completed, including any prerequisites.
  - **Estimated Timeframes**: If applicable, provide estimated timeframes for each task to assist with project planning.
- Whenever a task is started, completed or otherwise changes state, update the overview table to show the new status. If applicable include Git Branch and/or commit hash.
- Whenever an implementation plan is edited, update the status section to reflect the current state of the feature and the date of the last update.
- Ensure that the implementation plan is clear, actionable, and provides a comprehensive guide for developers to follow.

## Coding

- When asked to carry out a task in the Implementation Plan follow the steps outlined in the implementation plan.
- Ensure that the code adheres to the coding standards and best practices of the project.
- Write clear, maintainable code with appropriate comments
- All new features must be documented from both the user and developer perspective.
- Use TODO, OPTIMIZATION, and FIXME comments to indicate areas that need further work or optimization.
- Follow the project's branching strategy for version control, creating a new branch for each task. Branch names should take the form `[SpecificationID]-[TaskID]-[keywords]`, where SpecificationID is the number used in the specification filename, TaskID is the three digit number assigned to the task, and keywords are 1-3 CamelCase words describing the task
- Commit changes frequently with clear commit messages that describe the changes made.
- Before attempting to commit changes ensure all checks in the pre-commit hook are passing.
- When a task or subtask is completed (that is acceptance criteria have been met), update the status in the implementation plan.
- Use pull requests to facilitate code reviews and discussions about the changes made.
- Document any changes made to the codebase in the project's documentation, including updates to the README or other relevant files.
- If a task requires additional information or clarification, ask for it before proceeding with the implementation.
- Save the code in the appropriate directory as specified in the implementation plan, following the project's directory structure.

### Testing

- Write unit tests for new features and ensure that existing tests are updated and continue to pass as necessary.
- Ensure that the code is tested thoroughly before merging, including unit tests, integration tests, and any other relevant tests.
- Always run tests in CI mode using `CI=true` environment variable.

## Project Management

- Use the project's issue tracker to manage tasks, bugs, and feature requests.
- When updating an implementation plan check the status of tasks by looking at linked issues in the issue tracker.
- When instructed to create an issue for a task in the implementation plan, ensure that the issue is created using the content from the implementation plan. Complete with links to the plan and the specification.
- When an issue is created add a link to the issue in the implementation plan, both in the summary table and in the task details



### Version management

- When starting work on a new feature or task, ensure that the latest version of the codebase is pulled from the main branch.
- If asked to restart work as part of a review merge any changes from the main branch into your feature branch before continuing work.
- When creating a new feature branch update the project version in the `package.json` file or equivalent, following the semantic versioning (tasks are patch versions, features are minor versions).
- If a task is in development the version number should note this with a suffix of `-dev` (e.g., `1.0.0-dev`).
- Use Git for version control, following the project's branching strategy.
- Commit changes frequently with clear commit messages that describe the changes made.
- Always ensure that the Pull Request has instructions on what needs to be done to update the project (e.g. is `npm install` required)
- When a task is completed (all acceptance criteria are checked off with peer review), remove the `-dev` version postfix an open a pull request to merge the changes into the main branch.
- Ensure that the main branch is always in a deployable state, adding automation where applicable.