package main

import (
	"fmt"
	"github.com/heeeeeng/node_stewpot/types"
	"strings"
)

var (
	DefaultDelay int64 = 100
)

type Node struct {
	db *NodeCacheDB

	ip        string
	bandwidth Bandwidth
	location  types.Location
	perf      int

	// peer related
	peerNumIn  int
	peerNumOut int
	maxIn      int
	maxOut     int
	blackList  map[string]struct{}
	peers      map[string]types.Peer

	CpuLocked bool

	close chan bool
}

type Bandwidth struct {
	Upload   types.FileSize
	Download types.FileSize
}

type NodeConfig struct {
	IP       string
	Upload   types.FileSize
	Download types.FileSize
	MaxIn    int
	MaxOut   int
}

func NewNode(config NodeConfig, loc types.Location, perf int) *Node {
	n := &Node{}

	n.db = newNodeCacheDB()

	n.ip = config.IP
	n.bandwidth = Bandwidth{Upload: config.Upload, Download: config.Download}
	n.location = loc
	n.perf = perf

	n.maxIn = config.MaxIn
	n.maxOut = config.MaxOut
	n.peerNumIn = 0
	n.peerNumOut = 0

	n.blackList = make(map[string]struct{})
	n.blackList[n.ip] = struct{}{}

	n.peers = make(map[string]types.Peer)

	n.CpuLocked = false

	closeC := make(chan bool)
	n.close = closeC

	return n
}

func (n *Node) Close() {
	close(n.close)
}

func (n *Node) IP() string {
	return n.ip
}

func (n *Node) Location() types.Location {
	return n.location
}

func (n *Node) Perf() int {
	return n.perf
}

func (n *Node) Bandwidth() types.FileSize {
	return n.bandwidth.Upload
}

func (n *Node) BandwidthInMillisecond() types.FileSize {
	return n.bandwidth.Upload / 1000
}

func (n *Node) Peers() map[string]types.Peer {
	return n.peers
}

func (n *Node) GetDelay(loc types.Location) int64 {
	if delay, ok := n.location.Delays[loc.Name]; ok {
		return delay
	}
	return DefaultDelay
}

func (n *Node) MsgExists(msg types.Message) bool {
	return n.db.Exist(msg.ID)
}

func (n *Node) StoreMsg(msg types.Message) {
	n.db.Insert(msg.ID)
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

func (n *Node) AddPeer(p *Peer) {
	n.peers[p.ipRemote] = p
	if p.out {
		n.peerNumOut++
	} else {
		n.peerNumIn++
	}
}

func (n *Node) ConnectIn(remoteNode types.Node) (bool, []types.Node) {
	fmt.Println(fmt.Sprintf("[ConnectIn]\t node %s connect in %s", remoteNode.IP(), n.ip))
	var connected bool

	if n.peerNumIn < n.maxIn {
		peer := NewPeer(n.ip, remoteNode.IP(), false, remoteNode)
		n.AddPeer(peer)
		n.blackList[remoteNode.IP()] = struct{}{}
		connected = true
	}

	fmt.Println(fmt.Sprintf("start return neighbor: %s", n.String()))
	var neighbors []types.Node
	for _, p := range n.peers {
		neighbors = append(neighbors, p.GetNode())
	}
	return connected, neighbors
}

func (n *Node) TryConnect(remoteNode types.Node) {
	fmt.Println(fmt.Sprintf("[TryConnect]\t node %s try to connect: %s", n.ip, remoteNode.IP()))

	if _, inBlackList := n.blackList[remoteNode.IP()]; inBlackList {
		return
	}
	n.blackList[remoteNode.IP()] = struct{}{}

	if n.peerNumOut >= n.maxOut {
		return
	}
	connected, neighbors := remoteNode.ConnectIn(n)
	if connected {
		peer := NewPeer(n.ip, remoteNode.IP(), true, remoteNode)
		n.AddPeer(peer)
	}

	for _, neighbor := range neighbors {
		if n.peerNumOut >= n.maxOut {
			return
		}
		n.TryConnect(neighbor)
	}
}

func (n *Node) String() string {
	var neighbors []string
	for _, p := range n.peers {
		neighborInfo := p.RemoteIP()
		if p.Out() {
			neighborInfo += "_out"
		} else {
			neighborInfo += "_in"
		}
		neighbors = append(neighbors, neighborInfo)
	}
	return fmt.Sprintf("{ ip: %s, neighbors: %s }", n.ip, strings.Join(neighbors, ", "))
}

type NodeCacheDB struct {
	trimNum  int
	hotData  []bool
	coldData int64
}

func newNodeCacheDB() *NodeCacheDB {
	return &NodeCacheDB{
		trimNum:  1000,
		hotData:  make([]bool, 0),
		coldData: 0,
	}
}

func (db *NodeCacheDB) Exist(data int64) bool {
	if db.coldData >= data {
		return true
	}
	if data-db.coldData-1 >= int64(len(db.hotData)) {
		return false
	}
	return db.hotData[data-db.coldData-1]
}

func (db *NodeCacheDB) Insert(data int64) {
	if db.coldData >= data {
		return
	}
	delta := data - db.coldData
	if delta <= int64(len(db.hotData)) {
		db.hotData[delta-1] = true
		return
	}
	db.hotData = append(db.hotData, make([]bool, int(delta)-len(db.hotData))...)
	db.hotData[delta-1] = true

	// trim if too large
	if len(db.hotData) >= 2*db.trimNum {
		db.hotData = db.hotData[db.trimNum:]
		db.coldData += int64(db.trimNum)
	}
}
