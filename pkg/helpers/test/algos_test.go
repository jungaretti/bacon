package helpers

import (
	"bacon/pkg/helpers"
	"fmt"
	"testing"
)

func TestDifference(t *testing.T) {
	a := []int{1, 2, 3, 4, 5, 6, 7}
	b := []int{6, 7, 8, 9}

	expected := []int{1, 2, 3, 4, 5}
	actual := helpers.Difference(a, b)

	for index, value := range expected {
		if value != actual[index] {
			t.Log("expected", value, "found", actual[index])
			t.Fail()
		}
	}
}

func TestDifferenceEmpty(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{7, 8, 9}

	expected := []int{}
	actual := helpers.Difference(a, b)

	for index, value := range expected {
		if value != actual[index] {
			t.Log("expected", value, "found", actual[index])
			t.Fail()
		}
	}
}

func TestDifferenceByHasher(t *testing.T) {
	a := []int{1, 2, 3, 4, 5, 6, 7}
	b := []int{6, 7, 8, 9}

	expected := []int{1, 2, 3, 4, 5}
	actual := helpers.DifferenceByHasher(a, b, func(t int) string { return fmt.Sprint(t) })

	for index, value := range expected {
		if value != actual[index] {
			t.Log("expected", value, "found", actual[index])
			t.Fail()
		}
	}
}

func TestDifferenceByHasherEmpty(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{7, 8, 9}

	expected := []int{}
	actual := helpers.DifferenceByHasher(a, b, func(t int) string { return fmt.Sprint(t) })

	for index, value := range expected {
		if value != actual[index] {
			t.Log("expected", value, "found", actual[index])
			t.Fail()
		}
	}
}
