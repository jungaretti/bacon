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

type RecordOperation struct {
	Type   OperationType
	Record porkbun.Record
}

type RecordOperationResult struct {
	Status OperationStatus `json:"status"`
	Error  error           `json:"error,omitempty"`
	Type   OperationType   `json:"operationType"`
	Record porkbun.Record  `json:"record"`
}

func (op RecordOperation) Execute(client *porkbun.Client, domain string) RecordOperationResult {
	var result RecordOperationResult
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

type RecordDeployment struct {
	Operations []RecordOperation
}

func NewRecordDeployment(added, removed, updated, unchanged []porkbun.Record) RecordDeployment {
	var operations []RecordOperation
	for _, record := range removed {
		operations = append(operations, RecordOperation{Type: Delete, Record: record})
	}
	for _, record := range updated {
		operations = append(operations, RecordOperation{Type: Update, Record: record})
	}
	for _, record := range added {
		operations = append(operations, RecordOperation{Type: Create, Record: record})
	}
	for _, record := range unchanged {
		operations = append(operations, RecordOperation{Type: Keep, Record: record})
	}
	return RecordDeployment{Operations: operations}
}

func (deployment RecordDeployment) Preview() []RecordOperationResult {
	var results []RecordOperationResult
	for _, operation := range deployment.Operations {
		results = append(results, RecordOperationResult{
			Status: Planned,
			Type:   operation.Type,
			Record: operation.Record,
		})
	}
	return results
}

func (deployment RecordDeployment) Execute(client *porkbun.Client, domain string) []RecordOperationResult {
	var results []RecordOperationResult
	for _, operation := range deployment.Operations {
		result := operation.Execute(client, domain)
		results = append(results, result)
		if result.Status == Failure {
			break
		}
	}
	return results
}
