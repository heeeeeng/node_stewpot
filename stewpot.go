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
	rate := 1000
	timeout := 1000

	ncCN := NodeConfig{
		IP:       "CN1",
		Upload:   1024,
		Download: 1024,
		MaxIn:    4,
		MaxOut:   8,
	}
	nodeCN := NewNode(rate, ncCN, LocCN)
	recvCN := make(chan PureMsg)

	ncNA := NodeConfig{
		IP:       "NA1",
		Upload:   1024,
		Download: 1024,
		MaxIn:    4,
		MaxOut:   8,
	}
	nodeNA := NewNode(rate, ncNA, LocNA)
	recvNA := make(chan PureMsg)

	nodeCN.ConnOut(nodeNA, timeout, recvCN, recvNA)
	nodeNA.ConnIn(nodeCN, timeout, recvNA, recvCN)

	s.nodes = append(s.nodes, nodeCN, nodeNA)

	for _, n := range s.nodes {
		n.Start()
	}
}

// Stew start the simulating of the nodes network.
func (s *Stewpot) Start() {
	msgSize := 100
	sendNum := 100

	for i := 0; i < sendNum; i++ {
		index := rand.Intn(len(s.nodes))
		n := s.nodes[index]

		msg := PureMsg{
			ID:   i,
			Data: msgSize,
		}
		n.BroadcastMsg(msg)
	}

}

type Protocol interface {
	GetPackages() [][]byte
}
