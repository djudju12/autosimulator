package afdMachine

import (
	"autosimulator/src/machine"
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

		_currentState string
	}

	Transition struct {
		Symbol      string `json:"symbol"`
		ResultState string `json:"resultState"`
	}
)

func New() *Machine {
	return &Machine{}
}

func (m *Machine) Init() {
	m._currentState = m.InitialState
}

func (m *Machine) IsLastState() bool {
	return utils.Contains(m.FinalStates, m._currentState)
}

func (m *Machine) PossibleTransitions() []machine.Transition {
	transitions := m.Transitions[m._currentState]
	result := make([]machine.Transition, len(transitions))

	// Aparentemente interaces possuem diferentes layouts in
	// memory que concrete types. Segundo stack overflow o
	// compilador nao faz automaticamente por causa da complexidade
	// O(n)
	for i := range transitions {
		result[i] = &transitions[i]
	}

	return result
}

func (t *Transition) GetSymbol() string {
	return t.Symbol
}

func (t *Transition) MakeTransition(m machine.Machine) bool {
	// TODO handle ok properly
	v, ok := m.(*Machine)
	v._currentState = t.ResultState
	return ok
}

func (t *Transition) UnmarshalJSON(data []byte) error {
	parsed, err := utils.ParseTransition((string(data)))
	if err != nil {
		fmt.Print(err.Error())
		panic(1)
	}

	if len(parsed) > 2 {
		fmt.Print(
			`transições de maquinas AFD devem seguir o padrao:
				"<estadoAtual>":[
					"(<simbolo>, <proximoEstado>)"
					],`)
		panic(1)
	}

	*t = Transition{
		Symbol:      parsed[0],
		ResultState: parsed[1],
	}

	return nil
}
