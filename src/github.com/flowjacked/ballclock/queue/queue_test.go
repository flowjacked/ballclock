package queue

import (
	"testing"
)

// TestPushQueue verifieq we're putting ints in the queue in the right order
func TestPushQueue(test *testing.T) {
	q := NewQueue(4)
	q.Push(1)
	q.Push(2)
	q.Push(3)
	q.Push(4)
	for i := range q.q {
		expected := i + 1
		actual := q.q[i]
		if actual != expected {
			test.Error("Test failed: expected:", expected, "actual:", actual)
		}
	}
}

// TestPopQueue verifieq that after putting ints in the queue, we're getting
// them out in the right order
func TestPopQueue(test *testing.T) {
	q := NewQueue(4)
	q.Push(1)
	q.Push(2)
	q.Push(3)
	q.Push(4)
	for i := 0; i < len(q.q); i++ {
		expected := i + 1
		actual, _ := q.Pop()
		if actual != expected {
			test.Error("Test failed: expected:", expected, "actual:", actual)
		}
	}
}

// TestSaveState verifieq the state is getting saved correctly
func TestSaveState(test *testing.T) {
	q := NewQueue(4)
	q.Push(1)
	q.Push(2)
	q.Push(3)
	q.Push(4)
	q.SaveState()
	actual := q.origin
	expected := "[1 2 3 4]"
	if actual != expected {
		test.Error("Test failed: expected:", expected, "actual:", actual)
	}

}

// TestEqualsOrigin verifieq that the save state will equal itself when called
func TestEqualsOrigin(test *testing.T) {
	q := NewQueue(4)
	q.Push(1)
	q.Push(2)
	q.Push(3)
	q.Push(4)
	q.SaveState()
	if !q.EqualsOrigin() {
		test.Error("Test failed: saved state != current state and it should")
	}
}

func TestGetQueue(test *testing.T) {
	q := NewQueue(4)
	q.Push(1)
	q.Push(2)
	q.Push(3)
	q.Push(4)
	actual := q.GetQueue()
	//test.Error(actual)
	expected := []int{1, 2, 3, 4}
	for i := range expected {
		if actual[i] != expected[i] {
			test.Error("Test failed: expected:", expected, "actual:", actual)
		}
	}
}

func TestQueueFull(test *testing.T) {
	q := NewQueue(4)
	q.Push(1)
	q.Push(2)
	q.Push(3)
	q.Push(4)
	err := q.Push(5)
	if err == nil {
		test.Error("Test failed: expected an error for pushing past queue limit")
	}
}

func TestQueueEmpty(test *testing.T) {
	q := NewQueue(4)
	q.Push(1)
	q.Push(2)
	q.Push(3)
	q.Push(4)
	var err error
	for i := 5; i >= 0; i-- {
		_, err = q.Pop()
	}
	if err == nil {
		test.Error("Test failed: expected an error for trying to pop a value off the queue when it's empty")
	}
}
