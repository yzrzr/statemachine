package statemachine

import (
	"fmt"
	"strings"
)

type stateMachine[S, E ID, C any] struct {
	machineId    string
	stateMap     stateMap[S, E, C]
	ready        bool
	failCallback FailCallback[S, E, C]
	err          error
}

func newStateMachine[S, E ID, C any](stateMap stateMap[S, E, C]) *stateMachine[S, E, C] {
	return &stateMachine[S, E, C]{
		stateMap: stateMap,
	}
}

func (s *stateMachine[S, E, C]) GetMachineId() string {
	return s.machineId
}

func (s *stateMachine[S, E, C]) FireEvent(stateId S, event E, ctx C) (r S, err error) {
	if !s.ready {
		return r, NewError("状态机尚未构建，不能工作")
	}
	transition := s.routeTransition(stateId, event, ctx)
	// 没有找到对应的transition，可能是没定义，也可能是条件满足
	if transition == nil {
		if s.failCallback != nil {
			s.failCallback(stateId, event, ctx)
		}
		return stateId, nil
	}
	state, err := transition.transit(ctx, false)
	if err != nil {
		return r, err
	}
	return state.id, nil
}

func (s *stateMachine[S, E, C]) Verify(stateId S, event E) bool {
	transitions := s.getEventTransitions(stateId, event)
	return len(transitions) != 0
}

func (s *stateMachine[S, E, C]) ShowStateMachine() {
	builder := strings.Builder{}
	builder.WriteString("-----StateMachine:" + s.machineId + "-------")
	for stateId, state := range s.stateMap {
		builder.WriteString(fmt.Sprintf("State: %v\n", stateId))
		for _, transition := range state.getAllEventTransitions() {
			builder.WriteString(fmt.Sprintf("    Transition:%s\n", transition))
		}
	}
	builder.WriteString("------------------------")
	fmt.Println(builder.String())
}

func (s *stateMachine[S, E, C]) GeneratePlantUML() string {
	builder := strings.Builder{}
	builder.WriteString("@startuml\n")
	for _, state := range s.stateMap {
		for _, transition := range state.getAllEventTransitions() {
			builder.WriteString(fmt.Sprintf("%v --> %v : %v\n", transition.source.id, transition.target.id, transition.event))
		}
	}
	builder.WriteString("@enduml")
	return builder.String()
}

func (s *stateMachine[S, E, C]) routeTransition(stateId S, event E, ctx C) *Transition[S, E, C] {
	transitions := s.getEventTransitions(stateId, event)
	if len(transitions) == 0 {
		return nil
	}
	var transit *Transition[S, E, C]
	for _, transition := range transitions {
		if transition.condition == nil {
			transit = transition
		} else if transition.condition(ctx) {
			transit = transition
			break
		}
	}
	return transit
}

func (s *stateMachine[S, E, C]) createAndGetState(stateId S) *state[S, E, C] {
	return s.stateMap.createAndGet(stateId)
}

func (s *stateMachine[S, E, C]) getEventTransitions(stateId S, event E) []*Transition[S, E, C] {
	sourceState := s.stateMap.createAndGet(stateId)
	return sourceState.getEventTransitions(event)
}

var _ StateMachine[int, int, int] = (*stateMachine[int, int, int])(nil)
