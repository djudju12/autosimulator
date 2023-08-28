package machine

import (
	"autosimulator/src/collections"
	"fmt"
)

const (
	STATE_CHANGE         = iota
	STATE_NOT_CHANGE     = iota
	STATE_INPUT_ACCEPTED = iota
	STATE_INPUT_REJECTED = iota
)

const (
	SIMPLE_MACHINE    = iota
	ONE_STACK_MACHINE = iota
	TWO_STACK_MACHINE = iota
)

const TAIL_FITA = "?"
const PALAVRA_VAZIA = "&"

type (
	Machine interface {
		Init()
		Type() int
		Stacks() []*collections.Stack
		GetStates() []string
		GetTransitions(state string) []Transition
		IsLastState() bool
		PossibleTransitions() []Transition
		CurrentState() string
	}

	Transition interface {
		GetSymbol() string
		GetResultState() string
		MakeTransition(machine Machine) bool
	}

	BaseMachine struct {
		States       []string `json:"states"`
		InitialState string   `json:"initialState"`
		FinalStates  []string `json:"finalStates"`
		Alfabet      []string `json:"alfabet"`
		Type         string   `json:"type"`
	}
)

func Execute(m Machine, fita *collections.Fita, channel chan int) bool {
	// init seta o estado inicial
	m.Init()

	isAccepted := true
	var s string
	for {
		s, _ = fita.Read()

		if ok := NextTransition(m, s); !ok {
			fmt.Println("entrada rejeitada")
			channel <- STATE_INPUT_REJECTED
			return !isAccepted
		}

		// Se o seguinte for o final da fita
		// retorna par achecar se esta no ultimo
		// estado. Apenas uma mundaça visual!
		channel <- STATE_CHANGE

		if s == TAIL_FITA {
			break
		}
	}

	isAccepted = m.IsLastState()
	if isAccepted {
		fmt.Println("entrada aceita")
		channel <- STATE_INPUT_ACCEPTED
	} else {
		fmt.Println("entrada rejeitada")
		channel <- STATE_INPUT_REJECTED
	}

	return isAccepted
}

func NextTransition(m Machine, symbol string) bool {
	possibleTransitions := m.PossibleTransitions()
	if possibleTransitions == nil {
		return false
	}

	for _, t := range possibleTransitions {
		if t.GetSymbol() == symbol {
			return t.MakeTransition(m)
		}
	}

	return false

}
