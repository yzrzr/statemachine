package statemachine

// Builder 用来构建状态机
type Builder[S, E ID, C any] struct {
	stateMachine *stateMachine[S, E, C]
	failCallback FailCallback[S, E, C]
}

// ExternalTransition 外部流转，不同状态之间的流转
func (b *Builder[S, E, C]) ExternalTransition() ExternalTransitionBuilder[S, E, C] {
	return newTransitionBuilder[S, E, C](b.stateMachine, EXTERNAL)
}

// InternalTransition 内部流转，相同状态的流转，只触发事件执行动作
func (b *Builder[S, E, C]) InternalTransition() InternalTransitionBuilder[S, E, C] {
	return newTransitionBuilder[S, E, C](b.stateMachine, INTERNAL)
}

// SetFailCallback 设置失败回调
func (b *Builder[S, E, C]) SetFailCallback(failCallback FailCallback[S, E, C]) {
	b.failCallback = failCallback
}

// Build 构建状态机
func (b *Builder[S, E, C]) Build(machineId string) (StateMachine[S, E, C], error) {
	if b.stateMachine.err != nil {
		return nil, b.stateMachine.err
	}
	b.stateMachine.machineId = machineId
	b.stateMachine.ready = true
	b.stateMachine.failCallback = b.failCallback
	err := registerStateMachine[S, E, C](b.stateMachine)
	if err != nil {
		return nil, err
	}
	return b.stateMachine, nil
}

// NewBuilder 创建一个状态机构建器
func NewBuilder[S, E ID, C any]() *Builder[S, E, C] {
	return &Builder[S, E, C]{
		stateMachine: newStateMachine[S, E, C](make(map[S]*state[S, E, C])),
	}
}
