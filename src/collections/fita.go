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

func NewFita() *Fita {
	return &Fita{
		first:   nil,
		last:    nil,
		current: nil,
		len:     0,
	}
}

func (f *Fita) Read() (string, bool) {
	if f.current == nil {
		f.Reset()
		return "", false
	}

	value, _ := f.current.value.(string)
	f.current = f.current.next
	return value, true
}

func (f *Fita) Reset() {
	f.current = f.first
}

func (f *Fita) Peek(amount int) []string {
	if amount < 0 {
		fmt.Printf("Amount  nao pode ser menor que 0: f.BulkRead()")
		os.Exit(1)
	}

	var result []string
	var value string
	node := f.current
	for i := 0; i < amount; i++ {
		if node == nil {
			break
		}
		value, _ = node.value.(string)
		result = append(result, value)
		node = node.next
	}

	return result
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

// func (f *Fita) WriteLast(item string) {
// 	newNode := &node{
// 		value: item,
// 		next:  nil,
// 	}

// 	if f.first == nil {
// 		f.first = newNode
// 		f.current = newNode
// 		f.last = newNode
// 	} else {

// 	}

// }

func (f Fita) ToArray() []string {
	var result []string
	for s, ok := f.Read(); ok; f.Read() {
		result = append(result, s)
	}

	return result
}

func FitaFromArray(value []string) *Fita {
	fita := NewFita()
	for _, item := range value {
		fita.Write(item)
	}

	return fita
}

func (f *Fita) Length() int {
	return f.len
}
