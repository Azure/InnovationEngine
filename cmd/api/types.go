package main

type DeploymentStep struct {
	Name    string `json:"name"`
	Command string `json:"command"`
}

type DeploymentResponse struct {
	AgentWebsocketUrl string           `json:"agentWebsocketUrl"`
	Steps             []DeploymentStep `json:"steps"`
}
