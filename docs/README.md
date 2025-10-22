# Welcome to Innovation Engine

Innovation Engine is a tool for rapid innovation and simplification. Innovation Engine is 
a CLI tool known as ie that enables execution and testing of Executable Documentation.

Executable Documentation is a shell script, leveraging any tools available in the shell and embedding it within documentation. That is, it takes standard markdown language and amplifies it by allowing the code commands within the document to be executed interacted with an executed.
This means that for the first time documentation is also code. 

Using Innovation Engine you can:

  * Describe the intent and expected behaviour of your shell scripts in markdown rather than comments. This means you documentation can contain hyperlinks, images, formatting etc. It can be rendered as standard markdown, e.g. as a `README.md` or a wiki page in GitHub, or as a web page. It also means that there is no need to keep two separate documents in sync. Editing code and documentation is now done in a single file.
  * Execute the code within your documentation just like any other shell script. The Innovation Engine CLI tool will parse out your script and execute it for you, as if it were a standard shell script.
  * Execute in "learn mode" onboarding new team members can be hard. Telling them to learn from a script is often going too deep too quickly, while starting from documentation presents the challenge of finding the right starting point for all skill levels. Innovation Engine allows individuals to work through the documented script at their own pace. Telling the engine to execute up to the point that they can follow and then working through step by step guided by the documentation.
  * Test the intended results of a script through the inclusion of self-documenting results blocks. This allows you to test your documentation/scripts in the CLI before merging, or in your CI/CD environment using, for example, GitHub Actions.
  * Extract the executable script from the documentation for use without Innovation Engine in the workflow.

Innovation Engine is designed to be reused in custom user experiences. For example, Microsoft Azure uses Innovation Engine to provide documentation on their Learn site, which can also be executed in the Azure Portal. This allows users to explore "good practice" documentation at the pace they prefer. They can simply read the documentation, they can interactively work through it in a customer Portal interface or they can simply go ahead and run it in order to deploy the architecture described within the document.
  
# Hello World

There's nothing magical about a document that Innovation Engine will execute. If you are writing your content in Markdown and you are marking up your bash code blocks correctly then it will be executed. For example:

```bash
echo $"Hello World"
```

In order to test that the output is as expected there's one tiny little thing you need to do, that is add an `<!-- expected_similarity=1.0 -->` comment before your code block showing the expected results:

<!-- expected_similarity=1.0 -->
```text
Hello World
```

There are multiple ways of testing the output against the actual results. Here we are saying the output should be identical to the example text (a 100%, or 1.0, match in this case). We'll cover other ways of validating the results in subsequent content.

As an excercise, if you have checked out the Innovation Engine code you could edit the echo statement above, whilst not updating the results block. When you run `ie test scripts/README.md` this code will then fail. Not very useful, but a good first example of how you can improve the quality of your documentation and scripts through testing with Innovation Engine.

# Next Steps (TBD)

1. [Modes of operation](modesOfOperation.md)
2. [Hello World Demo](helloWorldDemo.md)

<!--
TODO: port relevant content from SimDem to here and update to cover IE
  
  2. 
  3. [Build a Hello World script](tutorial/README.md)
  4. [Write SimDem documents](syntax/README.md)
  5. [Special Commands](special_commands/README.md)
  6. [Configure your scripts through variables](variables/README.md)
  7. [Write multi-part documents](multipart/README.md)
  8. [Use your documents as interactive tutorials or demos](running/README.md)
  9. [Use your documents as automated tests](test/README.md)
 10. [Build an SimDem container](building/README.md)
-->

  
