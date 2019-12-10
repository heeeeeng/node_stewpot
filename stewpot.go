package main

import (
	"fmt"
	"math/rand"
)

type Stewpot struct {
	bootstrap *Node
	nodes     []*Node
}

func NewStewpot() *Stewpot {
	s := &Stewpot{}

	return s
}

func (s *Stewpot) InitNetwork() {
	nodeNum := 100

	loc := ConstLocations[rand.Intn(len(ConstLocations))]
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
		loc := ConstLocations[rand.Intn(len(ConstLocations))]
		conf := NodeConfig{
			IP:       fmt.Sprintf("%d_%s", i, loc.Name),
			Upload:   1024,
			Download: 1024,
			MaxIn:    8,
			MaxOut:   4,
		}
		perf := 1024
		node := NewNode(conf, loc, perf)
		s.nodes = append(s.nodes, node)

		node.TryConnect(s.bootstrap)
	}
}

// Stew start the simulating of the nodes network.
func (s *Stewpot) Start() {
	//msgSize := 100
	//sendNum := 100
	//
	//for i := 0; i < sendNum; i++ {
	//	index := rand.Intn(len(s.nodes))
	//	n := s.nodes[index]
	//
	//	msg := PureMsg{
	//		ID:   i,
	//		Data: msgSize,
	//	}
	//	n.BroadcastMsg(msg)
	//}

}

func (s *Stewpot) PrintOutNodes() {
	for _, n := range s.nodes {
		fmt.Println(n.String())
	}
}
