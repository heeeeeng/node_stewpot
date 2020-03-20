package stewpot

import "github.com/heeeeeng/node_stewpot/types"

type Peer struct {
	ipLocal  string
	ipRemote string
	out      bool

	node types.Node

	close chan bool
}

func NewPeer(ipLocal, ipRemote string, out bool, remoteNode types.Node) *Peer {
	p := &Peer{}

	p.ipLocal = ipLocal
	p.ipRemote = ipRemote
	p.out = out
	p.node = remoteNode
	p.close = make(chan bool)

	//conn := NewConn(ipLocal, ipRemote, timeout, delay, bandwidth, pkgLoss, p.close)
	//p.conn = conn

	return p
}

func (p *Peer) Close() {
	close(p.close)
}

//func (p *Peer) GetDelay() int64 {
//	return p.conn.delay
//}

func (p *Peer) GetNode() types.Node {
	return p.node
}

func (p *Peer) RemoteIP() string {
	return p.ipRemote
}

func (p *Peer) Out() bool {
	return p.out
}

//func (p *Peer) LockConn(endTime int64) {
//	p.conn.lock = true
//	p.conn.rlsTime = endTime
//}
//
//func (p *Peer) ReleaseConn() {
//	p.conn.lock = false
//	p.conn.rlsTime = 0
//}

//// Conn simulates a tcp link between two nodes.
//type Conn struct {
//	ipLocal  string
//	ipRemote string
//
//	timeout   int   // in milliseconds
//	delay     int64 // in milliseconds
//	bandwidth int   // in Bytes
//	pkgLoss   int   // in percentage
//
//	lock    bool
//	rlsTime int64
//
//	close chan bool
//}
//
//func NewConn(ipLocal, ipRemote string, timeout int, delay int64, bandwidth, pkgLoss int, closeC chan bool) *Conn {
//	c := &Conn{}
//
//	c.ipLocal = ipLocal
//	c.ipRemote = ipRemote
//	c.timeout = timeout
//	c.delay = delay
//	c.bandwidth = bandwidth
//	c.pkgLoss = pkgLoss
//	c.lock = false
//
//	c.close = closeC
//
//	return c
//}
