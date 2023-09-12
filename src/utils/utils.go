package utils

import (
	"errors"
	"fmt"
)

func ParseTransition(s string) ([]string, error) {
	// retira os double quotes do json
	s = s[1 : len(s)-1]

	if s[0] != '(' || s[len(s)-1] != ')' {
		return []string{}, errors.New("transicao deve começar com '(' e terminar com ')'")
	}

	i, j := 2, 1
	result := []string{}
	for {
		// remove espaços à esquerda
		if s[j] == ' ' {
			j++
		}

		switch currentChar := s[i]; currentChar {
		case ')':
			result = appendWithVoidWord(result, s[j:i])
			return result, nil
		case ',':
			result = appendWithVoidWord(result, s[j:i])
			i++
			j = i
		default:
			i++
		}
	}
}

// append com validacao para palavra vazia
func appendWithVoidWord(slice []string, symbol string) []string {
	if symbol == "" {
		slice = append(slice, "&")
	} else {
		slice = append(slice, symbol)
	}
	return slice
}

func Contains(slice []string, symbol string) bool {
	for _, s := range slice {
		if s == symbol {
			return true
		}
	}

	return false
}

func DebugFita(fita []string, index int) {
	for i, item := range fita {
		if i == index {
			fmt.Printf("%s* ", item)
		} else {
			fmt.Printf("%s ", item)
		}
	}
	fmt.Println()

}

func Reverse(arr []string) []string {
	if len(arr) == 0 {
		return []string{}
	}

	result := make([]string, len(arr))
	j := 0
	for i := (len(arr) - 1); i >= 0; i-- {
		result[j] = arr[i]
		j++
	}

	return result
}

func AjustMaxLen(buffer []string, index, maxLen int) []string {
	if len(buffer) == 0 {
		return []string{}
	}
	// Se há apenas 1 caracter na fita
	if len(buffer)-index < 1 {
		return []string{buffer[len(buffer)-1]}
	}

	// Se a fita é menor que o tamanho da estrutura
	if (len(buffer) - index) < maxLen {
		return buffer[index:]
	}

	return buffer[index : index+maxLen]
}
