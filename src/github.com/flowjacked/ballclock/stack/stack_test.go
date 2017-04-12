package stack

import (
	"testing"
)

// TestPushStack verifies we're putting ints in the stack in the right order
func TestPushStack(test *testing.T) {
	s := NewStack(4)
	s.Push(1)
	s.Push(2)
	s.Push(3)
	s.Push(4)
	for i := len(s.s); i > 0; i-- {
		expected := i
		actual := s.s[i-1]
		if actual != expected {
			test.Error("Test failed: expected:", expected, "actual:", actual)
		}
	}
}

// TestPopStack verifies that after putting ints in the stack, we're getting
// them out in the right order
func TestPopStack(test *testing.T) {
	s := NewStack(4)
	s.Push(1)
	s.Push(2)
	s.Push(3)
	s.Push(4)
	for i := len(s.s); i > 0; i-- {
		expected := i
		actual, _ := s.Pop()
		if actual != expected {
			test.Error("Test failed: expected:", expected, "actual:", actual)
		}
	}
}

func TestGetStack(test *testing.T) {
	s := NewStack(4)
	s.Push(1)
	s.Push(2)
	s.Push(3)
	s.Push(4)
	actual := s.GetStack()
	//test.Error(actual)
	expected := []int{1, 2, 3, 4}
	for i := range expected {
		if actual[i] != expected[i] {
			test.Error("Test failed: expected:", expected, "actual:", actual)
		}
	}
}

func TestStackFull(test *testing.T) {
	s := NewStack(4)
	s.Push(1)
	s.Push(2)
	s.Push(3)
	s.Push(4)
	err := s.Push(5)
	if err == nil {
		test.Error("Test failed: expected an error for pushing past stack limit")
	}
}

func TestStackEmpty(test *testing.T) {
	s := NewStack(4)
	s.Push(1)
	s.Push(2)
	s.Push(3)
	s.Push(4)
	var err error
	for i := 5; i >= 0; i-- {
		_, err = s.Pop()
	}
	if err == nil {
		test.Error("Test failed: expected an error for trying to pop a value off the stack when it's empty")
	}
}
