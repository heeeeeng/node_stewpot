package main

import (
	"fmt"
	"math/rand"
)

type Stewpot struct {
	bootstrap *Node
	nodes     []*Node
	timeline  *Timeline
}

func NewStewpot() *Stewpot {
	db := newMemDB()
	timeline := newTimeline(db)

	s := &Stewpot{}
	s.timeline = timeline

	return s
}

func (s *Stewpot) InitNetwork(nodeNum int) {
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
