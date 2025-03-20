# Modes of Operation

Innovation Engine provides a number of modes of operation. You can view a summary of these with `ie --help`, this document provides more detail about each mode:

  * `execute` - Execute the commands in an executable document without interaction - ideal for unattended execution.
  * `interactive` - Execute a document in interactive mode - ideal for learning.
  * `test` - Execute the commands in a document and test output against the expected. Abort if a test fails. 
  * `to-bash` - Convert the commands in a document into a bash script for standalone execution.
  * `inspect` - Deprecated
  
## Interactive Mode

In Innovation Engine parses the document and presents it one chunk at a time. The the console displays the descriptive text along with the commands to be run and pauses for the user to indicate they are ready to progress. The user can look forward, or backward in the document and can execute the command being displayed (including any outstanding commands up until that point).

This mode is ideal for learning or teaching scenarios as it presents full context and descriptive text. If, however, you would prefer to simply run the commands without interactions use the `execute` mode instead.

## Execute Mode

Execute mode allows for unnatended execution of the document. Unless the script in the document requires user interaction the user can simply leave the script to run in this mode. However, they are also not given the opportunity to review commands before they are executed. If manual review is important use the `interactive` mode instead.

## Test Mode

Test mode runs the commands and then verifies that the output is sufficiently similar to the expected results (recorded in the markdown file) to be considered correct. This mode is similar to `execute` mode but provides more useful output in the event of a test failure.

## To-bash mode

`to-bash` mode does not execute any of the commands, instead is outputs a bash script that can be run independently of Innovation Engine. Generally you will want to send the outputs of this command to a file, e.g. `ie to-bash coolmd > cool.sh`.

## Inspect mode

This mode is deprecated and should not be used.

# Next Steps

<!--
TODO: port relevant content from SimDem to here and update to cover IE
  1. [Hello World Demo](../demo/README.md)
  2. [SimDem Index](../README.md)
  3. [Write SimDem documents](../syntax/README.md)
-->