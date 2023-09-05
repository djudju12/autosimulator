package twoStackMachine

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

		stackA        *collections.Stack
		stackB        *collections.Stack
		stackAHistory [][]string
		stackBHistory [][]string
		currentState  string
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
		stackA:        a,
		stackB:        b,
		stackAHistory: [][]string{},
		stackBHistory: [][]string{},
	}
}

func (m *Machine) GetInitialState() string {
	return m.InitialState
}

func (m *Machine) GetFinalStates() []string {
	return m.FinalStates
}

func (m *Machine) Init(input *collections.Fita) {
	m.Input = input
	m.currentState = m.InitialState
	m.stackA = collections.NewStack()
	m.stackB = collections.NewStack()
	m.stackAHistory = [][]string{}
	m.stackBHistory = [][]string{}
	m.backupStacks()
}

func (m *Machine) Type() int {
	return machine.TWO_STACK_MACHINE
}

func (m *Machine) InLastState() bool {
	return utils.Contains(m.FinalStates, m.currentState)
}

func (m *Machine) CurrentState() string {
	return m.currentState
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
	return m.GetTransitions(m.currentState)
}

func (m *Machine) Stacks() []*collections.Stack {
	return []*collections.Stack{m.stackA, m.stackB}
}

func (t *Transition) MakeTransition(m machine.Machine) bool {
	stackMachine, ok := m.(*Machine)
	if !ok {
		fmt.Printf("tipo invalido de maquina: %v\n", m)
		os.Exit(1)
	}

	a := stackMachine.stackA
	b := stackMachine.stackB

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
		stackMachine.currentState = t.ResultState
	}

	// Salva o historico do stack
	stackMachine.backupStacks()

	return ok
}

func (m *Machine) StackHistory() ([][]string, [][]string) {
	return m.stackAHistory, m.stackBHistory
}

func (m *Machine) backupStacks() {
	m.stackAHistory = append(m.stackAHistory, utils.Reserve(m.stackA.Peek(m.stackA.Length())))
	m.stackBHistory = append(m.stackBHistory, utils.Reserve(m.stackB.Peek(m.stackB.Length())))
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
