package test

// Represents the structure of

type TestedCodeBlocks struct {
	CodeBlock       string  `json:"codeBlock"`
	ExpectedOutput  string  `json:"expectedOutput"`
	ActualOutput    string  `json:"actualOutput"`
	ComparisonScore float64 `json:"score"`
	Success         bool    `json:"success"`
	Error           string  `json:"error"`
}

type TestedStep struct {
	Header      string             `json:"header"`
	Description string             `json:"description"`
	CodeBlocks  []TestedCodeBlocks `json:"codeBlocks"`
}

// A generated report from the execution of `ie test` on a markdown document.
type Report struct {
	Name                 string            `json:"name"`
	Properties           map[string]string `json:"properties"`
	EnvironmentVariables map[string]string `json:"environmentVariables"`
	Success              bool              `json:"success"`
	Error                string            `json:"error"`
	FailedAt             int               `json:"failedAt"`
	Steps                []TestedStep      `json:"steps"`
}

// TODO(vmarcella): Build out the rest of the test reporting JSON.
func BuildReport() Report {
	return Report{
		Name:                 "",
		Properties:           make(map[string]string),
		EnvironmentVariables: make(map[string]string),
	}
}
