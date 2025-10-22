# Prerequisites: Installing Terraform on Linux

This document provides instructions to download and install Terraform on a Linux system using a bash script.

## Install Terraform

Define the version of Terraform to install

```bash
export TERRAFORM_VERSION="1.10.5"
```

Download, install and configure Terraform

```bash
if ! command -v terraform &> /dev/null
then
    curl -O https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip

    mkdir -p ~/bin
    unzip -j terraform_${TERRAFORM_VERSION}_linux_amd64.zip terraform -d ~/bin
    
    if [[ ":$PATH:" != *":$HOME/bin:"* ]]; then
        export PATH="$HOME/bin:$PATH"
        echo 'export PATH="$HOME/bin:$PATH"' >> ~/.bashrc
    fi
    
    terraform -v

    rm terraform_${TERRAFORM_VERSION}_linux_amd64.zip

    echo "Terraform has been installed"
else
    echo "Terraform is already installed"
fi
```

