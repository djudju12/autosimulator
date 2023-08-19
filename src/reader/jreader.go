package reader

import (
	"autosimulator/src/machine/afdMachine"
	"autosimulator/src/machine/stackMachine"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ReadMachine(path string) *afdMachine.Machine {
	m := afdMachine.New()
	content := readFileContent(path)
	err := json.Unmarshal(content, &m)
	if err != nil {
		unmarshalError(path, err)
	}
	return m
}

func ReadStackMachine(path string) *stackMachine.Machine {
	m := stackMachine.New()
	content := readFileContent(path)
	err := json.Unmarshal(content, &m)
	if err != nil {
		unmarshalError(path, err)
	}
	return m
}

// func ReadInputs(path string) []afdMachine.Input {
// 	var in []afdMachine.Input
// 	content := readFileContent(path)
// 	err := json.Unmarshal(content, &in)
// 	if err != nil {
// 		unmarshalError(path, err)
// 	}
// 	return in
// }

func unmarshalError(path string, err error) {
	fmt.Printf("Erro ao tentar fazer o unmarshal do arquivo %s. Error: %s\n", path, err)
	os.Exit(1)
}

func readFileContent(path string) []byte {
	file, err := os.Open(path)

	if strings.ToLower(filepath.Ext(path)) != ".json" {
		fmt.Printf("Arquivo deve ser em formato JSON\n")
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Erro ao tentar abrir o arquivo: %s\n", path)
		os.Exit(1)
	}
	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Erro ao tentar ler o conteudo do arquivo: %s. Error: %s\n", path, err)
		os.Exit(1)
	}

	return content
}
