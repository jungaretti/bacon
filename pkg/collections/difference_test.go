package collections

import "testing"

const (
	hello = "hello"
	world = "world"
)

func TestSetDifferenceEmpty(t *testing.T) {
	from := []string{hello, world}
	to := []string{hello, world}

	diff := SetDifferenceByHash(from, to, func(thing string) string { return thing })

	if len(diff) != 0 {
		t.Error("expected", 0, "actual", len(diff))
	}
}

func TestSetDifferenceAdded(t *testing.T) {
	from := []string{hello}
	to := []string{hello, world}

	diff := SetDifferenceByHash(from, to, func(thing string) string { return thing })

	if len(diff) != 0 {
		t.Error("expected", 0, "actual", len(diff))
	}
}

func TestSetDifferenceRemoved(t *testing.T) {
	from := []string{hello, world}
	to := []string{hello}

	diff := SetDifferenceByHash(from, to, func(thing string) string { return thing })

	if len(diff) != 1 {
		t.Error("expected", 1, "actual", len(diff))
	}
	if diff[0] != world {
		t.Error("expected", world, "actual", diff[0])
	}
}

func TestAddedRemovedEmpty(t *testing.T) {
	from := []string{hello, world}
	to := []string{hello, world}

	added, removed := AddedRemovedByHash(from, to, func(thing string) string { return thing })

	if len(added) != 0 {
		t.Error("expected", 0, "actual", len(added))
	}
	if len(removed) != 0 {
		t.Error("expected", 0, "actual", len(added))
	}
}

func TestAddedRemovedAdd(t *testing.T) {
	from := []string{hello}
	to := []string{hello, world}

	added, removed := AddedRemovedByHash(from, to, func(thing string) string { return thing })

	if len(added) != 1 {
		t.Error("expected", 0, "actual", len(added))
	}
	if len(removed) != 0 {
		t.Error("expected", 0, "actual", len(added))
	}
}

func TestAddedRemovedRemove(t *testing.T) {
	from := []string{hello, world}
	to := []string{hello}

	added, removed := AddedRemovedByHash(from, to, func(thing string) string { return thing })

	if len(added) != 0 {
		t.Error("expected", 0, "actual", len(added))
	}
	if len(removed) != 1 {
		t.Error("expected", 0, "actual", len(added))
	}
}
