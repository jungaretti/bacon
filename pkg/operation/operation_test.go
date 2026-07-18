package operation

import (
	"errors"
	"testing"
)

func TestRunExecutesRealOperation(t *testing.T) {
	manager := NewManager()

	calls := 0
	result := manager.Run(RecordOperation{
		Action:  Create,
		Execute: func() error { calls++; return nil },
	})

	if calls != 1 {
		t.Error("expected", 1, "actual", calls)
	}
	if result.Status() != Succeeded {
		t.Error("expected", Succeeded, "actual", result.Status())
	}
}

func TestRunDoesNotExecuteDryRun(t *testing.T) {
	manager := NewManager()

	calls := 0
	result := manager.Run(RecordOperation{
		Action:  Delete,
		DryRun:  true,
		Execute: func() error { calls++; return nil },
	})

	if calls != 0 {
		t.Error("expected", 0, "actual", calls)
	}
	if result.Status() != Planned {
		t.Error("expected", Planned, "actual", result.Status())
	}
}

func TestRunDoesNotExecuteSkip(t *testing.T) {
	manager := NewManager()

	result := manager.Run(RecordOperation{
		Action: Skip,
	})

	if result.Status() != Unchanged {
		t.Error("expected", Unchanged, "actual", result.Status())
	}
}

func TestRunCapturesFailure(t *testing.T) {
	manager := NewManager()

	failure := errors.New("api error")
	result := manager.Run(RecordOperation{
		Action:  Create,
		Execute: func() error { return failure },
	})

	if result.Status() != Failed {
		t.Error("expected", Failed, "actual", result.Status())
	}
	if result.Err != failure {
		t.Error("expected", failure, "actual", result.Err)
	}
}

func TestRunContinuesAfterFailure(t *testing.T) {
	manager := NewManager()

	manager.Run(RecordOperation{
		Action:  Delete,
		Execute: func() error { return errors.New("api error") },
	})

	calls := 0
	manager.Run(RecordOperation{
		Action:  Create,
		Execute: func() error { calls++; return nil },
	})

	if calls != 1 {
		t.Error("expected", 1, "actual", calls)
	}
	if len(manager.Results()) != 2 {
		t.Error("expected", 2, "actual", len(manager.Results()))
	}
}

func TestSummaryCounts(t *testing.T) {
	manager := NewManager()

	succeed := func() error { return nil }
	fail := func() error { return errors.New("api error") }

	manager.Run(RecordOperation{Action: Create, Execute: succeed})
	manager.Run(RecordOperation{Action: Delete, Execute: succeed})
	manager.Run(RecordOperation{Action: Create, DryRun: true, Execute: succeed})
	manager.Run(RecordOperation{Action: Delete, DryRun: true, Execute: succeed})
	manager.Run(RecordOperation{Action: Skip})
	manager.Run(RecordOperation{Action: Create, Execute: fail})

	summary := manager.Summary()
	expected := Summary{
		Created:        1,
		Deleted:        1,
		PlannedCreates: 1,
		PlannedDeletes: 1,
		Unchanged:      1,
		Failed:         1,
	}

	if summary != expected {
		t.Error("expected", expected, "actual", summary)
	}
}
