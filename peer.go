package main

import (
	"fmt"
	"sync"
	"time"
)

type Peer struct {
	speedRate int

	ipLocal  string
	ipRemote string
	out      bool

	conn *Conn

	close chan bool
}

func NewPeer(rate int, ipLocal, ipRemote string, out bool, timeout int, delay, bandwidth, pkgLoss int, recvC, sendC chan PureMsg) *Peer {
	p := &Peer{}
	p.speedRate = rate

	p.ipLocal = ipLocal
	p.ipRemote = ipRemote
	p.out = out
	p.close = make(chan bool)

	conn := NewConn(rate, ipLocal, ipRemote, timeout, delay, bandwidth, pkgLoss, recvC, sendC, p.close)
	p.conn = conn

	return p
}

func (p *Peer) Close() {
	close(p.close)
}

func (p *Peer) SendMsg(msg PureMsg) error {
	dumper.Log("{MSG-SEND} [%s] send msg %d to [%s]", p.ipLocal, msg.ID, p.ipRemote)

	err := p.conn.SendMsg(msg)
	if err != nil {
		dumper.Log("{MSG-SEND-FAIL} [%s] send msg %d to [%s], err: %v", p.ipLocal, msg.ID, p.ipRemote, err)
	} else {
		dumper.Log("{MSG-SEND-SUCC} [%s] send msg %d to [%s]", p.ipLocal, msg.ID, p.ipRemote)
	}
	return err
}

func (p *Peer) ReceiveMsg(receiver chan NodeMsg) {
	for {
		data, err := p.conn.RecvMsg()
		if err != nil {
			p.Close()
			return
		}

		dumper.Log("{MSG-RECV} [%s] receive msg %d from [%s]", p.ipLocal, data.ID, p.ipRemote)

		msg := NodeMsg{
			IP:        p.ipRemote,
			Timestamp: time.Now().Unix(),
			Data:      data,
		}
		select {
		case receiver <- msg:
			continue
		case <-p.close:
			return
		}
	}
}

// Conn simulates a tcp link between two nodes.
type Conn struct {
	speedRate int

	ipLocal  string
	ipRemote string

	timeout   int // in milliseconds
	delay     int // in milliseconds
	bandwidth int // in Bytes
	pkgLoss   int // in percentage

	// recvC and sendC simulates a general tcp connection
	recvC chan PureMsg
	sendC chan PureMsg

	close chan bool

	mu sync.RWMutex
}

func NewConn(rate int, ipLocal, ipRemote string, timeout int, delay, bandwidth, pkgLoss int,
	recvC, sendC chan PureMsg, closeC chan bool) *Conn {
	c := &Conn{}
	c.speedRate = rate

	c.ipLocal = ipLocal
	c.ipRemote = ipRemote
	c.timeout = timeout
	c.delay = delay
	c.bandwidth = bandwidth
	c.pkgLoss = pkgLoss
	c.recvC = recvC
	c.sendC = sendC
	c.close = closeC

	return c
}

func (c *Conn) SendMsg(msg PureMsg) error {
	//timer := time.NewTimer(time.Millisecond * time.Duration(c.timeout) / time.Duration(c.speedRate))

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.delay > 0 {
		time.Sleep(time.Millisecond * time.Duration(c.delay) / time.Duration(c.speedRate))
	}

	for msg.Data > 0 {
		select {
		case <-c.close:
			return fmt.Errorf("conn closed")

		//case <-timer.C:
		//	return fmt.Errorf("conn timeout")

		case c.sendC <- msg:
			msg.Data -= c.bandwidth
			continue
		}
	}

	return nil
}

func (c *Conn) RecvMsg() (PureMsg, error) {
	select {
	case <-c.close:
		return PureMsg{}, fmt.Errorf("conn closed")

	case msg := <-c.recvC:
		return msg, nil

	}
}
