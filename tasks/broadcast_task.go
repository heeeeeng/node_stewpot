package tasks

import (
	"encoding/json"
	"github.com/heeeeeng/node_stewpot/types"
)

type BroadcastTask struct {
	taskType  TaskType
	startTime int64
	endTime   int64
	node      types.Node
	msg       types.Message
}

func NewBroadcastTask(startTime int64, node types.Node, msg types.Message) *BroadcastTask {
	return &BroadcastTask{
		taskType:  TaskTypeBroadcast,
		startTime: startTime,
		endTime:   startTime,
		node:      node,
		msg:       msg,
	}
}

func (t *BroadcastTask) Process(tl types.Timeline) {
	if t.node == nil {
		return
	}
	n := t.node

	var destinations []types.Node
	for _, p := range t.node.Peers() {
		if t.msg.Source != nil && p.RemoteIP() == t.msg.Source.IP() {
			continue
		}

		destinations = append(destinations, p.GetNode())
	}

	// calculate bandwidth usage time.
	totalMsgSize := t.msg.Size * int64(len(destinations))
	bandwidthTimeUsage := totalMsgSize / int64(n.BandwidthInMillisecond())

	for _, dest := range destinations {
		msg := t.msg
		msg.Source = t.node
		task := NewMsgTransmitUploadTask(t.startTime, bandwidthTimeUsage, t.node, dest, msg)
		tl.ImportTask(task.StartTime(), task)
	}

}

func (t *BroadcastTask) Type() int        { return int(t.taskType) }
func (t *BroadcastTask) StartTime() int64 { return t.startTime }
func (t *BroadcastTask) EndTime() int64   { return t.endTime }

type BroadcastTaskJson struct {
	Type      int    `json:"type"`
	TypeDef   string `json:"type_def"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	Node      string `json:"node"`
	Msg       int64  `json:"msg"`
}

func (t *BroadcastTask) MarshalJSON() ([]byte, error) {
	j := BroadcastTaskJson{}

	j.Type = int(t.taskType)
	j.TypeDef = t.taskType.String()
	j.StartTime = t.startTime
	j.EndTime = t.endTime
	j.Node = t.node.IP()
	j.Msg = t.msg.ID

	return json.Marshal(j)
}
