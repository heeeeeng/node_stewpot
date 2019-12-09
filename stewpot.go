package main

import "math/rand"

type Stewpot struct {
	protocol Protocol

	nodes []*Node
}

func NewStewpot() *Stewpot {
	s := &Stewpot{}

	return s
}

func (s *Stewpot) InitNetwork() {
	nodeNum := 100

	loc := ConstLocations[rand.Intn(len(ConstLocations))]
	conf := NodeConfig{
		IP:       loc.Name,
		Upload:   1024,
		Download: 1024,
		MaxIn:    4,
		MaxOut:   8,
	}
	perf := 1024
	bootstrap := NewNode(conf, loc, perf)
	s.nodes = append(s.nodes, bootstrap)

	for i := 1; i < nodeNum; i++ {
		loc := ConstLocations[rand.Intn(len(ConstLocations))]
		conf := NodeConfig{
			IP:       loc.Name,
			Upload:   1024,
			Download: 1024,
			MaxIn:    4,
			MaxOut:   8,
		}
		perf := 1024
		node := NewNode(conf, loc, perf)
		s.nodes = append(s.nodes, node)

		node.TryConnect(bootstrap)
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

type Protocol interface {
	GetPackages() [][]byte
}
