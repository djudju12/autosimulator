package machine

import (
	"autosimulator/src/collections"
	"fmt"
)

const (
	ACCEPTED = "[V]"
	REJECTED = "[X]"
	RUNNING  = "[ ]"
	INITIAL  = "[I]"
)

const (
	SIMPLE_MACHINE    = iota
	ONE_STACK_MACHINE = iota
	TWO_STACK_MACHINE = iota
)

type (
	Machine interface {
		Init(input *collections.Fita)
		Type() int
		Stacks() []*collections.Stack
		IsLastState() bool
		PossibleTransitions() []Transition
		CurrentState() string
		GetTransitions(state string) []Transition
		GetStates() []string
	}

	Transition interface {
		GetSymbol() string
		GetResultState() string
		MakeTransition(machine Machine) bool
		Stringfy() string
	}

	BaseMachine struct {
		States       []string `json:"states"`
		InitialState string   `json:"initialState"`
		FinalStates  []string `json:"finalStates"`
		Alfabet      []string `json:"alfabet"`
		Type         string   `json:"type"`
		Input        *collections.Fita
	}

	Computation struct {
		History []ComputationRecord
	}

	ComputationRecord struct {
		lastState    string
		currentState string
		result       string
	}
)

func Execute(m Machine, fita *collections.Fita) *Computation {
	// Seta o estado inicial
	m.Init(fita)

	// Criar um registro para salvar o historico da computação
	comp := newComputation(m)

	ok := true
	for ok {
		// Lê o proximo input
		symbol := fita.Read()

		// Salva o estado atual
		stateBefore := m.CurrentState()

		// Faz a transição de estados
		ok = NextTransition(m, symbol)

		// Salva o histórico da transição
		comp.add(stateBefore, m.CurrentState())
	}

	// Marca se foi aceita a entrada
	comp.setResult(m)

	// Pritna o histórico da computação
	fmt.Println(comp.Stringfy())
	return comp
}

func NextTransition(m Machine, symbol string) bool {
	if symbol == "" {
		return false
	}

	possibleTransitions := m.PossibleTransitions()
	if possibleTransitions == nil {
		return false
	}

	for _, t := range possibleTransitions {
		if symbol == t.GetSymbol() {
			return t.MakeTransition(m)
		}
	}

	return false
}

func newComputation(m Machine) *Computation {
	record := newComputationRecord(m)
	return &Computation{
		History: []ComputationRecord{*record},
	}
}

func newComputationRecord(m Machine) *ComputationRecord {
	record := ComputationRecord{
		lastState:    m.CurrentState(),
		currentState: m.CurrentState(),
		result:       INITIAL,
	}
	return &record
}

func (c *Computation) setResult(m Machine) {
	if m.IsLastState() {
		c.History[len(c.History)-1].result = ACCEPTED
	} else {
		c.History[len(c.History)-1].result = REJECTED
	}
}

func (c *Computation) add(lastState string, currentState string) {
	record := ComputationRecord{
		lastState:    lastState,
		currentState: currentState,
		result:       RUNNING,
	}

	c.History = append(c.History, record)
}

func (c *Computation) Stringfy() string {
	var s string
	for i, v := range c.History {
		s += fmt.Sprintf("%d: %s\n", i, v.Stringfy())
	}

	return s
}

func (cr *ComputationRecord) Stringfy() string {
	if cr.result != RUNNING {
		return fmt.Sprintf("%s %s", cr.lastState, cr.result)
	}

	return fmt.Sprintf("%s -> %s", cr.lastState, cr.currentState)
}

func (cr *ComputationRecord) Details() map[string]string {
	details := make(map[string]string)
	details["LAST_STATE"] = cr.lastState
	details["NEXT_STATE"] = cr.currentState
	details["RESULT"] = cr.result
	return details
}
