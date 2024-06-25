package common

import "github.com/Azure/InnovationEngine/internal/parsers"

// State for the codeblock in interactive mode. Used to keep track of the
// state of each codeblock.
type StatefulCodeBlock struct {
	CodeBlock       parsers.CodeBlock
	CodeBlockNumber int
	Error           error
	StdErr          string
	StdOut          string
	StepName        string
	StepNumber      int
	Success         bool
}

// Checks if a codeblock was executed by looking at the
// output, errors, and if success is true.
func (s StatefulCodeBlock) WasExecuted() bool {
	return s.StdOut != "" || s.StdErr != "" || s.Error != nil || s.Success
}
