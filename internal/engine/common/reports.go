package common

// Reports are summaries of the execution of a markdown document. They will
// include the name of the document, the properties found in the yaml header of
// the document, the environment variables set by the document, and general
// information about the execution.
type Report struct {
	Name                 string                 `json:"name"`
	Properties           map[string]interface{} `json:"properties"`
	EnvironmentVariables map[string]string      `json:"environmentVariables"`
	Success              bool                   `json:"success"`
	Error                string                 `json:"error"`
	FailedAt             int                    `json:"failedAt"`
	codeBlocks           []StatefulCodeBlock    `json:"codeBlocks"`
}

func (report *Report) WithProperties(properties map[string]interface{}) *Report {
	report.Properties = properties
	return report
}

func (report *Report) WithEnvironmentVariables(envVars map[string]string) *Report {
	report.EnvironmentVariables = envVars
	return report
}

func (report *Report) WithCodeBlocks(codeBlocks []StatefulCodeBlock) *Report {
	report.codeBlocks = codeBlocks
	return report
}

// TODO(vmarcella): Implement this to write the report to JSON.
func (report *Report) WriteToJSONFile(outputPath string) error {
	return nil
}

func BuildReport(name string) Report {
	return Report{
		Name:                 name,
		Properties:           make(map[string]interface{}),
		EnvironmentVariables: make(map[string]string),
	}
}
