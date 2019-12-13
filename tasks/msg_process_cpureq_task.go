package tasks

import (
	"encoding/json"
	"github.com/heeeeeng/node_stewpot/types"
)

type MsgProcessCPUReqTask struct {
	startTime int64
	endTime   int64
	node      types.Node
	msg       types.Message
}

func NewMsgProcessCPUReqTask(startTime int64, n types.Node, msg types.Message) *MsgProcessCPUReqTask {
	return &MsgProcessCPUReqTask{
		startTime: startTime,
		endTime:   startTime,
		node:      n,
		msg:       msg,
	}
}

func (t *MsgProcessCPUReqTask) Process(tl types.Timeline) {
	if !t.node.LockCpu() {
		t.endTime = tl.NextTime()
		tl.ImportTask(t.StartTime(), t)
		return
	}

	task := NewMsgProcessCalcTask(tl.NextTime(), t.node, t.msg)
	tl.ImportTask(task.StartTime(), task)
}

func (t *MsgProcessCPUReqTask) StartTime() int64 { return t.startTime }
func (t *MsgProcessCPUReqTask) EndTime() int64   { return t.endTime }

type MsgProcessCPUReqTaskJson struct {
	Type      int    `json:"type"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	Node      string `json:"node"`
	Msg       int64  `json:"msg"`
}

func (t *MsgProcessCPUReqTask) MarshalJSON() ([]byte, error) {
	j := MsgProcessCPUReqTaskJson{}

	j.Type = TaskTypeMsgProcessCPUReq
	j.StartTime = t.startTime
	j.EndTime = t.endTime
	j.Node = t.node.IP()
	j.Msg = t.msg.ID

	return json.Marshal(j)
}
