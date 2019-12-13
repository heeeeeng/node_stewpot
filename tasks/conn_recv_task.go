package tasks

import (
	"encoding/json"
	"github.com/heeeeeng/node_stewpot/types"
)

type ConnRecvTask struct {
	startTime int64
	endTime   int64
	sender    types.Node
	recver    types.Node
	msg       types.Message
}

func NewConnRecvTask(startTime int64, sender, recver types.Node, msg types.Message) *ConnRecvTask {
	return &ConnRecvTask{
		startTime: startTime,
		endTime:   startTime,
		sender:    sender,
		recver:    recver,
		msg:       msg,
	}
}

func (t *ConnRecvTask) Process(tl types.Timeline) {
	task := NewMsgProcessTask(t.startTime, t.recver, t.msg)
	tl.ImportTask(t.startTime, task)
}

func (t *ConnRecvTask) StartTime() int64 { return t.startTime }
func (t *ConnRecvTask) EndTime() int64   { return t.endTime }

type ConnRecvTaskJson struct {
	Type      int    `json:"type"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	Sender    string `json:"sender"`
	Recver    string `json:"recver"`
	Msg       int64  `json:"msg"`
}

func (t *ConnRecvTask) MarshalJSON() ([]byte, error) {
	j := ConnRecvTaskJson{}
	j.Type = TaskTypeConnRecv
	j.StartTime = t.startTime
	j.EndTime = t.endTime
	j.Sender = t.sender.IP()
	j.Recver = t.recver.IP()
	j.Msg = t.msg.ID

	return json.Marshal(j)
}
