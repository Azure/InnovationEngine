# Azure Static Web Apps Quickstart: Building Your First Static Site Using the Azure CLI

Azure Static Web Apps publishes websites to production by building apps from a code repository. In this quickstart, you deploy a web application to Azure Static Web Apps using the Azure CLI.

## Prerequisites

- Azure account. If you don't have an Azure subscription, you can create a free trial account.
- Azure CLI installed (version 2.29.0 or higher)
- A Git setup and a GitHub account (optional)

## Create a Repository (optional)

(Optional) This article uses a GitHub template repository as another way to make it easy for you to get started. The template features a starter app to deploy to Azure Static Web Apps.

- Navigate to the following location to create a new repository: https://github.com/staticwebdev/vanilla-basic/generate
- Name your repository `my-first-static-web-app`

> **Note:** Azure Static Web Apps requires at least one HTML file to create a web app. The repository you create in this step includes a single `index.html` file.

Select `Create repository`.

## Deploy a Static Web App

You can deploy the app as a static web app from the Azure CLI.

1. Create a resource group.

   ```bash
   az group create \
     --name my-swa-group \
     --location "eastus2"
   ```

Results:

<!-- expected_similarity=0.3 -->
```json
{
  "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/my-swa-group",
  "location": "eastus2",
  "managedBy": null,
  "name": "my-swa-group",
  "properties": {
    "provisioningState": "Succeeded"
  },
  "tags": null,
  "type": "Microsoft.Resources/resourceGroups"
}
```

2. Deploy a new static web app from your repository.

   ```bash
   az staticwebapp create \
       --name my-first-static-web-app \
       --resource-group my-swa-group \
       --location "eastus2" 
   ```

There are two aspects to deploying a static app. The first operation creates the underlying Azure resources that make up your app. The second is a workflow that builds and publishes your application.

Before you can go to your new static site, the deployment build must first finish running.

3. Return to your console window and run the following command to list the URLs associated with your app.

   ```bash
   az staticwebapp show \
     --name  my-first-static-web-app 
   ```

Results:

<!-- expected_similarity=0.3 -->
```json
{
    "branch": "main",
    "buildProperties": {
        "appLocation": "/",
        "apiLocation": null,
        "outputLocation": "dist"
    },
    "customDomains": [],
    "defaultHostname": "my-first-static-web-app.azurestaticapps.net",
    "hostnameSslStates": [
        {
            "name": "my-first-static-web-app.azurestaticapps.net",
            "sslState": "Disabled",
            "thumbprint": null,
            "toUpdate": null,
            "ipBasedSslResult": null,
            "hostType": "Standard"
        }
    ],
    "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/my-resource-group/providers/Microsoft.Web/staticSites/my-first-static-web-app",
    "location": "centralus",
    "name": "my-first-static-web-app",
    "repositoryUrl": "https://github.com/my-account/my-repo",
    "sku": "Free",
    "tags": {},
    "type": "Microsoft.Web/staticSites",
    "updateDomain": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
    "userDomainVerificationError": null
}
```

## Next Steps

Congratulations! You have successfully deployed a static web app to Azure Static Web Apps using the Azure CLI. Now that you have a basic understanding of how to deploy a static web app, you can explore more advanced features and functionality of Azure Static Web Apps.

In case you want to use the GitHub template repository, follow the additional steps below.

Go to https://github.com/login/device and enter the user code 329B-3945 to activate and retrieve your GitHub personal access token.

1. Go to https://github.com/login/device.
2. Enter the user code as displayed your console's message.
3. Select `Continue`.
4. Select `Authorize AzureAppServiceCLI`.

### View the Website via Git

1. As you get the repository URL while running the script, copy the repository URL and paste it into your browser.
2. Select the `Actions` tab.

   At this point, Azure is creating the resources to support your static web app. Wait until the icon next to the running workflow turns into a check mark with green background ( ). This operation may take a few minutes to complete.

3. Once the success icon appears, the workflow is complete and you can return back to your console window.
4. Run the following command to query for your website's URL.

   ```bash
   az staticwebapp show \
     --name my-first-static-web-app \
     --query "defaultHostname"
   ```

5. Copy the URL into your browser to go to your website.

## Feedback

We would love to hear your feedback on this quickstart. If you have any comments or suggestions, please let us know by opening an issue on the [GitHub repository](https://github.com/Azure/static-web-apps-docs/issues/new/choose). Thank you for your feedback!

---

That's it! I hope this markdown document is helpful to you. Let me know if you have any other questions or if there's anything else I can help you with.