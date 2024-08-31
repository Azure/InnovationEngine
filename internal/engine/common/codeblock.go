package common

import "github.com/Azure/InnovationEngine/internal/parsers"

const (
	STATUS_SUCCESS = "success"
	STATUS_FAILURE = "failure"
	STATUS_PENDING = "pending"
)

// State for the codeblock in interactive mode. Used to keep track of the
// state of each codeblock.
type StatefulCodeBlock struct {
	CodeBlock       parsers.CodeBlock `json:"codeBlock"`
	CodeBlockNumber int               `json:"codeBlockNumber"`
	Error           error             `json:"error"`
	StdErr          string            `json:"stdErr"`
	StdOut          string            `json:"stdOut"`
	StepName        string            `json:"stepName"`
	StepNumber      int               `json:"stepNumber"`
	Status          string            `json:"success"`
	SimilarityScore float64           `json:"similarityScore"`
}

// Checks if a codeblock was executed by looking at the
// output, errors, and if success is true.
func (s StatefulCodeBlock) WasExecuted() bool {
	return s.Status != STATUS_PENDING
}

func (s StatefulCodeBlock) Succeeded() bool {
	return s.Status == STATUS_SUCCESS
}

func (s StatefulCodeBlock) Failed() bool {
	return s.Status == STATUS_FAILURE
}
