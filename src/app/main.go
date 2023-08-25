package main

import (
	"autosimulator/src/collections"
	"autosimulator/src/graphics"
	"autosimulator/src/reader"
)

func main() {
	machine := reader.ReadMachine("/home/jonathan/programacao/autosimulator/src/machine/afdMachine/machine_example.json")
	window := graphics.NewSDLWindow()
	machine.Init()
	environment := graphics.PopulateEnvironment(window, machine)
	fita := collections.FitaFromArray([]string{"a", "b", "d", "d"})
	environment.Input(fita)
	graphics.Mainloop(environment)
}
