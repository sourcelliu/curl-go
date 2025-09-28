package tool

import (
	"reflect"
	"testing"
)

func TestStringList(t *testing.T) {
	t.Run("new list is empty", func(t *testing.T) {
		list := NewStringList()
		if list == nil {
			t.Fatal("NewStringList() returned nil")
		}
		if len(list.Strings()) != 0 {
			t.Errorf("New list should be empty, but has length %d", len(list.Strings()))
		}
	})

	t.Run("append adds items correctly", func(t *testing.T) {
		list := NewStringList()

		// Append first item
		list.Append("first")
		expected1 := []string{"first"}
		if !reflect.DeepEqual(list.Strings(), expected1) {
			t.Errorf("After 1 append, list = %v; want %v", list.Strings(), expected1)
		}

		// Append second item
		list.Append("second")
		expected2 := []string{"first", "second"}
		if !reflect.DeepEqual(list.Strings(), expected2) {
			t.Errorf("After 2 appends, list = %v; want %v", list.Strings(), expected2)
		}

		// Append third item
		list.Append("") // Append empty string
		expected3 := []string{"first", "second", ""}
		if !reflect.DeepEqual(list.Strings(), expected3) {
			t.Errorf("After 3 appends, list = %v; want %v", list.Strings(), expected3)
		}
	})
}