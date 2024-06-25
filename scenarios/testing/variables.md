# Variable declaration and usage

## Simple declaration

```bash
export MY_VAR="Hello, World!"
echo $MY_VAR
```

<!-- expected_similarity=1.0 -->

```text
Hello, World!
```

## Double variable declaration

```bash
export NEXT_VAR="Hello" && export OTHER_VAR="Hello, World!"
echo $NEXT_VAR
```

<!-- expected_similarity=1.0 -->

```text
Hello
```

## Double declaration with semicolon

```bash
export THIS_VAR="Hello"; export THAT_VAR="Hello, World!"
echo $THAT_VAR
```

<!-- expected_similarity=1.0 -->

```text
Hello, World!
```

## Declaration with subshell value

```bash
export SUBSHELL_VARIABLE=$(echo "Hello, World!")
echo $SUBSHELL_VARIABLE
```

<!-- expected_similarity=1.0 -->

```text
Hello, World!
```

## Declaration with other variable in value

```bash
export VAR1="Hello"
export VAR2="$VAR1, World!"
echo $VAR2
```

<!-- expected_similarity=1.0 -->

```text
Hello, World!
```
