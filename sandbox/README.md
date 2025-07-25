# Sandbox

The `sandbox` directory is dedicated to experimental code, prototypes, and early-stage integrations for the Innovation Engine project. It is intended as a safe space for trying out new ideas and features before they are considered for inclusion in the main codebase. In particular it includes experiments that focus on integration of Innovation Engine with Headlamp, a UI for Kubernetes.

## What is Headlamp?

[Headlamp](https://headlamp.dev/) is a modern, extensible Kubernetes UI that allows users to manage and visualize their Kubernetes clusters. It supports a plugin system, enabling developers to extend its functionality with custom features and integrations.

### Running Headlamp 

A complete copy of Headlamp Desktop application is included in the sandbox directory. This is included for convenience and to ensure that we are developing against a consistent version of Headlamp. However, you should be able to run these plugins in any recent version of Headlamp (if not please file a bug).

#### Install Docker

1. [Install the binaries](https://docs.docker.com/engine/install/ubuntu/)
2. Configure the docker user:
```bash
sudo groupadd docker
sudo usermod -aG docker $(whoami)
newgrp docker
``` 

#### Installing Headlamp

1. Download the [latest Headlamp release](https://github.com/kubernetes-sigs/headlamp/releases) as a tarball
2. Install from the tarball
```bash
tar xvzf ./Headlamp-0.30.0-linux-x64.tar.gz
```
3. If using Linux prepare the executable:
```bash
cd Headlamp-0.30.0-linux-x64
sudo chown root headlamp
sudo chmod 4755 headlamp
```
4. Run Headlamp
```
./headlamp
```

#### Install Minikube plugin for local K8s

1. Locate Minikube in the Headlamp plugin catalog
2. Click install
3. Restart Headlamp
4. Create a minikube cluster in Headlamp `Home -> Load Cluster -> Minikube Add`

## Innovation Engine Headlamp Plugin

The `innovation-engine-headlamp` subfolder contains a Headlamp plugin that integrates the Innovation Engine with the Headlamp UI. This plugin provides:

- A sidebar entry for accessing Innovation Engine features directly from Headlamp.
- A "Getting Started" page to help users begin using the Innovation Engine.
- A shell execution interface for running allowlisted commands (such as `ie execute ...`) and viewing their output within the Headlamp UI.

For more details on the plugin's features and how to run it, see the [README](innovation-engine-headlamp/README.md).