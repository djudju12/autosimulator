package main

import (
	"autosimulator/src/collections"
	"autosimulator/src/graphics"
	"autosimulator/src/reader"
)

func main() {
	// machine := reader.ReadMachine("/home/jonathan/programacao/autosimulator/src/machine/afdMachine/machine_example.json")
	machine := reader.ReadStackMachine("/home/jonathan/programacao/autosimulator/src/machine/stackMachine/stack_machine_example.json")
	window := graphics.NewSDLWindow()
	machine.Init()
	environment := graphics.PopulateEnvironment(window, machine)
	fita := collections.FitaFromArray([]string{"a", "a", "b", "b"})
	environment.Input(fita)
	graphics.Mainloop(environment)
}

// func main() {
// m := reader.ReadMachine("/home/jonathan/programacao/autosimulator/src/machine/afdMachine/machine_example.json")
// m := reader.ReadStackMachine("/home/jonathan/programacao/autosimulator/src/machine/stackMachine/stack_machine_example.json")
// fita := collections.FitaFromArray([]string{"a", "a", "b"})
// computation := machine.Execute(m, fita)
// fmt.Println(computation.Stringfy())
// }
