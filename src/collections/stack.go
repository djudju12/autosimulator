package collections

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
		value string
		next  *node
	}
)

func NewStack() *Stack {
	firstN := &node{"?", nil}
	return &Stack{firstN, 0}
}

func (s *Stack) Length() int {
	return s.len
}

// remove o topo e retorna o valor
func (s *Stack) Pop() string {
	if s.first == nil {
		fmt.Println("Pop() em uma pilha vazia!")
		os.Exit(1)
	}
	temp := s.first
	s.first = s.first.next
	return temp.value
}

func (s *Stack) Push(value string) {
	n := node{
		value,
		s.first,
	}
	s.first = &n
}

func (s *Stack) IsEmpty() bool {
	return s.first == nil
}

func (s *Stack) Peek(amout int) []string {
	var result []string
	current := s.first
	for i := 0; i < amout; i++ {
		if current == nil {
			break
		}
		result = append(result, current.value)
		current = current.next
	}

	return result
}
