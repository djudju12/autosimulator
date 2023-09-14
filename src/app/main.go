package main

import (
	"autosimulator/src/collections"
	"autosimulator/src/graphics"
	"autosimulator/src/machine"
	"autosimulator/src/machine/afdMachine"
)

func main() {
	m := &afdMachine.Machine{
		BaseMachine: machine.BaseMachine{
			Type:         "simple_machine",
			States:       []string{"Q0"},
			FinalStates:  []string{"Q0"},
			InitialState: "Q0",
			Input:        collections.FitaFromArray([]string{"1"}),
		},
	}

	window := graphics.NewSDLWindow()
	environment := graphics.PopulateEnvironment(window, m)
	graphics.Mainloop(environment)
}
