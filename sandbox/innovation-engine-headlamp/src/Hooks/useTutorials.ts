import React from 'react';

export const useTutorials = () => {
    const [tutorials, setTutorials] = React.useState<Array<any>>([]);

    React.useEffect(() => {
        setTutorials([
            { key: "azure-docs/articles/aks/learn/quick-kubernetes-deploy-cli.md", text: "Deploy an Azure Kubernetes Service (AKS) cluster" },
            { key: "azure-databases-docs/articles/mysql/flexible-server/tutorial-deploy-wordpress-on-aks.md", text: "Tutorial: Deploy WordPress on AKS cluster by using Azure CLI" },
            { key: "CreateAKSWebApp/README.md", text: "Deploy a Scalable & Secure Azure Kubernetes Service cluster using the Azure CLI" }
        ]);
        // update to retrieve tutorials from storage account once scalability effort is complete
    }, []);

    return { tutorials };
}