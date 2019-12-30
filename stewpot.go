package main

import (
	"fmt"
	"github.com/heeeeeng/node_stewpot/types"
	"math/rand"
	"sync"
)

type Stewpot struct {
	msg       int64
	msgLocker sync.RWMutex

	bootstrap *Node
	nodes     []*Node
	timeline  *Timeline
}

func NewStewpot() *Stewpot {
	db := newMemDB()
	timeline := newTimeline(db)

	s := &Stewpot{}
	s.msg = 1
	s.timeline = timeline

	return s
}

func (s *Stewpot) RestartNetwork(nodeNum, maxIn, maxOut int, bandwidth types.FileSize) {
	s.Stop()

	db := newMemDB()
	timeline := newTimeline(db)

	s.msg = 1
	s.timeline = timeline
	s.nodes = make([]*Node, 0)
	s.bootstrap = nil

	s.InitNetwork(nodeNum, maxIn, maxOut, bandwidth)
	s.Start()
}

func (s *Stewpot) InitNetwork(nodeNum, maxIn, maxOut int, bandwidth types.FileSize) {
	locSet := []types.Location{types.LocCN, types.LocSEA, types.LocJP, types.LocRU, types.LocNA, types.LocEU}
	loc := locSet[rand.Intn(len(locSet))]

	conf := NodeConfig{
		IP:       fmt.Sprintf("%d_%s", 0, loc.Name),
		Upload:   bandwidth,
		Download: bandwidth,
		MaxIn:    maxIn,
		MaxOut:   maxOut,
	}
	perf := 1024
	s.bootstrap = NewNode(conf, loc, perf)
	s.nodes = append(s.nodes, s.bootstrap)

	for i := 1; i < nodeNum; i++ {
		loc = locSet[rand.Intn(len(locSet))]
		conf.IP = fmt.Sprintf("%d_%s", i, loc.Name)
		perf := 1024
		node := NewNode(conf, loc, perf)

		node.TryConnect(s.bootstrap)
		s.nodes = append(s.nodes, node)
	}
}

// Stew start the simulating of the nodes network.
func (s *Stewpot) Start() {
	s.timeline.Start()
}

func (s *Stewpot) Stop() {
	s.timeline.Stop()
}

func (s *Stewpot) PrintOutNodes() {
	for _, n := range s.nodes {
		fmt.Println(n.String())
	}
}

func (s *Stewpot) GenerateMsg() types.Message {
	s.msgLocker.Lock()
	defer s.msgLocker.Unlock()

	// TODO msg size should not be hard coded here
	msg := types.NewMessage(nil, 1, s.msg, 256*types.Byte)
	s.msg++

	return msg
}

func (s *Stewpot) SendNewMsg() {
	node := s.nodes[rand.Intn(len(s.nodes))]
	msg := s.GenerateMsg()
	timestamp := s.timeline.SendNewMsg(node, msg)
	fmt.Println("send msg at time: ", timestamp)
}
