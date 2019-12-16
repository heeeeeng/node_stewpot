package tasks

import (
	"encoding/json"
	"github.com/heeeeeng/node_stewpot/types"
)

type MsgProcessCalcTask struct {
	taskType  TaskType
	startTime int64
	endTime   int64
	node      types.Node
	msg       types.Message
}

func NewMsgProcessCalcTask(startTime int64, n types.Node, msg types.Message) *MsgProcessCalcTask {
	timeUsage := int64(msg.Difficulty / n.Perf())
	endTime := startTime + timeUsage

	return &MsgProcessCalcTask{
		taskType:  TaskTypeMsgProcessCalc,
		startTime: startTime,
		endTime:   endTime,
		node:      n,
		msg:       msg,
	}
}

func (t *MsgProcessCalcTask) Process(tl types.Timeline) {
	if t.endTime > tl.CurrentTime() {
		tl.ImportTask(tl.NextTime(), t)
		return
	}

	t.node.ReleaseCpu()
	t.node.StoreMsg(t.msg)

	broadcastTask := NewBroadcastTask(tl.NextTime(), t.node, t.msg)
	tl.ImportTask(broadcastTask.StartTime(), broadcastTask)
}

func (t *MsgProcessCalcTask) Type() int        { return int(t.taskType) }
func (t *MsgProcessCalcTask) StartTime() int64 { return t.startTime }
func (t *MsgProcessCalcTask) EndTime() int64   { return t.endTime }

type MsgProcessCalcTaskJson struct {
	Type      int    `json:"type"`
	TypeDef   string `json:"type_def"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	Node      string `json:"node"`
	Msg       int64  `json:"msg"`
}

func (t *MsgProcessCalcTask) MarshalJSON() ([]byte, error) {
	j := MsgProcessCalcTaskJson{}

	j.Type = int(t.taskType)
	j.TypeDef = t.taskType.String()
	j.StartTime = t.startTime
	j.EndTime = t.endTime
	j.Node = t.node.IP()
	j.Msg = t.msg.ID

	return json.Marshal(j)
}
