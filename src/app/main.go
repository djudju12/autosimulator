package main

import (
	"autosimulator/src/collections"
	"autosimulator/src/graphics"
	"autosimulator/src/reader"
)

func main() {
	m, _ := reader.ReadMachine("examples/machine_example.json")
	window := graphics.NewSDLWindow()
	environment := graphics.PopulateEnvironment(window, m)
	fita := collections.FitaFromArray([]string{"a", "a", "a", "b", "b", "b"})
	environment.Input(fita)
	graphics.Mainloop(environment)
}
