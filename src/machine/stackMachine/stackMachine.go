package stackMachine

import (
	"autosimulator/src/stack"
	"errors"
	"fmt"
)

type (
	Transition struct {
		Symbol      string `json:"symbol"`
		ReadA       string `json:"readA"`
		WriteA      string `json:"writeA"`
		ReadB       string `json:"readB"`
		WriteB      string `json:"writeB"`
		ResultState string `json:"resultState"`
	}

	Machine struct {
		States       []string                `json:"states"`
		InitialState string                  `json:"initialState"`
		FinalStates  []string                `json:"finalStates"`
		Alfabet      []string                `json:"alfabet"`
		Transitions  map[string][]Transition `json:"transitions"`

		_stackA       stack.Stack
		_stackB       stack.Stack
		_currentState string
	}
)

func New() *Machine {
	a := stack.New()
	b := stack.New()
	return &Machine{_stackA: *a, _stackB: *b}
}

func (m *Machine) Execute(fita []string) bool {
	m._currentState = m.InitialState
	fita = append(fita, "?")
	for _, s := range fita {
		if !m.nextTransition(s) {
			return false
		}
	}
	fmt.Printf("last state %s\n", m._currentState)
	return true
}

func (m *Machine) nextTransition(symbol string) bool {
	possibleTransitions := m.Transitions[m._currentState]
	if possibleTransitions == nil {
		return false
	}

	for _, t := range possibleTransitions {
		if t.Symbol == symbol {
			return t.execute(m)
		}
	}

	return true
}

func (t *Transition) execute(m *Machine) bool {
	// TODO: devo fazer sempre readA->writeA->readB->writeB?
	a := &m._stackA
	b := &m._stackB

	check := func(s *stack.Stack, read string, write string) bool {
		if read != "&" {
			if s.IsEmpty() {
				return false
			}

			if s.Pop() != read {
				return false
			}
		}
		if write != "&" {
			s.Push(write)
		}

		return true
	}

	m._currentState = t.ResultState
	// checa os dois stacks
	return check(a, t.ReadA, t.WriteA) || check(b, t.ReadB, t.WriteB)
}

func (t *Transition) UnmarshalJSON(data []byte) error {
	parsed, err := parse(string(data))
	if err != nil {
		fmt.Print(err.Error())
		panic(1)
	}

	*t = Transition{
		Symbol:      parsed[0],
		ReadA:       parsed[1],
		WriteA:      parsed[2],
		ReadB:       parsed[3],
		WriteB:      parsed[4],
		ResultState: parsed[5],
	}

	return nil
}

// append com validacao para palavra vazia
func appendWithVoidWord(slice []string, symbol string) []string {
	if symbol == "" {
		slice = append(slice, "&")
	} else {
		slice = append(slice, symbol)
	}
	return slice
}

func parse(s string) ([]string, error) {
	// TODO: proper handling

	// retira os double quotes do json
	s = s[1 : len(s)-1]

	if s[0] != '(' || s[len(s)-1] != ')' {
		return []string{}, errors.New("transicao deve começar com '(' e terminar com ')'")
	}

	i, j := 2, 1
	result := []string{}
	for {
		// remove espaços à esquerda
		if s[j] == ' ' {
			j++
		}

		switch currentChar := s[i]; currentChar {
		case ')':
			result = appendWithVoidWord(result, s[j:i])
			return result, nil
		case ',':
			result = appendWithVoidWord(result, s[j:i])
			i++
			j = i
		default:
			i++
		}
	}
}
