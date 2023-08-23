package afdMachine

import (
	"autosimulator/src/machine"
	"autosimulator/src/utils"
	"errors"
)

type (
	Machine struct {
		machine.BaseMachine
		Transitions map[string][]Transition `json:"transitions"`

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

func (m *Machine) CurrentState() string {
	return m._currentState
}

func (m *Machine) IsLastState() bool {
	return utils.Contains(m.FinalStates, m._currentState)
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

func (m *Machine) GetStates() []string {
	return m.States
}

func (t *Transition) GetSymbol() string {
	return t.Symbol
}

func (t *Transition) GetResultState() string {
	return t.ResultState
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
		return err
	}

	if len(parsed) > 2 {
		err = errors.New(
			`transições de maquinas AFD devem seguir o padrao:
				"<estadoAtual>":[
					"(<simbolo>, <proximoEstado>)"
					],`)

		return err
	}

	*t = Transition{
		Symbol:      parsed[0],
		ResultState: parsed[1],
	}

	return nil
}
