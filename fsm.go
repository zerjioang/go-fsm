package fsm

import (
	"bytes"
	"encoding/json"
	"fmt"
)

/*
implementation of a concurrent safe finite state machine

In a replicated state machine, if a transaction is valid,
a set of inputs will cause the state of the system to
transition to the next state. A transaction is an atomic
operation on a database. This means the operations either
complete in full or never complete at all. The set of
transactions maintained in a replicated state machine is
known as a “transaction log.”
*/

type FiniteStateMachine struct {
	namedStates        map[string]StateEvents
	namedTransitions   map[string]StateTransition
	currentState       string
	currentStateEvents StateEvents
}

func (machine *FiniteStateMachine) AddState(stateName string, stateEvent StateEvents) {
	stateEvent.StateType = State
	machine.namedStates[stateName] = stateEvent
}

func (machine *FiniteStateMachine) DeleteState(stateName string) {
	delete(machine.namedStates, stateName)
}

// HasState returns if the FSM has a StateEvents added with given name
func (machine FiniteStateMachine) HasState(name string) bool {
	_, hasKey := machine.namedStates[name]
	return hasKey
}

// adds new transaction to the Finite state machine specifing the transaction name, origin and destination
func (machine *FiniteStateMachine) AddTransaction(transactioName string, fromState string, toState string) {
	tx := StateTransition{TransactionName: transactioName, FromState: fromState, ToState: toState}
	key := fromState + "-" + toState
	machine.namedTransitions[key] = tx
}

// deletes an existing transaction from the Finite state machine specifing the transaction name
func (machine *FiniteStateMachine) DeleteTransaction(transactioName string) {
	delete(machine.namedTransitions, transactioName)
}

func (machine *FiniteStateMachine) SetInitialState(stateName string) {
	if machine.currentState == "" {
		//an initial state setup has been requested
		//fmt.Println("setting FSM initial state to", stateName)
		machine.currentState = stateName
		//trigger requested state onEnter event
		state, ok := machine.namedStates[stateName]
		if ok {
			machine.currentStateEvents = state
			if state.OnEnter != nil {
				//fmt.Println("triggering ", stateName, "OnEnter() event")
				state.OnEnter()
			}
			state.StateType = StartState
			machine.namedStates[stateName] = state
		}
	}
}

func (machine *FiniteStateMachine) SetFinalState(stateName string) {
	//an initial state setup has been requested
	//fmt.Println("setting FSM final state to", stateName)
	state, ok := machine.namedStates[stateName]
	if ok {
		state.StateType = EndState
		machine.namedStates[stateName] = state
	}
}

// changes the current state of the Finite state machine to requested state if allowed by transactions
func (machine *FiniteStateMachine) ChangeStateTo(stateName string) {
	//a typical FSM state change
	// check if we have a valid transaction from current state to requested state
	_, valid := machine.HasValidTransaction(machine.currentState, stateName)
	if valid {
		//fmt.Println("changing state from", machine.currentState, "to", stateName, "using transaction", txData.TransactionName)
		//trigger previous state exit event
		if machine.currentStateEvents.OnExit != nil {
			machine.currentStateEvents.OnExit()
		}
		//trigger requested state onEnter event
		state, ok := machine.namedStates[stateName]
		if ok {
			machine.currentStateEvents = state
			if state.OnEnter != nil {
				//fmt.Println("triggering ", stateName, "OnEnter() event")
				state.OnEnter()
			}
		}
		machine.currentState = stateName
	} else {
		//fmt.Println("there is no a direct transaction from", machine.currentState, "to", stateName)
	}
}
func (machine FiniteStateMachine) HasValidTransaction(from string, to string) (StateTransition, bool) {
	//fmt.Println("checkinf if valid transaction exist from", from, "to", to)
	key := from + "-" + to
	tx, found := machine.namedTransitions[key]
	return tx, found
}

// Visualize outputs a visualization of a FSM in Graphviz format.
func (machine FiniteStateMachine) DotGraph() string {
	var buf bytes.Buffer
	states := make(map[string]int)

	buf.WriteString(`digraph fsm {
	size ="4,4";
	node [shape=circle,fontsize=12,fixedsize=true,width=0.8];
	edge [fontsize=6];
	rankdir=LR;
`)

	// make sure the initial state is at top
	for k, v := range machine.namedTransitions {
		states[k]++
		buf.WriteString(fmt.Sprintf(`    "%s" -> "%s" [ label = "%s" ];`, v.FromState, v.ToState, v.TransactionName))
		buf.WriteString("\n")
	}

	buf.WriteString("\n")

	for k, v := range machine.namedStates {
		buf.WriteString(fmt.Sprintf(`    "%s" %s;`, k, v.dot()))
		buf.WriteString("\n")
	}
	buf.WriteString(fmt.Sprintln("}"))
	return buf.String()
}

// return machine content encoded as Json
func (machine FiniteStateMachine) Json() ([]byte, error) {
	type internalMachine struct {
		NamedStates       map[string]StateEvents     `json:"states"`
		NamedTransactions map[string]StateTransition `json:"transitions"`
		CurrentState      string                     `json:"current"`
	}
	var internal internalMachine
	internal = internalMachine{
		NamedStates:       machine.namedStates,
		NamedTransactions: machine.namedTransitions,
		CurrentState:      machine.currentState,
	}
	return json.Marshal(internal)
}

//buils the machine from json input
func (machine *FiniteStateMachine) Load(raw []byte) error {
	type internalMachine struct {
		NamedStates      map[string]StateEvents     `json:"states"`
		NamedTransitions map[string]StateTransition `json:"transitions"`
		CurrentState     string                     `json:"current"`
	}
	var internal internalMachine
	internal = internalMachine{}
	err := json.Unmarshal(raw, &internal)
	if err != nil {
		return err
	} else {
		machine.namedStates = internal.NamedStates
		machine.namedTransitions = internal.NamedTransitions
		machine.currentState = "start"
		machine.currentStateEvents = StateEvents{}
		return nil
	}
}

// State returns current state of the FSM
func (machine *FiniteStateMachine) State() string {
	return machine.currentState
}

// Create a new Finite state machine and returns it as struct
func New() FiniteStateMachine {
	m := FiniteStateMachine{}
	m.namedStates = make(map[string]StateEvents, 0)
	m.namedTransitions = make(map[string]StateTransition, 0)
	return m
}

// Create a new Finite state machine and returns it as pointer
func NewPtr() *FiniteStateMachine {
	m := new(FiniteStateMachine)
	m.namedStates = make(map[string]StateEvents, 0)
	m.namedTransitions = make(map[string]StateTransition, 0)
	return m
}
