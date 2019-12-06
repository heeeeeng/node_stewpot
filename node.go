package main

var (
	DefaultDelay int64 = 100
)

type Node struct {
	db *CacheDB

	IP        string
	Bandwidth Bandwidth
	Loc       Location
	Perf      int

	peers      map[string]*Peer
	peerNumIn  int
	peerNumOut int
	maxIn      int
	maxOut     int

	CpuLocked bool

	close chan bool
}

type Bandwidth struct {
	Upload   int
	Download int
}

type NodeConfig struct {
	IP       string
	Upload   int
	Download int
	MaxIn    int
	MaxOut   int
}

func NewNode(config NodeConfig, loc Location, perf int) *Node {
	n := &Node{}

	n.db = newCacheDB()

	n.IP = config.IP
	n.Bandwidth = Bandwidth{Upload: config.Upload, Download: config.Download}
	n.Loc = loc
	n.Perf = perf
	n.peers = make(map[string]*Peer)
	n.maxIn = config.MaxIn
	n.maxOut = config.MaxOut
	n.CpuLocked = false

	closeC := make(chan bool)
	n.close = closeC

	return n
}

func (n *Node) Close() {
	close(n.close)
}

func (n *Node) GetDelay(remote *Node) int64 {
	if delay, ok := n.Loc.Delays[remote.Loc.Name]; ok {
		return delay
	}
	return DefaultDelay
}

func (n *Node) MsgExists(msg Message) bool {
	return n.db.Exist(msg.ID)
}

func (n *Node) LockCpu() bool {
	if n.CpuLocked {
		return false
	}
	n.CpuLocked = true
	return true
}

func (n *Node) ReleaseCpu() {
	n.CpuLocked = false
}

//func (n *Node) minBandwidth(rn *Node) int {
//	min := n.Bandwidth.Upload
//	if n.Bandwidth.Download < min {
//		min = n.Bandwidth.Download
//	}
//	if rn.Bandwidth.Upload < min {
//		min = rn.Bandwidth.Upload
//	}
//	if rn.Bandwidth.Download < min {
//		min = rn.Bandwidth.Download
//	}
//	return min
//}

type CacheDB struct {
	trimNum  int
	hotData  []bool
	coldData int64
}

func newCacheDB() *CacheDB {
	return &CacheDB{
		trimNum:  1000,
		hotData:  make([]bool, 0),
		coldData: 0,
	}
}

func (db *CacheDB) Exist(data int64) bool {
	if db.coldData > data {
		return true
	}
	if data-db.coldData > int64(len(db.hotData)) {
		return false
	}
	return db.hotData[data-db.coldData-1]
}

func (db *CacheDB) Insert(data int64) {
	if db.coldData > data {
		return
	}
	if data-db.coldData <= int64(len(db.hotData)) {
		db.hotData[data-db.coldData-1] = true
		return
	}
	delta := data - db.coldData
	db.hotData = append(db.hotData, make([]bool, delta)...)
	db.hotData[data-db.coldData-1] = true

	// trim if too large
	if len(db.hotData) >= 2*db.trimNum {
		db.hotData = db.hotData[db.trimNum:]
		db.coldData += int64(db.trimNum)
	}
}
