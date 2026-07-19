package deployment

type Formatter interface {
	FormatStart(domain string, dryRun bool) string
	FormatResult(deploymentResult DeploymentResult) string
}
