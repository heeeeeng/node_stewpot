package main

type Task interface {
	Process(tl *Timeline)
}

type CPUCalculateTask struct {
	startTime int64
	node      *Node
}

func (t *CPUCalculateTask) Process(tl *Timeline) {

}

type MsgTransmitTask struct {
	startTime int64
	from      *Node
	to        *Node
}

func newMsgTransmitTask(t int64, from, to *Node) *MsgTransmitTask {
	return &MsgTransmitTask{
		startTime: t,
		from:      from,
		to:        to,
	}
}

func (t *MsgTransmitTask) Process(tl *Timeline) {
	delay := t.from.GetDelay(t.to)
	endtime := t.startTime + delay

	lockTask := newConnLockTask(t.startTime, endtime, t.from, t.to)
	tl.ImportTask(t.startTime, lockTask)

	rlsTask := newConnReleaseTask(endtime, t.from, t.to)
	tl.ImportTask(endtime, rlsTask)
}

type BroadcastTask struct {
	startTime int64
	source    *Node
	node      *Node
	msg       int64
}

func (t *BroadcastTask) Process(tl *Timeline) {
	if t.node == nil {
		return
	}

	for _, p := range t.node.peers {
		if t.source != nil && p.ipRemote == t.source.IP {
			continue
		}
		task := newMsgTransmitTask(t.startTime, t.node, p.node)
		tl.ImportTask(t.startTime, task)
	}
}

type MsgCheckTask struct {
	startTime int64
	node      *Node
	msg       int64
}

func (t *MsgCheckTask) Process(tl *Timeline) {

}

type ConnLockTask struct {
	startTime int64
	endTime   int64
	n1        *Node
	n2        *Node
}

func newConnLockTask(t int64, endTime int64, n1, n2 *Node) *ConnLockTask {
	return &ConnLockTask{
		startTime: t,
		endTime:   endTime,
		n1:        n1,
		n2:        n2,
	}
}

func (t *ConnLockTask) Process(tl *Timeline) {
	t.n1.LockConn(t.n2.IP, t.endTime)
	t.n2.LockConn(t.n1.IP, t.endTime)
}

type ConnReleaseTask struct {
	startTime int64
	n1        *Node
	n2        *Node
}

func newConnReleaseTask(t int64, n1, n2 *Node) *ConnReleaseTask {
	return &ConnReleaseTask{
		startTime: t,
		n1:        n1,
		n2:        n2,
	}
}

func (t *ConnReleaseTask) Process(tl *Timeline) {
	t.n1.ReleaseConn(t.n2.IP)
	t.n2.ReleaseConn(t.n1.IP)
}
