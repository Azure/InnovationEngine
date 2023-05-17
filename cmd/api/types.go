package main

import "github.com/Azure/InnovationEngine/internal/engine"

type ExecuteResponse struct {
	RunnerID string        `json:"runnerID"`
	Steps    []engine.Step `json:"steps"`
}
