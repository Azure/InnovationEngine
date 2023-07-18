This document is to show the hierarchy of environment variables 


<!---
```variables
export MY_RESOURCE_GROUP=setInComments
export MY_VARIABLE_NAME=commentVariable
```
--->

```bash
echo $MY_RESOURCE_GROUP
echo $MY_VARIABLE_NAME
```

# The following will now declare variables locally which will overwrite comment variables

```bash
export MY_RESOURCE_GROUP=RGSetLocally
export MY_VARIABLE_NAME=LocallySetVariable
```

```bash
echo $MY_RESOURCE_GROUP
echo $MY_VARIABLE_NAME
```
