package stewpot

import (
	"encoding/json"
	"fmt"
	"github.com/heeeeeng/node_stewpot/protocols"
	"github.com/heeeeeng/node_stewpot/tasks"
	"github.com/heeeeeng/node_stewpot/types"
	"math/rand"
	"sync"
)

type Stewpot struct {
	db *MemDB

	msgID     int64
	msgLocker sync.RWMutex

	bootstrap *Node
	nodes     *NodeSet
	timeline  *Timeline
	protocol  types.Protocol
}

func NewStewpot() *Stewpot {
	s := &Stewpot{}
	s.msgID = 1

	return s
}

func (s *Stewpot) RestartNetwork(nodeNum, maxIn, maxOut, maxBest int, bandwidth int64, callback func(task types.Task)) {
	s.Stop()

	db := newMemDB()
	s.db = db
	timeline := newTimeline(db, s.protocol, callback)

	s.msgID = 1
	s.timeline = timeline
	s.nodes = NewNodeSet()
	s.bootstrap = nil

	s.InitNetwork(nodeNum, maxIn, maxOut, maxBest, bandwidth)
	s.Start()
}

func (s *Stewpot) InitNetwork(nodeNum, maxIn, maxOut, maxBest int, bandwidth int64) {
	locSet := []types.Location{types.LocCN, types.LocSEA, types.LocJP, types.LocRU, types.LocNA, types.LocEU}
	loc := locSet[rand.Intn(len(locSet))]

	conf := NodeConfig{
		IP:       fmt.Sprintf("%d_%s", 0, loc.Name),
		Upload:   bandwidth,
		Download: bandwidth,
		MaxIn:    maxIn,
		MaxOut:   maxOut,
		MaxBest:  maxBest,
	}
	perf := 1024
	s.bootstrap = NewNode(conf, loc, perf)
	s.nodes.Add(s.bootstrap)

	for i := 1; i < nodeNum; i++ {
		loc = locSet[rand.Intn(len(locSet))]
		conf.IP = fmt.Sprintf("%d_%s", i, loc.Name)
		perf := 1024
		node := NewNode(conf, loc, perf)

		node.TryConnect(s.bootstrap)
		s.nodes.Add(node)
	}

	s.protocol = protocols.NewFakeProtocol()

	db := newMemDB()
	s.db = db
	s.timeline = newTimeline(db, s.protocol, nil)
}

// Starts the simulating of the nodes network.
func (s *Stewpot) Start() {
	s.timeline.Start()
}

func (s *Stewpot) Stop() {
	s.timeline.Stop()
}

func (s *Stewpot) PrintOutNodes() {
	for _, key := range s.nodes.Keys() {
		node := s.nodes.Get(key)
		fmt.Println(node.String())
	}
}

func (s *Stewpot) Node(id string) types.Node {
	return s.nodes.Get(id)
}

func (s *Stewpot) Nodes() map[string]types.Node {
	return s.nodes.NodesMap()
}

// generate a new message without source node.
func (s *Stewpot) GenerateMsg(difficulty int64, msgSize int64, content string) types.Message {
	s.msgLocker.Lock()
	defer s.msgLocker.Unlock()

	msg := types.NewMessage(nil, difficulty, s.msgID, msgSize, content)
	s.msgID++
	s.db.InsertMessage(&msg)

	return msg
}

func (s *Stewpot) SendMsg(source types.Node, msg types.Message) int64 {
	timestamp := s.timeline.SendNewMsg(source, msg)
	return timestamp
}

func (s *Stewpot) GetMsg(msgID int64) *types.Message {
	return s.db.GetMessage(msgID)
}

func (s *Stewpot) GetTimeUnitTasks(t int64) []types.Task {
	timeUnit := s.timeline.GetTimeUnit(t)
	if timeUnit == nil {
		return nil
	}
	return timeUnit.tasks
}

func (s *Stewpot) MsgProducer(src types.Node, difficulty int64, size int64, content string) {
	msg := s.GenerateMsg(difficulty, size, content)
	s.timeline.SendNewMsg(src, msg)
}

type SimConfig struct {
	IterNum   int
	MsgSize   int64
	NodeNum   int
	Bandwidth int64
	MaxIn     int
	MaxOut    int
}

type cbData struct {
	nodeIP    string
	msgID     int64
	timestamp int64
}

func (s *Stewpot) MultiSimulate(conf SimConfig) int64 {
	callbackChan := make(chan cbData)

	callback := func(task types.Task) {
		if task.Type() != int(tasks.TaskTypeMsgProcessCPUReq) {
			return
		}
		t := task.(*tasks.MsgProcessCPUReqTask)
		data := cbData{
			nodeIP:    t.Node().IP(),
			msgID:     t.Msg(),
			timestamp: t.StartTime(),
		}
		callbackChan <- data
	}
	// restart network
	s.RestartNetwork(conf.NodeNum, conf.MaxIn, conf.MaxOut, 4, conf.Bandwidth, callback)

	totalTime := int64(0)
	for i := 0; i < conf.IterNum; i++ {
		//fmt.Println(fmt.Sprintf("iter %d, totalTime: %d", i, totalTime))
		totalTime += s.simulate(conf.MsgSize, callbackChan)
	}
	return totalTime / int64(conf.IterNum)
}

func (s *Stewpot) simulate(msgSize int64, cbChan chan cbData) int64 {
	nodes := make(map[string]struct{})
	for _, ip := range s.nodes.Keys() {
		nodes[ip] = struct{}{}
	}

	msg := s.GenerateMsg(1*types.SizeBit, msgSize, "")

	srcNode := s.nodes.RandomNode()
	startTime := s.SendMsg(srcNode, msg)
	endTime := int64(0)

	for {
		if len(nodes) == 0 {
			break
		}

		select {
		case data := <-cbChan:
			if data.msgID != msg.ID {
				continue
			}
			if _, ok := nodes[data.nodeIP]; !ok {
				continue
			}
			delete(nodes, data.nodeIP)
			endTime = data.timestamp
		}
	}

	return endTime - startTime
}

func (s *Stewpot) MarshalNodes() []byte {
	var nodesJson []types.GraphNode
	var linksJson []types.GraphLink
	linksMap := make(map[string][]string)

	for _, n := range s.nodes.NodesMap() {
		//n := s.nodes.GetTimeUnit(key)

		graphNode := types.GraphNode{
			Name: n.IP(),
		}
		nodesJson = append(nodesJson, graphNode)

		for _, neighbor := range n.Peers() {
			if neighbor.Out() {
				linksMap[n.IP()] = append(linksMap[n.IP()], neighbor.RemoteIP())
			}
		}
	}
	for source, targets := range linksMap {
		for _, target := range targets {
			graphLink := types.GraphLink{
				Source: source,
				Target: target,
			}
			linksJson = append(linksJson, graphLink)
		}
	}

	graph := types.Graph{
		Nodes: nodesJson,
		Links: linksJson,
	}
	graphData, _ := json.Marshal(graph)
	return graphData
}
