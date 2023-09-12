package main

import (
	"autosimulator/src/graphics"
	"autosimulator/src/reader"
)

func main() {
	m, _ := reader.ReadMachine("examples/[dfa]simple_example.json")
	// fmt.Printf("%+v", machine.DefaultInput.Stringfy())
	window := graphics.NewSDLWindow()
	environment := graphics.PopulateEnvironment(window, m)
	// fita := collections.FitaFromArray([]string{"a", "a", "a", "b", "b", "b"})
	graphics.Mainloop(environment)
}
