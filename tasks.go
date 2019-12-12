package main

type Task interface {
	StartTime() int64
	EndTime() int64
	Process(tl *Timeline)
}

type MsgTransmitTask struct {
	startTime int64
	endTime   int64
	from      *Node
	to        *Node
	msg       Message
}

func newMsgTransmitTask(startTime int64, from, to *Node, msg Message) *MsgTransmitTask {
	delay := from.GetDelay(to)
	endTime := startTime + delay

	return &MsgTransmitTask{
		startTime: startTime,
		endTime:   endTime,
		from:      from,
		to:        to,
		msg:       msg,
	}
}

func (t *MsgTransmitTask) Process(tl *Timeline) {
	if t.endTime > tl.CurrentTime() {
		tl.ImportTask(tl.NextTime(), t)
		return
	}

	recvTask := newConnRecvTask(tl.CurrentTime(), t.from, t.to, t.msg)
	tl.ImportTask(recvTask.StartTime(), recvTask)
}

func (t *MsgTransmitTask) StartTime() int64 { return t.startTime }
func (t *MsgTransmitTask) EndTime() int64   { return t.endTime }

type BroadcastTask struct {
	startTime int64
	endTime   int64
	node      *Node
	msg       Message
}

func newBroadcastTask(startTime int64, node *Node, msg Message) *BroadcastTask {
	return &BroadcastTask{
		startTime: startTime,
		endTime:   startTime,
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
		tl.ImportTask(task.StartTime(), task)
	}
}

func (t *BroadcastTask) StartTime() int64 { return t.startTime }
func (t *BroadcastTask) EndTime() int64   { return t.endTime }

type MsgProcessTask struct {
	startTime int64
	endTime   int64
	node      *Node
	msg       Message
}

func newMsgProcessTask(startTime int64, n *Node, msg Message) *MsgProcessTask {
	return &MsgProcessTask{
		startTime: startTime,
		endTime:   startTime,
		node:      n,
		msg:       msg,
	}
}

func (t *MsgProcessTask) Process(tl *Timeline) {
	if t.node.MsgExists(t.msg) {
		return
	}
	cpuTask := newMsgProcessCPUReqTask(tl.CurrentTime(), t.node, t.msg)
	tl.ImportTask(cpuTask.StartTime(), cpuTask)
}

func (t *MsgProcessTask) StartTime() int64 { return t.startTime }
func (t *MsgProcessTask) EndTime() int64   { return t.endTime }

type MsgProcessCPUReqTask struct {
	startTime int64
	endTime   int64
	node      *Node
	msg       Message
}

func newMsgProcessCPUReqTask(startTime int64, n *Node, msg Message) *MsgProcessCPUReqTask {
	return &MsgProcessCPUReqTask{
		startTime: startTime,
		endTime:   startTime,
		node:      n,
		msg:       msg,
	}
}

func (t *MsgProcessCPUReqTask) Process(tl *Timeline) {
	if !t.node.LockCpu() {
		t.endTime = tl.NextTime()
		tl.ImportTask(t.StartTime(), t)
		return
	}

	task := newMsgProcessCalcTask(tl.NextTime(), t.node, t.msg)
	tl.ImportTask(task.StartTime(), task)
}

func (t *MsgProcessCPUReqTask) StartTime() int64 { return t.startTime }
func (t *MsgProcessCPUReqTask) EndTime() int64   { return t.endTime }

type MsgProcessCalcTask struct {
	startTime int64
	endTime   int64
	node      *Node
	msg       Message
}

func newMsgProcessCalcTask(startTime int64, n *Node, msg Message) *MsgProcessCalcTask {
	timeUsage := int64(msg.Difficulty / n.Perf)
	endTime := startTime + timeUsage

	return &MsgProcessCalcTask{
		startTime: startTime,
		endTime:   endTime,
		node:      n,
		msg:       msg,
	}
}

func (t *MsgProcessCalcTask) Process(tl *Timeline) {
	if t.endTime > tl.CurrentTime() {
		tl.ImportTask(tl.NextTime(), t)
		return
	}

	t.node.ReleaseCpu()
	t.node.StoreMsg(t.msg)

	broadcastTask := newBroadcastTask(tl.NextTime(), t.node, t.msg)
	tl.ImportTask(broadcastTask.StartTime(), broadcastTask)
}

func (t *MsgProcessCalcTask) StartTime() int64 { return t.startTime }
func (t *MsgProcessCalcTask) EndTime() int64   { return t.endTime }

type ConnRecvTask struct {
	startTime int64
	endTime   int64
	sender    *Node
	recver    *Node
	msg       Message
}

func newConnRecvTask(startTime int64, sender, recver *Node, msg Message) *ConnRecvTask {
	return &ConnRecvTask{
		startTime: startTime,
		endTime:   startTime,
		sender:    sender,
		recver:    recver,
		msg:       msg,
	}
}

func (t *ConnRecvTask) Process(tl *Timeline) {
	task := newMsgProcessTask(t.startTime, t.recver, t.msg)
	tl.ImportTask(t.startTime, task)
}

func (t *ConnRecvTask) StartTime() int64 { return t.startTime }
func (t *ConnRecvTask) EndTime() int64   { return t.endTime }
