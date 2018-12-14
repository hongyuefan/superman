package mavg

import (
	"github.com/hongyuefan/superman/krang"
	"github.com/hongyuefan/superman/logs"
	"github.com/hongyuefan/superman/protocol"
	"github.com/hongyuefan/superman/strategy"
)

/*
 mavg策略的正常状态
 正常状态下，不会加仓
 1. 盈亏状态，当前仓位的盈亏曲线变化
 2. MA信号
 3. 上一个仓位结束后，新建仓位的要冷静一段时间
*/

type normalState struct {
	spm map[string]*symbolParam // key:symbol
}

func NewNormalState() strategy.FSMState {
	ltcParam := &symbolParam{
		klkind:         protocol.KL5Min,
		stopLoseRate:   -0.13,
		stopProfitRate: 0.17,
		minVol:         0.1,
		maxVol:         0.5,
		stepRate:       0.1,
		marketSt:       true,
		level:          10,
	}
	etcParam := &symbolParam{
		klkind:         protocol.KL5Min,
		stopLoseRate:   -0.18,
		stopProfitRate: 0.20,
		minVol:         1,
		maxVol:         5,
		stepRate:       0.1,
		marketSt:       true,
		level:          10,
	}

	st := &normalState{
		spm: make(map[string]*symbolParam),
	}
	st.spm["ltc_usd"] = ltcParam
	st.spm["etc_usd"] = etcParam
	return st
}

/////////////////////////////////////////////////////

func (t *normalState) Name() string {
	return STATE_NAME_NORMAL
}

func (t *normalState) Init() {
}

func (t *normalState) Enter(ctx krang.Context) {
	logs.Info("进入状态[%s]", t.Name())

	// 重新进入normal state，增加限制
	UpLossTimesLimit()

	// 重新读取头寸信息
	mavg.queryAllPos(ctx)
}

/*
 1. 判断当前是否有仓位
 2. 如果没有仓位，且信号是买，买多
 3. 如果没有仓位，信号是卖，做空
 4. 如果有仓位，normal state不加仓，可以考虑从这里跳转到radical state
 5. 如果有仓位，是否超出止盈止损范围，如果超出，平仓
 6. 如果有仓位，信号是买，空头头寸平仓
 7. 如果有仓位，信号是卖，多头头寸平仓

 平仓或者建仓后，应先更新头寸，因为下单到查询头寸，更新头寸信息，这里有一个时间差
 在这个时间差内应该按照最新的头寸信息操作
*/
func (t *normalState) Decide(ctx krang.Context, tick *krang.Tick, evc *strategy.EventCompose) string {
	// 超过亏损限制，进入保守状态
	if IsOverLossLimit() {
		return STATE_NAME_DEFENSE
	}

	t.handleLongPart(ctx, tick, evc)
	t.handleShortPart(ctx, tick, evc)

	if evc.HasEmergency() {
		logs.Info("紧急情况, 策略暂时关闭")
		return STATE_NAME_SHUTDOWN
	}
	return t.Name()
}

/////////////////////////////////////////////////////

func (t *normalState) getSymbolParam(symbol string) *symbolParam {
	v, ok := t.spm[symbol]
	if !ok {
		return nil
	}
	return v
}

func (t *normalState) handleLongPart(ctx krang.Context, tick *krang.Tick, evc *strategy.EventCompose) {
	sp := t.getSymbolParam(tick.Symbol)
	if sp == nil {
		return
	}

	s, ok := evc.Macd.Signals[sp.klkind]
	if !ok || evc.Pos == nil {
		return
	}

	if evc.Pos.LongAvai <= 0 {
		if s == strategy.SIGNAL_BUY {
			reason := "买入信号"
			ArcherOpenPos(ctx, tick, evc, protocol.ORDERTYPE_OPENLONG, reason, sp)
		}
		return
	}

	// 有多头头寸情况
	if s == strategy.SIGNAL_EMERGENCY {
		reason := "紧急情况"
		ArcherClosePos(ctx, tick, evc, protocol.ORDERTYPE_CLOSELONG, reason, sp)
		return
	}

	if evc.Pos.LongFloatPRate <= sp.stopLoseRate || evc.Pos.LongFloatPRate >= sp.stopProfitRate {
		reason := "超出止盈止损范围"
		ArcherClosePos(ctx, tick, evc, protocol.ORDERTYPE_CLOSELONG, reason, sp)
		return
	}

	if s == strategy.SIGNAL_SELL {
		reason := "卖出信号，平多"
		ArcherClosePos(ctx, tick, evc, protocol.ORDERTYPE_CLOSELONG, reason, sp)
	}
}

func (t *normalState) handleShortPart(ctx krang.Context, tick *krang.Tick, evc *strategy.EventCompose) {
	sp := t.getSymbolParam(tick.Symbol)
	if sp == nil {
		return
	}

	s, ok := evc.Macd.Signals[sp.klkind]
	if !ok || evc.Pos == nil {
		return
	}

	if evc.Pos.ShortAvai <= 0 {
		if s == strategy.SIGNAL_SELL {
			reason := "卖出信号"
			ArcherOpenPos(ctx, tick, evc, protocol.ORDERTYPE_OPENSHORT, reason, sp)
		}
		return
	}

	// 有空头头寸情况
	if s == strategy.SIGNAL_EMERGENCY {
		reason := "紧急情况"
		ArcherClosePos(ctx, tick, evc, protocol.ORDERTYPE_CLOSESHORT, reason, sp)
		return
	}

	if evc.Pos.ShortFloatPRate <= sp.stopLoseRate || evc.Pos.ShortFloatPRate >= sp.stopProfitRate {
		reason := "超出止盈止损范围"
		ArcherClosePos(ctx, tick, evc, protocol.ORDERTYPE_CLOSESHORT, reason, sp)
		return
	}

	if s == strategy.SIGNAL_BUY {
		reason := "买入信号，平空"
		ArcherClosePos(ctx, tick, evc, protocol.ORDERTYPE_CLOSESHORT, reason, sp)
	}
}
