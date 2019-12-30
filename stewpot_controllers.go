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

func (s *Stewpot) StaticController(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir("./"))
	fs.ServeHTTP(w, r)
}

func (s *Stewpot) RestartController(w http.ResponseWriter, r *http.Request) {
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

	s.RestartNetwork(nodeNum, maxIn, maxOut, types.FileSize(bandwidth))
}

func (s *Stewpot) GetNetworkGraph(w http.ResponseWriter, r *http.Request) {
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

		for _, neighbor := range n.peers {
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

func (s *Stewpot) SendMsg(w http.ResponseWriter, r *http.Request) {
	node := s.nodes[rand.Intn(len(s.nodes))]
	msg := s.GenerateMsg()
	timestamp := s.timeline.SendNewMsg(node, msg)
	fmt.Println(timestamp)
	w.Write([]byte(strconv.FormatInt(timestamp, 10)))
}

type Tasks []types.Task

func (s *Stewpot) GetTimeUnit(w http.ResponseWriter, r *http.Request) {
	timestamps := r.URL.Query()["time"]
	if len(timestamps) == 0 {
		return
	}
	timestamp, err := strconv.Atoi(timestamps[0])
	if err != nil {
		return
	}
	//fmt.Println("GetTimeUnit request: ", timestamp)
	timeUnit := s.timeline.GetTimeUnit(int64(timestamp))
	if timeUnit == nil {
		return
	}
	data, _ := json.Marshal(timeUnit.tasks)
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
