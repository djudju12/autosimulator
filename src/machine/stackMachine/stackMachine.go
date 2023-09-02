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

		_stackA        *collections.Stack
		_stackB        *collections.Stack
		_stackAHistory [][]string
		_stackBHistory [][]string
		_currentState  string
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
	return &Machine{
		_stackA:        a,
		_stackB:        b,
		_stackAHistory: [][]string{},
		_stackBHistory: [][]string{},
	}
}

func (m *Machine) Init(input *collections.Fita) {
	m.Input = input
	m._currentState = m.InitialState
	m._stackA = collections.NewStack()
	m._stackB = collections.NewStack()
	m._stackAHistory = [][]string{}
	m._stackBHistory = [][]string{}
	m.backupStacks()
}

func (m *Machine) IsLastState() bool {
	return utils.Contains(m.FinalStates, m._currentState)
}

func (m *Machine) CurrentState() string {
	return m._currentState
}

func (m *Machine) Type() int {
	return machine.TWO_STACK_MACHINE
}

func (m *Machine) GetStates() []string {
	return m.States
}

func (m *Machine) GetTransitions(state string) []machine.Transition {
	transitions := m.Transitions[state]
	result := make([]machine.Transition, len(transitions))

	// Necessarios pois interfaces possuem diferentes layouts in
	// memory que concrete types.
	for i := range transitions {
		result[i] = &transitions[i]
	}

	return result
}

func (m *Machine) PossibleTransitions() []machine.Transition {
	return m.GetTransitions(m._currentState)
}

func (m *Machine) Stacks() []*collections.Stack {
	return []*collections.Stack{m._stackA, m._stackB}
}

func (t *Transition) MakeTransition(m machine.Machine) bool {
	stackMachine, ok := m.(*Machine)
	if !ok {
		fmt.Printf("tipo invalido de maquina: %v\n", m)
		os.Exit(1)
	}

	a := stackMachine._stackA
	b := stackMachine._stackB

	check := func(stack *collections.Stack, read string, write string) bool {
		if read != collections.PALAVRA_VAZIA {
			if stack.IsEmpty() {
				return false
			}

			current := stack.Peek(1)[0]
			if current != read {
				return false
			}

			stack.Pop()
		}

		if write != collections.PALAVRA_VAZIA {
			stack.Push(write)
		}

		return true
	}

	// checa os dois stacks
	if ok = (check(a, t.ReadA, t.WriteA) && check(b, t.ReadB, t.WriteB)); ok {
		stackMachine._currentState = t.ResultState
	}

	// Salva o historico do stack
	v, _ := m.(*Machine)
	v.backupStacks()

	return ok
}

func (m *Machine) StackHistory() ([][]string, [][]string) {
	return m._stackAHistory, m._stackBHistory
}

func (m *Machine) backupStacks() {
	m._stackAHistory = append(m._stackAHistory, utils.Reserve(m._stackA.Peek(m._stackA.Length())))
	m._stackBHistory = append(m._stackBHistory, utils.Reserve(m._stackB.Peek(m._stackB.Length())))
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

func (t *Transition) Stringfy() string {
	return fmt.Sprintf("(%s, %s, %s, %s, %s, %s)", t.Symbol, t.ReadA, t.WriteA, t.ReadB, t.WriteB, t.ResultState)
}
