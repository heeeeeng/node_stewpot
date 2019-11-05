package main

import (
	"fmt"
	"time"
)

type Peer struct {
	ip  string
	out bool

	conn *Conn

	close chan bool
}

func NewPeer(ip string, out bool, timeout int, recvC, sendC chan []byte) *Peer {
	p := &Peer{}
	p.ip = ip
	p.out = out
	p.close = make(chan bool)

	conn := NewConn(ip, timeout, recvC, sendC, p.close)
	p.conn = conn

	return p
}

func (p *Peer) Close() {
	close(p.close)
}

func (p *Peer) SendMsg(msg []byte) error {
	err := p.conn.SendMsg(msg)
	// TODO dump the result of this sending.

	//callback <- err
	return err
}

func (p *Peer) ReceiveMsg(receiver chan NodeMsg) {
	for {
		data, err := p.conn.RecvMsg()
		if err != nil {
			p.Close()
			return
		}

		msg := NodeMsg{
			IP:        p.ip,
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
	ip string

	timeout   int // in milliseconds
	delay     int // in milliseconds
	bandwidth int // in KiB
	pkgLoss   int // in percentage

	// recvC and sendC simulates a general tcp connection
	recvC chan []byte
	sendC chan []byte

	close chan bool
}

func NewConn(ip string, timeout int, recvC, sendC chan []byte, closeC chan bool) *Conn {
	c := &Conn{}

	c.ip = ip
	c.timeout = timeout
	c.recvC = recvC
	c.sendC = sendC
	c.close = closeC

	return c
}

func (c *Conn) SendMsg(msg []byte) error {
	timer := time.NewTimer(time.Millisecond * time.Duration(c.timeout))
	if c.delay > 0 {
		time.Sleep(time.Millisecond * time.Duration(c.delay))
	}

	select {
	case <-c.close:
		return fmt.Errorf("conn closed")

	case c.sendC <- msg:
		return nil

	case <-timer.C:
		return fmt.Errorf("conn timeout")

	}
}

func (c *Conn) RecvMsg() ([]byte, error) {
	select {
	case <-c.close:
		return nil, fmt.Errorf("conn closed")

	case msg := <-c.recvC:
		return msg, nil

	}
}
