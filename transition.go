package fsm

type StateTransition struct {
	TransactionName string `json:"name"`
	FromState       string `json:"from"`
	ToState         string `json:"to"`
}
