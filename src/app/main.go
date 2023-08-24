package main

import (
	"autosimulator/src/graphics"
	"autosimulator/src/reader"
)

func main() {
	machine := reader.ReadMachine("/home/jonathan/programacao/autosimulator/src/machine/afdMachine/machine_example.json")
	window := graphics.NewSDLWindow()
	environment := graphics.PopulateEnvironment(window, machine)
	graphics.Mainloop(environment)
}
