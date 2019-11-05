package main

import "sync"

// TODO
// Node config

type Node struct {
	IP        string
	Bandwidth bandwidth
	Protocol  Protocol
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

type bandwidth struct {
	upload   int
	download int
}

func NewNode() *Node {
	// TODO
	// config

	n := &Node{}

	receiver := make(chan NodeMsg)
	n.receiver = receiver

	closeC := make(chan bool)
	n.close = closeC

	return n
}

func (n *Node) Close() {
	close(n.close)
}

func (n *Node) ConnOut(ip string, timeout int, recvC, sendC chan []byte) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.peers[ip] != nil {
		return
	}
	if n.peerNumOut == n.maxOut {
		return
	}

	p := NewPeer(ip, true, timeout, recvC, sendC)
	n.peers[ip] = p
	n.peerNumOut++

	go p.ReceiveMsg(n.receiver)
}

func (n *Node) ConnIn(ip string, timeout int, recvC, sendC chan []byte) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.peers[ip] != nil {
		return
	}
	if n.peerNumIn == n.maxIn {
		return
	}

	p := NewPeer(ip, false, timeout, recvC, sendC)
	n.peers[ip] = p
	n.peerNumIn++

	go p.ReceiveMsg(n.receiver)
}

func (n *Node) BroadcastMsg(msg []byte) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	n.broadcastMsg(msg)
}

func (n *Node) broadcastMsg(msg []byte) {
	for _, peer := range n.peers {
		go peer.SendMsg(msg)
	}
}

func (n *Node) loop() {
	for {
		select {
		case nodeMsg := <-n.receiver:
			// TODO
			// deal with node

		case <-n.close:
			for _, p := range n.peers {
				p.Close()
			}
			return
		}
	}
}

type NodeMsg struct {
	IP        string
	Timestamp int64
	Data      []byte
}
