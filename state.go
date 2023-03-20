package statemachine

import "fmt"

func newState[S, E ID, C any](stateId S) *state[S, E, C] {
	return &state[S, E, C]{
		id:               stateId,
		eventTransitions: newEventTransitions[S, E, C](),
	}
}

type state[S, E ID, C any] struct {
	id               S
	eventTransitions *eventTransitions[S, E, C]
}

func (s *state[S, E, C]) addTransition(event E, target *state[S, E, C], transitionType TransitionType) (*Transition[S, E, C], error) {
	transition := &Transition[S, E, C]{
		source: s,
		target: target,
		event:  event,
		ty:     transitionType,
	}
	err := s.eventTransitions.put(event, transition)
	if err != nil {
		return nil, err
	}
	return transition, nil
}

func (s *state[S, E, C]) getEventTransitions(event E) []*Transition[S, E, C] {
	return s.eventTransitions.get(event)
}

func (s *state[S, E, C]) getAllEventTransitions() []*Transition[S, E, C] {
	return s.eventTransitions.all()
}

func (s *state[S, E, C]) String() string {
	return fmt.Sprintf("%v", s.id)
}

func (s *state[S, E, C]) equals(o *state[S, E, C]) bool {
	return s.id == o.id
}
