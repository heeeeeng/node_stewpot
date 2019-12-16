package tasks

import (
	"encoding/json"
	"github.com/heeeeeng/node_stewpot/types"
)

type MsgTransmitTask struct {
	taskType  TaskType
	startTime int64
	endTime   int64
	from      types.Node
	to        types.Node
	msg       types.Message
}

func NewMsgTransmitTask(startTime int64, from, to types.Node, msg types.Message) *MsgTransmitTask {
	delay := from.GetDelay(to.Location())
	endTime := startTime + delay

	return &MsgTransmitTask{
		taskType:  TaskTypeMsgTransmit,
		startTime: startTime,
		endTime:   endTime,
		from:      from,
		to:        to,
		msg:       msg,
	}
}

func (t *MsgTransmitTask) Process(tl types.Timeline) {
	if t.endTime > tl.CurrentTime() {
		tl.ImportTask(tl.NextTime(), t)
		return
	}

	recvTask := NewConnRecvTask(tl.CurrentTime(), t.from, t.to, t.msg)
	tl.ImportTask(recvTask.StartTime(), recvTask)
}

func (t *MsgTransmitTask) Type() int        { return int(t.taskType) }
func (t *MsgTransmitTask) StartTime() int64 { return t.startTime }
func (t *MsgTransmitTask) EndTime() int64   { return t.endTime }

type MsgTransmitTaskJson struct {
	Type      int    `json:"type"`
	TypeDef   string `json:"type_def"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	Source    string `json:"source"`
	Target    string `json:"target"`
	Msg       int64  `json:"msg"`
}

func (t *MsgTransmitTask) MarshalJSON() ([]byte, error) {
	j := MsgTransmitTaskJson{}

	j.Type = int(t.taskType)
	j.TypeDef = t.taskType.String()
	j.StartTime = t.startTime
	j.EndTime = t.endTime
	j.Source = t.from.IP()
	j.Target = t.to.IP()
	j.Msg = t.msg.ID

	return json.Marshal(j)
}
