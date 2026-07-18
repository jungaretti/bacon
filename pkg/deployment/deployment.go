package deployment

import (
	"bacon/pkg/porkbun"
	"fmt"
	"io"
)

type Action string

const (
	Delete Action = "delete"
	Update Action = "update"
	Create Action = "create"
	Keep   Action = "keep"
)

type Operation struct {
	Action Action
	DryRun bool
	Record porkbun.Record
}

type DeploymentRenderer interface {
	// Preview operations before they are applied.
	Preview(ops []Operation) error
	// Report results of completed operations.
	Report(ops []Operation) error
}

func NewDeploymentRenderer(format string, writer io.Writer) (DeploymentRenderer, error) {
	switch format {
	case "table":
		return &tableReporter{writer: writer}, nil
	case "json":
		return &jsonReporter{writer: writer}, nil
	}

	return nil, fmt.Errorf("unknown output format: %v", format)
}
