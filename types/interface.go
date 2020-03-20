package types

//type StewpotInterface interface {
//	Nodes() map[string]Node
//	GenerateMsg(difficulty int64, msgSize int64, content string) Message
//	SendMsg(source Node, msg Message) int64
//	GetTimeUnitTasks(t int64) []Task
//	MarshalNodes() []byte
//	RestartNetwork(nodeNum, maxIn, maxOut, maxBest int, bandwidth int64, callback func(task Task))
//}

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
	String() string
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
