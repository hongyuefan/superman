package mavg

import (
	"github.com/hongyuefan/superman/krang"
	"github.com/hongyuefan/superman/strategy"
)

/*
 mavg策略的激进状态
 激进状态，止盈止损幅度，下单信号选择，加仓策略上
 都会比较激进
*/

type radicalState struct {
}

func NewRadicalState() strategy.FSMState {
	return &radicalState{}
}

func (t *radicalState) Name() string {
	return STATE_NAME_RADICAL
}

func (t *radicalState) Init() {
}

func (t *radicalState) Enter(ctx krang.Context) {
}

func (t *radicalState) Decide(ctx krang.Context, tick *krang.Tick, evc *strategy.EventCompose) string {
	return t.Name()
}
