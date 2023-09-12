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

func (m *Machine) GetInput() *collections.Fita {
	return m.Input
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
	isPopable := func(stack *collections.Stack, read string) bool {
		if stack.IsEmpty() {
			return false
		}

		// Olha o primeiro elemento do stack
		head := stack.Peek(1)[0]

		// Retorna true se for igual o elemento que deverá ser lido
		return head == read
	}

	// Nota: Eu escolhi uma maneira um tanto estranha para verificar a transição
	// isPopable retorna true se for necessário ler uma palavra e essa palavra
	// esta no topo do stack. Se ela não estiver vai retornar FALSE e, como
	// não foi possível fazer a transição, retorna da função.
	// Se esta sendo lida a palavra vazia então não fará o pop e seguirá normalmente.
	popA := false
	if t.ReadA != collections.PALAVRA_VAZIA {
		if popA = isPopable(a, t.ReadA); !popA {
			return false
		}
	}

	popB := false
	if t.ReadB != collections.PALAVRA_VAZIA {
		if popB = isPopable(b, t.ReadB); !popB {
			return false
		}
	}

	if popA {
		a.Pop()
	}

	if popB {
		b.Pop()
	}

	if t.WriteA != collections.PALAVRA_VAZIA {
		a.Push(t.WriteA)
	}

	if t.WriteB != collections.PALAVRA_VAZIA {
		b.Push(t.WriteB)
	}

	// Avança para o proximo estado
	stackMachine.currentState = t.GetResultState()

	// Salva o historico do stack
	stackMachine.backupStacks()

	return true
}

func (m *Machine) StackHistory() ([][]string, [][]string) {
	return m.stackAHistory, m.stackBHistory
}

func (m *Machine) backupStacks() {
	m.stackAHistory = append(m.stackAHistory, utils.Reverse(m.stackA.Peek(m.stackA.Length())))
	m.stackBHistory = append(m.stackBHistory, utils.Reverse(m.stackB.Peek(m.stackB.Length())))
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
