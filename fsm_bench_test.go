package fsm_test

import (
	"testing"

	"github.com/zerjioang/go-fsm"
)

/*
* Benchmark functions start with Benchmark not Test.

* Benchmark functions are run several times by the testing package.
  The value of b.N will increase each time until the benchmark runner
  is satisfied with the stability of the benchmark. This has some important
  ramifications which we’ll investigate later in this article.

* Each benchmark is run for a minimum of 1 second by default.
  If the second has not elapsed when the Benchmark function returns,
  the value of b.N is increased in the sequence 1, 2, 5, 10, 20, 50, …
  and the function run again.

* the for loop is crucial to the operation of the benchmark driver
  it must be: for n := 0; n < b.N; n++

* Add b.ReportAllocs() at the beginning of each benchmark to know allocations
* Add b.SetBytes(1) InitConstants()to measure byte transfers

  Optimization info: https://stackimpact.com/blog/practical-golang-benchmarks/
*/

func BenchmarkFsm(b *testing.B) {
	b.Run("instantiation", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(1)
		for n := 0; n < b.N; n++ {
			_ = fsm.New()
		}
	})
	b.Run("instantiation-ptr", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(1)
		for n := 0; n < b.N; n++ {
			_ = fsm.NewPtr()
		}
	})
	b.Run("add-state", func(b *testing.B) {
		m := fsm.NewPtr()
		b.SetBytes(1)
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			m.AddState("start", fsm.NoStateEvents)
		}
	})
	b.Run("add-state-no-name", func(b *testing.B) {
		m := fsm.NewPtr()
		b.SetBytes(1)
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			m.AddState("", fsm.NoStateEvents)
		}
	})
	b.Run("add-state-no-event", func(b *testing.B) {
		m := fsm.NewPtr()
		b.SetBytes(1)
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			m.AddState("start", fsm.NoStateEvents)
		}
	})
	b.Run("change-state-empty", func(b *testing.B) {
		m := fsm.NewPtr()
		b.SetBytes(1)
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			m.ChangeStateTo("")
		}
	})
	b.Run("change-state", func(b *testing.B) {
		m := fsm.NewPtr()
		m.AddState("start", fsm.NoStateEvents)
		m.AddState("finish", fsm.NoStateEvents)
		m.AddState("a", fsm.NoStateEvents)
		m.AddState("b", fsm.NoStateEvents)
		m.AddState("c", fsm.NoStateEvents)

		m.AddTransaction("toA", "start", "a")
		m.AddTransaction("toFinish", "a", "finish")
		m.SetInitialState("start")
		m.SetFinalState("finish")

		b.SetBytes(1)
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			m.ChangeStateTo("a")
		}
	})
	b.Run("has-state-false", func(b *testing.B) {
		m := fsm.NewPtr()
		b.ReportAllocs()
		b.SetBytes(1)
		for n := 0; n < b.N; n++ {
			m.HasState("start")
		}
	})
	b.Run("has-state-true", func(b *testing.B) {
		m := fsm.NewPtr()
		m.AddState("start", fsm.NoStateEvents)
		m.AddState("finish", fsm.NoStateEvents)
		m.AddState("a", fsm.NoStateEvents)
		m.AddState("b", fsm.NoStateEvents)
		m.AddState("c", fsm.NoStateEvents)
		b.ReportAllocs()
		b.SetBytes(1)
		for n := 0; n < b.N; n++ {
			m.HasState("start")
		}
	})

	b.Run("to-dot", func(b *testing.B) {
		m := fsm.NewPtr()
		m.AddState("start", fsm.NoStateEvents)
		m.AddState("finish", fsm.NoStateEvents)
		m.AddState("a", fsm.NoStateEvents)
		m.AddState("b", fsm.NoStateEvents)
		m.AddState("c", fsm.NoStateEvents)

		m.AddTransaction("toA", "start", "a")
		m.AddTransaction("toFinish", "a", "finish")
		m.SetInitialState("start")
		m.SetFinalState("finish")

		b.SetBytes(1)
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			m.DotGraph()
		}
	})
}
