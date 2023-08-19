package main

import (
	"autosimulator/src/machine"
	"autosimulator/src/reader"
	"fmt"
	"os"
)

func main() {
	m := reader.ReadStackMachine("/home/jonathan/programacao/autosimulator/src/machine/stackMachine/stack_machine_example.json")
	input := []string{"a", "a", "b", "b"}
	fmt.Printf("%t\n", m.Execute(input))
	// inputs := reader.ReadInputs("/home/jonathan/programacao/autosimulator/src/reader/fita_example.json")

	// for i, in := range inputs {
	// 	fmt.Printf("%d: ", i)
	// 	check(m, in.Fita, in.ExpectedResult)
	// }
}

func check(m *machine.Machine, fita []string, expected bool) {
	actual := m.Execute(fita)
	fmt.Printf("%-6t<=> %t\n", actual, expected)
	if actual != expected {
		os.Exit(1)
	}
}
