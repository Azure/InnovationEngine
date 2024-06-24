# Testing multi Line code block

```azurecli-interactive
echo "Hello World"
```

This is what the expected output should be

<!--expected_similarity=0.8-->

```text
Hello world
```

# Testing multi Line code block

```azurecli-interactive
echo "Hello \
world"
```

# Output Should Fail

<!--expected_similarity=0.9-->

```text
Hello world
```

# Code block

```azurecli-interactive
echo "Hello \
world"
```

# Output Should Pass

<!--expected_similarity=1.0-->

```text
Hello world
```

# Code block

```azurecli-interactive
echo "Hello \
world"
```

# Bad similarity - should fail

<!--expected_similarity=9.0-->

```text
Hello world
```
