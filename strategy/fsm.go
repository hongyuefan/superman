package strategy

import (
	"github.com/hongyuefan/superman/krang"
	"github.com/hongyuefan/superman/logs"
)

const (
	SIGNAL_BUY       = 1
	SIGNAL_SELL      = 2
	SIGNAL_EMERGENCY = 3
)

/*
  策略都是有限状态机
  这里是状态机的接口，具体状态机有哪些状态，由各个策略决定
*/

/*
 状态接口
 状态里有一整套决策参数，不同状态参数不一样，
 状态依据参数和ticker的变化集合来做决策，并决定
 跳到哪个状态
*/
type FSMState interface {
	Name() string
	Init()
	Enter(ctx krang.Context)
	Decide(ctx krang.Context, tick *krang.Tick, e *EventCompose) string
}

/*
 处理器接口
 处理器计算ticker引起的不同变化
 将结果写到组合里
*/
type FSMHandler interface {
	Name() string
	OnTick(ctx krang.Context, tick *krang.Tick, e *EventCompose)
}

/*
 变化结果组合
*/
type EventCompose struct {
	Exchange     string
	Symbol       string
	ContractType string

	Pos *krang.Pos

	Money struct {
		Balance float32
	}

	Macd struct {
		Signals map[int32]int32 // key: k线种类， value: 买入卖出信号
	}
}

func newEventCompose() *EventCompose {
	evc := &EventCompose{}
	evc.Macd.Signals = make(map[int32]int32)
	return evc
}

func (e *EventCompose) reset() {
	e.Exchange = ""
	e.Symbol = ""
	e.ContractType = ""
	e.Pos = nil
	e.Money.Balance = 0

	for k, _ := range e.Macd.Signals {
		e.Macd.Signals[k] = 0
	}
}

func (e *EventCompose) HasEmergency() bool {
	for _, v := range e.Macd.Signals {
		if v == SIGNAL_EMERGENCY {
			return true
		}
	}
	return false
}

/*
 状态机
 状态机只需管理handlers和各个状态的跳转
*/
type FSM struct {
	name     string
	state    FSMState // 当前状态
	handlers []FSMHandler
	evc      *EventCompose
	states   map[string]FSMState // 全部的状态
}

func NewFSM(name string) *FSM {
	return &FSM{
		name:     name,
		state:    nil,
		handlers: make([]FSMHandler, 0),
		evc:      newEventCompose(),
		states:   make(map[string]FSMState),
	}
}

func (t *FSM) GetState() FSMState {
	return t.state
}

func (t *FSM) SetState(stn string) {
	st, ok := t.states[stn]
	if !ok {
		panic("SetState param invalid")
	}
	t.state = st
}

func (t *FSM) AddHandler(h FSMHandler) {
	if h == nil {
		panic("AddHandler param nil")
	}
	for _, v := range t.handlers {
		if v.Name() == h.Name() {
			panic("AddHandler repeat handler")
		}
	}
	t.handlers = append(t.handlers, h)
}

func (t *FSM) AddState(st FSMState) {
	if st == nil {
		panic("AddState param nil")
	}
	_, dup := t.states[st.Name()]
	if dup {
		panic("AddState repeat state")
	}
	t.states[st.Name()] = st
}

func (t *FSM) Call(ctx krang.Context, tick *krang.Tick) {
	t.evc.reset()
	t.evc.Exchange = tick.Exchange
	t.evc.Symbol = tick.Symbol
	t.evc.ContractType = tick.ContractType

	for _, v := range t.handlers {
		v.OnTick(ctx, tick, t.evc)
	}

	oldst := t.GetState()
	newStname := oldst.Decide(ctx, tick, t.evc)
	t.SetState(newStname)

	if oldst.Name() != newStname {
		logs.Info("[%s]fsm 从[%s]状态跳转到[%s]状态", t.name, oldst.Name(), newStname)
		t.GetState().Enter(ctx)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////
