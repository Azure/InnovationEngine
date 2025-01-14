# Hello World Innovation Engine Demo

This document is intended to be used in interactive mode to demonstrate the key features of
Innovation Engine. If you are reading this as a document please consider working through it in Innoation Engine.

When running Innovation Engine in the terminal the user interface is split vertically into three sections, The information area is at the top, the output section is in the middle and the command list is at the bottom.

The information are shows the title of the current document, a progress meter showing your position in the document, the heading for the current step and the content of that step, up to and including the next command to be run.

Initially the output area will be empty, but when you execute a command it will display the results of that command.

The bottom section provides information about the options the user has at that moment.

If you are running this document in IE right now you will be able to see that the next command to be run is `echo "Hello World"`. You can hit the `e` key (for execute). Go ahead and hit the `e` key.

```bash
echo "Hello World"
```

<!-- expected_similarity=1.0 -->

That's great, but so what? Isn't it just outputing what is in the document illustrating the expected resutls? Well, no. The commands are actually run in a shell. For example, the date command will output the actual time at your location rather than what is written in the document.

```bash
date
```

<!-- expected_similarity=0.3 -->
```text
Wed Apr 20 15:35:31 PDT 2022
```

You can run almost any shell command this way.

<!-- TODO: add examples of major features -->

# Next Steps

  1. [Prerequisites](prerequesites.md)

<!--
TODO: port relevant content from SimDem to here and update to cover IE

It's possible to provide a branching point a the end of a script. The
user can select one of a selection of options or they can enter "quit"
(or just "q") to exit SimDem.

  1. [Write SimDem documents](../syntax/README.md)
  2. [SimDem Index](../README.md)
  3. [Modes of operation](../modes/README.md)

-->
