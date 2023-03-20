package statemachine

import "fmt"

func newEventTransitions[S, E ID, C any]() *eventTransitions[S, E, C] {
	return &eventTransitions[S, E, C]{
		eventTransitions: make(map[E][]*Transition[S, E, C]),
	}
}

type eventTransitions[S, E ID, C any] struct {
	eventTransitions map[E][]*Transition[S, E, C]
}

func (e *eventTransitions[S, E, C]) put(event E, transition *Transition[S, E, C]) error {
	if li, ok := e.eventTransitions[event]; ok {
		if err := e.verify(li, transition); err != nil {
			return err
		}
		e.eventTransitions[event] = append(li, transition)
	} else {
		e.eventTransitions[event] = []*Transition[S, E, C]{transition}
	}
	return nil
}

func (e *eventTransitions[S, E, C]) get(event E) []*Transition[S, E, C] {
	if li, ok := e.eventTransitions[event]; ok {
		return li
	}
	return nil
}

func (e *eventTransitions[S, E, C]) all() []*Transition[S, E, C] {
	res := make([]*Transition[S, E, C], 0, 8)
	for _, transitions := range e.eventTransitions {
		res = append(res, transitions...)
	}
	return res
}

func (e *eventTransitions[S, E, C]) verify(existingTransitions []*Transition[S, E, C], transition *Transition[S, E, C]) error {
	for _, v := range existingTransitions {
		if v.equals(transition) {
			return NewError(fmt.Sprintf("%v already Exist, you can not add another one", transition))
		}
	}
	return nil
}
