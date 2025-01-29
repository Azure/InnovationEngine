# Overview

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

## Install Innovation Engine CLI

To install the Innovation Engine CLI, run the following commands. To install a specific version, set VERSION to the desired release number, such as "v0.1.3".
You can find all releases [here](https://github.com/Azure/InnovationEngine/releases).

```bash
VERSION="latest"
wget -q -O ie https://github.com/Azure/InnovationEngine/releases/download/$VERSION/ie

# Setup permissions & move to the local bin
chmod +x ie
mkdir -p ~/.local/bin
mv ie ~/.local/bin
```

## Build Innovation Engine from Source

Paste the following commands into the shell. This will 
clone the Innovation Engine repo, install the requirements, and build out the 
Innovation Engine executable.

```bash
git clone https://github.com/Azure/InnovationEngine;
cd InnovationEngine;
make build-ie;
```

Now you can run the Innovation Engine tutorial with the following 
command:

```bash
./bin/ie execute tutorial.md
```
### Building a Container from Source

```bash
docker build -t ie .
```

Once built you can run the container and connect to it. Innovation Engine will automatically run an introductory
document when you execute this command.

```bash
docker run -it ie
```

You can override the start command if you want to take control immediately with:

```bash
docker run -it ie /bin/sh
```

## Testing Innovation Engine

Innovation Engine is self-documenting, that is all our documentation is written to be executable. Since Innovation Engine can test the results of an execution against the intended results this means our documentation is also part of our test suite. Testing against all our documentation is easy as:

```bash
make test-docs
```

If you make any changes to the IE code (see Contributing below) we would encourage you to tun the full test suite before issuing a PR.

To manual test a document it is best to run in `interactive` mode (see below). This mode provides an interactive console for reading and executing the content of Executable Documentation.

# How to Use Innovation Engine
The general format to run an executable document is: 
`ie <MODE_OF_OPERATION> <MARKDOWN_FILE>`

## Modes of Operation
Today, executable documentation can be run in 3 modes of operation:

Interactive: Displays the descriptive text of the tutorial and pauses at code 
blocks and headings to allow user interaction 
`ie interactive tutorial.md`

Test: Runs the commands and then verifies that the output is sufficiently 
similar to the expected results (recorded in the markdown file) to be 
considered correct. `ie test tutorial.md`

Execute: Reads the document and executes all of the code blocks not pausing for 
input or testing output. Essentially executes a markdown file as a script. 
`ie execute tutorial.md`

## Use Innovation Engine with any URL

Documentation does not need to be stored locally in order to run IE with it. With v0.1.3 and greater, you can run `ie execute`, `ie interactive`, and `ie test` with any URL that points to a public markdown file, including raw GitHub URLs. See the below demo:

https://github.com/Azure/InnovationEngine/assets/55719566/ce37f53c-9876-42b9-a033-1e4acaeb9d50

## Use Executable documentation for Automated Testing
One of the core benefits of executable documentation is the ability to run 
automated testing on markdown file. This can be used to ensure freshness of 
content.

In order to do this one will need to combine innovation engine executable 
documentation syntax with GitHub actions. 

In order to test if a command or action ran correctly executable documentation 
needs something to compare the results against. This requirement is met with 
result blocks.

### Result Blocks
Result blocks are distinguished in Executable documentation by a custom 
expected_similarity comment tag followed by a code block. For example

\<!--expected_similarity=0.8-->
<!--expected_similarity=0.8-->
```text
Hello world
```
In the above example we have escaped the comment syntax so that it shows up in 
markdown. Otherwise, the tag of expected_similarity is completely invisible.

The expected similarity value is a floating point number between 0 and 1 which 
specifies how closely the output needs to match the results block. 0 being no 
similarity, 1 being an exact match.

>**Note** It may take a little bit of trial and error to find the exact value for expected_similarity.

### Environment Variables

You can pass in variable declarations as an argument to the ie CLI command using the 'var' parameter. For example:
```bash
ie execute tutorial.md --var REGION=eastus
```

CLI argument variables override environment variables declared within the markdown document,
which override preexisting environment variables.

Local variables declared within the markdown document will override CLI argument variables.

Local variables (ex: `REGION=eastus`) will not persist across code blocks. It is recommended
to instead use environment variables (ex: `export REGION=eastus`).

### Setting Up GitHub Actions to use Innovation Engine

After documentation is set up to take advantage of automated testing a github 
action will need to be created to run testing on a recurring basis. The action 
will simply create a basic Linux container, install Innovation Engine 
Executable Documentation and run Executable documentation in the Test mode on 
whatever markdown files are specified.

It is important to note that if you require any specific access or cli tools 
not included in standard bash that will need to be installed in the container. 
The following example is how this may be done for a document which runs Azure 
commands.

```yml
name: 00-testing

on:
  push:
    branches:
    - main

  workflow_dispatch:

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - uses: azure/login@v1
      with:
        creds: ${{ secrets.AZURE_CREDENTIALS }}
    - name: Deploy
      env:
        AZURE_CREDENTIALS: ${{ secrets.AZURE_CREDENTIALS }}
        GITHUB_SHA: ${{ github.sha }}
      run: |
        cd $GITHUB_WORKSPACE/
        git clone https://github.com/Azure/InnovationEngine/tree/ParserAndExecutor
        cd innovationEngine
        pip3 install -r requirements.txt
        cp ../../articles/quick-create-cli.md README.md
        python3 main.py test README.md
```

# Authoring Documents

Authoring documents for use in Innovation Engine is no different from writing high quality documentation for reading. However, it does force you to follow good practice and therefore can sometimes feel a little too involved. That is  every edge case needs to be accounted for so that automated testing will reliably pass. We are therefore working on tools to help you in the authoring process.

These tools are independent of Innovation Engine, however, if you build a container from source they will be included in that container. To use them you will need an Azure OpenAI key (you can use an OpenAI key if you prefer) - be sure to add them in the command below.

```bash
docker run -it \
  -e AZURE_OPENAI_API_KEY=$AZURE_OPENAI_API_KEY \
  -e AZURE_OPENAI_ENDPOINT=$AZURE_OPENAI_ENDPOINT \
  ie /bin/sh -c "python AuthoringTools/ada.py"
```



## Contributing

This is an open source project. Don't keep your code improvements,
features and cool ideas to yourself. Please issue pull requests
against our [GitHub repo](https://github.com/Azure/innovationengine).

Be sure to use our Git pre-commit script to test your contributions
before committing, simply run the following command: `make test-docs`

This project welcomes contributions and suggestions.  Most
contributions require you to agree to a Contributor License Agreement
(CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit
https://cla.microsoft.com.

When you submit a pull request, a CLA-bot will automatically determine
whether you need to provide a CLA and decorate the PR appropriately
(e.g., label, comment). Simply follow the instructions provided by the
bot. You will only need to do this once across all repos using our
CLA.

This project has adopted
the
[Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see
the
[Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with
any additional questions or comments.


## Trademarks

This project may contain trademarks or logos for projects, products, or 
services. Authorized use of Microsoft trademarks or logos is subject to and 
must follow [Microsoft's Trademark & Brand Guidelines](https://www.microsoft.com/en-us/legal/intellectualproperty/trademarks/usage/general).
Use of Microsoft trademarks or logos in modified versions of this project must 
not cause confusion or imply Microsoft sponsorship. Any use of third-party 
trademarks or logos are subject to those third-party's policies.
