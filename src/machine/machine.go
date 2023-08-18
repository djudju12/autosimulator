package machine

type Input struct {
	Fita           []string
	ExpectedResult bool
}

type Transition struct {
	State       string `json:"state"`
	Symbol      string `json:"symbol"`
	ResultState string `json:"resultState"`
}

type Machine struct {
	States       []string     `json:"states"`
	InitialState string       `json:"initialState"`
	FinalStates  []string     `json:"finalStates"`
	Alfabet      []string     `json:"alfabet"`
	Transitions  []Transition `json:"transitions"`

	_currentState string
}

type MachineOperations interface {
	New() *Machine
	Execute(fita []string) bool
}

func New() *Machine {
	return &Machine{}
}

func (m *Machine) Execute(fita []string) bool {

	m._currentState = m.InitialState
	for _, s := range fita {
		if !m.nextTransition(s) {
			return false
		}
	}

	return m.isInLastState()
}

func (m *Machine) isInLastState() bool {
	return contains(m._currentState, m.FinalStates)
}

func (m *Machine) nextTransition(symbol string) bool {
	// TODO: usar um map para procurar a transição
	for _, transition := range m.Transitions {
		if transition.State == m._currentState && transition.Symbol == symbol {
			m._currentState = transition.ResultState
			return true
		}
	}

	return false
}

func contains(state string, states []string) bool {
	for _, s := range states {
		if state == s {
			return true
		}
	}
	return false
}
