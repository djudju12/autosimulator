package machine

import (
	"fmt"
)

const (
	STATE_CHANGE         = iota
	STATE_NOT_CHANGE     = iota
	STATE_INPUT_ACCEPTED = iota
	STATE_INPUT_REJECTED = iota
)

const TAIL_FITA = "?"
const PALAVRA_VAZIA = "&"

type (
	Machine interface {
		Init()
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
	}
)

func Execute(m Machine, fita []string, channel chan int) bool {
	m.Init()

	// init seta o estado inicial
	channel <- STATE_CHANGE

	fita = append(fita, TAIL_FITA)
	i := 0
	s := fita[i]
	for s != TAIL_FITA {
		if ok := NextTransition(m, s); !ok {
			fmt.Printf("entrada: %s rejeitada", fita)
			channel <- STATE_INPUT_REJECTED
			return false
		}
		channel <- STATE_CHANGE

		i++
		s = fita[i]
	}

	isAccepted := m.IsLastState()
	if isAccepted {
		channel <- STATE_INPUT_ACCEPTED
	} else {
		channel <- STATE_INPUT_REJECTED
	}

	fmt.Printf("last Stage: %s\n", m.CurrentState())
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
