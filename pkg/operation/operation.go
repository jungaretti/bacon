package operation

import (
	"bacon/pkg/porkbun"
	"encoding/json"
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

func (r RecordOperationResult) Status() Status {
	if r.Operation.Action == Skip {
		return Unchanged
	}
	if r.Operation.DryRun {
		return Planned
	}
	if r.Err != nil {
		return Failed
	}
	return Succeeded
}

func (r RecordOperationResult) MarshalJSON() ([]byte, error) {
	result := struct {
		Action Action         `json:"action"`
		Status Status         `json:"status,omitempty"`
		Record porkbun.Record `json:"record"`
		Error  string         `json:"error,omitempty"`
	}{
		Action: r.Operation.Action,
		Record: r.Operation.Record,
		Status: r.Status(),
	}
	if r.Err != nil {
		result.Error = r.Err.Error()
	}

	return json.Marshal(result)
}

// Manager runs operations and collects their results.
type Manager struct {
	results []RecordOperationResult
}

func NewManager() *Manager {
	return &Manager{}
}

// Run executes an operation, unless it is a skip or a dry run, and records
// the result.
func (m *Manager) Run(op RecordOperation) RecordOperationResult {
	result := RecordOperationResult{Operation: op}
	if op.Action != Skip && !op.DryRun && op.Execute != nil {
		result.Err = op.Execute()
	}

	m.results = append(m.results, result)
	return result
}

func (m *Manager) Results() []RecordOperationResult {
	return m.results
}

func (m *Manager) Summary() Summary {
	summary := Summary{}
	for _, result := range m.results {
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
