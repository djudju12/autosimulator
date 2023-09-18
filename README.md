# Relatório: Simulador de Autômato 

Aluno: Jonathan Willian dos Santos 

Disciplina: Compubilidade 

Professor: Ivan L. Suptitz

--- 

### Descrição do Funcionamento
O presente trabalho tem como objetivo a implementação de um simulador de autômato. No simulador será possível modelar três tipos de autômatos:
- Autômato Finito Determinístico (AFD)
- Autômato de uma Pilha
- Autômato de duas Pilhas

A linguagem escolhida para a implementação da lógica do simulador foi o **GoLang**. A escolha da linguagem se deu pelos seguintes motivos:
1. A linguagem possui uma sintaxe simples e de fácil compreensão;
2. Eu, o autor, não possuo conhecimento prévio da linguagem, o que me possibilitou aprender uma nova linguagem de programação;


O modelo de dados escolhido para representar os autômatos foi o **JSON** pelos seguintes motivos:
1. o formato é amplamente utilizado e conhecido;
2. Golang possui suporte nativo para JSON, o que facilita a implementação do simulador;

Segue exemplos de JSON's que definem cada um dos autômatos que poodem ser simulados:

<br>

_Autômato Finito Determinístico (AFD)_ 

```json
{
   "type": "simple_machine",
   "states": [
      "q0",
      "q1",
      "q2"
   ],
   "initialState": "q0",
   "finalStates": [
      "q2"
   ],
   "alfabet": [
      "a",
      "b"
   ],
   "defaultInput": [
      "a",
      "a",
      "b"
   ],
   "transitions": {
      "q0": [
         "(a, q1)"
      ],
      "q1": [
         "(a, q1)",
         "(b, q2)"
      ]
   }
}
```

<br>

_Autômato de uma Pilha_

```json
{
   "type": "1_stack_machine",
   "states": [
      "q0",
      "q1",
      "qf"
   ],
   "initialState": "q0",
   "finalStates": [
      "qf"
   ],
   "alfabet": [
      "a",
      "b"
   ],
   "defaultInput": [
      "a",
      "a",
      "b"
   ],
   "transitions": {
      "q0": [
         "(a, &, b, q0)",
         "(b, b, &, q1)"
      ],
      "q1": [
         "(b, b, &, q1)",
         "(?, ?, &, qf)"
      ]
   }
}
```
<br>

_Autômato de duas Pilha_

```json
{
   "type": "2_stack_machine",
   "states": [
      "q0",
      "q1",
      "qf"
   ],
   "initialState": "q0",
   "finalStates": [
   	"qf"
   ],
   "alfabet": [
      "a",
      "b"
   ],
   "defaultInput": [
      "a",
      "a",
      "b"
   ],
   "transitions": {
      "q0": [
         "(a, &, b, &, &, q0)",
         "(b, b, &, &, &, q1)",
         "(?, ?, &, ?, &, qf)"
      ],
      "q1": [
         "(b, b, &, &, &, q1)",
         "(?, ?, &, ?, &, qf)"
      ]
   }
}
```

Campos:
- **type**: tipo do autômato que será simulado. Pode ser: "_simple_machine_", "_1_stack_machine_" ou "_2_stack_machine_";
- **states**: lista de estados do autômato;
- **initialState**: estado inicial do autômato;
- **finalStates**: lista de estados finais do autômato;
- **alfabet**: alfabeto do autômato;
- **defaultInput**: uma entrada padrão para o autômato;
- **transitions**: lista de transições do autômato. Cada estado pode possuir uma lista de transições. As transições são representadas diferentemente para cada tipo de autômato.
  - **simple_machine**: as transições são representadas como uma lista de tuplas, onde o primeiro elemento da tupla é o símbolo de entrada e o segundo elemento é o estado de destino;
  - **1_stack_machine**: as transições são representadas como uma lista de tuplas, onde o primeiro elemento da tupla é o símbolo de entrada, o segundo elemento é o símbolo que será lido da pilha, o terceiro será escrito na pilha e o quarto elemento é o estado de destino;
  - **2_stack_machine**: funciona de forma análoga ao _1_stack_machine_, porém, com duas pilhas. (Entrada, Ler 1 stack, Escreve 1 stack, Ler 2 stack, Escreve 2 stack, Estado de destino);
- **?**: símbolo que representa o final da fita de entrada ou da pilha;
- **&**: símbolo que representa a palavra vazia.

### Estruturas de Dados Utilizadas

Para a implementação do simulador foram utilizadas as seguintes estruturas de dados:
- **fita**: utilizado para representar a fita de entrada do autômato. A fita é uma lista de simbolos e um ponteiro que indica a posição atual da fita;
- **pilha**: utilizado para representar a pilha do autômato. A pilha é uma lista de simbolos e um ponteiro que indica a posição atual da pilha;

As estruturas de dados acima foram implementadas do zero. O restante das estruturas de dados utilizadas foram as nativas da linguagem GoLang (maps, slices...).

### Print das Telas

Para a visualização do funcionamento do simulador, foi implementado uma interface gráfica utilizando SDL2 (Simple DirectMedia Layer). A interface gráfica é por uma tela principal onde os estados do autômato, as pilhas, a fita e o histórico da computação estão dispostos. Além disso, o usuário conta com um menu para alterar o autômato que está sendo simulado, a fita de entrada e possui opções para salvar e carregar uma entrada de um arquivo.

Tela principal:
![Tela Principal](./img/tela_principal.png)

Mais sobre o funcionamento da interface gráfica pode ser visto no vídeo de demonstração do simulador: [Vídeo de Demonstração](https://placeholder.com)

### Lógica de Funcionamento
No módulo **Machine** é possível averiguar a implementação do simulador. O simulador da máquina possuí uma interface principal que todas as variantes de autômato implementam. A interface é a seguinte:

```go
type Machine interface {
      Type() int
		GetInitialState() string
		GetFinalStates() []string
		GetInput() *collections.Fita
		Init(input *collections.Fita)
		Stacks() []*collections.Stack
		InLastState() bool
		PossibleTransitions() []Transition
		CurrentState() string
		GetTransitions(state string) []Transition
		GetStates() []string
}
``` 

Além disso, as máquinas possuem transições e tais transições também precisam seguir esta interface:

```go
type Transition interface {
		GetSymbol() string
		GetResultState() string
		MakeTransition(machine Machine) bool
		Stringfy() string
}
```

Essas interfaces buscam separar a execução da máquina da sua representação. Dessa forma, podemos ter uma única implementação da execução da máquina e várias implementações de representação da máquina.

A seguir, a função que executa uma maquina que implementa as interfaces anteriormente citadas:

```go
func Execute(m Machine, fita *collections.Fita) *Computation {
	// Seta o estado inicial
	m.Init(fita)

	// Criar um registro para salvar o historico da computação
	comp := newComputation(m)

	var ok bool = true
	for {
		// Lê o proximo input
		symbol := fita.Read()

		// Salva o estado atual
		stateBefore := m.CurrentState()

		// Faz a transição de estados
		ok = NextTransition(m, symbol)
		if !ok {
			break
		}

		// Salva o histórico da transição
		comp.add(stateBefore, m.CurrentState())
	}

	// Marca se foi aceita a entrada
	comp.setResult(m)

	// Printa o histórico da computação
	fmt.Printf("Fita: %s\nResultado:\n%s\n", fita.Stringfy(), comp.Stringfy())
	return comp
}
```

Note que a função acima recebe uma máquina e uma fita de entrada. A função então executa a máquina até que não seja possível fazer mais transições. Ao final, a função retorna um registro da computação que foi executada. O registro da computação é uma estrutura de dados que armazena o histórico de transições e o resultado da computação, isso será utilizado para dar log da computação e para a interface gráfica.