package collections

import (
	"encoding/json"
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
		return TAIL_FITA
	}

	value := f.current.value
	f.current = f.current.next
	return value
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
	return f.current == nil
}

func (f *Fita) Peek(amount int) []string {
	if amount < 0 {
		fmt.Printf("Amount nao pode ser menor que 0: f.BulkRead()")
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

func (f *Fita) PeekAll() []string {
	return f.Peek(f.Length())
}

func (f *Fita) ToArray() []string {
	var result []string
	var value interface{}
	node := f.first
	for node != nil {
		value = node.value
		result = append(result, value.(string))
		node = node.next
	}

	return result
}

func (f *Fita) Stringfy() string {
	var s string
	current := f.first
	for current != nil {
		s += fmt.Sprintf("%s ", current.value)
		current = current.next
	}

	return s
}

func (f *Fita) UnmarshalJSON(data []byte) error {
	arr := []string{}
	err := json.Unmarshal(data, &arr)
	if err != nil {
		return err
	}

	*f = *FitaFromArray(arr)
	return nil
}
