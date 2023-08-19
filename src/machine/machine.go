package machine

import "fmt"

const TAIL_FITA = "?"
const PALAVRA_VAZIA = "&"

type (
	Machine interface {
		Init()
		IsLastState() bool
		PossibleTransitions() []Transition
	}

	Transition interface {
		GetSymbol() string
		MakeTransition(machine Machine) bool
	}
)

func Execute(m Machine, fita []string) bool {
	m.Init()
	fita = append(fita, TAIL_FITA)
	for _, s := range fita {
		if s == TAIL_FITA {
			break
		}

		if ok := NextTransition(m, s); !ok {
			fmt.Printf("entrada: %s rejeitada", fita)
			return false
		}
	}

	return m.IsLastState()
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
