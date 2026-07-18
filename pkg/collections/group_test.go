package collections

import (
	"strings"
	"testing"
)

func hashFirstTerm(thing string) string {
	first, _, _ := strings.Cut(thing, "-")
	return first
}

func TestGroupByHashAllGrouped(t *testing.T) {
	from := []string{"hello-world", "other-thing"}
	to := []string{"hello-dog", "other-stuff"}

	groups, ungroupedFrom, ungroupedTo := GroupByHash(from, to, hashFirstTerm)

	if len(groups) != 2 {
		t.Error("expected", 2, "actual", len(groups))
	}

	if len(ungroupedFrom) != 0 {
		t.Error("expected", 0, "actual", len(ungroupedFrom))
	}

	if len(ungroupedTo) != 0 {
		t.Error("expected", 0, "actual", len(ungroupedTo))
	}
}

func TestGroupByHashSingleGroup(t *testing.T) {
	from := []string{"hello-world", "other-thing"}
	to := []string{"hello-dog", "goose-egg"}

	groups, ungroupedFrom, ungroupedTo := GroupByHash(from, to, hashFirstTerm)

	if len(groups) != 1 {
		t.Error("expected", 1, "actual", len(groups))
	}
	if groups[0].From != "hello-world" || groups[0].To != "hello-dog" {
		t.Error("expected", "hello-world/hello-dog", "actual", groups[0].From+"/"+groups[0].To)
	}

	if len(ungroupedFrom) != 1 {
		t.Error("expected", 1, "actual", len(ungroupedFrom))
	}
	if ungroupedFrom[0] != "other-thing" {
		t.Error("expected", "other-thing", "actual", ungroupedFrom[0])
	}

	if len(ungroupedTo) != 1 {
		t.Error("expected", 1, "actual", len(ungroupedTo))
	}
	if ungroupedTo[0] != "goose-egg" {
		t.Error("expected", "goose-egg", "actual", ungroupedTo[0])
	}
}

func TestGroupByHashMultipleGroups(t *testing.T) {
	from := []string{"hello-world", "hello-dog", "hello-cat"}
	to := []string{"hello-moon", "hello-star"}

	groups, ungroupedFrom, ungroupedTo := GroupByHash(from, to, hashFirstTerm)

	if len(groups) != 2 {
		t.Error("expected", 2, "actual", len(groups))
	}
	if groups[0].From != "hello-world" || groups[0].To != "hello-moon" {
		t.Error("expected", "hello-world/hello-moon", "actual", groups[0].From+"/"+groups[0].To)
	}
	if groups[1].From != "hello-dog" || groups[1].To != "hello-star" {
		t.Error("expected", "hello-dog/hello-star", "actual", groups[1].From+"/"+groups[1].To)
	}

	if len(ungroupedFrom) != 1 {
		t.Error("expected", 1, "actual", len(ungroupedFrom))
	}
	if ungroupedFrom[0] != "hello-cat" {
		t.Error("expected", "hello-cat", "actual", ungroupedFrom[0])
	}

	if len(ungroupedTo) != 0 {
		t.Error("expected", 0, "actual", len(ungroupedTo))
	}
}
