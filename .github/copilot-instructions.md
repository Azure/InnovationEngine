## Version Management

- Use semantic versioning (e.g., v1.0.0) for releases.
- Maintain a changelog to document changes, improvements, and bug fixes.
- Use Git tags to mark releases in the repository.
- When asked to increment a version number (e.g., "bump version to 1.0.1"), change the version number in package.json and the changelog.
- When asked to push git commits follow these steps
  - Make sure that the version number is updated in package.json and the changelog and ensure that the changelog is up to date
  - Push the changes, taking note of the commit hash
  - Record the commit hash against the version number in the changelog