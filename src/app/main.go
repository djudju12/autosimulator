package main

import (
	"autosimulator/src/graphics"
	"autosimulator/src/reader"
)

func main() {
	m, _ := reader.ReadMachine("/home/jonathan/hd/programacao/autosimulator/examples/machine_example.json")
	window := graphics.NewSDLWindow()
	environment := graphics.PopulateEnvironment(window, m)
	environment.Input([]string{"a", "a", "a", "b"})
	graphics.Mainloop(environment)
}
