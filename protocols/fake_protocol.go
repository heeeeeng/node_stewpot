package protocols

import "github.com/heeeeeng/node_stewpot/types"

type FakeProtocol struct {
	producer types.MsgProducer
}

func NewFakeProtocol() *FakeProtocol {
	return &FakeProtocol{}
}

func (p *FakeProtocol) RegisterProducer(producer types.MsgProducer) {
	p.producer = producer
}

func (p *FakeProtocol) ConsumeMsg(source types.Node, content string) {
	return
}
