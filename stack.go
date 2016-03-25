package main

import "time"

type stack []string

func (s *stack) Push(who string, what string) {
	x := append(*s, time.Now().Format(time.RFC3339Nano)+"\t"+who+"\t"+what)
	*s = x
}
