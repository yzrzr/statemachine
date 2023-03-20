package statemachine

import "fmt"

var registerMap = make(map[string]any)

func registerStateMachine[S, E ID, C any](stateMachine StateMachine[S, E, C]) error {
	machineId := stateMachine.GetMachineId()
	if _, ok := registerMap[machineId]; ok {
		return NewError(fmt.Sprintf("状态机 [%s] 已经构建, 不需要重新构建", machineId))
	}
	stateMachine.GetMachineId()
	registerMap[machineId] = stateMachine
	return nil
}

func getStateMachine[S, E ID, C any](machineId string) StateMachine[S, E, C] {
	if v, ok := registerMap[machineId]; ok {
		return v.(StateMachine[S, E, C])
	}
	return nil
}
