# Prerequisites: Installing Terraform on Linux

This document provides instructions to download and install Terraform on a Linux system using a bash script.

## Install Terraform

Define the version of Terraform to install

```bash
export TERRAFORM_VERSION="1.10.5"
```

Download Terraform

```bash
curl -O https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip
```

Unzip the downloaded file

```bash
unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip
```

Move the Terraform binary to /usr/local/bin

```bash
mkdir -p ~/bin
mv terraform ~/bin/
```

Verify the installation

```bash
terraform -v
```

Cleanup

```bash
rm terraform_${TERRAFORM_VERSION}_linux_amd64.zip
rm LICENSE
```

