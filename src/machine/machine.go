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
		Init()
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
	}

	Computation struct {
		History []ComputationRecord
	}

	ComputationRecord struct {
		currentState string
		transition   Transition
		result       string
	}
)

var Comp *Computation

func Execute(m Machine, fita *collections.Fita) *Computation {
	// init seta o estado inicial
	m.Init()
	Comp = newComputation(m)
	var s string
	for {
		s, _ = fita.Read()
		if s == collections.TAIL_FITA {
			break
		}

		if ok := NextTransition(m, s); !ok {
			Comp.setResult(REJECTED)
			return Comp
		}
	}

	if m.IsLastState() {
		Comp.setResult(ACCEPTED)
	} else {
		Comp.setResult(REJECTED)
	}

	return Comp
}

func NextTransition(m Machine, symbol string) bool {
	possibleTransitions := m.PossibleTransitions()
	if possibleTransitions == nil {
		return false
	}

	for _, t := range possibleTransitions {
		if t.GetSymbol() == symbol {
			Comp.add(m.CurrentState(), t)
			return t.MakeTransition(m)
		}
	}

	return false

}

func newComputationRecord(m Machine) *ComputationRecord {
	record := ComputationRecord{
		currentState: m.CurrentState(),
		result:       INITIAL,
	}
	return &record
}

func newComputation(m Machine) *Computation {
	record := newComputationRecord(m)
	return &Computation{
		History: []ComputationRecord{*record},
	}
}

func (c *Computation) setResult(result string) {
	c.History[len(c.History)-1].result = result
}

func (c *Computation) add(currentState string, trans Transition) {
	record := ComputationRecord{
		currentState: currentState,
		transition:   trans,
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
	var s string
	if cr.transition == nil {
		s += fmt.Sprintf("%s, ( ) %s", cr.currentState, cr.result)
	} else {
		s += fmt.Sprintf("%s -> %s, %s %s", cr.currentState, cr.transition.GetResultState(), cr.transition.Stringfy(), cr.result)
	}

	return s
}

func (cr *ComputationRecord) Details() map[string]string {
	details := make(map[string]string)
	details["CURRENT_STATE"] = cr.currentState
	if cr.transition == nil {
		details["NEXT_STATE"] = ""
	} else {
		details["NEXT_STATE"] = cr.transition.GetResultState()
	}

	details["TRANSITION"] = cr.Stringfy()
	details["RESULT"] = cr.result

	return details
}
