package oneStackMachine

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

		stack        *collections.Stack
		stackHistory [][]string
		currentState string
	}

	Transition struct {
		Symbol      string `json:"symbol"`
		Read        string `json:"read"`
		Write       string `json:"write"`
		ResultState string `json:"resultState"`
	}
)

func New() *Machine {
	return &Machine{
		stack:        collections.NewStack(),
		stackHistory: [][]string{},
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
	m.stack = collections.NewStack()
	m.stackHistory = [][]string{}
	m.backupStacks()
}

func (m *Machine) Type() int {
	return machine.ONE_STACK_MACHINE
}

func (m *Machine) InLastState() bool {
	return utils.Contains(m.FinalStates, m.currentState)
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
	return []*collections.Stack{m.stack}
}

func (m *Machine) CurrentState() string {
	return m.currentState
}

func (m *Machine) GetStates() []string {
	return m.States
}

func (t *Transition) MakeTransition(m machine.Machine) bool {
	stackMachine, ok := m.(*Machine)
	if !ok {
		fmt.Printf("tipo invalido de maquina: %v\n", m)
		os.Exit(1)
	}

	// checa stack
	if ok = check(stackMachine.stack, t.Read, t.Write); ok {
		stackMachine.currentState = t.ResultState
	}

	// Salva o historico do stack
	stackMachine.backupStacks()

	return ok
}

func check(stack *collections.Stack, read string, write string) bool {
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

func (m *Machine) StackHistory() [][]string {
	return m.stackHistory
}

func (m *Machine) backupStacks() {
	m.stackHistory = append(m.stackHistory, utils.Reverse(m.stack.Peek(m.stack.Length())))
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

	if len(parsed) > 4 {
		err = fmt.Errorf("transição mal formada: %v", parsed)
		return err
	}

	*t = Transition{
		Symbol:      parsed[0],
		Read:        parsed[1],
		Write:       parsed[2],
		ResultState: parsed[3],
	}

	return nil
}

func (t *Transition) Stringfy() string {
	return fmt.Sprintf("(%s, %s, %s, %s)", t.Symbol, t.Read, t.Write, t.ResultState)
}
