package main

import (
	"autosimulator/src/graphics"
	"autosimulator/src/reader"
)

func main() {
	m, _ := reader.ReadMachine("machines/[dfa]simple_example.json")
	window := graphics.NewSDLWindow()
	environment := graphics.PopulateEnvironment(window, m)
	graphics.Mainloop(environment)
}
