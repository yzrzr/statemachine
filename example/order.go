package main

import (
	"fmt"

	"github.com/yzrzr/statemachine"
)

type OrderStatus int

const (
	None OrderStatus = iota + 1
	WaitPayment
	WaitDeliver
	WaitConfirm
	WaitEvaluation
	Complete
	CancelOrder
)

func (o OrderStatus) String() string {
	switch o {
	case None:
		return "None"
	case WaitPayment:
		return "WaitPayment"
	case WaitDeliver:
		return "WaitDeliver"
	case WaitConfirm:
		return "WaitConfirm"
	case WaitEvaluation:
		return "WaitEvaluation"
	case Complete:
		return "Complete"
	case CancelOrder:
		return "CancelOrder"
	default:
		return ""
	}
}

type OrderEvent int

const (
	CreateEvent OrderEvent = iota + 1
	ChangePriceEvent
	PaymentEvent
	DeliverEvent
	ConfirmEvent
	EvaluationEvent
	CancelEvent
)

func (o OrderEvent) String() string {
	switch o {
	case CreateEvent:
		return "CreateEvent"
	case ChangePriceEvent:
		return "ChangePriceEvent"
	case PaymentEvent:
		return "PaymentEvent"
	case DeliverEvent:
		return "DeliverEvent"
	case ConfirmEvent:
		return "ConfirmEvent"
	case EvaluationEvent:
		return "EvaluationEvent"
	case CancelEvent:
		return "CancelEvent"
	default:
		return ""
	}
}

type Order struct {
	Status OrderStatus
}

func main() {
	machine := createOrderStateMachine()
	order := &Order{
		Status: None,
	}
	target, err := machine.FireEvent(None, CreateEvent, order)
	fmt.Println(target, err)
	target, err = machine.FireEvent(order.Status, PaymentEvent, order)
	fmt.Println(target, err)
	target, err = machine.FireEvent(order.Status, DeliverEvent, order)
	fmt.Println(target, err)
	target, err = machine.FireEvent(order.Status, ConfirmEvent, order)
	fmt.Println(target, err)
	target, err = machine.FireEvent(order.Status, EvaluationEvent, order)
	fmt.Println(target, err)
	// 状态转移失败，状态不变
	target, err = machine.FireEvent(order.Status, EvaluationEvent, order)
	fmt.Println(target, err)

	uml := machine.GeneratePlantUML()
	fmt.Println(uml)
}

func createOrderStateMachine() statemachine.StateMachine[OrderStatus, OrderEvent, *Order] {
	builder := statemachine.NewBuilder[OrderStatus, OrderEvent, *Order]()
	builder.SetFailCallback(func(sourceState OrderStatus, event OrderEvent, ctx *Order) {
		fmt.Println("状态转移失败")
	})
	// 创建订单，触发创建事件，状态转移到等待支付
	builder.ExternalTransition().From(None).To(WaitPayment).On(CreateEvent).
		When(func(ctx *Order) bool {
			return ctx.Status == None
		}).Perform(func(from OrderStatus, to OrderStatus, event OrderEvent, ctx *Order) error {
		fmt.Println("订单创建成功，等待支付")
		ctx.Status = to
		return nil
	})
	// 商户改价，触发改价事件，状态不变
	builder.InternalTransition().Within(WaitPayment).On(ChangePriceEvent).
		When(func(ctx *Order) bool {
			return ctx.Status == WaitPayment
		}).Perform(func(from OrderStatus, to OrderStatus, event OrderEvent, ctx *Order) error {
		fmt.Println("商户改价成功，等待支付")
		return nil
	})
	// 支付，触发支付事件，状态转移到等待发货
	builder.ExternalTransition().From(WaitPayment).To(WaitDeliver).On(PaymentEvent).
		When(func(ctx *Order) bool {
			return ctx.Status == WaitPayment
		}).Perform(func(from OrderStatus, to OrderStatus, event OrderEvent, ctx *Order) error {
		fmt.Println("订单支付成功，等待发货")
		ctx.Status = to
		return nil
	})
	// 取消订单，触发取消事件，状态转移到交易关闭
	builder.ExternalTransition().From(WaitPayment).To(CancelOrder).On(CancelEvent).
		When(func(ctx *Order) bool {
			return ctx.Status == WaitPayment
		}).Perform(func(from OrderStatus, to OrderStatus, event OrderEvent, ctx *Order) error {
		fmt.Println("用户取消订单，交易关闭")
		ctx.Status = to
		return nil
	})
	// 发货，触发发货事件，状态转移到等待收货
	builder.ExternalTransition().From(WaitDeliver).To(WaitConfirm).On(DeliverEvent).
		When(func(ctx *Order) bool {
			return ctx.Status == WaitDeliver
		}).Perform(func(from OrderStatus, to OrderStatus, event OrderEvent, ctx *Order) error {
		fmt.Println("订单发货成功，等待用户确认收货")
		ctx.Status = to
		return nil
	})
	// 用户确认发货，触发收货事件，状态转移到等待评价
	builder.ExternalTransition().From(WaitConfirm).To(WaitEvaluation).On(ConfirmEvent).
		When(func(ctx *Order) bool {
			return ctx.Status == WaitConfirm
		}).Perform(func(from OrderStatus, to OrderStatus, event OrderEvent, ctx *Order) error {
		fmt.Println("用户确认发货成功，等待用户评价")
		ctx.Status = to
		return nil
	})
	// 用户评价，触发评价事件，状态转移到交易完成
	builder.ExternalTransition().From(WaitEvaluation).To(Complete).On(EvaluationEvent).
		When(func(ctx *Order) bool {
			return ctx.Status == WaitEvaluation
		}).Perform(func(from OrderStatus, to OrderStatus, event OrderEvent, ctx *Order) error {
		fmt.Println("用户评价成功，交易完成")
		ctx.Status = to
		return nil
	})

	// 构建状态机，定义状态机的名称
	machine, err := builder.Build("stateMachine-order")
	if err != nil {
		panic(err)
	}
	return machine
}
