package tasks

import (
	"encoding/json"
	"github.com/heeeeeng/node_stewpot/types"
)

type MsgProcessTask struct {
	startTime int64
	endTime   int64
	node      types.Node
	msg       types.Message
}

func NewMsgProcessTask(startTime int64, n types.Node, msg types.Message) *MsgProcessTask {
	return &MsgProcessTask{
		startTime: startTime,
		endTime:   startTime,
		node:      n,
		msg:       msg,
	}
}

func (t *MsgProcessTask) Process(tl types.Timeline) {
	if t.node.MsgExists(t.msg) {
		return
	}
	cpuTask := NewMsgProcessCPUReqTask(tl.CurrentTime(), t.node, t.msg)
	tl.ImportTask(cpuTask.StartTime(), cpuTask)
}

func (t *MsgProcessTask) StartTime() int64 { return t.startTime }
func (t *MsgProcessTask) EndTime() int64   { return t.endTime }

type MsgProcessTaskJson struct {
	Type      int    `json:"type"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	Node      string `json:"node"`
	Msg       int64  `json:"msg"`
}

func (t *MsgProcessTask) MarshalJSON() ([]byte, error) {
	j := MsgProcessTaskJson{}

	j.Type = TaskTypeMsgProcess
	j.StartTime = t.startTime
	j.EndTime = t.endTime
	j.Node = t.node.IP()
	j.Msg = t.msg.ID

	return json.Marshal(j)
}
