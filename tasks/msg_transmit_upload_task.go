package tasks

import (
	"encoding/json"
	"github.com/heeeeeng/node_stewpot/types"
)

type MsgTransmitUploadTask struct {
	taskType  TaskType
	startTime int64
	endTime   int64
	from      types.Node
	to        types.Node
	msg       types.Message
}

func NewMsgTransmitUploadTask(startTime int64, uploadTime int64, from, to types.Node, msg types.Message) *MsgTransmitUploadTask {
	return &MsgTransmitUploadTask{
		taskType:  TaskTypeMsgTransmitUpload,
		startTime: startTime,
		endTime:   startTime + uploadTime,
		from:      from,
		to:        to,
		msg:       msg,
	}
}

func (t *MsgTransmitUploadTask) Process(tl types.Timeline) {
	if t.endTime > tl.CurrentTime() {
		tl.ImportTask(tl.NextTime(), t)
		return
	}

	delayTask := NewMsgTransmitDelayTask(tl.CurrentTime(), t.from, t.to, t.msg)
	tl.ImportTask(delayTask.StartTime(), delayTask)
}

func (t *MsgTransmitUploadTask) Type() int        { return int(t.taskType) }
func (t *MsgTransmitUploadTask) StartTime() int64 { return t.startTime }
func (t *MsgTransmitUploadTask) EndTime() int64   { return t.endTime }

type MsgTransmitUploadTaskJson struct {
	Type      int    `json:"type"`
	TypeDef   string `json:"type_def"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	Source    string `json:"source"`
	Target    string `json:"target"`
	Msg       int64  `json:"msg"`
}

func (t *MsgTransmitUploadTask) MarshalJSON() ([]byte, error) {
	j := MsgTransmitUploadTaskJson{}

	j.Type = int(t.taskType)
	j.TypeDef = t.taskType.String()
	j.StartTime = t.startTime
	j.EndTime = t.endTime
	j.Source = t.from.IP()
	j.Target = t.to.IP()
	j.Msg = t.msg.ID

	return json.Marshal(j)
}
