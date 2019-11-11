package main

import (
	"fmt"
	"sync"
)

var (
	DefaultDelay int64 = 100
)

type Node struct {
	SpeedRate int

	IP        string
	Bandwidth Bandwidth
	Loc       Location

	peers      map[string]*Peer
	peerNumIn  int
	peerNumOut int
	maxIn      int
	maxOut     int

	receiver chan NodeMsg
	close    chan bool

	mu sync.RWMutex
}

type Bandwidth struct {
	Upload   int
	Download int
}

type NodeConfig struct {
	IP       string
	Upload   int
	Download int
	MaxIn    int
	MaxOut   int
}

func NewNode(rate int, config NodeConfig, loc Location) *Node {
	n := &Node{}
	n.SpeedRate = rate

	n.IP = config.IP
	n.Bandwidth = Bandwidth{Upload: config.Upload, Download: config.Download}
	n.Loc = loc
	n.peers = make(map[string]*Peer)
	n.maxIn = config.MaxIn
	n.maxOut = config.MaxOut
	//n.Protocol = protocol

	receiver := make(chan NodeMsg)
	n.receiver = receiver

	closeC := make(chan bool)
	n.close = closeC

	return n
}

//func (n *Node) Start() {
//	go n.loop()
//}

func (n *Node) Close() {
	close(n.close)
}

func (n *Node) ConnOut(remoteNode *Node, timeout int) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.peers[remoteNode.IP] != nil {
		return fmt.Errorf("peer [%s] already exists", remoteNode.IP)
	}
	if n.peerNumOut == n.maxOut {
		return fmt.Errorf("max connect out peers")
	}

	delay := n.GetDelay(remoteNode)
	bw := n.minBandwidth(remoteNode)
	pkgLoss := 0

	p := NewPeer(n.SpeedRate, n.IP, remoteNode.IP, true, timeout, delay, bw, pkgLoss)
	//go p.ReceiveMsg(n.receiver)

	n.peers[remoteNode.IP] = p
	n.peerNumOut++

	return nil
}

func (n *Node) ConnIn(remoteNode *Node, timeout int) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.peers[remoteNode.IP] != nil {
		return fmt.Errorf("peer [%s] already exists", remoteNode.IP)
	}
	if n.peerNumIn == n.maxIn {
		return fmt.Errorf("max connect in peers")
	}

	delay := n.GetDelay(remoteNode)
	bw := n.minBandwidth(remoteNode)
	pkgLoss := 0

	p := NewPeer(n.SpeedRate, n.IP, remoteNode.IP, false, timeout, delay, bw, pkgLoss)
	//go p.ReceiveMsg(n.receiver)

	n.peers[remoteNode.IP] = p
	n.peerNumIn++

	return nil
}

//func (n *Node) BroadcastMsg(msg PureMsg) {
//	n.mu.RLock()
//	defer n.mu.RUnlock()
//
//	n.broadcastMsg(msg)
//}
//
//func (n *Node) broadcastMsg(msg PureMsg) {
//	for _, peer := range n.peers {
//		go peer.SendMsg(msg)
//	}
//}

func (n *Node) GetDelay(remote *Node) int64 {
	if delay, ok := n.Loc.Delays[remote.Loc.Name]; ok {
		return delay
	}
	return DefaultDelay
}

func (n *Node) minBandwidth(rn *Node) int {
	min := n.Bandwidth.Upload
	if n.Bandwidth.Download < min {
		min = n.Bandwidth.Download
	}
	if rn.Bandwidth.Upload < min {
		min = rn.Bandwidth.Upload
	}
	if rn.Bandwidth.Download < min {
		min = rn.Bandwidth.Download
	}
	return min
}

func (n *Node) LockConn(remote string, endTime int64) {
	peer := n.peers[remote]
	if peer == nil {
		return
	}
	peer.LockConn(endTime)
}

func (n *Node) ReleaseConn(remote string) {
	peer := n.peers[remote]
	if peer == nil {
		return
	}
	peer.ReleaseConn()
}

//func (n *Node) loop() {
//	for {
//		select {
//		case <-n.receiver:
//		//case nodeMsg := <-n.receiver:
//		// TODO
//		// deal with node
//
//		//msg := nodeMsg.Data
//		//fmt.Println(fmt.Sprintf("node 【 %s 】 receive msg id: %d", n.IP, msg.ID))
//
//		case <-n.close:
//			for _, p := range n.peers {
//				p.Close()
//			}
//			return
//		}
//	}
//}

type NodeMsg struct {
	IP        string
	Timestamp int64
	Data      PureMsg
}

//func (m *NodeMsg) Decode() *PureMsg {
//	id := BytesToInt64(m.Data[:8])
//	data := m.Data[8:]
//
//	return &PureMsg{ID: id, Data: data}
//}

type PureMsg struct {
	ID   int
	Data int
}
