# Environment Variables

Executable Documents should use environment variables extensively and should use the following naming conventions to maximize reuse across documents.

## Prerequisites

Prerequisites can set environment variables. For example, lets set a nil value for a test variable:

```bash
export ENV_VAR_TEST=
```

Now, if we run the [Environment Variables From Prerequisite](Common/environmentVariablesFromPrerequisites.md) we should find that the value is no longer nil.

```bash
echo $ENV_VAR_TEST
```

<!-- expected_results=1.0 -->
```text
Value set in Environment Variables From Prerequisite.
```

## Naming Conventions

In general each Enviroment Variable declared in a primary Executable Document (not a prerequisite document) should us a consistent prefix. This meakes it possible to print (to the console) all variables used by that document. This can be useful in faciliating further work with the resources created. For example, here are three variables that both use the `EV_` prefix.

```bash
export EV_VAR_ONE=1
export EV_VAR_TWO=2
export EV_VAR_THREE=3
```

Now we can dump all the values set in this document with the following code:

```bash
printenv | grep '^EV_'
```