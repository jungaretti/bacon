package deployment

import (
	"bacon/pkg/config"
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
	Error  string          `json:"error,omitempty"`
	Type   OperationType   `json:"operationType"`
	Record config.Record   `json:"record"`
}

func (op Operation) Execute(client *porkbun.Client, domain string) OperationResult {
	var result OperationResult
	result.Type = op.Type

	record, err := config.RecordFromPorkbun(op.Record)
	if err != nil {
		result.Status = Failure
		result.Error = err.Error()
		return result
	}
	result.Record = record

	switch op.Type {
	case Create:
		_, err := client.CreateRecord(domain, op.Record)
		if err != nil {
			result.Status = Failure
			result.Error = err.Error()
		} else {
			result.Status = Success
		}
	case Update:
		err := client.EditRecord(domain, op.Record)
		if err != nil {
			result.Status = Failure
			result.Error = err.Error()
		} else {
			result.Status = Success
		}
	case Delete:
		err := client.DeleteRecord(domain, op.Record)
		if err != nil {
			result.Status = Failure
			result.Error = err.Error()
		} else {
			result.Status = Success
		}
	case Keep:
		result.Status = Success
	default:
		result.Status = Failure
		result.Error = fmt.Sprintf("unknown operation type: %s", op.Type)
	}
	return result
}
