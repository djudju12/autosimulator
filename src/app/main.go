package main

import (
	"autosimulator/src/graphics"
	"autosimulator/src/reader"
)

func main() {
	m := reader.ReadMachine("/home/jonathan/programacao/autosimulator/src/machine/afdMachine/machine_example.json")
	graphics.Run(m)
}
