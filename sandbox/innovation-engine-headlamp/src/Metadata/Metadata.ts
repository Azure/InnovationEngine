import { CloudNativeResources } from "./CloudnativeResources";

export type MetaData = {
    status: "active" | "inactive";
    key: string;
    title: string;
    description: string;
    stackDetails: string[];
    sourceUrl: string;
    nextSteps: Array<{ title: string, url: string }>;
    documentationUrl: string;
    configurations?: any;
};

const tutorialMetaData: MetaData[] = [
    {
        "status": "active",
        "key": "azure-docs/articles/aks/learn/quick-kubernetes-deploy-cli.md",
        "title": CloudNativeResources.quickKubernetesTitle,
        "description": CloudNativeResources.quickKubernetesDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-docs/articles/aks/learn/quick-kubernetes-deploy-cli.md",
        "documentationUrl": "https://learn.microsoft.com/azure/aks/learn/quick-kubernetes-deploy-cli",
        "nextSteps": [
            {
                "title": CloudNativeResources.quickKubernetesNextStepOne,
                "url": "https://learn.microsoft.com/azure/aks/tutorial-kubernetes-prepare-app"
            },
            {
                "title": CloudNativeResources.quickKubernetesNextStepTwo,
                "url": "https://learn.microsoft.com/azure/architecture/reference-architectures/containers/aks-start-here?toc=/azure/aks/toc.json&bc=/azure/aks/breadcrumb/toc.json"
            }
        ],
        "configurations": {
            "permissions": [
                "Microsoft.Resources/resourceGroups/write",
                "Microsoft.Resources/resourceGroups/read",
                "Microsoft.Network/virtualNetworks/write",
                "Microsoft.Network/virtualNetworks/read",
                "Microsoft.Network/publicIPAddresses/write",
                "Microsoft.Network/publicIPAddresses/read",
                "Microsoft.Network/networkSecurityGroups/write",
                "Microsoft.Network/networkSecurityGroups/read",
                "Microsoft.Network/networkSecurityGroups/securityRules/write",
                "Microsoft.Network/networkSecurityGroups/securityRules/read",
                "Microsoft.Network/networkInterfaces/write",
                "Microsoft.Network/networkInterfaces/read",
                "Microsoft.Network/networkInterfaces/ipConfigurations/write",
                "Microsoft.Network/networkInterfaces/ipConfigurations/read",
                "Microsoft.Network/networkInterfaces/ipConfigurations/dnsSettings/update",
                "Microsoft.Network/networkInterfaces/ipConfigurations/dnsSettings/delete",
                "Microsoft.Storage/storageAccounts/write",
                "Microsoft.Storage/storageAccounts/read",
                "Microsoft.Provider/register",
                "Microsoft.Compute/virtualMachines/write",
                "Microsoft.Compute/virtualMachines/read",
                "Microsoft.Compute/virtualMachines/extensions/write",
                "Microsoft.Compute/virtualMachines/extensions/read",
                "Microsoft.ContainerService/managedClusters/write",
                "Microsoft.ContainerService/managedClusters/read",
                "Microsoft.ContainerService/managedClusters/listClusterUserCredential/action",
                "Microsoft.ContainerService/managedClusters/listClusterAdminCredential/action",
                "Microsoft.ContainerService/managedClusters/GetAccessProfiles/action",
                "Microsoft.Network/publicIPAddresses/list/action"
            ],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "MY_RESOURCE_GROUP_NAME",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "MY_AKS_CLUSTER_NAME",
                    "title": CloudNativeResources.aksClusterNameParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "azure-databases-docs/articles/mysql/flexible-server/tutorial-deploy-wordpress-on-aks.md",
        "title": CloudNativeResources.flexibleServerTitle,
        "description": CloudNativeResources.flexibleServerDescription,
        "stackDetails": [
            CloudNativeResources.flexibleServerStackItemOne,
            CloudNativeResources.flexibleServerStackItemTwo,
            CloudNativeResources.flexibleServerStackItemThree,
            CloudNativeResources.flexibleServerStackItemFive,
            CloudNativeResources.flexibleServerStackItemSix,
            CloudNativeResources.flexibleServerStackItemSeven,
            CloudNativeResources.flexibleServerStackItemEight
        ],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-databases-docs/articles/mysql/flexible-server/tutorial-deploy-wordpress-on-aks.md",
        "documentationUrl": "https://learn.microsoft.com/azure/mysql/flexible-server/tutorial-deploy-wordpress-on-aks",
        "nextSteps": [
            {
                "title": CloudNativeResources.flexibleServerNextStepOne,
                "url": "https://learn.microsoft.com/azure/aks/kubernetes-dashboard"
            },
            {
                "title": CloudNativeResources.flexibleServerNextStepTwo,
                "url": "https://learn.microsoft.com/azure/aks/tutorial-kubernetes-scale"
            },
            {
                "title": CloudNativeResources.flexibleServerNextStepThree,
                "url": "https://learn.microsoft.com/azure/mysql/flexible-server/quickstart-create-server-cli"
            },
            {
                "title": CloudNativeResources.flexibleServerNextStepFour,
                "url": "https://learn.microsoft.com/azure/mysql/flexible-server/how-to-configure-server-parameters-cli"
            }
        ],
        "configurations": {
            "permissions": [
                "Microsoft.Resources/resourceGroups/write",
                "Microsoft.Network/virtualNetworks/write",
                "Microsoft.Network/publicIPAddresses/write",
                "Microsoft.Network/networkSecurityGroups/write",
                "Microsoft.Network/networkSecurityGroups/securityRules/write",
                "Microsoft.Network/networkInterfaces/write",
                "Microsoft.Network/networkInterfaces/ipConfigurations/write",
                "Microsoft.Storage/storageAccounts/write",
                "Microsoft.Network/privateDnsZones/write",
                "Microsoft.Network/privateDnsZones/virtualNetworkLinks/write",
                "Microsoft.Network/privateDnsZones/privateDnsRecordSets/write",
                "Microsoft.Network/privateDnsZones/privateDnsRecordSets/A/write",
                "Microsoft.Network/privateDnsZones/privateDnsRecordSets/TXT/write",
                "Microsoft.Network/privateDnsZones/privateDnsRecordSets/SRV/write",
                "Microsoft.Network/privateDnsZones/privateDnsRecordSets/CNAME/write",
                "Microsoft.Network/privateDnsZones/privateDnsRecordSets/MX/write",
                "Microsoft.Network/privateDnsZones/privateDnsRecordSets/AAAA/write",
                "Microsoft.Network/privateDnsZones/privateDnsRecordSets/PTR/write",
                "Microsoft.Network/privateDnsZones/privateDnsRecordSets/CERT/write",
                "Microsoft.Network/privateDnsZones/privateDnsRecordSets/NS/write",
                "Microsoft.Network/privateDnsZones/privateDnsRecordSets/SOA/write",
                "Microsoft.Network/privateDnsZones/privateDnsRecordSets/CAA/write",
                "Microsoft.Network/privateDnsZones/privateDnsRecordSets/ANY/write",
                "Microsoft.Network/privateDnsZones/privateDnsRecordSets/SSHFP/write",
                "Microsoft.Network/privateDnsZones/privateDnsRecordSets/SPF/write",
                "Microsoft.Network/privateDnsZones/privateDnsRecordSets/DNSKEY/write",
                "Microsoft.Network/privateDnsZones/privateDnsRecordSets/DS/write",
                "Microsoft.Network/privateDnsZones/privateDnsRecordSets/NAPTR/write",
                "Microsoft.Compute/virtualMachines/write",
                "Microsoft.Compute/virtualMachines/extensions/write",
                "Microsoft.Compute/virtualMachines/read",
                "Microsoft.Authorization/roleAssignments/write",
                "Microsoft.Authorization/roleAssignments/read",
                "Microsoft.Authorization/roleDefinitions/read",
                "Microsoft.Authorization/roleDefinitions/write"
            ]
        }
    },
    {
        "status": "active",
        "key": "azure-docs/articles/static-web-apps/get-started-cli.md",
        "title": CloudNativeResources.staticWebAppTitle,
        "description": CloudNativeResources.staticWebAppDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-docs/articles/static-web-apps/get-started-cli.md",
        "documentationUrl": "https://learn.microsoft.com/azure/static-web-apps/get-started-cli",
        "nextSteps": [
            {
                "title": CloudNativeResources.staticWebAppNextStepOne,
                "url": "https://learn.microsoft.com/azure/static-web-apps/add-api"
            }
        ],
        "configurations": {
            "permissions": []
        }
    },
    {
        "status": "active",
        "key": "azure-docs/articles/virtual-machine-scale-sets/flexible-virtual-machine-scale-sets-cli.md",
        "title": CloudNativeResources.flexibleVMTitle,
        "description": CloudNativeResources.flexibleVMDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-docs/articles/virtual-machine-scale-sets/flexible-virtual-machine-scale-sets-cli.md",
        "documentationUrl": "https://learn.microsoft.com/azure/virtual-machine-scale-sets/flexible-virtual-machine-scale-sets-cli",
        "nextSteps": [
            {
                "title": CloudNativeResources.flexibleVMNextStepOne,
                "url": "https://learn.microsoft.com/azure/virtual-machine-scale-sets/flexible-virtual-machine-scale-sets-portal"
            },
            {
                "title": CloudNativeResources.flexibleVMNextStepTwo,
                "url": "https://learn.microsoft.com/azure/virtual-machine-scale-sets/overview"
            },
            {
                "title": CloudNativeResources.flexibleVMNextStepThree,
                "url": "https://learn.microsoft.com/azure/virtual-machine-scale-sets/tutorial-autoscale-cli"
            }
        ],
        "configurations": {
            "permissions": []
        }
    },
    {
        "status": "active",
        "key": "azure-docs/articles/virtual-machines/linux/quick-create-cli.md",
        "title": CloudNativeResources.quickCreateVMTitle,
        "description": CloudNativeResources.quickCreateVMDescription,
        "stackDetails": [
            CloudNativeResources.quickCreateVMStackItemOne,
            CloudNativeResources.quickCreateVMStackItemTwo,
            CloudNativeResources.quickCreateVMStackItemThree,
            CloudNativeResources.quickCreateVMStackItemFour
        ],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-docs/articles/virtual-machines/linux/quick-create-cli.md",
        "documentationUrl": "https://learn.microsoft.com/azure/virtual-machines/linux/quick-create-cli",
        "nextSteps": [
            {
                "title": CloudNativeResources.quickCreateVMNextStepOne,
                "url": "https://learn.microsoft.com/azure/virtual-machines/"
            },
            {
                "title": CloudNativeResources.quickCreateVMNextStepTwo,
                "url": "https://learn.microsoft.com/azure/virtual-machines/linux/tutorial-automate-vm-deployment"
            },
            {
                "title": CloudNativeResources.quickCreateVMNextStepThree,
                "url": "https://learn.microsoft.com/azure/virtual-machines/linux/tutorial-custom-images"
            },
            {
                "title": CloudNativeResources.quickCreateVMNextStepFour,
                "url": "https://learn.microsoft.com/azure/load-balancer/quickstart-load-balancer-standard-public-cli"
            }
        ],
        "configurations": {
            "permissions": [
                "Microsoft.Resources/resourceGroups/write",
                "Microsoft.Compute/virtualMachines/write",
                "Microsoft.Network/publicIPAddresses/write",
                "Microsoft.Network/networkInterfaces/write",
                "Microsoft.Storage/storageAccounts/listKeys/action",
                "Microsoft.Authorization/roleAssignments/write",
                "Microsoft.Compute/virtualMachines/extensions/write",
                "Microsoft.Compute/virtualMachines/read",
                "Microsoft.Network/publicIPAddresses/read",
                "Microsoft.Compute/virtualMachines/instanceView/read"
            ]
        }
    },
    {
        "status": "active",
        "key": "azure-docs/articles/virtual-machines/linux/tutorial-lemp-stack.md",
        "title": CloudNativeResources.lempStackTitle,
        "description": CloudNativeResources.lempStackDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-docs/articles/virtual-machines/linux/tutorial-lemp-stack.md",
        "documentationUrl": "https://learn.microsoft.com/azure/virtual-machines/linux/tutorial-lemp-stack",
        "nextSteps": [
            {
                "title": CloudNativeResources.lempStackNextStepOne,
                "url": "https://learn.microsoft.com/azure/virtual-machines/"
            },
            {
                "title": CloudNativeResources.lempStackNextStepTwo,
                "url": "https://learn.microsoft.com/azure/virtual-machines/linux/tutorial-manage-vm"
            },
            {
                "title": CloudNativeResources.lempStackNextStepThree,
                "url": "https://learn.microsoft.com/azure/virtual-machines/linux/tutorial-secure-web-server"
            }
        ],
        "configurations": {
            "permissions": []
        }
    },
    {
        "status": "active",
        "key": "DeployIGonAKS/README.md",
        "title": CloudNativeResources.igonTitle,
        "description": CloudNativeResources.igonDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/DeployIGonAKS/README.md",
        "documentationUrl": "https://learn.microsoft.com/troubleshoot/azure/azure-kubernetes/logs/capture-system-insights-from-aks",
        "nextSteps": [
            {
                "title": CloudNativeResources.igonNextStepOne,
                "url": "https://go.microsoft.com/fwlink/p/?linkid=2260402#use-cases"
            },
            {
                "title": CloudNativeResources.igonNextStepTwo,
                "url": "https://go.microsoft.com/fwlink/p/?linkid=2260070"
            },
            {
                "title": CloudNativeResources.igonNextStepThree,
                "url": "https://go.microsoft.com/fwlink/p/?linkid=2259865"
            }
        ],
        "configurations": {
            "permissions": []
        }
    },
    {
        "status": "active",
        "key": "CreateAKSWebApp/README.md",
        "title": CloudNativeResources.aksWebAppTitle,
        "description": CloudNativeResources.aksWebAppDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/CreateAKSWebApp/README.md",
        "documentationUrl": "",
        "nextSteps": [
            {
                "title": CloudNativeResources.aksWebAppNextStepOne,
                "url": "https://learn.microsoft.com/azure/aks/"
            },
            {
                "title": CloudNativeResources.aksWebAppNextStepTwo,
                "url": "https://learn.microsoft.com/azure/aks/tutorial-kubernetes-prepare-acr?tabs=azure-cli"
            },
            {
                "title": CloudNativeResources.aksWebAppNextStepThree,
                "url": "https://learn.microsoft.com/azure/aks/tutorial-kubernetes-scale?tabs=azure-cli"
            },
            {
                "title": CloudNativeResources.aksWebAppNextStepFour,
                "url": "https://learn.microsoft.com/azure/aks/tutorial-kubernetes-app-update?tabs=azure-cli"
            }
        ],
        "configurations": {
            "permissions": []
        }
    },
    {
        "status": "active",
        "key": "CreateRHELVMAndSSH/create-rhel-vm-ssh.md",
        "title": CloudNativeResources.rhelVmTitle,
        "description": CloudNativeResources.rhelVmDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/CreateRHELVMAndSSH/create-rhel-vm-ssh.md",
        "documentationUrl": "",
        "nextSteps": [
            {
                "title": CloudNativeResources.rhelVmNextStepOne,
                "url": "https://learn.microsoft.com/azure/virtual-machines/"
            },
            {
                "title": CloudNativeResources.rhelVmNextStepTwo,
                "url": "https://learn.microsoft.com/azure/virtual-machines/linux/quick-create-cli"
            },
            {
                "title": CloudNativeResources.rhelVmNextStepThree,
                "url": "https://learn.microsoft.com/azure/virtual-machines/linux/tutorial-custom-images"
            },
            {
                "title": CloudNativeResources.rhelVmNextStepFour,
                "url": "https://learn.microsoft.com/azure/load-balancer/quickstart-load-balancer-standard-public-cli"
            }
        ],
        "configurations": {
            "permissions": []
        }
    },
    {
        "status": "active",
        "key": "AttachDataDiskLinuxVM/attach-data-disk-linux-vm.md",
        "title": CloudNativeResources.attachDataDiskTitle,
        "description": CloudNativeResources.attachDataDiskDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/AttachDataDiskLinuxVM/attach-data-disk-linux-vm.md",
        "documentationUrl": "",
        "nextSteps": [
            {
                "title": CloudNativeResources.attachDataDiskNextStepOne,
                "url": "https://learn.microsoft.com/azure/virtual-machines/"
            },
            {
                "title": CloudNativeResources.attachDataDiskNextStepTwo,
                "url": "https://learn.microsoft.com/azure/virtual-machines/linux/tutorial-automate-vm-deployment"
            },
            {
                "title": CloudNativeResources.attachDataDiskNextStepThree,
                "url": "https://learn.microsoft.com/azure/virtual-machines/linux/tutorial-custom-images"
            },
            {
                "title": CloudNativeResources.attachDataDiskNextStepFour,
                "url": "https://learn.microsoft.com/azure/load-balancer/quickstart-load-balancer-standard-public-cli"
            }
        ],
        "configurations": {
            "permissions": []
        }
    },
    {
        "status": "active",
        "key": "ObtainPerformanceMetricsLinuxSustem/obtain-performance-metrics-linux-system.md",
        "title": CloudNativeResources.linuxSystemTitle,
        "description": CloudNativeResources.linuxSystemDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/ObtainPerformanceMetricsLinuxSustem/obtain-performance-metrics-linux-system.md",
        "documentationUrl": "https://learn.microsoft.com/troubleshoot/azure/virtual-machines/linux/collect-performance-metrics-from-a-linux-system",
        "nextSteps": [],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "MY_RESOURCE_GROUP_NAME",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "MY_VM_NAME",
                    "title": CloudNativeResources.vmNameParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "ConfigurePythonContainer/configure-python-container.md",
        "title": CloudNativeResources.pythonContainerTitle,
        "description": CloudNativeResources.pythonContainerDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/ConfigurePythonContainer/configure-python-container.md",
        "documentationUrl": "https://learn.microsoft.com/azure/app-service/configure-language-python",
        "nextSteps": [],
        "configurations": {
            "permissions": []
        }
    },
    {
        "status": "active",
        "key": "CreateSpeechService/create-speech-service.md",
        "title": CloudNativeResources.speechServiceTitle,
        "description": CloudNativeResources.speechServiceDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/CreateSpeechService/create-speech-service.md",
        "documentationUrl": "https://learn.microsoft.com/azure/ai-services/speech-service/spx-basics?tabs=windowsinstall%2Cterminal",
        "nextSteps": [],
        "configurations": {
            "permissions": []
        }
    },
    {
        "status": "active",
        "key": "DeployPremiumSSDV2/deploy-premium-ssd-v2.md",
        "title": CloudNativeResources.ssdv2Title,
        "description": CloudNativeResources.ssdv2Description,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/DeployPremiumSSDV2/deploy-premium-ssd-v2.md",
        "documentationUrl": "https://learn.microsoft.com/azure/virtual-machines/disks-deploy-premium-v2?tabs=azure-cli",
        "nextSteps": [],
        "configurations": {
            "permissions": []
        }
    },
    {
        "status": "active",
        "key": "PostgresRagLlmDemo/README.md",
        "title": CloudNativeResources.postgresRagTitle,
        "description": CloudNativeResources.postgresRagDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/PostgresRagLlmDemo/README.md",
        "documentationUrl": "",
        "nextSteps": [],
        "configurations": {
            "permissions": []
        }
    },
    {
        "status": "active",
        "key": "CreateAOAIDeployment/create-aoai-deployment.md",
        "title": CloudNativeResources.aoaDeployTitle,
        "description": CloudNativeResources.aoaDeployDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/CreateAOAIDeployment/create-aoai-deployment.md",
        "documentationUrl": "https://learn.microsoft.com/azure/ai-services/openai/how-to/create-resource?pivots=cli",
        "nextSteps": [],
        "configurations": {
            "permissions": []
        }
    },
    {
        "status": "active",
        "key": "AksKaito/README.md",
        "title": CloudNativeResources.aksKaitoTitle,
        "description": CloudNativeResources.aksKaitoDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/AksKaito/README.md",
        "documentationUrl": "https://learn.microsoft.com/azure/aks/ai-toolchain-operator",
        "nextSteps": [
            {
                "title": CloudNativeResources.aksKaitoNextStepOne,
                "url": "https://github.com/Azure/kaito"
            }
        ],
        "configurations": {
            "permissions": [],
        }
    },
    {
        "status": "active",
        "key": "azure-docs/articles/confidential-computing/confidential-enclave-nodes-aks-get-started.md",
        "title": CloudNativeResources.enclaveNodesTitle,
        "description": CloudNativeResources.enclaveNodesDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-docs/articles/confidential-computing/confidential-enclave-nodes-aks-get-started.md",
        "documentationUrl": "https://learn.microsoft.com/azure/confidential-computing/confidential-enclave-nodes-aks-get-started",
        "nextSteps": [
            {
                "title": CloudNativeResources.enclaveNodesNextStepOne,
                "url": "https://github.com/Azure-Samples/confidential-container-samples"
            },
            {
                "title": CloudNativeResources.enclaveNodesNextStepTwo,
                "url": "https://github.com/Azure-Samples/confidential-computing/blob/main/containersamples/"
            }
        ],
        "configurations": {
            "permissions": []
        }
    },
    {
        "status": "active",
        "key": "azure-management-docs/articles/azure-linux/quickstart-azure-cli.md",
        "title": CloudNativeResources.quickAzureCliTitle,
        "description": CloudNativeResources.quickAzureCliDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-management-docs/articles/azure-linux/quickstart-azure-cli.md",
        "documentationUrl": "https://learn.microsoft.com/azure/azure-linux/quickstart-azure-cli",
        "nextSteps": [
            {
                "title": CloudNativeResources.quickAzureCliNextStepOne,
                "url": "https://learn.microsoft.com/azure/azure-linux/tutorial-azure-linux-create-cluster"
            }
        ],
        "configurations": {
            "permissions": []
        }
    },
    {
        "status": "active",
        "key": "azure-docs/articles/virtual-machine-scale-sets/tutorial-use-custom-image-cli.md",
        "title": CloudNativeResources.customImageCliTitle,
        "description": CloudNativeResources.customImageCliDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-docs/articles/virtual-machine-scale-sets/tutorial-use-custom-image-cli.md",
        "documentationUrl": "https://learn.microsoft.com/azure/virtual-machine-scale-sets/tutorial-use-custom-image-cli",
        "nextSteps": [
            {
                "title": CloudNativeResources.customImageNextStepOne,
                "url": "https://learn.microsoft.com/azure/virtual-machine-scale-sets/tutorial-install-apps-cli"
            }
        ],
        "configurations": {
            "permissions": []
        }
    },
    {
        "status": "active",
        "key": "azure-docs/articles/virtual-network/create-virtual-machine-accelerated-networking.md",
        "title": CloudNativeResources.vmAccNetworkTitle,
        "description": CloudNativeResources.vmAccNetworkDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-docs/articles/virtual-network/create-virtual-machine-accelerated-networking.md",
        "documentationUrl": "https://learn.microsoft.com/azure/virtual-network/create-virtual-machine-accelerated-networking?tabs=cli",
        "nextSteps": [
            {
                "title": CloudNativeResources.vmAccNetworkNextStepOne,
                "url": "https://learn.microsoft.com/azure/virtual-network/accelerated-networking-how-it-works"
            }
        ],
        "configurations": {
            "permissions": []
        }
    },
    {
        "status": "active",
        "key": "azure-aks-docs/articles/aks/workload-identity-migrate-from-pod-identity.md",
        "title": CloudNativeResources.podIdentityTitle,
        "description": CloudNativeResources.podIdentityDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-aks-docs/articles/aks/workload-identity-migrate-from-pod-identity.md",
        "documentationUrl": "https://learn.microsoft.com/azure/aks/workload-identity-migrate-from-pod-identity",
        "nextSteps": [
            {
                "title": CloudNativeResources.podIdentityNextStepOne,
                "url": "https://learn.microsoft.com/azure/aks/workload-identity-overview"
            }
        ],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "MY_AKS_RESOURCE_GROUP",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "MY_AKS_CLUSTER_NAME",
                    "title": CloudNativeResources.aksClusterNameParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "DeployTensorflowOnAKS/deploy-tensorflow-on-aks.md",
        "title": CloudNativeResources.tensorFlowTitle,
        "description": CloudNativeResources.tensorFlowDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/DeployTensorflowOnAKS/deploy-tensorflow-on-aks.md",
        "documentationUrl": "",
        "nextSteps": [],
        "configurations": {
            "permissions": []
        }
    },
];

export const deepLinkedScenarios: MetaData[] = [
    {
        "status": "active",
        "key": "FixFstabIssuesRepairVM/fix-fstab-issues-repair-vm.md",
        "title": CloudNativeResources.fstabTitle,
        "description": CloudNativeResources.fstabDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/FixFstabIssuesRepairVM/fix-fstab-issues-repair-vm.md",
        "documentationUrl": "https://learn.microsoft.com/troubleshoot/azure/virtual-machines/linux/linux-virtual-machine-cannot-start-fstab-errors#use-azure-linux-auto-repair-alar",
        "nextSteps": [
            {
                "title": CloudNativeResources.createSupportRequestTitle,
                "url": "https://portal.azure.com/#view/Microsoft_Azure_Support/HelpAndSupportBlade/~/overview"
            }
        ],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "MY_RESOURCE_GROUP_NAME",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "MY_VM_NAME",
                    "title": CloudNativeResources.vmNameParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "KernelBootIssuesRepairVM/kernel-related-boot-issues-repairvm.md",
        "title": CloudNativeResources.kernelPanicTitle,
        "description": CloudNativeResources.kernelPanicDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/KernelBootIssuesRepairVM/kernel-related-boot-issues-repairvm.md",
        "documentationUrl": "https://learn.microsoft.com/troubleshoot/azure/virtual-machines/linux/kernel-related-boot-issues#missing-initramfs-alar",
        "nextSteps": [
            {
                "title": CloudNativeResources.createSupportRequestTitle,
                "url": "https://portal.azure.com/#view/Microsoft_Azure_Support/HelpAndSupportBlade/~/overview"
            }
        ],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "MY_RESOURCE_GROUP_NAME",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "MY_VM_NAME",
                    "title": CloudNativeResources.vmNameParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "TroubleshootVMGrubError/troubleshoot-vm-grub-error-repairvm.md",
        "title": CloudNativeResources.grubRescueTitle,
        "description": CloudNativeResources.grubRescueDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/TroubleshootVMGrubError/troubleshoot-vm-grub-error-repairvm.md",
        "documentationUrl": "https://learn.microsoft.com/troubleshoot/azure/virtual-machines/linux/troubleshoot-vm-boot-error",
        "nextSteps": [
            {
                "title": CloudNativeResources.createSupportRequestTitle,
                "url": "https://portal.azure.com/#view/Microsoft_Azure_Support/HelpAndSupportBlade/~/overview"
            }
        ],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "MY_RESOURCE_GROUP_NAME",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "MY_VM_NAME",
                    "title": CloudNativeResources.vmNameParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "azure-docs/articles/batch/quick-create-cli.md",
        "title": CloudNativeResources.batchQuickstartTitle,
        "description": CloudNativeResources.batchQuickstartDescription,
        "stackDetails": [
        ],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-docs/articles/batch/quick-create-cli.md",
        "documentationUrl": "https://learn.microsoft.com/azure/batch/quick-create-cli",
        "nextSteps": [
            {
                "title": CloudNativeResources.batchNextStepOne,
                "url": "https://learn.microsoft.com/azure/batch/tutorial-parallel-python"
            }
        ],
        "configurations": {
        }
    },
    {
        "status": "active",
        "key": "azure-compute-docs/articles/virtual-machines/linux/tutorial-manage-vm.md",
        "title": CloudNativeResources.manageLinuxVmTitle,
        "description": CloudNativeResources.manageLinuxVmDescription,
        "stackDetails": [
        ],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-compute-docs/articles/virtual-machines/linux/tutorial-manage-vm.md",
        "documentationUrl": "https://learn.microsoft.com/azure/virtual-machines/linux/tutorial-manage-vm",
        "nextSteps": [
            {
                "title": CloudNativeResources.manageLinuxVmNextStepOne,
                "url": "https://learn.microsoft.com/azure/virtual-machines/linux/tutorial-manage-disks"
            }
        ],
        "configurations": {
        }
    },
    {
        "status": "active",
        "key": "azure-compute-docs/articles/virtual-machine-scale-sets/tutorial-autoscale-cli.md",
        "title": CloudNativeResources.autoscaleVmssTitle,
        "description": CloudNativeResources.autoscaleVmssDescription,
        "stackDetails": [
        ],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-compute-docs/articles/virtual-machine-scale-sets/tutorial-autoscale-cli.md",
        "documentationUrl": "https://learn.microsoft.com/azure/virtual-machine-scale-sets/tutorial-autoscale-cli?tabs=Ubuntu",
        "nextSteps": [
        ],
        "configurations": {
        }
    },
    {
        "status": "active",
        "key": "azure-compute-docs/articles/virtual-machine-scale-sets/tutorial-modify-scale-sets-cli.md",
        "title": CloudNativeResources.modifyVmssTitle,
        "description": CloudNativeResources.modifyVmssDescription,
        "stackDetails": [
        ],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-compute-docs/articles/virtual-machine-scale-sets/tutorial-modify-scale-sets-cli.md",
        "documentationUrl": "https://learn.microsoft.com/azure/virtual-machine-scale-sets/tutorial-modify-scale-sets-cli",
        "nextSteps": [
            {
                "title": CloudNativeResources.modifyVmssNextStepOne,
                "url": "https://learn.microsoft.com/azure/virtual-machine-scale-sets/tutorial-use-disks-powershell"
            }
        ],
        "configurations": {
        }
    },
    {
        "status": "active",
        "key": "azure-compute-docs/articles/virtual-machines/disks-enable-performance.md",
        "title": CloudNativeResources.diskPerfTitle,
        "description": CloudNativeResources.diskPerfDescription,
        "stackDetails": [
        ],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-compute-docs/articles/virtual-machines/disks-enable-performance.md",
        "documentationUrl": "https://learn.microsoft.com/azure/virtual-machines/disks-enable-performance?tabs=azure-cli",
        "nextSteps": [
            {
                "title": CloudNativeResources.diskPerfNextStepOne,
                "url": "https://learn.microsoft.com/azure/virtual-machines/disks-incremental-snapshots"
            },
            {
                "title": CloudNativeResources.diskPerfNextStepTwo,
                "url": "https://learn.microsoft.com/azure/virtual-machines/linux/expand-disks"
            }
        ],
        "configurations": {
        }
    },
    {
        "status": "active",
        "key": "azure-compute-docs/articles/container-instances/container-instances-vnet.md",
        "title": CloudNativeResources.aciVnetTitle,
        "description": CloudNativeResources.aciVnetDescription,
        "stackDetails": [
        ],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-compute-docs/articles/container-instances/container-instances-vnet.md",
        "documentationUrl": "https://learn.microsoft.com/azure/container-instances/container-instances-vnet",
        "nextSteps": [
            {
                "title": CloudNativeResources.aciVnetNextStepOne,
                "url": "https://github.com/Azure/azure-quickstart-templates/tree/master/quickstarts/microsoft.containerinstance/aci-vnet"
            },
            {
                "title": CloudNativeResources.aciVnetNextStepTwo,
                "url": "https://learn.microsoft.com/azure/container-instances/using-azure-container-registry-mi"
            }
        ],
        "configurations": {
        }
    },
    {
        "status": "active",
        "key": "azure-compute-docs/articles/virtual-machines/linux/multiple-nics.md",
        "title": CloudNativeResources.multiNicVmTitle,
        "description": CloudNativeResources.multiNicVmDescription,
        "stackDetails": [
        ],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-compute-docs/articles/virtual-machines/linux/multiple-nics.md",
        "documentationUrl": "https://learn.microsoft.com/azure/virtual-machines/linux/multiple-nics",
        "nextSteps": [
            {
                "title": CloudNativeResources.multiNicVmNextStepOne,
                "url": "https://learn.microsoft.com/azure/virtual-machines/sizes"
            },
            {
                "title": CloudNativeResources.multiNicVmNextStepTwo,
                "url": "https://learn.microsoft.com/azure/defender-for-cloud/just-in-time-access-usage"
            }
        ],
        "configurations": {
        }
    },
    {
        "status": "inactive",
        "key": "azure-compute-docs/articles/virtual-machines/linux/quick-create-terraform/quick-create-terraform.md",
        "title": CloudNativeResources.terraformVmTitle,
        "description": CloudNativeResources.terraformVmDescription,
        "stackDetails": [
        ],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-compute-docs/articles/virtual-machines/linux/quick-create-terraform/quick-create-terraform.md",
        "documentationUrl": "https://learn.microsoft.com/azure/virtual-machines/linux/quick-create-terraform?tabs=azure-cli",
        "nextSteps": [
            {
                "title": CloudNativeResources.terraformVmNextStepOne,
                "url": "https://learn.microsoft.com/azure/developer/terraform/troubleshoot"
            },
            {
                "title": CloudNativeResources.terraformVmNextStepTwo,
                "url": "https://learn.microsoft.com/azure/virtual-machines/linux/tutorial-manage-vm"
            }
        ],
        "configurations": {
        }
    },
    {
        "status": "active",
        "key": "upstream/FlatcarOnAzure/flatcar-on-azure.md",
        "title": CloudNativeResources.flatcarTitle,
        "description": CloudNativeResources.flatcarDescription,
        "stackDetails": [
        ],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/upstream/FlatcarOnAzure/flatcar-on-azure.md",
        "documentationUrl": "https://www.flatcar.org/docs/latest/installing/cloud/azure/",
        "nextSteps": [],
        "configurations": {
        }
    },
    {
        "status": "active",
        "key": "azure-management-docs/articles/azure-linux/tutorial-azure-linux-migration.md",
        "title": CloudNativeResources.azLinuxMigrateTitle,
        "description": CloudNativeResources.azLinuxMigrateDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-management-docs/articles/azure-linux/tutorial-azure-linux-migration.md",
        "documentationUrl": "https://learn.microsoft.com/azure/azure-linux/tutorial-azure-linux-migration?tabs=azure-cli",
        "nextSteps": [
            {
                "title": CloudNativeResources.azLinuxMigrateNextStepOne,
                "url": "https://learn.microsoft.com/azure/azure-linux/tutorial-azure-linux-telemetry-monitor"
            }
        ],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "RESOURCE_GROUP",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "CLUSTER_NAME",
                    "title": CloudNativeResources.aksClusterNameParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "azure-management-docs/articles/azure-linux/tutorial-azure-linux-create-cluster.md",
        "title": CloudNativeResources.azLinuxCreateClusterTitle,
        "description": CloudNativeResources.azLinuxCreateClusterDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-management-docs/articles/azure-linux/tutorial-azure-linux-create-cluster.md",
        "documentationUrl": "https://learn.microsoft.com/azure/azure-linux/tutorial-azure-linux-create-cluster",
        "nextSteps": [
            {
                "title": CloudNativeResources.azLinuxCreateClusterNextStepOne,
                "url": "https://learn.microsoft.com/azure/azure-linux/tutorial-azure-linux-add-nodepool"
            }
        ],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "RESOURCE_GROUP_NAME",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "CLUSTER_NAME",
                    "title": CloudNativeResources.aksClusterNameParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "azure-management-docs/articles/azure-linux/tutorial-azure-linux-add-nodepool.md",
        "title": CloudNativeResources.azLinuxAddNodepoolTitle,
        "description": CloudNativeResources.azLinuxAddNodepoolDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-management-docs/articles/azure-linux/tutorial-azure-linux-add-nodepool.md",
        "documentationUrl": "https://learn.microsoft.com/azure/azure-linux/tutorial-azure-linux-add-nodepool",
        "nextSteps": [
            {
                "title": CloudNativeResources.azLinuxAddNodepoolNextStepOne,
                "url": "https://learn.microsoft.com/azure/azure-linux/tutorial-azure-linux-migration"
            }
        ],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "RESOURCE_GROUP",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "CLUSTER_NAME",
                    "title": CloudNativeResources.aksClusterNameParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "azure-management-docs/articles/azure-linux/tutorial-azure-linux-upgrade.md",
        "title": CloudNativeResources.azLinuxUpgradeTitle,
        "description": CloudNativeResources.azLinuxUpgradeDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-management-docs/articles/azure-linux/tutorial-azure-linux-upgrade.md",
        "documentationUrl": "https://learn.microsoft.com/azure/azure-linux/tutorial-azure-linux-upgrade",
        "nextSteps": [
            {
                "title": CloudNativeResources.azLinuxUpgradeNextStepOne,
                "url": "https://learn.microsoft.com/azure/azure-linux/intro-azure-linux"
            }
        ],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "AZ_LINUX_RG",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "AZ_LINUX_CLUSTER",
                    "title": CloudNativeResources.aksClusterNameParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "azure-management-docs/articles/azure-linux/tutorial-azure-linux-telemetry-monitor.md",
        "title": CloudNativeResources.azLinuxTelemetryTitle,
        "description": CloudNativeResources.azLinuxTelemetryDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-management-docs/articles/azure-linux/tutorial-azure-linux-telemetry-monitor.md",
        "documentationUrl": "https://learn.microsoft.com/azure/azure-linux/tutorial-azure-linux-telemetry-monitor",
        "nextSteps": [
            {
                "title": CloudNativeResources.azLinuxTelemetryNextStepOne,
                "url": "https://learn.microsoft.com/azure/azure-linux/tutorial-azure-linux-upgrade"
            }
        ],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "RESOURCE_GROUP",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "CLUSTER_NAME",
                    "title": CloudNativeResources.aksClusterNameParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "azure-stack-docs/azure-stack/user/azure-stack-quick-create-vm-linux-cli.md",
        "title": CloudNativeResources.azStackVmTitle,
        "description": CloudNativeResources.azStackVmDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-stack-docs/azure-stack/user/azure-stack-quick-create-vm-linux-cli.md",
        "documentationUrl": "https://learn.microsoft.com/azure-stack/user/azure-stack-quick-create-vm-linux-cli?view=azs-2501",
        "nextSteps": [
            {
                "title": CloudNativeResources.azStackVmNextStepOne,
                "url": "https://learn.microsoft.com/azure-stack/user/azure-stack-vm-considerations"
            }
        ],
        "configurations": {
            "permissions": [],
            "configurableParams": []
        }
    },
    {
        "status": "active",
        "key": "azure-aks-docs/articles/aks/azure-cni-powered-by-cilium.md",
        "title": CloudNativeResources.cniCiliumTitle,
        "description": CloudNativeResources.cniCiliumDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-aks-docs/articles/aks/azure-cni-powered-by-cilium.md",
        "documentationUrl": "https://learn.microsoft.com/azure/aks/azure-cni-powered-by-cilium",
        "nextSteps": [
            {
                "title": CloudNativeResources.cniCiliumNextStepOne,
                "url": "https://learn.microsoft.com/azure/aks/upgrade-azure-cni"
            },
            {
                "title": CloudNativeResources.cniCiliumNextStepTwo,
                "url": "https://learn.microsoft.com/azure/aks/static-ip"
            },
            {
                "title": CloudNativeResources.cniCiliumNextStepThree,
                "url": "https://learn.microsoft.com/azure/aks/internal-lb"
            },
            {
                "title": CloudNativeResources.cniCiliumNextStepFour,
                "url": "https://learn.microsoft.com/azure/aks/ingress-basic"
            }
        ],
        "configurations": {
            "permissions": [],
            "configurableParams": []
        }
    },
    {
        "status": "active",
        "key": "azure-compute-docs/articles/virtual-machines/linux/tutorial-automate-vm-deployment.md",
        "title": CloudNativeResources.cloudInitVmTitle,
        "description": CloudNativeResources.cloudInitVmDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-compute-docs/articles/virtual-machines/linux/tutorial-automate-vm-deployment.md",
        "documentationUrl": "https://learn.microsoft.com/azure/virtual-machines/linux/tutorial-automate-vm-deployment",
        "nextSteps": [
            {
                "title": CloudNativeResources.cloudInitVmNextStepOne,
                "url": "https://learn.microsoft.com/azure/virtual-machines/linux/tutorial-custom-images"
            }
        ],
        "configurations": {
            "permissions": [],
            "configurableParams": []
        }
    },
    {
        "status": "active",
        "key": "azure-compute-docs/articles/virtual-machines/linux/multiple-nics.md",
        "title": CloudNativeResources.multiNicVmTitle,
        "description": CloudNativeResources.multiNicVmDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-compute-docs/articles/virtual-machines/linux/multiple-nics.md",
        "documentationUrl": "https://learn.microsoft.com/azure/virtual-machines/linux/multiple-nics",
        "nextSteps": [
            {
                "title": CloudNativeResources.multiNicVmNextStepOne,
                "url": "https://learn.microsoft.com/azure/virtual-machines/sizes"
            },
            {
                "title": CloudNativeResources.multiNicVmNextStepTwo,
                "url": "https://github.com/MicrosoftDocs/azure-compute-docs/blob/main/azure/security-center/security-center-just-in-time"
            }
        ],
        "configurations": {
            "permissions": [],
            "configurableParams": []
        }
    },
    {
        "status": "active",
        "key": "azure-compute-docs/articles/virtual-machines/disks-enable-performance.md",
        "title": CloudNativeResources.diskPerfTitle,
        "description": CloudNativeResources.diskPerfDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-compute-docs/articles/virtual-machines/disks-enable-performance.md",
        "documentationUrl": "https://learn.microsoft.com/azure/virtual-machines/disks-enable-performance?tabs=azure-cli",
        "nextSteps": [
            {
                "title": CloudNativeResources.diskPerfNextStepOne,
                "url": "https://learn.microsoft.com/azure/virtual-machines/disks-incremental-snapshots"
            },
            {
                "title": CloudNativeResources.diskPerfNextStepTwo,
                "url": "https://learn.microsoft.com/azure/virtual-machines/linux/expand-disks"
            },
            {
                "title": CloudNativeResources.diskPerfNextStepThree,
                "url": "https://learn.microsoft.com/azure/virtual-machines/windows/expand-os-disk"
            }
        ],
        "configurations": {
            "permissions": [],
            "configurableParams": []
        }
    },
    {
        "status": "active",
        "key": "azure-compute-docs/articles/virtual-machine-scale-sets/tutorial-modify-scale-sets-cli.md",
        "title": CloudNativeResources.modifyVmssTitle,
        "description": CloudNativeResources.modifyVmssDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-compute-docs/articles/virtual-machine-scale-sets/tutorial-modify-scale-sets-cli.md",
        "documentationUrl": "https://learn.microsoft.com/azure/virtual-machine-scale-sets/tutorial-modify-scale-sets-cli",
        "nextSteps": [
            {
                "title": CloudNativeResources.modifyVmssNextStepOne,
                "url": "https://learn.microsoft.com/azure/virtual-machine-scale-sets/tutorial-use-disks-powershell"
            }
        ],
        "configurations": {
            "permissions": [],
            "configurableParams": []
        }
    },
    {
        "status": "active",
        "key": "azure-compute-docs/articles/virtual-machine-scale-sets/tutorial-autoscale-cli.md",
        "title": CloudNativeResources.autoscaleVmssTitle,
        "description": CloudNativeResources.autoscaleVmssDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-compute-docs/articles/virtual-machine-scale-sets/tutorial-autoscale-cli.md",
        "documentationUrl": "https://learn.microsoft.com/azure/virtual-machine-scale-sets/tutorial-autoscale-cli",
        "nextSteps": [
            {
                "title": CloudNativeResources.autoscaleVmssNextStepOne,
                "url": "https://learn.microsoft.com/azure/virtual-machine-scale-sets/virtual-machine-scale-sets-instance-protection"
            }
        ],
        "configurations": {
            "permissions": [],
            "configurableParams": []
        }
    },
    {
        "status": "active",
        "key": "azure-compute-docs/articles/virtual-machines/linux/tutorial-manage-vm.md",
        "title": CloudNativeResources.manageLinuxVmTitle,
        "description": CloudNativeResources.manageLinuxVmDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-compute-docs/articles/virtual-machines/linux/tutorial-manage-vm.md",
        "documentationUrl": "https://learn.microsoft.com/azure/virtual-machines/linux/tutorial-manage-vm",
        "nextSteps": [
            {
                "title": CloudNativeResources.manageLinuxVmNextStepOne,
                "url": "https://github.com/MicrosoftDocs/azure-compute-docs/blob/main/articles/virtual-machines/linux/tutorial-manage-disks.md"
            }
        ],
        "configurations": {
            "permissions": [],
            "configurableParams": []
        }
    },
    {
        "status": "active",
        "key": "azure-compute-docs/articles/virtual-machines/linux/tutorial-lamp-stack.md",
        "title": CloudNativeResources.lampWpVmTitle,
        "description": CloudNativeResources.lampWpVmDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-compute-docs/articles/virtual-machines/linux/tutorial-lamp-stack.md",
        "documentationUrl": "https://learn.microsoft.com/azure/virtual-machines/linux/tutorial-lamp-stack",
        "nextSteps": [
            {
                "title": CloudNativeResources.lampWpVmNextStepOne,
                "url": "https://learn.microsoft.com/azure/virtual-machines/linux/tutorial-secure-web-server"
            }
        ],
        "configurations": {
            "permissions": [],
            "configurableParams": []
        }
    },
    {
        "status": "active",
        "key": "azure-docs/articles/batch/quick-create-cli.md",
        "title": CloudNativeResources.batchQuickstartTitle,
        "description": CloudNativeResources.batchQuickstartDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-docs/articles/batch/quick-create-cli.md",
        "documentationUrl": "https://learn.microsoft.com/azure/batch/quick-create-cli",
        "nextSteps": [
            {
                "title": CloudNativeResources.batchNextStepOne,
                "url": "https://learn.microsoft.com/azure/batch/tutorial-parallel-python"
            }
        ],
        "configurations": {
            "permissions": [],
            "configurableParams": []
        }
    },
    {
        "status": "active",
        "key": "azure-aks-docs/articles/aks/node-image-upgrade.md",
        "title": CloudNativeResources.aksNodeImageUpgradeTitle,
        "description": CloudNativeResources.aksNodeImageUpgradeDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-aks-docs/articles/aks/node-image-upgrade.md",
        "documentationUrl": "https://learn.microsoft.com/azure/aks/node-image-upgrade",
        "nextSteps": [
            {
                "title": CloudNativeResources.aksNodeImageUpgradeNextStepOne,
                "url": "https://github.com/Azure/AKS/releases"
            },
            {
                "title": CloudNativeResources.aksNodeImageUpgradeNextStepTwo,
                "url": "https://learn.microsoft.com/azure/aks/upgrade-aks-cluster"
            },
            {
                "title": CloudNativeResources.aksNodeImageUpgradeNextStepThree,
                "url": "https://learn.microsoft.com/azure/aks/node-upgrade-github-actions"
            },
            {
                "title": CloudNativeResources.aksNodeImageUpgradeNextStepFour,
                "url": "https://learn.microsoft.com/azure/aks/create-node-pools"
            },
            {
                "title": CloudNativeResources.aksNodeImageUpgradeNextStepFive,
                "url": "https://learn.microsoft.com/azure/architecture/operator-guides/aks/aks-upgrade-practices"
            }
        ],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "AKS_RESOURCE_GROUP",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "AKS_CLUSTER",
                    "title": CloudNativeResources.aksClusterNameParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "AKS_NODEPOOL",
                    "title": CloudNativeResources.aksNodePoolNameTitle,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "azure-compute-docs/articles/virtual-machines/linux/tutorial-elasticsearch.md",
        "title": CloudNativeResources.elasticSearchVmTitle,
        "description": CloudNativeResources.elasticSearchVmDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-compute-docs/articles/virtual-machines/linux/tutorial-elasticsearch.md",
        "documentationUrl": "https://learn.microsoft.com/azure/virtual-machines/linux/tutorial-elasticsearch",
        "nextSteps": [
            {
                "title": CloudNativeResources.elasticSearchVmNextStepOne,
                "url": "https://learn.microsoft.com/azure/virtual-machines/linux/quick-create-cli"
            }
        ],
        "configurations": {
            "permissions": [],
            "configurableParams": []
        }
    },
    {
        "status": "active",
        "key": "azure-aks-docs/articles/aks/learn/quick-windows-container-deploy-cli.md",
        "title": CloudNativeResources.aksWindowsContainerTitle,
        "description": CloudNativeResources.aksWindowsContainerDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-aks-docs/articles/aks/learn/quick-windows-container-deploy-cli.md",
        "documentationUrl": "https://learn.microsoft.com/azure/aks/learn/quick-windows-container-deploy-cli?tabs=add-windows-node-pool",
        "nextSteps": [
            {
                "title": CloudNativeResources.aksWindowsContainerNextStepOne,
                "url": "https://learn.microsoft.com/azure/architecture/reference-architectures/containers/aks-start-here?toc=/azure/aks/toc.json&bc=/azure/aks/breadcrumb/toc.json"
            },
            {
                "title": CloudNativeResources.aksWindowsContainerNextStepTwo,
                "url": "https://learn.microsoft.com/azure/aks/tutorial-kubernetes-prepare-app"
            }
        ],
        "configurations": {
            "permissions": [],
            "configurableParams": []
        }
    },
    {
        "status": "active",
        "key": "azure-aks-docs/articles/aks/spot-node-pool.md",
        "title": CloudNativeResources.aksSpotNodepoolTitle,
        "description": CloudNativeResources.aksSpotNodepoolDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-aks-docs/articles/aks/spot-node-pool.md",
        "documentationUrl": "https://learn.microsoft.com/azure/aks/spot-node-pool",
        "nextSteps": [
            {
                "title": CloudNativeResources.aksSpotNodepoolNextStepOne,
                "url": "https://learn.microsoft.com/azure/aks/operator-best-practices-advanced-scheduler"
            }
        ],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "RESOURCE_GROUP",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "AKS_CLUSTER",
                    "title": CloudNativeResources.aksClusterNameParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "azure-aks-docs/articles/aks/auto-upgrade-cluster.md",
        "title": CloudNativeResources.aksAutoUpgradeClusterTitle,
        "description": CloudNativeResources.aksAutoUpgradeClusterDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-aks-docs/articles/aks/auto-upgrade-cluster.md",
        "documentationUrl": "https://learn.microsoft.com/azure/aks/auto-upgrade-cluster?tabs=azure-cli",
        "nextSteps": [
            {
                "title": CloudNativeResources.aksAutoUpgradeClusterNextStepOne,
                "url": "https://learn.microsoft.com/azure/architecture/operator-guides/aks/aks-upgrade-practices"
            }
        ],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "RESOURCE_GROUP",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "AKS_CLUSTER_NAME",
                    "title": CloudNativeResources.aksClusterNameParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "azure-aks-docs/articles/aks/auto-upgrade-node-os-image.md",
        "title": CloudNativeResources.aksAutoUpgradeNodeOsTitle,
        "description": CloudNativeResources.aksAutoUpgradeNodeOsDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-aks-docs/articles/aks/auto-upgrade-node-os-image.md",
        "documentationUrl": "https://learn.microsoft.com/azure/aks/auto-upgrade-node-os-image?tabs=azure-cli",
        "nextSteps": [
            {
                "title": CloudNativeResources.aksAutoUpgradeNodeOsNextStepOne,
                "url": "https://learn.microsoft.com/azure/architecture/operator-guides/aks/aks-upgrade-practices"
            }
        ],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "RESOURCE_GROUP",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "AKS_CLUSTER",
                    "title": CloudNativeResources.aksClusterNameParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "azure-aks-docs/articles/aks/cost-analysis.md",
        "title": CloudNativeResources.aksCostAnalysisTitle,
        "description": CloudNativeResources.aksCostAnalysisDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-aks-docs/articles/aks/cost-analysis.md",
        "documentationUrl": "https://learn.microsoft.com/azure/aks/cost-analysis",
        "nextSteps": [
            {
                "title": CloudNativeResources.aksCostAnalysisNextStepOne,
                "url": "https://learn.microsoft.com/azure/aks/understand-aks-costs"
            }
        ],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "RESOURCE_GROUP",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "CLUSTER_NAME",
                    "title": CloudNativeResources.aksClusterNameParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "azure-aks-docs/articles/aks/istio-deploy-addon.md",
        "title": CloudNativeResources.aksIstioAddonTitle,
        "description": CloudNativeResources.aksIstioAddonDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-aks-docs/articles/aks/istio-deploy-addon.md",
        "documentationUrl": "https://learn.microsoft.com/azure/aks/istio-deploy-addon",
        "nextSteps": [
            {
                "title": CloudNativeResources.aksIstioAddonNextStepOne,
                "url": "https://learn.microsoft.com/azure/aks/istio-deploy-ingress"
            },
            {
                "title": CloudNativeResources.aksIstioAddonNextStepTwo,
                "url": "https://learn.microsoft.com/azure/aks/istio-scale#scaling"
            },
            {
                "title": CloudNativeResources.aksIstioAddonNextStepThree,
                "url": "https://learn.microsoft.com/azure/aks/istio-metrics-managed-prometheus"
            }
        ],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "RESOURCE_GROUP",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "CLUSTER",
                    "title": CloudNativeResources.aksClusterNameParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "azure-aks-docs/articles/aks/istio-scale.md",
        "title": CloudNativeResources.istioServiceMeshTitle,
        "description": CloudNativeResources.istioServiceMeshScalingDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-aks-docs/articles/aks/istio-scale.md",
        "documentationUrl": "https://learn.microsoft.com/azure/aks/istio-scale",
        "nextSteps": [],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "RESOURCE_GROUP",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "AKS_CLUSTER",
                    "title": CloudNativeResources.aksClusterNameParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "SupportArticles-docs/support/azure/azure-kubernetes/availability-performance/node-not-ready-custom-script-extension-errors.md",
        "title": CloudNativeResources.nodeNotReadyCseErrorsTitle,
        "description": CloudNativeResources.nodeNotReadyCseErrorsDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/SupportArticles-docs/support/azure/azure-kubernetes/availability-performance/node-not-ready-custom-script-extension-errors.md",
        "documentationUrl": "https://learn.microsoft.com/troubleshoot/azure/azure-kubernetes/availability-performance/node-not-ready-custom-script-extension-errors",
        "nextSteps": [],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "RG_NAME",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "CLUSTER_NAME",
                    "title": CloudNativeResources.aksClusterNameParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "AVAILABILITY_SET_VM",
                    "title": CloudNativeResources.availabilitySetVmNameParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "SupportArticles-docs/support/azure/azure-kubernetes/availability-performance/node-not-ready-after-being-healthy.md",
        "title": CloudNativeResources.nodeNotReadyHealthyStateTitle,
        "description": CloudNativeResources.nodeNotReadyHealthyStateDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/NodeNotReadyAKS/node-not-ready-aks.md",
        "documentationUrl": "https://learn.microsoft.com/troubleshoot/azure/azure-kubernetes/availability-performance/node-not-ready-after-being-healthy",
        "nextSteps": [],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "RESOURCE_GROUP",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "AKS_CLUSTER",
                    "title": CloudNativeResources.aksClusterNameParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "SupportArticles-docs/support/azure/azure-kubernetes/connectivity/tcp-timeouts-dial-tcp-nodeip-10250-io-timeout.md",
        "title": CloudNativeResources.tcpTimeoutKubeletTitle,
        "description": CloudNativeResources.tcpTimeoutKubeletDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/SupportArticles-docs/support/azure/azure-kubernetes/connectivity/tcp-timeouts-dial-tcp-nodeip-10250-io-timeout.md",
        "documentationUrl": "https://learn.microsoft.com/troubleshoot/azure/azure-kubernetes/connectivity/tcp-timeouts-dial-tcp-nodeip-10250-io-timeout",
        "nextSteps": [],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "POD_NAME",
                    "title": CloudNativeResources.podNameParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "RESOURCE_GROUP",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "CLUSTER_NAME",
                    "title": CloudNativeResources.clusterNameParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "azure-aks-docs/articles/aks/kubelet-logs.md",
        "title": CloudNativeResources.viewKubeletLogsTitle,
        "description": CloudNativeResources.viewKubeletLogsDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-aks-docs/articles/aks/kubelet-logs.md",
        "documentationUrl": "https://learn.microsoft.com/azure/aks/kubelet-logs",
        "nextSteps": [
            {
                "title": CloudNativeResources.sshIntoAksNodesNextStep,
                "url": "https://learn.microsoft.com/azure/aks/ssh"
            }
        ],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "RESOURCE_GROUP_NAME",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "AKS_CLUSTER_NAME",
                    "title": CloudNativeResources.aksClusterNameParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "NODE_NAME",
                    "title": CloudNativeResources.nodeNameParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "azure-aks-docs/articles/aks/delete-cluster.md",
        "title": CloudNativeResources.deleteAksClusterTitle,
        "description": CloudNativeResources.deleteAksClusterDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-aks-docs/articles/aks/delete-cluster.md",
        "documentationUrl": "https://learn.microsoft.com/azure/aks/delete-cluster",
        "nextSteps": [
            {
                "title": CloudNativeResources.stopAksClusterNextStep,
                "url": "https://learn.microsoft.com/azure/aks/start-stop-cluster?tabs=azure-cli"
            },
            {
                "title": CloudNativeResources.upgradeAksClusterNextStep,
                "url": "https://learn.microsoft.com/azure/aks/upgrade-cluster"
            }
        ],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "AKS_CLUSTER_NAME",
                    "title": CloudNativeResources.aksClusterNameParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "RESOURCE_GROUP_NAME",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "azure-aks-docs/articles/aks/access-control-managed-azure-ad.md",
        "title": CloudNativeResources.controlClusterAccessConditionalAccessTitle,
        "description": CloudNativeResources.controlClusterAccessConditionalAccessDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-aks-docs/articles/aks/access-control-managed-azure-ad.md",
        "documentationUrl": "https://learn.microsoft.com/azure/aks/access-control-managed-azure-ad",
        "nextSteps": [],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "RESOURCE_GROUP",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "AKS_CLUSTER",
                    "title": CloudNativeResources.aksClusterNameParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "azure-aks-docs/articles/aks/concepts-network-azure-cni-pod-subnet.md",
        "title": CloudNativeResources.azureCniPodSubnetTitle,
        "description": CloudNativeResources.azureCniPodSubnetDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-aks-docs/articles/aks/concepts-network-azure-cni-pod-subnet.md",
        "documentationUrl": "https://learn.microsoft.com/azure/aks/concepts-network-azure-cni-pod-subnet",
        "nextSteps": [],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "CLUSTER_NAME",
                    "title": CloudNativeResources.clusterNameParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "RESOURCE_GROUP_NAME",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "azure-aks-docs/articles/aks/aks-migration.md",
        "title": CloudNativeResources.migrateToAksTitle,
        "description": CloudNativeResources.migrateToAksDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-aks-docs/articles/aks/aks-migration.md",
        "documentationUrl": "https://learn.microsoft.com/azure/aks/aks-migration",
        "nextSteps": [],
        "configurations": {
            "permissions": [],
            "configurableParams": []
        }
    },
    {
        "status": "active",
        "key": "azure-aks-docs/articles/aks/use-etags.md",
        "title": CloudNativeResources.enhancingConcurrencyETagsTitle,
        "description": CloudNativeResources.enhancingConcurrencyETagsDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-aks-docs/articles/aks/use-etags.md",
        "documentationUrl": "https://learn.microsoft.com/azure/aks/use-etags",
        "nextSteps": [],
        "configurations": {
            "permissions": [],
            "configurableParams": []
        }
    },
    {
        "status": "active",
        "key": "azure-aks-docs/articles/aks/istio-meshconfig.md",
        "title": CloudNativeResources.istioMeshConfigTitle,
        "description": CloudNativeResources.istioMeshConfigDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-aks-docs/articles/aks/istio-meshconfig.md",
        "documentationUrl": "https://learn.microsoft.com/azure/aks/istio-meshconfig",
        "nextSteps": [],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "CLUSTER",
                    "title": CloudNativeResources.clusterNameParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "RESOURCE_GROUP",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "azure-aks-docs/articles/aks/access-private-cluster.md",
        "title": CloudNativeResources.accessPrivateAksClusterTitle,
        "description": CloudNativeResources.accessPrivateAksClusterDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-aks-docs/articles/aks/access-private-cluster.md",
        "documentationUrl": "https://learn.microsoft.com/azure/aks/access-private-cluster",
        "nextSteps": [],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "AKS_RESOURCE_GROUP",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "AKS_CLUSTER_NAME",
                    "title": CloudNativeResources.aksClusterNameParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "AKSDNSLookupFailError/aksdns-lookup-fail-error.md",
        "title": CloudNativeResources.troubleshootDnsLookupFailTitle,
        "description": CloudNativeResources.troubleshootDnsLookupFailDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/AKSDNSLookupFailError/aksdns-lookup-fail-error.md",
        "documentationUrl": "",
        "nextSteps": [
            {
                "title": CloudNativeResources.privateAksCustomDnsNextStep,
                "url": "https://github.com/Azure/terraform/tree/master/quickstart/301-aks-private-cluster"
            }
        ],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "CLUSTER_NAME",
                    "title": CloudNativeResources.aksClusterNameParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "azure-aks-docs/articles/aks/resize-cluster.md",
        "title": CloudNativeResources.resizeAksClustersTitle,
        "description": CloudNativeResources.resizeAksClustersDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-aks-docs/articles/aks/resize-cluster.md",
        "documentationUrl": "https://learn.microsoft.com/azure/aks/resize-cluster",
        "nextSteps": [],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "RESOURCE_GROUP",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "AKS_CLUSTER_NAME",
                    "title": CloudNativeResources.aksClusterNameParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "NUM_NODES",
                    "title": CloudNativeResources.numberOfNodesParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "NODE_POOL_NAME",
                    "title": CloudNativeResources.nodePoolNameParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "azure-aks-docs/articles/aks/concepts-preview-api-life-cycle.md",
        "title": CloudNativeResources.aksPreviewApiLifecycleTitle,
        "description": CloudNativeResources.aksPreviewApiLifecycleDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-aks-docs/articles/aks/concepts-preview-api-life-cycle.md",
        "documentationUrl": "https://learn.microsoft.com/azure/aks/concepts-preview-api-life-cycle",
        "nextSteps": [
            {
                "title": CloudNativeResources.aksPreviewCliExtensionNextStep,
                "url": "https://github.com/Azure/azure-cli-extensions/tree/main/src/aks-preview"
            },
            {
                "title": CloudNativeResources.newerVersionSdkNextStep,
                "url": "https://azure.github.io/azure-sdk/releases/latest/index.html?search=containerservice"
            },
            {
                "title": CloudNativeResources.terraformProviderClientGoTitle,
                "url": "https://github.com/hashicorp/terraform-provider-azurerm/blob/main/internal/services/containers/client/client.go"
            }
        ],
        "configurations": {
            "permissions": [],
            "configurableParams": []
        }
    },
    {
        "status": "active",
        "key": "azure-aks-docs/articles/aks/use-labels.md",
        "title": CloudNativeResources.useLabelsInAksTitle,
        "description": CloudNativeResources.useLabelsInAksDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-aks-docs/articles/aks/use-labels.md",
        "documentationUrl": "https://learn.microsoft.com/azure/aks/use-labels",
        "nextSteps": [],
        "configurations": {
            "permissions": [],
            "configurableParams": []
        }
    },
    {
        "status": "active",
        "key": "SupportArticles-docs/support/azure/azure-kubernetes/extensions/aks-cost-analysis-add-on-issues.md",
        "title": CloudNativeResources.aksCostAnalysisIssuesTitle,
        "description": CloudNativeResources.aksCostAnalysisIssuesDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/SupportArticles-docs/support/azure/azure-kubernetes/extensions/aks-cost-analysis-add-on-issues.md",
        "documentationUrl": "https://learn.microsoft.com/troubleshoot/azure/azure-kubernetes/extensions/aks-cost-analysis-add-on-issues",
        "nextSteps": [],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "RESOURCE_GROUP",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "AKS_CLUSTER",
                    "title": CloudNativeResources.aksClusterNameParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "SupportArticles-docs/support/azure/azure-kubernetes/availability-performance/cluster-service-health-probe-mode-issues.md",
        "title": CloudNativeResources.troubleshootHealthProbeModeTitle,
        "description": CloudNativeResources.troubleshootHealthProbeModeDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/SupportArticles-docs/support/azure/azure-kubernetes/availability-performance/cluster-service-health-probe-mode-issues.md",
        "documentationUrl": "https://learn.microsoft.com/troubleshoot/azure/azure-kubernetes/availability-performance/cluster-service-health-probe-mode-issues",
        "nextSteps": [],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "RESOURCE_GROUP",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "AKS_CLUSTER_NAME",
                    "title": CloudNativeResources.aksClusterNameParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "SupportArticles-docs/support/azure/azure-kubernetes/create-upgrade-delete/error-code-cnidownloadtimeoutvmextensionerror.md",
        "title": CloudNativeResources.troubleshootCniDownloadFailuresTitle,
        "description": CloudNativeResources.troubleshootCniDownloadFailuresDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/SupportArticles-docs/support/azure/azure-kubernetes/create-upgrade-delete/error-code-cnidownloadtimeoutvmextensionerror.md",
        "documentationUrl": "https://learn.microsoft.com/troubleshoot/azure/azure-kubernetes/create-upgrade-delete/error-code-cnidownloadtimeoutvmextensionerror",
        "nextSteps": [],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "RESOURCE_GROUP",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "AKS_CLUSTER_NAME",
                    "title": CloudNativeResources.aksClusterNameParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "SupportArticles-docs/support/azure/azure-kubernetes/connectivity/tcp-timeouts-kubetctl-third-party-tools-connect-api-server.md",
        "title": CloudNativeResources.tcpTimeoutsKubectlApiTitle,
        "description": CloudNativeResources.tcpTimeoutsKubectlApiDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/SupportArticles-docs/support/azure/azure-kubernetes/connectivity/tcp-timeouts-kubetctl-third-party-tools-connect-api-server.md",
        "documentationUrl": "https://learn.microsoft.com/troubleshoot/azure/azure-kubernetes/connectivity/tcp-timeouts-kubetctl-third-party-tools-connect-api-server",
        "nextSteps": [],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "ResourceGroupName",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "AKSClusterName",
                    "title": CloudNativeResources.aksClusterNameParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "azure-aks-docs/articles/aks/enable-host-encryption.md",
        "title": CloudNativeResources.enableHostBasedEncryptionTitle,
        "description": CloudNativeResources.enableHostBasedEncryptionDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-aks-docs/articles/aks/enable-host-encryption.md",
        "documentationUrl": "https://learn.microsoft.com/azure/aks/enable-host-encryption",
        "nextSteps": [],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "MY_AKS_CLUSTER",
                    "title": CloudNativeResources.aksClusterNameParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "MY_RESOURCE_GROUP",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "azure-aks-docs/articles/aks/nat-gateway.md",
        "title": CloudNativeResources.natGatewayAksTitle,
        "description": CloudNativeResources.natGatewayAksDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-aks-docs/articles/aks/nat-gateway.md",
        "documentationUrl": "https://learn.microsoft.com/azure/aks/nat-gateway",
        "nextSteps": [],
        "configurations": {
            "permissions": [],
            "configurableParams": []
        }
    },
    {
        "status": "active",
        "key": "azure-aks-docs/articles/aks/free-standard-pricing-tiers.md",
        "title": CloudNativeResources.aksPricingTiersTitle,
        "description": CloudNativeResources.aksPricingTiersDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-aks-docs/articles/aks/free-standard-pricing-tiers.md",
        "documentationUrl": "https://learn.microsoft.com/azure/aks/free-standard-pricing-tiers",
        "nextSteps": [],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "RESOURCE_GROUP",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "CLUSTER_NAME",
                    "title": CloudNativeResources.aksClusterNameParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "azure-aks-docs/articles/aks/events.md",
        "title": CloudNativeResources.useKubernetesEventsTitle,
        "description": CloudNativeResources.useKubernetesEventsDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/azure-aks-docs/articles/aks/events.md",
        "documentationUrl": "https://learn.microsoft.com/azure/aks/events",
        "nextSteps": [],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "RESOURCE_GROUP",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "AKS_CLUSTER",
                    "title": CloudNativeResources.aksClusterNameParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "SupportArticles-docs/support/azure/azure-kubernetes/create-upgrade-delete/upgrading-or-scaling-does-not-succeed.md",
        "title": CloudNativeResources.troubleshootClusterUpgradingTitle,
        "description": CloudNativeResources.troubleshootClusterUpgradingDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/SupportArticles-docs/support/azure/azure-kubernetes/create-upgrade-delete/upgrading-or-scaling-does-not-succeed.md",
        "documentationUrl": "https://learn.microsoft.com/troubleshoot/azure/azure-kubernetes/create-upgrade-delete/upgrading-or-scaling-does-not-succeed",
        "nextSteps": [],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "RESOURCE_GROUP",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "CLUSTER_NAME",
                    "title": CloudNativeResources.clusterNameParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "SupportArticles-docs/support/azure/azure-kubernetes/connectivity/user-cannot-get-cluster-resources.md",
        "title": CloudNativeResources.troubleshootForbiddenErrorTitle,
        "description": CloudNativeResources.troubleshootForbiddenErrorDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/SupportArticles-docs/support/azure/azure-kubernetes/connectivity/user-cannot-get-cluster-resources.md",
        "documentationUrl": "https://learn.microsoft.com/troubleshoot/azure/azure-kubernetes/connectivity/user-cannot-get-cluster-resources",
        "nextSteps": [],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "RESOURCE_GROUP",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "CLUSTER_NAME",
                    "title": CloudNativeResources.clusterNameParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "SupportArticles-docs/support/azure/azure-kubernetes/connectivity/troubleshoot-cluster-connection-issues-api-server.md",
        "title": CloudNativeResources.troubleshootApiServerConnectionTitle,
        "description": CloudNativeResources.troubleshootApiServerConnectionDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/SupportArticles-docs/support/azure/azure-kubernetes/connectivity/troubleshoot-cluster-connection-issues-api-server.md",
        "documentationUrl": "https://learn.microsoft.com/troubleshoot/azure/azure-kubernetes/connectivity/troubleshoot-cluster-connection-issues-api-server",
        "nextSteps": [],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "RESOURCE_GROUP",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "AKS_CLUSTER",
                    "title": CloudNativeResources.aksClusterNameParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "SupportArticles-docs/support/azure/azure-kubernetes/connectivity/client-ip-address-cannot-access-api-server.md",
        "title": CloudNativeResources.clientIpCannotAccessApiTitle,
        "description": CloudNativeResources.clientIpCannotAccessApiDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/SupportArticles-docs/support/azure/azure-kubernetes/connectivity/client-ip-address-cannot-access-api-server.md",
        "documentationUrl": "https://learn.microsoft.com/troubleshoot/azure/azure-kubernetes/connectivity/client-ip-address-cannot-access-api-server",
        "nextSteps": [],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "RG_NAME",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "CLUSTER_NAME",
                    "title": CloudNativeResources.clusterNameParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "SupportArticles-docs/support/azure/azure-kubernetes/create-upgrade-delete/error-code-badrequest-or-invalidclientsecret.md",
        "title": CloudNativeResources.aadstsBadRequestErrorTitle,
        "description": CloudNativeResources.aadstsBadRequestErrorDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/SupportArticles-docs/support/azure/azure-kubernetes/create-upgrade-delete/error-code-badrequest-or-invalidclientsecret.md",
        "documentationUrl": "https://learn.microsoft.com/troubleshoot/azure/azure-kubernetes/create-upgrade-delete/error-code-badrequest-or-invalidclientsecret",
        "nextSteps": [],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "RESOURCE_GROUP_NAME",
                    "title": CloudNativeResources.resourceGroupParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "AKS_CLUSTER_NAME",
                    "title": CloudNativeResources.aksClusterNameParam,
                    "defaultValue": ""
                }
            ]
        }
    },
    {
        "status": "active",
        "key": "SupportArticles-docs/support/azure/azure-kubernetes/create-upgrade-delete/cannot-scale-cluster-autoscaler-enabled-node-pool.md",
        "title": CloudNativeResources.clusterAutoscalerFailsToScaleTitle,
        "description": CloudNativeResources.clusterAutoscalerFailsToScaleDescription,
        "stackDetails": [],
        "sourceUrl": "https://raw.githubusercontent.com/MicrosoftDocs/executable-docs/main/scenarios/SupportArticles-docs/support/azure/azure-kubernetes/create-upgrade-delete/cannot-scale-cluster-autoscaler-enabled-node-pool.md",
        "documentationUrl": "https://learn.microsoft.com/troubleshoot/azure/azure-kubernetes/create-upgrade-delete/cannot-scale-cluster-autoscaler-enabled-node-pool",
        "nextSteps": [],
        "configurations": {
            "permissions": [],
            "configurableParams": [
                {
                    "inputType": "textInput",
                    "commandKey": "AKS_RG_NAME",
                    "title": CloudNativeResources.aksResourceGroupNameParam,
                    "defaultValue": ""
                },
                {
                    "inputType": "textInput",
                    "commandKey": "AKS_CLUSTER_NAME",
                    "title": CloudNativeResources.aksClusterNameParam,
                    "defaultValue": ""
                }
            ]
        }
    }
];

export const metadata: MetaData[] = [...tutorialMetaData, ...deepLinkedScenarios];