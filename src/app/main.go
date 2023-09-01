package main

import (
	"autosimulator/src/graphics"
	"autosimulator/src/reader"
)

func main() {
	m := reader.ReadMachine("/home/jonathan/programacao/autosimulator/src/machine/afdMachine/machine_example.json")
	// m := reader.ReadStackMachine("/home/jonathan/programacao/autosimulator/src/machine/stackMachine/stack_machine_example.json")
	window := graphics.NewSDLWindow()
	environment := graphics.PopulateEnvironment(window, m)
	environment.Input([]string{"a", "a", "a", "b"})
	graphics.Mainloop(environment)
}
