package collections

import (
	"slices"
	"testing"
)

func TestSetDifferenceEmpty(t *testing.T) {
	from := []string{"hello", "world"}
	to := []string{"hello", "world"}

	diff := SetDifferenceByHash(from, to, func(thing string) string { return thing })

	if len(diff) != 0 {
		t.Error("expected", 0, "actual", len(diff))
	}
}

func TestSetDifferenceAdded(t *testing.T) {
	from := []string{"hello"}
	to := []string{"hello", "world"}

	diff := SetDifferenceByHash(from, to, func(thing string) string { return thing })

	if len(diff) != 0 {
		t.Error("expected", 0, "actual", len(diff))
	}
}

func TestSetDifferenceRemoved(t *testing.T) {
	from := []string{"hello", "world"}
	to := []string{"hello"}

	diff := SetDifferenceByHash(from, to, func(thing string) string { return thing })

	if len(diff) != 1 {
		t.Error("expected", 1, "actual", len(diff))
	}
	if diff[0] != "world" {
		t.Error("expected", "world", "actual", diff[0])
	}
}

func TestSetIntersectionEmpty(t *testing.T) {
	from := []string{"hello", "world"}
	to := []string{"hello", "world"}

	intersection := SetIntersectionByHash(from, to, func(thing string) string { return thing })

	if len(intersection) != 2 {
		t.Error("expected", 2, "actual", len(intersection))
	}
	if slices.Contains(intersection, "hello") == false {
		t.Error("expected", true, "actual", false)
	}
	if slices.Contains(intersection, "world") == false {
		t.Error("expected", true, "actual", false)
	}
}

func TestSetIntersectionAdded(t *testing.T) {
	from := []string{"hello"}
	to := []string{"hello", "world"}

	intersection := SetIntersectionByHash(from, to, func(thing string) string { return thing })

	if len(intersection) != 1 {
		t.Error("expected", 1, "actual", len(intersection))
	}
	if intersection[0] != "hello" {
		t.Error("expected", "hello", "actual", intersection[0])
	}
}

func TestSetIntersectionRemoved(t *testing.T) {
	from := []string{"hello", "world"}
	to := []string{"hello"}

	intersection := SetIntersectionByHash(from, to, func(thing string) string { return thing })

	if len(intersection) != 1 {
		t.Error("expected", 1, "actual", len(intersection))
	}
	if intersection[0] != "hello" {
		t.Error("expected", "hello", "actual", intersection[0])
	}
}

func TestAddedRemovedUnchangedByHashEmpty(t *testing.T) {
	from := []string{"hello", "world"}
	to := []string{"hello", "world"}

	added, removed, unchanged := AddedRemovedUnchangedByHash(from, to, func(thing string) string { return thing })

	if len(added) != 0 {
		t.Error("expected", 0, "actual", len(added))
	}

	if len(removed) != 0 {
		t.Error("expected", 0, "actual", len(removed))
	}

	if len(unchanged) != 2 {
		t.Error("expected", 2, "actual", len(unchanged))
	}
	if slices.Contains(unchanged, "hello") == false {
		t.Error("expected", true, "actual", false)
	}
	if slices.Contains(unchanged, "world") == false {
		t.Error("expected", true, "actual", false)
	}
}

func TestAddedRemovedUnchangedByHashAdd(t *testing.T) {
	from := []string{"hello"}
	to := []string{"hello", "world"}

	added, removed, unchanged := AddedRemovedUnchangedByHash(from, to, func(thing string) string { return thing })

	if len(added) != 1 {
		t.Error("expected", 1, "actual", len(added))
	}
	if slices.Contains(added, "world") == false {
		t.Error("expected", true, "actual", false)
	}

	if len(removed) != 0 {
		t.Error("expected", 0, "actual", len(removed))
	}

	if len(unchanged) != 1 {
		t.Error("expected", 1, "actual", len(unchanged))
	}
}

func TestAddedRemovedUnchangedByHashRemove(t *testing.T) {
	from := []string{"hello", "world"}
	to := []string{"hello"}

	added, removed, unchanged := AddedRemovedUnchangedByHash(from, to, func(thing string) string { return thing })

	if len(added) != 0 {
		t.Error("expected", 0, "actual", len(added))
	}

	if len(removed) != 1 {
		t.Error("expected", 1, "actual", len(removed))
	}
	if slices.Contains(removed, "world") == false {
		t.Error("expected", true, "actual", false)
	}

	if len(unchanged) != 1 {
		t.Error("expected", 1, "actual", len(unchanged))
	}
	if slices.Contains(unchanged, "hello") == false {
		t.Error("expected", true, "actual", false)
	}
}
