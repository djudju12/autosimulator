package stack

import (
	"fmt"
	"os"
)

type (
	Stack struct {
		first *node
		len   int
	}

	node struct {
		value interface{}
		prev  *node
	}
)

func New() *Stack {
	firstN := &node{"?", nil}
	return &Stack{firstN, 0}
}

func (s *Stack) Len() int {
	return s.len
}

// remove o topo e retorna o valor
func (s *Stack) Pop() interface{} {
	if s.first == nil {
		fmt.Println("Pop() em uma pilha vazia!")
		os.Exit(1)
	}
	temp := s.first
	s.first = s.first.prev
	return temp.value
}

func (s *Stack) Push(value interface{}) {
	n := node{
		value,
		s.first,
	}
	s.first = &n
}

func (s *Stack) IsEmpty() bool {
	return s.first == nil
}

func (s *Stack) Peek() interface{} {
	return s.first.value
}
