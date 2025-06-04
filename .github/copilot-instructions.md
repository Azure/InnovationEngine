## Implementation Guidelines
- Follow the Headlamp plugin architecture and React best practices.
- Use functional components and hooks for state management.
- Ensure the UI is responsive and accessible, following WCAG guidelines.
- Use CSS modules or styled-components for styling to avoid global conflicts.
- Use TypeScript for type safety and better developer experience.
- Document all implementation details in the appropriate specification document
- Write tests for all major functionality, including unit tests for components and integration tests for user flows.
- Use Jest and React Testing Library for testing React components.
- Ensure all code is linted and formatted according to the project's style guide (e.g., ESLint, Prettier).
- Use Git for version control, following the Git Flow branching strategy.
- Do not commit anything to Git until it has been reviewed and approved by at least one other team member.
- Do not ask for a review until a unity of work is complete, tested, and ready for review.
- Ask for a review of work done as soon as possible, even if it is not complete, so that feedback can be incorporated early.

## Version Management

- Use semantic versioning (e.g., v1.0.0) for releases.
- Maintain a changelog to document changes, improvements, and bug fixes.
- Use Git tags to mark releases in the repository.
- When asked to increment a version number (e.g., "bump version to 1.0.1"), change the version number in package.json and the changelog.
- When asked to push git commits follow these steps
  - Make sure that the version number is updated in package.json and the changelog and ensure that the changelog is up to date
  - Push the changes, taking note of the commit hash
  - Record the commit hash against the version number in the changelog