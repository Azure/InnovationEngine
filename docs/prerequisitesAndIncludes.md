# Prerequisites and Includes

It is often useful to break down a large document into component pieces. Long and complex documents can be off-putting, especially when the reader already has some of the base knowledge needed. There are two ways to achieve this, prerequisites and includes. The difference between them is how they are handled and where they appear in the document. This document describes both approaches.

## Prerequisites

Prerequisites are documents that should be executed before the current document proceeds. They are used to ensure, for example, that the environment is correctly setup. When running in interactive mode the user is given the opportunity to run the prerequsites interactively or non-interactively. This allows the user to skip details they already understand or to step through concepts that are new to them.

```bash
if [ -z "$UNIQUE_HASH" ]; then
    echo "Unique hash has no value yet."
else
    echo "It looks like your environment already has a value for Unique Hash is '$UNIQUE_HASH'. This needs to be unset in order to test prerequisites correctly. Clearing the existing value."
    export UNIQUE_HASH=""
fi
```

This document defines a [prerequisite](Common/prerequisiteExample.md). In fact, if you are running in Innovation Engine you will already have seen it execute. In the following sections we'll explore how that happened. We can validate it ran by ensuring that the environment variable set in that document has a vlue.

### Check Prerequisites Ran

```bash
if [ -z "$UNIQUE_HASH" ]; then
    echo "Prerequisites didn't run since UNIQUE_HASH has no value."
else
    echo "Prerequisites ran, Unique Hash is '$UNIQUE_HASH'."
fi
```

<!-- expected_similarity=0.7 -->
```text
Prerequisites ran, Unique Hash is 'abcd1234'
```

### Prerequisites Syntax

The prerequisites section starts with a heading of `## prerequisites`.

The body of this section will contain 0 or more links to a document that should be executed ahead of the current one. When viewed in a rendered form, such as a web page, the link allows the user to click through to view the document. When interpreted by Innovation Engine the document will be loaded and executed within the same context as the current document.

<!-- The following documentation is from SimDem, this behavioud has not been implemented in IE at the time of writing

### Automatically validating Pre-requisites

Some pre-requisite steps can take a long time to execute. For this
reason it is possible to provide some validation checks to see if the
pre-requisite step has been completed. These are defined in a section
towards the end of the script, before the next steps section (if one
exists). The validation steps will be executed by SimDem *before*
running the pre-requisite steps, if the tests in that section pass
then there is no need to run the pre-requisites.

It's easier to explain with an example.

Imagine we have a prerequisite step that takes 5 seconds, we don't
want to wait 5 seconds only to find that we already completed that
pre-requisite (OK, we know 5 seconds is not long, but it's long enough
to serve for this demo). For this example we will merely sleep for 5
seconds then touch a file. To validate this prequisite has been
satisfied we will test the modified date of the file, if it has been
modified in the last 5 minutes then the pre-requisite has been
satisfied.

```bash
sleep 5
echo $SIMDEM_TEMP_DIR
mkdir -p $SIMDEM_TEMP_DIR
touch $SIMDEM_TEMP_DIR/this_file_must_be_modfied_every_minute.txt
```

Now we have a set of commands that should be executed as part of this
pre-requisite. In order to use them we simply add a reference to this
file in the pre-requisites section of any other script. 

Any code in a section headed with '# Validation' will be used by
SimDem to test whether the pre-requisites have been satisfied. If
validation tests pass the pre-requisite step will be skipped over,
otherwise the other commands in the script will be executed.

### Validation

In order to continue with our example we include some vlaidation steps
in this script. If you have not run through the commands above less
than one minute ago this validation stage will fail. If you are
working through this tutorial now you just executed the above
statements and so the tests here will pass, but if you include this
file as pre-requisite again it may well fail and thus automatically
execute this script.

For this pre-requisite we need to ensure that the test.txt file has
been updated in the last 5 minutes. If not then we need to run the
commands in this document. If you are running through this document in
SimDem itself then it might be worth going back to the page that calls
this as a pre-requisite, as long as you do this in the next five
minutes you won't come back here. You can do this by selecting
"Understanding SimDem Syntax" in the next steps section.

```bash
find $SIMDEM_TEMP_DIR -name "this_file_must_be_modfied_every_minute.txt" -newermt "1 minutes ago"
```

Results:

```
/home/<username>/.simdem/tmp/this_file_must_be_modfied_every_minute.txt
```

## Includes

Includes can appear anywhere in the document and are useful for including content that is shared across multiple documents. When an executable document contains includes the content of the included file is treated as if it were a part of the original file.

TODO: document the intended behavioud and implement it

# Next Steps

TODO: port relevant content from SimDem to here and update to cover IE

  1. [Understanding SimDem Syntax](../syntax/script.md)
  2. [Configure your scripts through variables](../variables/script.md)
  3. [SimDem Index](../script.md)
  4. [Write multi-part documents](../multipart/script.md)

-->
