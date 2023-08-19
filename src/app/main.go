package main

import (
	"autosimulator/src/machine"
	"autosimulator/src/reader"
)

func main() {
	m := reader.ReadMachine("/home/jonathan/programacao/autosimulator/src/machine/afdMachine/machine_example.json")
	// fmt.Printf("%+v", m)
	input := []string{"a", "a", "b"}

	machine.Execute(m, input)
	// for i, in := range inputs {
	// 	fmt.Printf("%d: ", i)
	// 	check(m, in.Fita, in.ExpectedResult)
	// }
}

// func check(m *machine.Machine, fita []string, expected bool) {
// 	actual := m.Execute(fita)
// 	fmt.Printf("%-6t<=> %t\n", actual, expected)
// 	if actual != expected {
// 		os.Exit(1)
// 	}
// }
