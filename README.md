# Overview

Innovation Engine is a tool for rapid innovation and simplification. Innovation Engine contains 
a CLI known as ie that enables execution and testing of Executable Documentation.

## What is Executable Documentation? 
Executable documentation takes standard markdown language and amplifies it by 
allowing the code commands within the document to be executed in full or step by step in an educational manner, and tested 
via automated CI/CD pipelines.

## Install Innovation Engine CLI
To install the Innovation Engine CLI, run the following commands. To install a specific version, set VERSION to the desired release number, such as "v0.1.3".

```bash
VERSION="latest"
wget -q -O ie https://github.com/Azure/InnovationEngine/releases/download/$VERSION/ie >/dev/null

# Setup permissions & move to the local bin
chmod +x ie >/dev/null
mkdir -p ~/.local/bin >/dev/null
mv ie ~/.local/bin >/dev/null

# Export the path to IE if it's not already available
if [[ !"$PATH" =~ "~/.local/bin" || !"$PATH" =~ "$HOME/.local/bin" ]]; then
	export PATH="$PATH:~/.local/bin"
fi
```

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

<!!--expected_similarity=0.8-->
<!--expected_similarity=0.8-->
```text
Hello world
```
This example purposely breaks the comment syntax so that it shows up in 
markdown. Otherwise, the tag of expected_similarity is completely invisible.

The expected similarity value is a floating point number between 0 and 1 which 
specifies how closely the output needs to match the results block. 0 being no 
similarity, 1 being an exact match.

>**Note** It may take a little bit of trial and error to find the exact value for expected_similarity.

### Environment Variables

Another barrier to automated testing is setting default values for test cases 
to use in running. This problem can be solved with command line variables in 
Executable documentation Syntax. 

Default environment variables can be set for executable documentation in a few 
different ways. 

1. A matching .ini file to the markdown
  - Upon running any document executable documentation will look for a 
  corresponding .ini file. For example if my markdown file is named tutorial.md 
  the corresponding ini file would be tutorial.ini. 
  - This file is a simple key value match for environment variable and value. 
  For example:
    ```ini
    MY_RESOURCE_GROUP_NAME = myResourceGroup
    MY_LOCATION = eastus
    MY_VM_NAME = myVM
    MY_VM_IMAGE = debian
    MY_ADMIN_USERNAME = azureuser
    ```
2. A comment at the beginning of the document containing a code blog with the 
tag 'variables'. This will be invisible to users unless they look at the raw 
markdown. For example:
    >**Note** The below example intentionally has broken comment syntax w/ two !'s.
  
    <!!-- 
    ```variables
    export MY_RESOURCE_GROUP_NAME=myResourceGroup22323
    export MY_LOCATION=eastus
    export MY_VM_NAME=myVM
    export MY_VM_IMAGE=debian
    export MY_ADMIN_USERNAME=azureuser
    ```
    -->
<!--TODO: 3. Pass in variable names as a paramter to ie CLI command. passed in vars are highest pri-->

Variables set in comments will override variables set in a .ini file. 
Consequently, locally declared variables in code samples will override 
variables set in comments. 



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

## Contributing

This is an open source project. Don't keep your code improvements,
features and cool ideas to yourself. Please issue pull requests
against our [GitHub repo](https://github.com/Azure/innovationengine).

Be sure to use our Git pre-commit script to test your contributions
before committing, simply run the following command: `python3 main.py test test`

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
