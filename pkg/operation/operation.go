package operation

import (
	"bacon/pkg/porkbun"
	"encoding/json"
	"errors"
)

type Action string

const (
	Create Action = "create"
	Delete Action = "delete"
	Skip   Action = "none"
)

type Status string

const (
	Succeeded Status = "succeeded"
	Failed    Status = "failed"
	Planned   Status = "planned"
	Unchanged Status = "unchanged"
)

// Planned unit of work on a DNS record.
type RecordOperation struct {
	Action  Action
	DryRun  bool
	Record  porkbun.Record
	Execute func() error
}

type RecordOperationResult struct {
	Operation RecordOperation
	Err       error
}

func (result RecordOperationResult) Status() Status {
	if result.Operation.Action == Skip {
		return Unchanged
	}
	if result.Operation.DryRun {
		return Planned
	}
	if result.Err != nil {
		return Failed
	}
	return Succeeded
}

func (result RecordOperationResult) MarshalJSON() ([]byte, error) {
	payload := struct {
		Action Action         `json:"action"`
		Status Status         `json:"status"`
		Record porkbun.Record `json:"record"`
		Error  string         `json:"error,omitempty"`
	}{
		Action: result.Operation.Action,
		Status: result.Status(),
		Record: result.Operation.Record,
	}
	if result.Err != nil {
		payload.Error = result.Err.Error()
	}

	return json.Marshal(payload)
}

// Manager runs operations and collects their results.
type Manager struct {
	results []RecordOperationResult
}

func NewManager() *Manager {
	return &Manager{}
}

// Executes an operation and records the result
func (manager *Manager) Run(operation RecordOperation) RecordOperationResult {
	result := RecordOperationResult{Operation: operation}

	isSkip := operation.Action == Skip
	isDryRun := operation.DryRun
	if !isSkip && !isDryRun {
		if operation.Execute == nil {
			result.Err = errors.New("operation is missing an execute function")
		} else {
			result.Err = operation.Execute()
		}
	}

	manager.results = append(manager.results, result)
	return result
}

func (manager *Manager) Results() []RecordOperationResult {
	return manager.results
}

func (manager *Manager) Summary() Summary {
	summary := Summary{}
	for _, result := range manager.results {
		switch result.Status() {
		case Succeeded:
			if result.Operation.Action == Create {
				summary.Created++
			} else {
				summary.Deleted++
			}
		case Failed:
			summary.Failed++
		case Planned:
			if result.Operation.Action == Create {
				summary.PlannedCreates++
			} else {
				summary.PlannedDeletes++
			}
		case Unchanged:
			summary.Unchanged++
		}
	}

	return summary
}
