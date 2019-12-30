package tasks

import (
	"encoding/json"
	"github.com/heeeeeng/node_stewpot/types"
)

type MsgTransmitDelayTask struct {
	taskType  TaskType
	startTime int64
	endTime   int64
	from      types.Node
	to        types.Node
	msg       types.Message
}

func NewMsgTransmitDelayTask(startTime int64, from, to types.Node, msg types.Message) *MsgTransmitDelayTask {
	delay := from.GetDelay(to.Location())
	endTime := startTime + delay

	return &MsgTransmitDelayTask{
		taskType:  TaskTypeMsgTransmitDelay,
		startTime: startTime,
		endTime:   endTime,
		from:      from,
		to:        to,
		msg:       msg,
	}
}

func (t *MsgTransmitDelayTask) Process(tl types.Timeline) {
	if t.endTime > tl.CurrentTime() {
		tl.ImportTask(tl.NextTime(), t)
		return
	}

	recvTask := NewConnRecvTask(tl.CurrentTime(), t.from, t.to, t.msg)
	tl.ImportTask(recvTask.StartTime(), recvTask)
}

func (t *MsgTransmitDelayTask) Type() int        { return int(t.taskType) }
func (t *MsgTransmitDelayTask) StartTime() int64 { return t.startTime }
func (t *MsgTransmitDelayTask) EndTime() int64   { return t.endTime }

type MsgTransmitDelayTaskJson struct {
	Type      int    `json:"type"`
	TypeDef   string `json:"type_def"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	Source    string `json:"source"`
	Target    string `json:"target"`
	Msg       int64  `json:"msg"`
}

func (t *MsgTransmitDelayTask) MarshalJSON() ([]byte, error) {
	j := MsgTransmitDelayTaskJson{}

	j.Type = int(t.taskType)
	j.TypeDef = t.taskType.String()
	j.StartTime = t.startTime
	j.EndTime = t.endTime
	j.Source = t.from.IP()
	j.Target = t.to.IP()
	j.Msg = t.msg.ID

	return json.Marshal(j)
}
