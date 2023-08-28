package main

import (
	"autosimulator/src/collections"
	"autosimulator/src/graphics"
	"autosimulator/src/reader"
)

func main() {
	// machine := reader.ReadMachine("/home/jonathan/programacao/autosimulator/src/machine/afdMachine/machine_example.json")
	m := reader.ReadStackMachine("/home/jonathan/programacao/autosimulator/src/machine/stackMachine/stack_machine_example.json")
	window := graphics.NewSDLWindow()
	m.Init()
	fita := collections.FitaFromArray([]string{"a", "a", "b"})
	// fmt.Printf("%+v\n", a[1])
	environment := graphics.PopulateEnvironment(window, m)
	environment.Input(fita)
	graphics.Mainloop(environment)
}
