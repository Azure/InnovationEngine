# Variable declaration and usage

## Simple declaration

```bash
export MY_VAR="Hello, World!"
echo $MY_VAR
```

## Double variable declaration

```bash
export NEXT_VAR="Hello" && export OTHER_VAR="Hello, World!"
echo $NEXT_VAR
```

## Double declaration with semicolon

```bash
export THIS_VAR="Hello"; export THAT_VAR="Hello, World!"
echo $OTHER_VAR
```

## Declaration with subshell value

```bash
export SUBSHELL_VARIABLE=$(echo "Hello, World!")
echo $SUBSHELL_VARIABLE
```

## Declaration with other variable in value

```bash
export VAR1="Hello"
export VAR2="$VAR1, World!"
echo $VAR2
```
