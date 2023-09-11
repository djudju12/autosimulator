package main

import (
	"autosimulator/src/collections"
	"autosimulator/src/graphics"
	"autosimulator/src/reader"
)

func main() {
	// m, _ := reader.ReadMachine("/home/jonathan/hd/programacao/autosimulator/examples/machine_example.json")
	// m, _ := reader.ReadMachine("/home/jonathan/hd/programacao/autosimulator/examples/count_words.json")
	m, _ := reader.ReadMachine("examples/one_stack_machine_example.json")
	window := graphics.NewSDLWindow()
	environment := graphics.PopulateEnvironment(window, m)
	// fita := collections.FitaFromArray([]string{"x", "x", "x"})
	fita := collections.FitaFromArray([]string{"a", "a", "a", "b", "b", "b"})
	environment.Input(fita)
	graphics.Mainloop(environment)
}
