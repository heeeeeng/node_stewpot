package main

import (
	"fmt"
	"github.com/heeeeeng/node_stewpot/tasks"
	"github.com/heeeeeng/node_stewpot/types"
	"math/rand"
	"sync"
)

type Stewpot struct {
	msg       int64
	msgLocker sync.RWMutex

	bootstrap *Node
	nodes     []*Node
	timeline  *Timeline
}

func NewStewpot() *Stewpot {
	db := newMemDB()
	timeline := newTimeline(db, nil)

	s := &Stewpot{}
	s.msg = 1
	s.timeline = timeline

	return s
}

func (s *Stewpot) RestartNetwork(nodeNum, maxIn, maxOut, maxBest int, bandwidth int64, callback func(task types.Task)) {
	s.Stop()

	db := newMemDB()
	timeline := newTimeline(db, callback)

	s.msg = 1
	s.timeline = timeline
	s.nodes = make([]*Node, 0)
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
	s.nodes = append(s.nodes, s.bootstrap)

	for i := 1; i < nodeNum; i++ {
		loc = locSet[rand.Intn(len(locSet))]
		conf.IP = fmt.Sprintf("%d_%s", i, loc.Name)
		perf := 1024
		node := NewNode(conf, loc, perf)

		node.TryConnect(s.bootstrap)
		s.nodes = append(s.nodes, node)
	}
}

// Stew start the simulating of the nodes network.
func (s *Stewpot) Start() {
	s.timeline.Start()
}

func (s *Stewpot) Stop() {
	s.timeline.Stop()
}

func (s *Stewpot) PrintOutNodes() {
	for _, n := range s.nodes {
		fmt.Println(n.String())
	}
}

func (s *Stewpot) GenerateMsg(msgSize int64) types.Message {
	s.msgLocker.Lock()
	defer s.msgLocker.Unlock()

	msg := types.NewMessage(nil, 1, s.msg, msgSize)
	s.msg++

	return msg
}

func (s *Stewpot) SendNewMsg() {
	msg := s.GenerateMsg(types.DefualtMsgSize)
	s.SendMsg(msg)
}

func (s *Stewpot) SendMsg(msg types.Message) int64 {
	node := s.nodes[rand.Intn(len(s.nodes))]
	timestamp := s.timeline.SendNewMsg(node, msg)
	//fmt.Println("send msg at time: ", timestamp)
	return timestamp
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
		fmt.Println(fmt.Sprintf("iter %d, totalTime: %d", i, totalTime))
		totalTime += s.simulate(conf.MsgSize, callbackChan)
	}
	return totalTime / int64(conf.IterNum)
}

func (s *Stewpot) simulate(msgSize int64, cbChan chan cbData) int64 {
	nodes := make(map[string]struct{})
	for _, node := range s.nodes {
		nodes[node.IP()] = struct{}{}
	}

	msg := s.GenerateMsg(msgSize)
	startTime := s.SendMsg(msg)
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
