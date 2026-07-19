package deployment

import "encoding/json"

type JSONFormatter struct{}

var _ Formatter = JSONFormatter{}

func (formatter JSONFormatter) Format(deploymentResult DeploymentResult) string {
	output, err := json.MarshalIndent(deploymentResult, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(output) + "\n"
}
