package deployment

type Formatter interface {
	Format(deploymentResult DeploymentResult) string
}
