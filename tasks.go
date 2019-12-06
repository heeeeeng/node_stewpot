package main

type Task interface {
	Process(tl *Timeline)
}

type MsgTransmitTask struct {
	startTime int64
	from      *Node
	to        *Node
	msg       Message
}

func newMsgTransmitTask(t int64, from, to *Node, msg Message) *MsgTransmitTask {
	return &MsgTransmitTask{
		startTime: t,
		from:      from,
		to:        to,
		msg:       msg,
	}
}

func (t *MsgTransmitTask) Process(tl *Timeline) {
	delay := t.from.GetDelay(t.to)
	endtime := t.startTime + delay

	recvTask := newConnRecvTask(endtime, t.from, t.to, t.msg)
	tl.ImportTask(endtime, recvTask)
}

type BroadcastTask struct {
	startTime int64
	node      *Node
	msg       Message
}

func newBroadcastTask(startTime int64, node *Node, msg Message) *BroadcastTask {
	return &BroadcastTask{
		startTime: startTime,
		node:      node,
		msg:       msg,
	}
}

func (t *BroadcastTask) Process(tl *Timeline) {
	if t.node == nil {
		return
	}

	for _, p := range t.node.peers {
		if t.msg.Source != nil && p.ipRemote == t.msg.Source.IP {
			continue
		}
		msg := t.msg
		msg.Source = t.node
		task := newMsgTransmitTask(t.startTime, t.node, p.node, msg)
		tl.ImportTask(t.startTime, task)
	}
}

type MsgProcessTask struct {
	startTime int64
	initTime  int64
	node      *Node
	msg       Message
	retryTime int64
}

func newMsgProcessTask(startTime, initTime int64, n *Node, msg Message) *MsgProcessTask {
	return &MsgProcessTask{
		startTime: startTime,
		initTime:  initTime,
		node:      n,
		msg:       msg,
		retryTime: 10, // TODO cpu retry every 10ms if it is locked
	}
}

func (t *MsgProcessTask) Process(tl *Timeline) {
	if t.node.MsgExists(t.msg) {
		return
	}
	timeUsage := t.msg.Difficulty / t.node.Perf
	if timeUsage == 0 {
		return
	}
	if !t.node.LockCpu() {
		t.startTime += t.retryTime
		tl.ImportTask(t.startTime, t)
		return
	}

	// broadcast
	newBroadcastTask(t.startTime+int64(timeUsage), t.node, t.msg)
}

type ConnRecvTask struct {
	startTime int64
	sender    *Node
	recver    *Node
	msg       Message
}

func newConnRecvTask(startTime int64, sender, recver *Node, msg Message) *ConnRecvTask {
	return &ConnRecvTask{
		startTime: startTime,
		sender:    sender,
		recver:    recver,
		msg:       msg,
	}
}

func (t *ConnRecvTask) Process(tl *Timeline) {
	task := newMsgProcessTask(t.startTime, t.startTime, t.recver, t.msg)
	tl.ImportTask(t.startTime, task)
}
