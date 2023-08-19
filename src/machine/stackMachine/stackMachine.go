package stackMachine

import (
	"autosimulator/src/stack"
	"autosimulator/src/utils"
	"fmt"
)

type (
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

	Transition struct {
		Symbol      string `json:"symbol"`
		ReadA       string `json:"readA"`
		WriteA      string `json:"writeA"`
		ReadB       string `json:"readB"`
		WriteB      string `json:"writeB"`
		ResultState string `json:"resultState"`
	}
)

func New() *Machine {
	a := stack.New()
	b := stack.New()
	return &Machine{_stackA: *a, _stackB: *b}
}

func (m *Machine) Init() {
	m._currentState = m.InitialState
}

// TODO: default impl?
func (m *Machine) IsLastState() bool {
	return utils.Contains(m.FinalStates, m._currentState)
}

func (m *Machine) NextTransition(symbol string) (Transition, bool) {
	possibleTransitions := m.Transitions[m._currentState]
	if possibleTransitions == nil {
		return Transition{}, false
	}

	for _, t := range possibleTransitions {
		if t.Symbol == symbol {
			return t, true
		}
	}

	return Transition{}, false
}

func (t *Transition) MakeTransition(m Machine) bool {
	// TODO: devo fazer sempre readA->writeA->readB->writeB?
	a := &m._stackA
	b := &m._stackB

	// TODO tenho que checar se esta no ultimo estado, não?
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
	parsed, err := utils.ParseTransition((string(data)))
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
