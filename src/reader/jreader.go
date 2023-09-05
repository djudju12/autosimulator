package reader

import (
	"autosimulator/src/machine"
	"autosimulator/src/machine/afdMachine"
	"autosimulator/src/machine/oneStackMachine"
	"autosimulator/src/machine/twoStackMachine"
	"autosimulator/src/utils"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ReadMachine(path string) (machine.Machine, error) {
	base := &machine.BaseMachine{}
	content, err := readFileContent(path)
	if err != nil {
		return nil, fmt.Errorf("não foi possível ler a maquina. err: %s", err)
	}

	err = json.Unmarshal(content, &base)
	if err != nil {
		return nil, fmt.Errorf("não foi possível ler a maquina. err: %s", err)
	}

	var m machine.Machine
	switch base.Type {
	case "simple_machine":
		m, err = ReadSimpleMachine(path)
	case "1_stack_machine":
		m, err = ReadOneStackMachine(path)
	case "2_stack_machine":
		m, err = ReadStackMachine(path)
	default:
		m, err = nil, fmt.Errorf("tipo de maquina não suportado: %s", base.Type)
	}

	if err != nil {
		return nil, err
	}

	if err = checkStates(m); err != nil {
		return nil, err
	}

	return m, nil
}

func ReadSimpleMachine(path string) (*afdMachine.Machine, error) {
	m := afdMachine.New()
	content, err := readFileContent(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(content, &m)
	if err != nil {
		return nil, unmarshalError(path, err)
	}

	return m, nil
}

func ReadStackMachine(path string) (*twoStackMachine.Machine, error) {
	m := twoStackMachine.New()
	content, err := readFileContent(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(content, &m)
	if err != nil {
		return nil, unmarshalError(path, err)
	}

	return m, nil
}

func ReadOneStackMachine(path string) (*oneStackMachine.Machine, error) {
	m := oneStackMachine.New()
	content, err := readFileContent(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(content, &m)
	if err != nil {
		return nil, unmarshalError(path, err)
	}

	return m, nil
}

func unmarshalError(path string, err error) error {
	return fmt.Errorf("erro ao tentar fazer o unmarshal do arquivo %s. err: %s", path, err)
}

func readFileContent(path string) ([]byte, error) {
	file, err := os.Open(path)

	if strings.ToLower(filepath.Ext(path)) != ".json" {
		return nil, fmt.Errorf("arquivo deve ser em formato JSON")
	}

	if err != nil {
		return nil, fmt.Errorf("erro ao tentar abrir o arquivo: %s", path)
	}
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("erro ao tentar ler o conteudo do arquivo: %s. err: %s", path, err)
	}

	return content, nil
}

func checkStates(machine machine.Machine) error {
	states := machine.GetStates()
	if len(states) == 0 {
		return fmt.Errorf("não há estados")
	}

	initialState := machine.GetInitialState()
	if initialState == "" {
		return fmt.Errorf("não há estado inicial")
	}

	if ok := utils.Contains(states, initialState); !ok {
		return fmt.Errorf("estado inicial {%s} não está presente nos estados da maquina", machine.GetInitialState())
	}

	finalStates := machine.GetFinalStates()
	if len(finalStates) == 0 {
		return fmt.Errorf("não há estados finais")
	}

	for _, state := range machine.GetFinalStates() {
		ok := false
		if ok = utils.Contains(states, state); !ok {
			return fmt.Errorf("estado final {%s} não está presente nos estados da maquina", state)
		}
	}

	return nil
}
