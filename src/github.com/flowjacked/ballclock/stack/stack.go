package stack

import (
	"errors"
	"sync"
)

/**
 * A basic stack. Code borrowed from:
 * http://stackoverflow.com/questions/28541609/looking-for-reasonable-stack-implementation-in-golang
 * I modified it to have a fixed length though I'm still using a slice of an undefined array
 **/
type stack struct {
	lock sync.Mutex
	s    []int
	l    int
}

/**
 * A simple wrapper around a slice that provides stack methods
 **/
func NewStack(length int) *stack {
	return &stack{sync.Mutex{}, []int{}, length}
}

/**
 * Lock the stack and push a value on
 **/
func (s *stack) Push(v int) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if len(s.s) == s.l {
		return errors.New("Stack full")
	}
	s.s = append(s.s, v)
	return err
}

/**
 * Lock the stack and give me the last value pushed onto the stack
 **/
func (s *stack) Pop() (int, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.s)
	if l == 0 {
		return 0, errors.New("Empty Stack")
	}

	res := s.s[l-1]
	s.s = s.s[:l-1]
	return res, nil
}

/**
 * When asking for an hour count, it will always reflect 1 smaller than it should
 * in all other cases, just return the length
 *
 **/
func (s *stack) Count(unit string) int {
	switch unit {
	case "hour":
		return len(s.s) + 1
	default:
		return len(s.s)
	}
}

func (s *stack) GetStack() []int {
	return s.s
}
