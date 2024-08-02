package common

import (
	"encoding/json"
	"os"

	"github.com/Azure/InnovationEngine/internal/logging"
)

type Report struct {
	Name                 string                 `json:"name"`
	Properties           map[string]interface{} `json:"properties"`
	EnvironmentVariables map[string]string      `json:"environmentVariables"`
	Success              bool                   `json:"success"`
	Error                string                 `json:"error"`
	FailedAtStep         int                    `json:"failedAtStep"`
	CodeBlocks           []StatefulCodeBlock    `json:"steps"`
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
	report.CodeBlocks = codeBlocks
	return report
}

func (report *Report) WithError(err error) *Report {
	if err == nil {
		return report
	}

	report.Error = err.Error()
	report.Success = false
	return report
}

// TODO(vmarcella): Implement this to write the report to JSON.
func (report *Report) WriteToJSONFile(outputPath string) error {
	jsonReport, err := json.MarshalIndent(report, "", "    ")
	if err != nil {
		return err
	}
	logging.GlobalLogger.Infof("Generated the test report:\n %s", jsonReport)

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write(jsonReport)
	if err != nil {
		return err
	}

	logging.GlobalLogger.Infof("Wrote the test report to %s", outputPath)

	return nil
}

func BuildReport(name string) Report {
	return Report{
		Name:                 name,
		Properties:           make(map[string]interface{}),
		EnvironmentVariables: make(map[string]string),
		Success:              true,
		Error:                "",
		FailedAtStep:         -1,
	}
}
