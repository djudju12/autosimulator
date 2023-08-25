package stackMachine

import (
	"autosimulator/src/collections"
	"autosimulator/src/machine"
	"autosimulator/src/utils"
	"fmt"
	"os"
)

type (
	Machine struct {
		machine.BaseMachine
		Transitions map[string][]Transition `json:"transitions"`

		_stackA       collections.Stack
		_stackB       collections.Stack
		_currentState string
		_lastRead     int
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
	a := collections.NewStack()
	b := collections.NewStack()
	return &Machine{_stackA: *a, _stackB: *b}
}

func (m *Machine) Init() {
	m._currentState = m.InitialState
}

func (m *Machine) IsLastState() bool {
	return utils.Contains(m.FinalStates, m._currentState)
}

func (m *Machine) GetStates() []string {
	return m.States
}

func (m *Machine) CurrentState() string {
	return m._currentState
}

func (m *Machine) GetTransitions(state string) []machine.Transition {
	transitions := m.Transitions[state]
	result := make([]machine.Transition, len(transitions))

	// Necessarios pois  interaces possuem diferentes layouts in
	// memory que concrete types.
	for i := range transitions {
		result[i] = &transitions[i]
	}

	return result
}
func (m *Machine) PossibleTransitions() []machine.Transition {
	return m.GetTransitions(m._currentState)
}

func (t *Transition) MakeTransition(m machine.Machine) bool {
	// TODO: devo fazer sempre readA->writeA->readB->writeB?
	stackMachine, ok := m.(*Machine)
	if !ok {
		fmt.Printf("invalido tipo de maquina: %v\n", m)
		os.Exit(1)
	}

	a := &stackMachine._stackA
	b := &stackMachine._stackB

	// TODO tenho que checar se esta no ultimo estado, não?
	check := func(s *collections.Stack, read string, write string) bool {
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

	stackMachine._currentState = t.ResultState

	// checa os dois stacks
	return check(a, t.ReadA, t.WriteA) || check(b, t.ReadB, t.WriteB)
}

func (t *Transition) GetSymbol() string {
	return t.Symbol
}

func (t *Transition) GetResultState() string {
	return t.ResultState
}

func (t *Transition) UnmarshalJSON(data []byte) error {
	parsed, err := utils.ParseTransition((string(data)))
	if err != nil {
		return err
	}

	if len(parsed) > 6 {
		err = fmt.Errorf("transição mal formada: %v", parsed)
		return err
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
