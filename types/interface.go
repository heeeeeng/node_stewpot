package types

type Timeline interface {
	CurrentTime() int64
	NextTime() int64
	ImportTask(startTime int64, task Task)
}

type Task interface {
	StartTime() int64
	EndTime() int64
	Process(tl Timeline)
	MarshalJSON() ([]byte, error)
}

type Node interface {
	IP() string
	Location() Location
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
