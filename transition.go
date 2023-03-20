package statemachine

import "fmt"

type TransitionType int

const (
	INTERNAL TransitionType = iota + 1
	LOCAL
	EXTERNAL
)

func (ty TransitionType) String() string {
	switch ty {
	case INTERNAL:
		return "INTERNAL"
	case LOCAL:
		return "LOCAL"
	case EXTERNAL:
		return "EXTERNAL"
	}
	return ""
}

type Transition[S, E ID, C any] struct {
	source    *state[S, E, C]
	target    *state[S, E, C]
	event     E
	ty        TransitionType
	condition Condition[C]
	action    Action[S, E, C]
}

func (t *Transition[S, E, C]) transit(ctx C, checkCondition bool) (*state[S, E, C], error) {
	err := t.verify()
	if err != nil {
		return nil, err
	}
	if !checkCondition || t.condition == nil || t.condition(ctx) {
		if t.action != nil {
			err = t.action(t.source.id, t.target.id, t.event, ctx)
			if err != nil {
				return nil, err
			}
		}
		return t.target, nil
	}
	return t.source, nil
}

func (t *Transition[S, E, C]) verify() error {
	// 内部流转，两个状态必须是同一个实例
	if t.ty == INTERNAL && t.source != t.target {
		return NewError(fmt.Sprintf("Internal transition source state '%s' and target state '%s' must be same.", t.source, t.target))
	}
	return nil
}

func (t *Transition[S, E, C]) equals(o *Transition[S, E, C]) bool {
	return t.event == o.event && t.source.equals(o.source) && t.target.equals(o.target)
}

func (t *Transition[S, E, C]) String() string {
	return fmt.Sprintf("%s-[%v, %s]->%s", t.source, t.event, t.ty, t.target)
}
