# EG

EG (meaning "for example") is a command line tool that assists in finding, customizing and executing Executable Docs. It uses Innovation Engine (`IE`) to execute the docs and Copilot to discover and customize documents.

## Setup

You will need an active Azure OpenAI deployment to use this tool locally. To create one follow the steps below.

## Prerequisites

The following prerequisites are required to complete this workshop:

- [Azure Subscription and Azure CLI](../Common/Prerequisite-AzureCLIAndSub.md)

### Environment Varaibles

In order to minimize the chance of errors and to facilitate reuse we will use Environment Variables for values we will use repeatedly in this document. For easy discovery we will use the prefix `EG_` on each variable name. The first time we encounter one of these variables in this document we will explain its purpose and a default value will be provided.

The first variable we need is `EG_HASH` this is a random string of 8 characters that will be used to create unique values when required.

```bash
export EG_HASH=$(head /dev/urandom | tr -dc A-Za-z0-9 | head -c 8)
```

### Create an Azure resource group

To create an Azure OpenAI resource, you need an Azure resource group. This will collect all of the Azure resources we need for this application to run. To create the Resource Group we will need a name and location:

```bash
export EG_RESOURCE_GROUP=EG_dev_${EG_HASH}
export EG_LOCATION=eastus2
```

To create the resource group us the `az group create` command:

```bash
az group create \
--name $EG_RESOURCE_GROUP \
--location $EG_LOCATION
```

## Create an Azure OpenAI Resource

We need an Azure OpenAI Resource, this will be configured with the following values:

```bash
export EG_RESOURCE_KIND=AIServices
export EG_OPENAI_RESOURCE_NAME=EG_dev_${EG_HASH}
export EG_OPENAI_RESOURCE_SKU=S0
```

Now you can use the [az cognitiveservices account create](/cli/azure/cognitiveservices/account?view=azure-cli-latest&preserve-view=true#az-cognitiveservices-account-create) command to create an Azure OpenAI resource in the resource group.

```bash
az cognitiveservices account create \
--name $EG_OPENAI_RESOURCE_NAME \
--resource-group $EG_RESOURCE_GROUP \
--location $EG_LOCATION \
--kind $EG_RESOURCE_KIND \
--sku $EG_OPENAI_RESOURCE_SKU
```

### Get the API key and endpoint URL

We will need the endpoint URL and API key in environment variables in order to communicate with the resource. These environment variables will not use the `EG_` prefix because this is used later to output the values of the variables and these should be conisdered secure information. Leaving the `EG_` off prevents them from being output but this script.

```bash
export OPENAI_API_ENDPOINT=$(az cognitiveservices account show \
--name $EG_OPENAI_RESOURCE_NAME \
--resource-group $EG_RESOURCE_GROUP \
| jq -r .properties.endpoint)

export OPENAI_API_KEY=$(az cognitiveservices account keys list \
--name $EG_OPENAI_RESOURCE_NAME \
--resource-group  $EG_RESOURCE_GROUP \
| jq -r .key1)
```

### Deploy a model

Now we can deploy a model into the Open AI resource. This requires a couple more variables to be defined:

```bash
export EG_DEPLOYMENT_NAME=EG_model_${EG_HASH}
export EG_MODEL=gpt-4o
export EG_MODEL_VERSION=2024-11-20
export EG_MODEL_FORMAT=OpenAI
export EG_SKU=GlobalStandard
export EG_SKU_CAPACITY=8
```

These settings will typicaly work, but be warned find the right combination can be quite a chore.

```bash
az cognitiveservices account deployment create \
--name $EG_OPENAI_RESOURCE_NAME \
--resource-group  $EG_RESOURCE_GROUP \
--deployment-name $EG_DEPLOYMENT_NAME \
--model-name $EG_MODEL \
--model-version $EG_MODEL_VERSION \
--model-format $EG_MODEL_FORMAT \
--sku $EG_SKU \
--capacity $EG_SKU_CAPACITY
```

Note that it can take a few minutes for your model to become available for use.

### Review Environment Variables

We now have an Azure OpenAI resource setup with a model deployed to it. Now would therefore be a good time to ensure that we have all the variables in one place for reference:

```bash
printenv | grep '^EG_'
```

### Install the EG CLI application from source

To install the CLI application from source you need to build the project from within the project root:

```bash
make
```

The EG command requires that the `OPENAI_API_KEY` and `OPENAI_ENDPOINT` are set. These can be retrieved using the comamnds above and thus, if you have been following along will already be set in your environment.

You can run the CLI application from source using:

```bash
./bin/eg --help
```

<!-- expected_results=".*EG is a Copilot for Executable Documentation.*" -->
```text
EG is a Copilot for Executable Documentation.

EG (meaning "for example") is a command line tool that assists in finding, customizing and executing Executable Docs.\n
Eg uses Copilot to interact with existing documentation in order to create custom executable docs.\n 
It then uses Innovation Engine (IE) to execute these docs.

[rest of help text]
```

## Usage



## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the MIT License. See the LICENSE file for details.