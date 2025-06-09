# Test Reports

**Status:** IN PROGRESS

## Summary

When users are testing their executable documentation using `ie test`, being
able to see the results of their execution is important, especially in instances
where the test fails. At the moment, there are only two ways to see what
happened during the test execution:

1. The user can look through the standard output from `ie test <scenario>`
   (Most common).
1. The user can look through `ie.log` (Not as common).

While these methods are effective for troubleshooting most issues, they also have
a few issues:

1. Storing the output of `ie test` in a file doesn't provide a good way to
   navigate the output and see the results of the test.
1. The log file `ie.log` is not user-friendly and can be difficult to navigate, especially
   so if the user is invoking multiple `ie test` commands as the log file is cumulative.
1. It's not easy to reproduce the execution of a specific scenario, as most of
   the variables declared by the scenario are randomized and don't have their
   values rendered in the output.
1. The output of `ie test` is not easily shareable with others.

To address these issues, we propose the introduction of test reports, a feature
for `ie test` that will generate a report of the scenario execution in JSON
format so that users can easily navigate the results of the test, reproduce
specific runs, and share the results with others.

## Requirements

- [x] The user can generate a test report by running `ie test <scenario> //report=<path>`
- [x] Reports capture the yaml metadata of the scenario.
- [x] Reports store the variables declared in the scenario and their values.
- [x] The report is generated in JSON format.
- [ ] Just like the scenarios that generated them, Reports are executable.
- [x] Outputs of the codeblocks executed are stored in the report.
- [x] Expected outputs for codeblocks are stored in the report.

## Technical specifications

- The report will be generated in JSON format, but in the future we may consider
  other formats like yaml or HTML. JSON format was chosen for v1 because it is
  easy to parse, and it is a common format for sharing data.
- Users must specify `//report=<path>` to generate a report. If the path is not
  specified, the report will not be generated.

### Report schema

The actual JSON schema is a work in progress, and will not be released with
the initial implementation, so we will list out the actual JSON with
documentation about each field until then.

```json
{
  // Name of the scenario
  "name": "Test reporting doc",
  // Properties found in the yaml header
  "properties": {
    "ms.author": "vmarcella",
    "otherProperty": "otherValue"
  },

  // Variables declared in the scenario
  "environmentVariables": {
    "NEW_VAR": "1"
  },
  // Whether the test was successful or not
  "success": true,
  // Error message if the test failed
  "error": "",
  // The step number where the test failed (-1 if successful)
  "failedAtStep": -1,
  "steps": [
    // The entire step
    {
      // The codeblock for the step
      "codeBlock": {
        // The language of the codeblock
        "language": "bash",
        // The content of the codeblock
        "content": "echo \"Hello, world!\"\n",
        // The header paired with the codeblock
        "header": "First step",
        // The paragraph paired with the codeblock
        "description": "This step will show you how to do something.",
        // The expected output for the codeblock
        "resultBlock": {
          // The language of the expected output
          "language": "text",
          // The content of the expected output
          "content": "Hello, world!\n",
          // The expected similarity score of the output (between 0 - 1)
          "expectedSimilarityScore": 1,
          // The expected regex pattern of the output
          "expectedRegexPattern": null
        }
      },
      // Codeblock number underneath the step (Should be ignored for now)
      "codeBlockNumber": 0,
      // Error message if the step failed (Would be same as top level error)
      "error": null,
      // Standard error output from executing the step
      "stdErr": "",
      // Standard output from executing the step
      "stdOut": "Hello, world!\n",
      // The name of the step
      "stepName": "First step",
      // The step number
      "stepNumber": 0,
      // Whether the step was successful or not
      "success": true,
      // The computed similarity score of the output (between 0 - 1)
      "similarityScore": 0
    },
    {
      "codeBlock": {
        "language": "bash",
        "content": "export NEW_VAR=1\n",
        "header": "Second step",
        "description": "This step will show you how to do something else.",
        "resultBlock": {
          "language": "",
          "content": "",
          "expectedSimilarityScore": 0,
          "expectedRegexPattern": null
        }
      },
      "codeBlockNumber": 0,
      "error": null,
      "stdErr": "",
      "stdOut": "",
      "stepName": "Second step",
      "stepNumber": 1,
      "success": true,
      "similarityScore": 0
    }
  ]
}
```

## Examples

Assuming you're running this command from the root of the repository:

```bash
ie test scenarios/testing/reporting.md --report=report.json >/dev/null && cat report.json
```

The output of the command above should look like this:

<!-- Need to increase this score once I fix issue #214 -->
<!-- expected_similarity=0.8 -->

```json
{
  "name": "Test reporting doc",
  "properties": {
    "ms.author": "vmarcella",
    "otherProperty": "otherValue"
  },
  "environmentVariables": {
    "NEW_VAR": "1"
  },
  "success": true,
  "error": "",
  "failedAtStep": -1,
  "steps": [
    {
      "codeBlock": {
        "language": "bash",
        "content": "echo \"Hello, world!\"\n",
        "header": "First step",
        "description": "This step will show you how to do something.",
        "resultBlock": {
          "language": "text",
          "content": "Hello, world!\n",
          "expectedSimilarityScore": 1,
          "expectedRegexPattern": null
        }
      },
      "codeBlockNumber": 0,
      "error": null,
      "stdErr": "",
      "stdOut": "Hello, world!\n",
      "stepName": "First step",
      "stepNumber": 0,
      "success": true,
      "similarityScore": 1
    },
    {
      "codeBlock": {
        "language": "bash",
        "content": "export NEW_VAR=1\n",
        "header": "Second step",
        "description": "This step will show you how to do something else.",
        "resultBlock": {
          "language": "",
          "content": "",
          "expectedSimilarityScore": 0,
          "expectedRegexPattern": null
        }
      },
      "codeBlockNumber": 0,
      "error": null,
      "stdErr": "",
      "stdOut": "",
      "stepName": "Second step",
      "stepNumber": 1,
      "success": true,
      "similarityScore": 1
    }
  ]
}
```
