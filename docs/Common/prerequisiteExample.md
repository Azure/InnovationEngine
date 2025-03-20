# Prerequisite Example

This document is a prerequisite example that is used by the [Prerequisites and Includes](prerequisitesAndIncludes.md) document. These two documents together describe and illustrate the use of Prerequisites in Innovation Engine.

## Environment Variables

Lets set an environment variable. This is a good use of pre-requisites because it allows document authors to use the same environment variables across multiple documents. This reduces the opportunity for errors and reduces the content that each author needs to create. Here we will create an 8 character hash that can be used in subsequent commands to ensure each run can create unique values for IDs.

```bash
export UNIQUE_HASH=$(head /dev/urandom | tr -dc A-Za-z0-9 | head -c 8)
```

Now we will echo this to the console. This will both serve to illustrate that this prerequisite has been executed but also allow the user to review the value.

```bash
echo "Unique hash: $UNIQUE_HASH"
```

