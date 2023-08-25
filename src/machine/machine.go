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

func Execute(m Machine, fita *collections.Fita, channel chan int) bool {
	// init seta o estado inicial
	m.Init()

	isAccepted := true
	var s string
	for i := 0; i < fita.Length()-1; i++ {
		s, _ = fita.Read()
		if ok := NextTransition(m, s); !ok {
			fmt.Printf("entrada: %+v rejeitada\n", fita.Peek(fita.Length()))
			break
		}
		channel <- STATE_CHANGE
	}

	isAccepted = m.IsLastState()
	if isAccepted {
		channel <- STATE_INPUT_ACCEPTED
	} else {
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
