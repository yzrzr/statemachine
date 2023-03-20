package statemachine

type stateMap[S, E ID, C any] map[S]*state[S, E, C]

func (s stateMap[S, E, C]) createAndGet(stateId S) *state[S, E, C] {
	if _, ok := s[stateId]; !ok {
		s[stateId] = newState[S, E, C](stateId)
	}
	return s[stateId]
}

func newTransitionBuilder[S, E ID, C any](stateMachine *stateMachine[S, E, C], transitionType TransitionType) *transitionBuilder[S, E, C] {
	return &transitionBuilder[S, E, C]{
		stateMachine:   stateMachine,
		transitionType: transitionType,
	}
}

type transitionBuilder[S, E ID, C any] struct {
	sources        []*state[S, E, C]
	target         *state[S, E, C]
	stateMachine   *stateMachine[S, E, C]
	transitions    []*Transition[S, E, C]
	transitionType TransitionType
}

func (t *transitionBuilder[S, E, C]) From(stateIds ...S) From[S, E, C] {
	for _, stateId := range stateIds {
		t.sources = append(t.sources, t.stateMachine.createAndGetState(stateId))
	}
	return t
}

func (t *transitionBuilder[S, E, C]) To(stateId S) To[S, E, C] {
	t.target = t.stateMachine.createAndGetState(stateId)
	return t
}

func (t *transitionBuilder[S, E, C]) Within(stateId S) To[S, E, C] {
	t.sources = []*state[S, E, C]{t.stateMachine.createAndGetState(stateId)}
	t.target = t.sources[0]
	return t
}

func (t *transitionBuilder[S, E, C]) On(event E) On[S, E, C] {
	for _, source := range t.sources {
		transition, err := source.addTransition(event, t.target, t.transitionType)
		if err != nil {
			t.stateMachine.err = err
			break
		}
		t.transitions = append(t.transitions, transition)
	}
	return t
}

func (t *transitionBuilder[S, E, C]) When(condition Condition[C]) When[S, E, C] {
	for _, transition := range t.transitions {
		transition.condition = condition
	}
	return t
}

func (t *transitionBuilder[S, E, C]) Perform(action Action[S, E, C]) {
	for _, transition := range t.transitions {
		transition.action = action
	}
}

var _ ExternalTransitionBuilder[int, int, int] = (*transitionBuilder[int, int, int])(nil)
var _ InternalTransitionBuilder[int, int, int] = (*transitionBuilder[int, int, int])(nil)
var _ From[int, int, int] = (*transitionBuilder[int, int, int])(nil)
var _ To[int, int, int] = (*transitionBuilder[int, int, int])(nil)
var _ On[int, int, int] = (*transitionBuilder[int, int, int])(nil)
var _ When[int, int, int] = (*transitionBuilder[int, int, int])(nil)
