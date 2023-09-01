package collections

import (
	"fmt"
	"os"
)

type Fita struct {
	first   *node
	current *node
	last    *node
	len     int
}

const (
	TAIL_FITA     = "?"
	PALAVRA_VAZIA = "&"
)

func NewFita() *Fita {
	return &Fita{
		first:   nil,
		last:    nil,
		current: nil,
		len:     0,
	}
}

func (f *Fita) Read() string {
	if f.current == nil {
		return ""
	}

	value := f.current.value
	f.current = f.current.next
	return value
}

func (f *Fita) Print() {
	fmt.Printf("%s\n", f.Peek(f.Length()))
}

func (f *Fita) Reset() {
	f.current = f.first
}

func (f *Fita) Write(item string) {
	newNode := &node{
		value: item,
		next:  nil,
	}

	if f.first == nil {
		f.first = newNode
		f.current = newNode
		f.last = newNode
	} else {
		f.last.next = newNode
		f.last = newNode
	}

	f.len++
}

func FitaFromArray(value []string) *Fita {
	fita := NewFita()
	for _, item := range value {
		fita.Write(item)
	}
	fita.Write(TAIL_FITA)
	return fita
}

func (f *Fita) Length() int {
	return f.len
}

func (f *Fita) IsLast() bool {
	if f.current == nil {
		return false
	}

	return f.current.next == nil
}

func (f *Fita) Peek(amount int) []string {
	if amount < 0 {
		fmt.Printf("Amount  nao pode ser menor que 0: f.BulkRead()")
		os.Exit(1)
	}

	var result []string
	var value interface{}
	node := f.current
	for i := 0; i < amount; i++ {
		if node == nil {
			break
		}
		value = node.value
		result = append(result, value.(string))
		node = node.next
	}

	return result
}
