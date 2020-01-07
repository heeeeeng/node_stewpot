package main

import (
	"encoding/json"
	"fmt"
	"github.com/heeeeeng/node_stewpot/types"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
)

func (s *Stewpot) MainPage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./static/index.html")
	t.Execute(w, "hello!")
}

func (s *Stewpot) SimPage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./static/multi_simulate.html")
	t.Execute(w, "hello!")

}

func (s *Stewpot) CtrlStatic(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir("./"))
	fs.ServeHTTP(w, r)
}

func (s *Stewpot) CtrlRestart(w http.ResponseWriter, r *http.Request) {
	nodeNum, err := urlParamToInt(r, "node_num")
	if err != nil {
		return
	}
	maxIn, err := urlParamToInt(r, "max_in")
	if err != nil {
		return
	}
	maxOut, err := urlParamToInt(r, "max_out")
	if err != nil {
		return
	}
	bandwidth, err := urlParamToInt(r, "bandwidth")
	if err != nil {
		return
	}
	maxBest := 4

	s.RestartNetwork(nodeNum, maxIn, maxOut, maxBest, int64(bandwidth), nil)
}

func (s *Stewpot) CtrlGetNetworkGraph(w http.ResponseWriter, r *http.Request) {
	w.Write(s.MarshalNodes())
}

func (s *Stewpot) MarshalNodes() []byte {
	var nodesJson []types.GraphNode
	var linksJson []types.GraphLink
	linksMap := make(map[string][]string)

	for _, n := range s.nodes {
		graphNode := types.GraphNode{
			Name: n.ip,
		}
		nodesJson = append(nodesJson, graphNode)

		for _, neighbor := range n.Peers() {
			if neighbor.Out() {
				linksMap[n.ip] = append(linksMap[n.ip], neighbor.RemoteIP())
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

func (s *Stewpot) CtrlSendMsg(w http.ResponseWriter, r *http.Request) {
	var size int64
	msgSize, err := urlParamToInt(r, "msg_size")
	if err != nil {
		fmt.Println("send msg param error: ", err)
		size = types.DefualtMsgSize
	} else {
		size = int64(msgSize) * types.Byte
	}

	node := s.nodes[rand.Intn(len(s.nodes))]
	msg := s.GenerateMsg(size)
	timestamp := s.timeline.SendNewMsg(node, msg)
	fmt.Println(timestamp)
	w.Write([]byte(strconv.FormatInt(timestamp, 10)))
}

type Tasks []types.Task

func (s *Stewpot) CtrlGetTimeUnit(w http.ResponseWriter, r *http.Request) {
	timestamps := r.URL.Query()["time"]
	if len(timestamps) == 0 {
		return
	}
	timestamp, err := strconv.Atoi(timestamps[0])
	if err != nil {
		return
	}
	//fmt.Println("CtrlGetTimeUnit request: ", timestamp)
	timeUnit := s.timeline.GetTimeUnit(int64(timestamp))
	if timeUnit == nil {
		return
	}
	data, _ := json.Marshal(timeUnit.tasks)
	w.Write(data)
}

type MultiSimResp struct {
	Xs []string `json:"xs"`
	Ys []int64  `json:"ys"`
}

func (s *Stewpot) CtrlMultiSimulate(w http.ResponseWriter, r *http.Request) {
	// load data
	iterNum, err := urlParamToInt(r, "iter_num")
	if err != nil {
		fmt.Println("iter_num param error: ", err)
		return
	}
	nodeNum, err := urlParamToInt(r, "node_num")
	if err != nil {
		fmt.Println("node_num param error: ", err)
		return
	}
	maxIn, err := urlParamToInt(r, "max_in")
	if err != nil {
		fmt.Println("max_in param error: ", err)
		return
	}
	maxOut, err := urlParamToInt(r, "max_out")
	if err != nil {
		fmt.Println("max_out param error: ", err)
		return
	}
	bandwidth, err := urlParamToInt(r, "bandwidth")
	if err != nil {
		fmt.Println("bandwidth param error: ", err)
		return
	}
	maxMsgSize, err := urlParamToInt(r, "max_msg_size")
	if err != nil {
		fmt.Println("max_msg_size param error: ", err)
		return
	}

	// init config
	conf := SimConfig{
		IterNum:   iterNum,
		NodeNum:   nodeNum,
		Bandwidth: int64(bandwidth),
		MaxIn:     maxIn,
		MaxOut:    maxOut,
	}

	// simulate loop
	var xs []string
	var ys []int64
	// TODO min msg size and enlarge number per iter should not be hard coded.
	minMsgSize := 256 * types.Byte
	sizeGap := 16 * types.KB
	for i := minMsgSize; i < int64(maxMsgSize); i += sizeGap {
		conf.MsgSize = i
		avgTimeUsage := s.MultiSimulate(conf)

		x := strconv.Itoa(int(i / types.Byte))
		xs = append(xs, x)

		ys = append(ys, avgTimeUsage)

		fmt.Println(fmt.Sprintf("x: %d KB", i/types.KB))
	}

	// response
	resp := MultiSimResp{
		Xs: xs,
		Ys: ys,
	}
	data, _ := json.Marshal(resp)
	w.Write(data)
}

func urlParamToInt(r *http.Request, key string) (int, error) {
	paramStr := r.URL.Query()[key]
	if len(paramStr) == 0 {
		return 0, fmt.Errorf("cannot find parameter: %s", key)
	}
	param, err := strconv.Atoi(paramStr[0])
	if err != nil {
		return param, fmt.Errorf("cannot format key %s from string to int: %v", key, err)
	}
	return param, nil
}
