package fsm

type StateType uint8

const (
	StartState StateType = iota
	EndState
	State
)

type StateEvents struct {
	OnEnter   func()
	OnUpdate  func()
	OnExit    func()
	StateType StateType
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
