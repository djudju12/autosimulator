package reader

import (
	"autosimulator/src/collections"
	"autosimulator/src/machine"
	"autosimulator/src/machine/afdMachine"
	"autosimulator/src/machine/oneStackMachine"
	"autosimulator/src/machine/twoStackMachine"
	"autosimulator/src/utils"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func ReadMachine(path string) (machine.Machine, error) {
	content, err := readFileContent(path)
	if err != nil {
		return nil, err
	}

	var m *machine.BaseMachine
	err = json.Unmarshal(content, &m)
	if err != nil {
		return nil, err
	}

	if m.Input == nil {
		return nil, fmt.Errorf("syntax error. Não há input default. Maquina: ", string(content))
	}

	var readedMachine machine.Machine
	switch m.Type {
	case "simple_machine":
		readedMachine, err = ReadSimpleMachine(path)
	case "1_stack_machine":
		readedMachine, err = ReadOneStackMachine(path)
	case "2_stack_machine":
		readedMachine, err = ReadTwoStackMachine(path)
	default:
		readedMachine, err = nil, fmt.Errorf("tipo de maquina não suportado: %s", m.Type)
	}

	if err != nil {
		return nil, err
	}

	if err = checkStates(readedMachine); err != nil {
		return nil, err
	}

	return readedMachine, nil
}

func ReadInputs(path string) ([]*collections.Fita, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	inputs, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	var result []*collections.Fita
	for _, input := range inputs {
		result = append(result, collections.FitaFromArray(input))
	}

	return result, nil

}

func ReadInput(path string) (*collections.Fita, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	input, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	if len(input) > 1 {
		return nil, fmt.Errorf("o input deve possuir apenas uma linha e seus elementos devem estar separados por vírgula e sem espaço entre eles. Input:", input)
	}

	fmt.Println(input)
	return collections.FitaFromArray(input[0]), nil
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

func ReadTwoStackMachine(path string) (*twoStackMachine.Machine, error) {
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

func readDir(path string) ([]fs.DirEntry, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func GetCsvList(path string) ([]string, error) {
	entries, err := readDir(path)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, entry := range entries {
		if !entry.IsDir() && isCsvExt(entry.Name()) {
			result = append(result, entry.Name())
		}
	}

	return result, nil
}

func GetJsonList(path string) ([]string, error) {
	entries, err := readDir(path)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, entry := range entries {
		if !entry.IsDir() && isJsonExt(entry.Name()) {
			result = append(result, entry.Name())
		}
	}

	return result, nil
}

func isJsonExt(fileName string) bool {
	return strings.ToLower(fileName[len(fileName)-5:]) == ".json"
}

func isCsvExt(fileName string) bool {
	return strings.ToLower(fileName[len(fileName)-4:]) == ".csv"
}

func WriteInput(input *collections.Fita, path string) error {
	f, err := os.Create(filepath.Join(path, time.Now().String()+".csv"))
	if err != nil {
		return err
	}

	defer f.Close()

	arr := input.ToArray()
	_, err = f.Write([]byte(arr[0]))
	if err != nil {
		return err
	}

	for _, word := range arr[1:] {
		if word == collections.TAIL_FITA {
			break
		}

		_, err = f.Write([]byte("," + word))
		if err != nil {
			return err
		}
	}

	return nil
}
