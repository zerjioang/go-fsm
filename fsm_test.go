package fsm_test

import (
	"fmt"
	"testing"

	"github.com/zerjioang/go-fsm"
)

const (
	expectedJson = `{"states":{"a":{"type":2},"b":{"type":2},"c":{"type":2},"finish":{"type":1},"start":{"type":0}},"transitions":{"a-b":{"name":"toB","from":"a","to":"b"},"a-c":{"name":"toC","from":"a","to":"c"},"b-c":{"name":"moveToC","from":"b","to":"c"},"b-finish":{"name":"toFinish","from":"b","to":"finish"},"c-a":{"name":"backToA","from":"c","to":"a"},"c-finish":{"name":"toFinish","from":"c","to":"finish"},"start-a":{"name":"toA","from":"start","to":"a"}},"current":"start"}`
)

func TestUnitFsm(t *testing.T) {
	t.Run("instantiation", func(t *testing.T) {
		fsm.New()
	})
	t.Run("instantiation-ptr", func(t *testing.T) {
		m := fsm.NewPtr()
		if m == nil {
			t.Error("failed to instantiate fsm via ptr")
		}
	})
	t.Run("add-state", func(t *testing.T) {
		m := fsm.NewPtr()
		m.AddState("start", fsm.NoStateEvents)
	})
	t.Run("add-state-no-name", func(t *testing.T) {
		m := fsm.NewPtr()
		m.AddState("", fsm.NoStateEvents)
	})
	t.Run("add-state-no-event", func(t *testing.T) {
		m := fsm.NewPtr()
		m.AddState("start", fsm.NoStateEvents)
	})
	t.Run("has-state-false", func(t *testing.T) {
		m := fsm.NewPtr()
		result := m.HasState("start")
		if result != false {
			t.Error("failing on hasState() function")
		}
	})
	t.Run("has-state-true", func(t *testing.T) {
		m := fsm.NewPtr()
		m.AddState("start", fsm.NoStateEvents)
		result := m.HasState("start")
		if result != true {
			t.Error("failing on hasState() function")
		}
	})
	t.Run("fsm-json", func(t *testing.T) {
		machine := fsm.New()

		machine.AddState("start", fsm.NoStateEvents)
		machine.AddState("a", fsm.NoStateEvents)
		machine.AddState("b", fsm.NoStateEvents)
		machine.AddState("c", fsm.NoStateEvents)
		machine.AddState("finish", fsm.NoStateEvents)

		machine.AddTransaction("toA", "start", "a")
		machine.AddTransaction("toB", "a", "b")
		machine.AddTransaction("toC", "a", "c")
		machine.AddTransaction("backToA", "c", "a")
		machine.AddTransaction("moveToC", "b", "c")
		machine.AddTransaction("toFinish", "b", "finish")
		machine.AddTransaction("toFinish", "c", "finish")

		machine.SetInitialState("start")
		machine.SetFinalState("finish")

		data, err := machine.Json()

		if err != nil {
			t.Error("failed to generate fsm json", err)
		}
		if string(data) != expectedJson {
			t.Error("unexpected json result")
			t.Log(string(data))
		}
	})

	t.Run("fsm-json-load", func(t *testing.T) {
		machine := fsm.New()

		err := machine.Load([]byte(expectedJson))
		if err != nil {
			t.Error("failed to load fsm from json", err)
		}

		data, err := machine.Json()
		if string(data) != expectedJson {
			t.Error("unexpected json result")
			t.Log(string(data))
		}
	})
}

func TestFsmExamples(t *testing.T) {
	t.Run("start-a-finish", func(t *testing.T) {
		machine := fsm.New()
		//define three example states
		machine.AddState("start", fsm.StateEvents{
			OnEnter: func() {
				fmt.Println("start state entered")
			},
			OnUpdate: nil,
			OnExit: func() {
				fmt.Println("start state exited")
			},
		})
		machine.AddState("a", fsm.StateEvents{
			OnEnter: func() {
				fmt.Println("a state entered")
			},
			OnUpdate: nil,
			OnExit: func() {
				fmt.Println("a state exited")
			},
		})
		machine.AddState("finish", fsm.StateEvents{
			OnEnter: func() {
				fmt.Println("finish state entered")
			},
			OnUpdate: nil,
			OnExit: func() {
				fmt.Println("finish state exited")
			},
		})
		//add new transaction from start->finish
		machine.AddTransaction("toA", "start", "a")
		machine.AddTransaction("toFinish", "a", "finish")

		// set our initial state to start
		machine.SetInitialState("start")

		//execute toFinish transaction
		machine.SetFinalState("finish")

		machine.ChangeStateTo("a")

		/*
			digraph fsm {
				"start" -> "finish" [ label = "toFinish" ];

				"start";
				"finish";
			}
		*/
		dot := machine.DotGraph()
		fmt.Println(dot)
	})
	t.Run("start-a-b-c-finish", func(t *testing.T) {
		machine := fsm.New()
		//define three example states
		machine.AddState("start", fsm.StateEvents{
			OnEnter: func() {
				fmt.Println("start state entered")
			},
			OnUpdate: nil,
			OnExit: func() {
				fmt.Println("start state exited")
			},
		})
		machine.AddState("a", fsm.StateEvents{
			OnEnter: func() {
				fmt.Println("a state entered")
			},
			OnUpdate: nil,
			OnExit: func() {
				fmt.Println("a state exited")
			},
		})
		machine.AddState("b", fsm.StateEvents{
			OnEnter: func() {
				fmt.Println("b state entered")
			},
			OnUpdate: nil,
			OnExit: func() {
				fmt.Println("b state exited")
			},
		})
		machine.AddState("c", fsm.StateEvents{
			OnEnter: func() {
				fmt.Println("c state entered")
			},
			OnUpdate: nil,
			OnExit: func() {
				fmt.Println("c state exited")
			},
		})
		machine.AddState("finish", fsm.StateEvents{
			OnEnter: func() {
				fmt.Println("finish state entered")
			},
			OnUpdate: nil,
			OnExit: func() {
				fmt.Println("finish state exited")
			},
		})
		//add new transaction from start->finish
		machine.AddTransaction("toA", "start", "a")
		machine.AddTransaction("toB", "a", "b")
		machine.AddTransaction("toC", "a", "c")
		machine.AddTransaction("backToA", "c", "a")
		machine.AddTransaction("moveToC", "b", "c")
		machine.AddTransaction("toFinish", "b", "finish")
		machine.AddTransaction("toFinish", "c", "finish")

		// set our initial state to start
		machine.SetInitialState("start")

		//execute toFinish transaction
		machine.SetFinalState("finish")

		machine.ChangeStateTo("a")

		/*
			digraph fsm {
				"start" -> "finish" [ label = "toFinish" ];

				"start";
				"finish";
			}
		*/
		dot := machine.DotGraph()
		fmt.Println(dot)
	})
}
