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

func (s *Stewpot) InitNetwork(nodeNum int) {
	locSet := []types.Location{types.LocCN, types.LocSEA, types.LocJP}
	loc := locSet[rand.Intn(len(locSet))]

	conf := NodeConfig{
		IP:       fmt.Sprintf("%d_%s", 0, loc.Name),
		Upload:   1024,
		Download: 1024,
		MaxIn:    8,
		MaxOut:   4,
	}
	perf := 1024
	s.bootstrap = NewNode(conf, loc, perf)
	s.nodes = append(s.nodes, s.bootstrap)

	for i := 1; i < nodeNum; i++ {
		loc := locSet[rand.Intn(len(locSet))]
		conf := NodeConfig{
			IP:       fmt.Sprintf("%d_%s", i, loc.Name),
			Upload:   1024,
			Download: 1024,
			MaxIn:    8,
			MaxOut:   4,
		}
		perf := 1024
		node := NewNode(conf, loc, perf)

		//node.TryConnect(s.nodes[rand.Intn(len(s.nodes))])
		node.TryConnect(s.bootstrap)
		s.nodes = append(s.nodes, node)
	}
}

// Stew start the simulating of the nodes network.
func (s *Stewpot) Start() {
	s.timeline.Start()
}

func (s *Stewpot) PrintOutNodes() {
	for _, n := range s.nodes {
		fmt.Println(n.String())
	}
}

func (s *Stewpot) GenerateMsg() types.Message {
	s.msgLocker.Lock()
	defer s.msgLocker.Unlock()

	msg := types.NewMessage(nil, 1, s.msg)
	s.msg++

	return msg
}

func (s *Stewpot) SendNewMsg() {
	node := s.nodes[rand.Intn(len(s.nodes))]
	msg := s.GenerateMsg()
	timestamp := s.timeline.SendNewMsg(node, msg)
	fmt.Println("send msg at time: ", timestamp)
}
