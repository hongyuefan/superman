package bows

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"

	"github.com/hongyuefan/superman/kfc"
	"github.com/hongyuefan/superman/logs"
	"github.com/hongyuefan/superman/protocol"
)

type bowLoop struct {
	m         map[string]chan *ArcherCmd
	exchanges []string
}

func InitKafkaClient(brokers []string) error {

	topics := []string{protocol.TOPIC_OKEX_ARCHER_REQ}

	kfc.InitClient(brokers)

	if err := kfc.TobeProducer(); err != nil {
		logs.Error("InitKafkaClient producer error ", err.Error())
		return err
	}

	if err := kfc.TobeConsumer(topics); err != nil {
		logs.Error("InitKafkaClient consumer error ", err.Error())
		return err
	}

	return nil
}

func InitBows() *bowLoop {
	return &bowLoop{
		m: make(map[string]chan *ArcherCmd),
	}
}

func StartExArcher(exchanges []string, bl *bowLoop) error {

	bl.exchanges = exchanges

	for _, ex := range exchanges {

		q := createArchers(ex)

		if q == nil {
			logs.Error("exchange [%s] is not supported !", ex)
			continue
		}

		if err := q.Init(); err != nil {
			logs.Error("exchange [%s] init fail, error:", ex, err.Error())
			return err
		}

		ch := make(chan *ArcherCmd)

		bl.m[ex] = ch

		go q.Run(ch)

		logs.Info("start exchange [%s] archer ok ...", ex)
	}
	return nil
}

func StartCmdLoop(bl *bowLoop) {

	logs.Info("wait for cmds .....")

	signals := make(chan os.Signal, 1)

	signal.Notify(signals, os.Interrupt)

	for {
		select {
		case msg := <-kfc.ReadMessages():
			if err := handleBrokerCmd(bl, msg); err != nil {
				logs.Error(err)
			}
		case <-signals:
			logs.Info("recv a break signal, exit archer...")
			doExit(bl)
			return
		}
	}
}

func dispatchCmd(exch chan<- *ArcherCmd, cmd *ArcherCmd) error {
	select {
	case exch <- cmd:
	default:
		return fmt.Errorf("Exchange %s Handler Chan Is Full", cmd.Exchange)
	}
	return nil
}

func handleBrokerCmd(bl *bowLoop, msg *sarama.ConsumerMessage) error {

	exch, ok := bl.m[string(msg.Key)]

	if !ok {
		return fmt.Errorf("recv msg with wrong key [%s]", string(msg.Key))
	}

	p := &protocol.FixPackage{}

	if err := p.ParseFromArray(msg.Value); err != nil {
		return fmt.Errorf("recv msg parse error, topic[%s], error: %s", msg.Topic, err.Error())
	}

	cmd := &ArcherCmd{}

	cmd.ReqSerial = p.GetReqSerial()
	cmd.Exchange = string(msg.Key)
	cmd.BusiId = p.GetBusinessId()
	cmd.PayLoad = p.GetPayload()

	return dispatchCmd(exch, cmd)
}

func doExit(bl *bowLoop) {

	cmd := &ArcherCmd{
		BusiId: protocol.CMD_EXIST,
	}

	for _, v := range bl.m {
		v <- cmd
	}

	kfc.ExitProducer()
	kfc.ExitConsumer()
}

func createArchers(ex string) Archer {
	switch ex {
	case "okex":
		return newOkexArcher()
	}
	return nil
}
