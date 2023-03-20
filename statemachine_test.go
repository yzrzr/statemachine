package statemachine

import (
	"fmt"
	"sync"
	"testing"
)

type States int

const (
	STATE1 = iota + 1
	STATE2
	STATE3
	STATE4
)

func (s States) String() string {
	switch s {
	case STATE1:
		return "STATE1"
	case STATE2:
		return "STATE2"
	case STATE3:
		return "STATE3"
	case STATE4:
		return "STATE4"
	default:
		return "UNKNOWN"
	}
}

type Events int

const (
	EVENT1 = iota + 1
	EVENT2
	EVENT3
	EVENT4
	INTERNAL_EVENT
)

func (e Events) String() string {
	switch e {
	case EVENT1:
		return "EVENT1"
	case EVENT2:
		return "EVENT2"
	case EVENT3:
		return "EVENT3"
	case EVENT4:
		return "EVENT4"
	case INTERNAL_EVENT:
		return "INTERNAL_EVENT"
	default:
		return "UNKNOWN"
	}
}

type Context1 struct {
	op string
	id int
}

func perform(from States, to States, event Events, ctx Context1) error {
	fmt.Printf("from: %v to: %v event: %v ctx: %v\n", from, to, event, ctx)
	return nil
}

func conditionFalse(ctx Context1) bool {
	return false
}

func conditionTrue(ctx Context1) bool {
	return true
}

func performInt(from States, to States, event Events, ctx int) error {
	fmt.Printf("from: %v to: %v event: %v ctx: %v\n", from, to, event, ctx)
	return nil
}

var context = Context1{"creat", 2}

func Test_external(t *testing.T) {
	builder := NewBuilder[States, Events, Context1]()
	builder.ExternalTransition().From(STATE1).To(STATE2).On(EVENT1).
		When(conditionTrue).Perform(perform)
	machine, err := builder.Build("TestStateMachine-external")
	if err != nil {
		t.Error(err)
	}
	target, err := machine.FireEvent(STATE1, EVENT1, context)
	if err != nil {
		t.Error(err)
	}
	if target != STATE2 {
		t.Errorf("FireEvent() = %v, want %v", target, STATE2)
	}
}

func Test_fail(t *testing.T) {
	builder := NewBuilder[States, Events, Context1]()
	builder.ExternalTransition().From(STATE1).To(STATE2).On(EVENT1).
		When(conditionTrue).Perform(perform)
	builder.SetFailCallback(func(sourceState States, event Events, ctx Context1) {
		fmt.Printf("当前状态：%v 无法触发事件：%v", sourceState, event)
	})
	machine, err := builder.Build("TestStateMachine-fail")
	if err != nil {
		t.Error(err)
	}
	_, err = machine.FireEvent(STATE2, EVENT1, context)
	if err != nil {
		t.Errorf("FireEvent error is = %v, want nil", err)
	}
}

func Test_verify(t *testing.T) {
	builder := NewBuilder[States, Events, Context1]()
	builder.ExternalTransition().From(STATE1).To(STATE2).On(EVENT1).
		When(conditionTrue).Perform(perform)
	machine, err := builder.Build("TestStateMachine-verify")
	if err != nil {
		t.Error(err)
	}
	b := machine.Verify(STATE1, EVENT1)
	if !b {
		t.Error("Verify() = false, want true")
	}
	b = machine.Verify(STATE1, EVENT2)
	if b {
		t.Error("Verify() = true, want false")
	}
}

func Test_externals(t *testing.T) {
	builder := NewBuilder[States, Events, Context1]()
	builder.ExternalTransition().From(STATE1, STATE2, STATE3).To(STATE4).On(EVENT1).
		When(conditionTrue).Perform(perform)
	machine, err := builder.Build("TestStateMachine-externals")
	if err != nil {
		t.Error(err)
	}
	target, err := machine.FireEvent(STATE2, EVENT1, context)
	if err != nil {
		t.Error(err)
	}
	if target != STATE4 {
		t.Errorf("FireEvent() = %v, want %v", target, STATE4)
	}
}

func Test_internal(t *testing.T) {
	builder := NewBuilder[States, Events, Context1]()
	builder.InternalTransition().Within(STATE1).On(INTERNAL_EVENT).
		When(conditionTrue).Perform(perform)
	machine, err := builder.Build("TestStateMachine-internal")
	if err != nil {
		t.Error(err)
	}
	target, err := machine.FireEvent(STATE1, EVENT1, context)
	if target != STATE1 {
		t.Errorf("FireEvent() = %v, want %v", target, STATE1)
	}
}

func Test_externalAndInternal(t *testing.T) {
	machine := buildStateMachine("TestStateMachine-externalAndInternal")
	target, err := machine.FireEvent(STATE1, EVENT1, context)
	if err != nil {
		t.Error(err)
	}
	if target != STATE2 {
		t.Errorf("FireEvent() = %v, want %v", target, STATE2)
	}
	target, err = machine.FireEvent(STATE2, INTERNAL_EVENT, context)
	if err != nil {
		t.Error(err)
	}
	if target != STATE2 {
		t.Errorf("FireEvent() = %v, want %v", target, STATE2)
	}
	target, err = machine.FireEvent(STATE2, EVENT2, context)
	if err != nil {
		t.Error(err)
	}
	if target != STATE1 {
		t.Errorf("FireEvent() = %v, want %v", target, STATE1)
	}
	target, err = machine.FireEvent(STATE1, EVENT3, context)
	if err != nil {
		t.Error(err)
	}
	if target != STATE3 {
		t.Errorf("FireEvent() = %v, want %v", target, STATE3)
	}
}

func Test_goroutine(t *testing.T) {
	group := sync.WaitGroup{}
	buildStateMachine("TestStateMachine-goroutine")
	group.Add(30)
	for i := 0; i < 10; i++ {
		go func() {
			defer group.Done()
			machine := getStateMachine[States, Events, Context1]("TestStateMachine-goroutine")
			target, err := machine.FireEvent(STATE1, EVENT1, context)
			if err != nil {
				t.Error(err)
			}
			if target != STATE2 {
				t.Errorf("FireEvent() = %v, want %v", target, STATE2)
			}
		}()
	}
	for i := 0; i < 10; i++ {
		go func() {
			defer group.Done()
			machine := getStateMachine[States, Events, Context1]("TestStateMachine-goroutine")
			target, err := machine.FireEvent(STATE1, EVENT4, context)
			if err != nil {
				t.Error(err)
			}
			if target != STATE4 {
				t.Errorf("FireEvent() = %v, want %v", target, STATE4)
			}
		}()
	}
	for i := 0; i < 10; i++ {
		go func() {
			defer group.Done()
			machine := getStateMachine[States, Events, Context1]("TestStateMachine-goroutine")
			target, err := machine.FireEvent(STATE1, EVENT3, context)
			if err != nil {
				t.Error(err)
			}
			if target != STATE3 {
				t.Errorf("FireEvent() = %v, want %v", target, STATE3)
			}
		}()
	}
	group.Wait()
}

func Test_plantUML(t *testing.T) {
	machine := buildStateMachine("TestStateMachine-plantUML")
	v := machine.GeneratePlantUML()
	fmt.Println(v)
}

func Test_conditionFalse(t *testing.T) {
	builder := NewBuilder[States, Events, Context1]()
	builder.ExternalTransition().
		From(STATE1).To(STATE2).On(EVENT1).
		When(conditionFalse).Perform(perform)

	machine, err := builder.Build("TestStateMachine-conditionFalse")
	if err != nil {
		t.Error(err)
	}
	target, err := machine.FireEvent(STATE1, EVENT1, context)
	if err != nil {
		t.Error(err)
	}
	if target != STATE1 {
		t.Errorf("FireEvent() = %v, want %v", target, STATE1)
	}
}

func Test_duplicateTransition(t *testing.T) {
	builder := NewBuilder[States, Events, Context1]()
	builder.ExternalTransition().
		From(STATE1).To(STATE2).On(EVENT1).
		When(conditionTrue).Perform(perform)

	builder.ExternalTransition().
		From(STATE1).To(STATE2).On(EVENT1).
		When(conditionTrue).Perform(perform)

	_, err := builder.Build("TestStateMachine-duplicateTransition")
	if !IsStateMachineError(err) {
		t.Errorf("Build err = %v, want StateMachineError", err)
	}
}

func Test_duplicateMachine(t *testing.T) {
	builder := NewBuilder[States, Events, Context1]()
	builder.ExternalTransition().
		From(STATE1).To(STATE2).On(EVENT1).
		When(conditionTrue).Perform(perform)

	_, err := builder.Build("TestStateMachine-duplicateMachine")
	if err != nil {
		t.Error(err)
	}
	_, err = builder.Build("TestStateMachine-duplicateMachine")
	if !IsStateMachineError(err) {
		t.Errorf("Build err = %v, want StateMachineError", err)
	}
}

func Test_choice(t *testing.T) {
	builder := NewBuilder[States, Events, int]()
	builder.InternalTransition().Within(STATE1).On(EVENT1).When(func(ctx int) bool {
		return ctx == 1
	}).Perform(performInt)
	builder.ExternalTransition().From(STATE1).To(STATE2).On(EVENT1).When(func(ctx int) bool {
		return ctx == 2
	}).Perform(performInt)
	builder.ExternalTransition().From(STATE1).To(STATE3).On(EVENT1).When(func(ctx int) bool {
		return ctx == 3
	}).Perform(performInt)
	machine, err := builder.Build("TestStateMachine-choice")
	if err != nil {
		t.Error(err)
	}
	target, err := machine.FireEvent(STATE1, EVENT1, 1)
	if err != nil {
		t.Error(err)
	}
	if target != STATE1 {
		t.Errorf("FireEvent() = %v, want %v", target, STATE1)
	}
	target, err = machine.FireEvent(STATE1, EVENT1, 2)
	if err != nil {
		t.Error(err)
	}
	if target != STATE2 {
		t.Errorf("FireEvent() = %v, want %v", target, STATE2)
	}
	target, err = machine.FireEvent(STATE1, EVENT1, 3)
	if err != nil {
		t.Error(err)
	}
	if target != STATE3 {
		t.Errorf("FireEvent() = %v, want %v", target, STATE3)
	}
}

func buildStateMachine(machineId string) StateMachine[States, Events, Context1] {
	builder := NewBuilder[States, Events, Context1]()
	builder.ExternalTransition().From(STATE1).To(STATE2).On(EVENT1).
		When(conditionTrue).Perform(perform)
	builder.InternalTransition().Within(STATE2).On(INTERNAL_EVENT).
		When(conditionTrue).Perform(perform)
	builder.ExternalTransition().From(STATE2).To(STATE1).On(EVENT2).
		When(conditionTrue).Perform(perform)
	builder.ExternalTransition().From(STATE1).To(STATE3).On(EVENT3).
		When(conditionTrue).Perform(perform)
	builder.ExternalTransition().From(STATE1, STATE2, STATE3).To(STATE4).On(EVENT4).
		When(conditionTrue).Perform(perform)

	_, err := builder.Build(machineId)
	if err != nil {
		panic(err)
	}
	machine := getStateMachine[States, Events, Context1](machineId)

	machine.ShowStateMachine()
	return machine
}
