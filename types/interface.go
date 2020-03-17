package types

type Timeline interface {
	CurrentTime() int64
	NextTime() int64
	Protocol() Protocol
	ImportTask(startTime int64, task Task)
}

type Task interface {
	Type() int
	StartTime() int64
	EndTime() int64
	Process(tl Timeline)
	MarshalJSON() ([]byte, error)
}

type Node interface {
	IP() string
	Location() Location
	Bandwidth() int64
	BandwidthInMillisecond() int64
	BestPeersLimit() int
	Perf() int
	Peers() map[string]Peer
	GetDelay(location Location) int64
	MsgExists(msg Message) bool
	StoreMsg(msg Message)
	LockCpu() bool
	ReleaseCpu()
	ConnectIn(remoteNode Node) (bool, []Node)
}

type Peer interface {
	RemoteIP() string
	Out() bool
	GetNode() Node
}

// Protocol should implements ConsumeMsg function to process the
// content of each message.
//
// A MsgProducer should be contained as an element of Protocol object.
// This is for producing new messages. e.g.
//
// type Protocol struct {
// 		producer MsgProducer
// }
type Protocol interface {
	RegisterProducer(producer MsgProducer)
	ConsumeMsg(source Node, content string)
}

type MsgProducer interface {
	ProduceMsg(source Node, difficulty int64, size int64, content string)
}
