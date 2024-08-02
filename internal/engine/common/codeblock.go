package common

import "github.com/Azure/InnovationEngine/internal/parsers"

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
	Success         bool              `json:"success"`
	SimilarityScore float64           `json:"similarityScore"`
}

// Checks if a codeblock was executed by looking at the
// output, errors, and if success is true.
func (s StatefulCodeBlock) WasExecuted() bool {
	return s.StdOut != "" || s.StdErr != "" || s.Error != nil || s.Success
}
