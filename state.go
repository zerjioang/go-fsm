package fsm

type StateType uint8

const (
	StartState StateType = iota
	EndState
	State
)

type StateEvents struct {
	OnEnter   func()    `json:",omitempty"`
	OnUpdate  func()    `json:",omitempty"`
	OnExit    func()    `json:",omitempty"`
	StateType StateType `json:"type"`
}

func (e StateEvents) dot() string {
	if e.StateType == StartState {
		return `[shape="doublecircle" color="blue"]`
	} else if e.StateType == EndState {
		return `[shape="doublecircle" color="red"]`
	} else if e.StateType == State {
		return `[shape="circle" color="black"]`
	} else {
		return ""
	}
}

var (
	NoStateEvents = StateEvents{}
)
