package queue

import (
	"errors"
	"fmt"
	"sync"
)

/**
 * Similar to the stack implementation but not
 **/
type Queue struct {
	lock   sync.Mutex
	q      []int
	l      int
	origin string
}

// Return a pointer to a Queue with a fixed length
func NewQueue(length int) *Queue {
	return &Queue{sync.Mutex{}, []int{}, length, ""}
}

// Push a value onto the Queue
func (q *Queue) Push(v int) (err error) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if len(q.q) == q.l {
		return errors.New("Queue is full")
	}
	q.q = append(q.q, v)
	return err
}

// Pop a value from the Queue
func (q *Queue) Pop() (ret int, err error) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if len(q.q) == 0 {
		return ret, errors.New("Queue is empty")
	}
	ret = q.q[0]
	q.q = q.q[1:]
	return
}

/**
 * Saves the current array of ints as a string for easy
 * comparison
 **/
func (q *Queue) SaveState() {
	q.origin = fmt.Sprintf("%v", q.q)
}

/**
 * Takes the current Queue and compares it to the origin
 **/
func (q *Queue) EqualsOrigin() bool {
	state := fmt.Sprintf("%v", q.q)
	if q.origin == state {
		return true
	}
	return false
}

/**
 * Returns the Queue
 **/
func (q *Queue) GetQueue() []int {
	return q.q
}
