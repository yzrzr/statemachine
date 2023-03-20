
# State machine

状态机的go语言实现




# Quickstart

构建一个简单的状态机

1、定义状态和事件
```go
type States int
const (
    STATE1 = iota + 1
    STATE2
)
type Events int
const (
    EVENT1 = iota + 1
)
```
2、定义一个实体
```go
type Entity struct{
	Name   string
	Status States
}
```
2、获取构建器
```go
builder := NewBuilder[States, Events, Entity]()
```
3、定义状态之间的流转关系
```go
builder.ExternalTransition().From(STATE1).To(STATE2).On(EVENT1).
    When(func(ctx Context1) bool {
    return true
    }).Perform(func(from States, to States, event Events, ctx Context1) error {
    return nil
    })
```
5、创建状态机
```go
machine, err := builder.Build("StateMachineName")
```
6、触发事件让状态机开始工作
```go
target, err := machine.FireEvent(STATE1, EVENT1, Entity{})
```

### PlantUML
状态机提供了接口，可以直接生成PlantUML
```go
v := machine.GeneratePlantUML()
fmt.Println(v)
```

# Demo

[一个复杂的订单状态例子](./example/order.go)
