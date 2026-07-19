package deployment

import (
	"bacon/pkg/porkbun"
	"fmt"
)

type OperationType string

const (
	Create OperationType = "create"
	Update OperationType = "update"
	Delete OperationType = "delete"
	Keep   OperationType = "keep"
)

type OperationStatus string

const (
	Planned OperationStatus = "planned"
	Success OperationStatus = "success"
	Failure OperationStatus = "failure"
)

type Operation struct {
	Type   OperationType
	Record porkbun.Record
}

type OperationResult struct {
	Status OperationStatus `json:"status"`
	Error  error           `json:"error,omitempty"`
	Type   OperationType   `json:"operationType"`
	Record porkbun.Record  `json:"record"`
}

func (op Operation) Execute(client *porkbun.Client, domain string) OperationResult {
	var result OperationResult
	result.Type = op.Type
	result.Record = op.Record

	switch op.Type {
	case Create:
		id, err := client.CreateRecord(domain, op.Record)
		if err != nil {
			result.Status = Failure
			result.Error = err
		} else {
			result.Status = Success
			// Set the newly created record's ID
			result.Record.Id = id
		}
	case Update:
		err := client.EditRecord(domain, op.Record)
		if err != nil {
			result.Status = Failure
			result.Error = err
		} else {
			result.Status = Success
		}
	case Delete:
		err := client.DeleteRecord(domain, op.Record)
		if err != nil {
			result.Status = Failure
			result.Error = err
		} else {
			result.Status = Success
		}
	case Keep:
		result.Status = Success
	default:
		result.Status = Failure
		result.Error = fmt.Errorf("unknown operation type: %s", op.Type)
	}
	return result
}

type Deployment struct {
	Operations []Operation
}

type DeploymentResult struct {
	Results []OperationResult
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
	return DeploymentResult{Results: results}
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
	return DeploymentResult{Results: results}
}
