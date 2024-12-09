# Test Code Blocks

This documetn should contain an near exhaustive set of code blocks and surronding content for test purposes.

If this document passes `ie test` then we are good to go.

## Simple Case

The simple case is some descriptive text before the code block, followed by the code block and its results.

```bash
echo "Hello, world!"
```

<!-- expected_similarity=1.0 -->

```text
Hello, world!
```

## Sandwhich Case

The sandwich case is like the simple case above, but there is more text after the code block. Execution should be no different but the output should include the content from both before and after the code block. Currently `ie test` does not validate this, so we will need to run `ie interactive` and eyeball the results of this one.

```bash
echo "Can I have a sandwich please."
```

<!-- expected_similarity=1.0 -->

```text
Can I have a sandwich please.
```

This is the content after the code block. As long as you can see this we are good to go.
