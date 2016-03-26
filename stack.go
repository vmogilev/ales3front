package main

import "time"

// Stack is used to save debug messages by all functions in a chain
// it should be passed into each function as (s *Stack) and used there
// via s.Push
type Stack []string

// NewStack creates a new stack
func NewStack() *Stack {
	return &Stack{}
}

// Push stacks the time.StampMilli + who + what
func (s *Stack) Push(who string, what string) {
	*s = append(*s, time.Now().Format(time.StampMilli)+"\t"+who+"\t"+what)
}
