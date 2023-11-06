# Setting the expectations around codeblock results

## Introduction

When ie executes a codeblock within a markdown document, the successful execution
of the command alone is usually an indication that we accomplished what the
document was aiming to achieve, however, it does not make any guarantees about
the results returned from the execution itself. For example, let's look at a
codeblock which performs an operation on virtual machines in azure. The
command may look something like this:

```bash
az vm create ...
```

And a successful output to that command may look like:

```json
{
  "fqdns": "",
  "id": "/subscriptions/<subscription-id>/resourceGroups/<resource-group>/providers/Microsoft.Compute/virtualMachines/<vm-name>",
  "identity": {
    "systemAssignedIdentity": "<managed-identity>",
    "userAssignedIdentities": {}
  },
  "location": "eastus",
  "macAddress": "00-0D-3A-1C-6B-66",
  "powerState": "VM running",
  "privateIpAddress": "10.0.0.4",
  "publicIpAddress": "172.178.12.226",
  "resourceGroup": "<resource-group>",
  "zones": ""
}
```

Looking throughout the output, we can see that there is a significant amount of
variation that can occur between mutliple successful runs of the previous
command. However though, most of the time we would only like to specify that the
output either looks like an output that we've seen before or that a
certain pattern can be found the output itself. In order for the
Innovation Engine to be able to solve this problem, we need to create a system
that allows authors to quantify what they expect the result of a command to
look like.

## Solution

Most of the codeblocks that the Innovation Engine currently executes are written
in bash and run executables to carry out tasks on behalf of a user. Instead of
only relying on the exit code to determine that a codeblock successfully
executed, we must look at the side effects (output to stdout/stderr) produced
by the command too. We will only perform comparisons when a command succeeds
and only from output sent to the standard output stream. If documentation
authors would like to make comparisons using the comamnd output from standard
error, they can manually combine the standard output/error stream into the
output stream within the codeblock itself.

To allow documentation authors to set their expectations around results,
we would like to introduce a new syntax to markdown documents in the form of
`expectation` tags. Expectation tags are HTML 5 comments that must be defined
inside of the markdown documents above a result block and below a codeblock.
Here's how that would look:

````markdown
# Example codeblock

```bash
echo "foo"
```

<!--expectation type="similarity" value="0.7"-->

```
foo
```
````

IE takes the standard output from the execution of the 
codeblock containing `echo "foo"` and compares it against the result codeblock
below our expectation tag using [Jaro Similarity](https://en.wikipedia.org/wiki/Jaro%E2%80%93Winkler_distance)
as the metric for determining the two strings similarity. The value provided by
the author sets the threshold for how similar the strings must be on a scale
of 0 to 1.

If a documentation author instead wanted to ensure that a specific pattern
existed within result of a command execution, they can change the parameters
`type` and `value` like so:

````markdown
# Example codeblock

```bash
echo "foo"

```

<!--expectation type="matches" value="foo" -->

```
foo
```
````

Instead of computing a similarity score, IE converts the `value` provided by
the author into a regular expression that is matched against the
output of `echo "foo"`.

This system should cover most scenarios in which a document author would like
to set expecatations around the value that is returned from the execution of
bash commands, but if more use cases arise then expanding on this implementation
only requires exposing new `type` and `value` parameter configurations

## Requirements

- [ ] Support for comparing the similarity of the actual command output
      against the document authors expected output codeblock given a threshold
      for how similar the strings should be.
- [ ] Support for checking if a pattern exists within the actual command output
      using regex.
- [ ] When an expectation fails, IE should give useful feedback about how to
      report the issue to the document author.
- [ ] Update at least one of the documents used within the testing pipeline
      to use the new expectation system

## Technical specifications

- If the codeblock that represents the expected result is marked as `JSON`, IE
  will attempt to parse the actual command output as JSON. If parsing the JSON
  fails, the error will include information about what to report and where to
  report it.
- When the similarity score is being computed for JSON codeblocks, the objects
  are sorted alphabetically by key before the comparison is made so that the
  similarity score that is computed is accurate as possible. We do this because
  we care more about the changes between the values inside the JSON objects
  than we do the difference in key ordering.

## Notes

- This will deprecate the old tag `expected_similarity`. We will initially
  retain support for it while we transition current documents towards using
  `expectation` but will add warnings in the log file for authors using it.
- A good but not absolute method to approximate what threshold should be used
  for a similarity score check is to measure the ratio of static text to dynamic
  text in the command output. For outputs where `static text > dynamic text`,
  you should set a higher threshold for the comparison because there is
  less text that changes between executions of the command. Conversely when
  `static text < dynamic text`, you should set a lower thresold.

## References

- [Jaro Similarity](https://en.wikipedia.org/wiki/Jaro%E2%80%93Winkler_distance)
