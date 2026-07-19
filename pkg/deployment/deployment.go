package deployment

import (
	"bacon/pkg/porkbun"
)

type Deployment struct {
	Operations []Operation
}

type DeploymentSummary struct {
	DryRun          bool                  `json:"dryRun"`
	OperationCounts map[OperationType]int `json:"operationCounts"`
}

type DeploymentResult struct {
	Summary DeploymentSummary `json:"summary"`
	Results []OperationResult `json:"results"`
}

func NewDeployment(added, removed, updated, unchanged []porkbun.Record) Deployment {
	var operations []Operation
	for _, record := range removed {
		operations = append(operations, Operation{Type: Delete, Record: record})
	}
	for _, record := range updated {
		operations = append(operations, Operation{Type: Update, Record: record})
	}
	for _, record := range added {
		operations = append(operations, Operation{Type: Create, Record: record})
	}
	for _, record := range unchanged {
		operations = append(operations, Operation{Type: Keep, Record: record})
	}
	return Deployment{Operations: operations}
}

func (deployment Deployment) Preview() DeploymentResult {
	var results []OperationResult
	for _, operation := range deployment.Operations {
		results = append(results, OperationResult{
			Status: Planned,
			Type:   operation.Type,
			Record: operation.Record,
		})
	}
	return newDeploymentResult(true, results)
}

func (deployment Deployment) Execute(client *porkbun.Client, domain string) DeploymentResult {
	var results []OperationResult
	for _, operation := range deployment.Operations {
		result := operation.Execute(client, domain)
		results = append(results, result)
		if result.Status == Failure {
			break
		}
	}
	return newDeploymentResult(false, results)
}

func newDeploymentResult(dryRun bool, results []OperationResult) DeploymentResult {
	operationCounts := make(map[OperationType]int)
	for _, result := range results {
		operationCounts[result.Type]++
	}

	return DeploymentResult{
		Summary: DeploymentSummary{
			DryRun:          dryRun,
			OperationCounts: operationCounts,
		},
		Results: results,
	}
}
