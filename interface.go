package statemachine

type ID interface {
	~string | ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type ExternalTransitionBuilder[S, E ID, C any] interface {
	From(stateId ...S) From[S, E, C]
}

type InternalTransitionBuilder[S, E ID, C any] interface {
	Within(stateId S) To[S, E, C]
}

type From[S, E ID, C any] interface {
	To(stateId S) To[S, E, C]
}

type To[S, E ID, C any] interface {
	On(event E) On[S, E, C]
}

type On[S, E ID, C any] interface {
	When(condition Condition[C]) When[S, E, C]
}

type When[S, E ID, C any] interface {
	Perform(action Action[S, E, C])
}

type Condition[C any] func(ctx C) bool

type Action[S, E ID, C any] func(from S, to S, event E, ctx C) error

type FailCallback[S, E ID, C any] func(sourceState S, event E, ctx C)

type StateMachine[S, E ID, C any] interface {
	// FireEvent 在状态 S 触发事件 E
	FireEvent(stateId S, event E, ctx C) (S, error)
	// GetMachineId 获取状态机id
	GetMachineId() string
	// Verify 验证状态 S 是否可以触发事件 E
	Verify(stateId S, event E) bool
	// ShowStateMachine 打印状态机结构
	ShowStateMachine()
	// GeneratePlantUML 生成PlantUML
	GeneratePlantUML() string
}
