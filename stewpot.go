package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
)

type Stewpot struct {
	bootstrap *Node
	nodes     []*Node
}

func NewStewpot() *Stewpot {
	s := &Stewpot{}

	return s
}

func (s *Stewpot) InitNetwork(nodeNum int) {
	loc := ConstLocations[rand.Intn(len(ConstLocations))]
	conf := NodeConfig{
		IP:       fmt.Sprintf("%d_%s", 0, loc.Name),
		Upload:   1024,
		Download: 1024,
		MaxIn:    8,
		MaxOut:   4,
	}
	perf := 1024
	s.bootstrap = NewNode(conf, loc, perf)
	s.nodes = append(s.nodes, s.bootstrap)

	for i := 1; i < nodeNum; i++ {
		loc := ConstLocations[rand.Intn(len(ConstLocations))]
		conf := NodeConfig{
			IP:       fmt.Sprintf("%d_%s", i, loc.Name),
			Upload:   1024,
			Download: 1024,
			MaxIn:    8,
			MaxOut:   4,
		}
		perf := 1024
		node := NewNode(conf, loc, perf)

		//node.TryConnect(s.nodes[rand.Intn(len(s.nodes))])
		node.TryConnect(s.bootstrap)
		s.nodes = append(s.nodes, node)
	}
}

// Stew start the simulating of the nodes network.
func (s *Stewpot) Start() {

}

func (s *Stewpot) PrintOutNodes() {
	for _, n := range s.nodes {
		fmt.Println(n.String())
	}
}

func (s *Stewpot) MainPage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./index.html")
	t.Execute(w, "hello!")
}

func (s *Stewpot) GetNetworkGraph(w http.ResponseWriter, r *http.Request) {
	w.Write(s.MarshalNodes())
}

func (s *Stewpot) MarshalNodes() []byte {
	var nodesJson []GraphNode
	var linksJson []GraphLink
	linksMap := make(map[string][]string)

	for _, n := range s.nodes {
		graphNode := GraphNode{
			Name: n.IP,
		}
		nodesJson = append(nodesJson, graphNode)

		for _, neighbor := range n.peers {
			if neighbor.out {
				linksMap[n.IP] = append(linksMap[n.IP], neighbor.ipRemote)
			} else {
				linksMap[neighbor.ipRemote] = append(linksMap[neighbor.ipRemote], n.IP)
			}
		}
	}
	for source, targets := range linksMap {
		for _, target := range targets {
			graphLink := GraphLink{
				Source: source,
				Target: target,
			}
			linksJson = append(linksJson, graphLink)
		}
	}

	graph := Graph{
		Nodes: nodesJson,
		Links: linksJson,
	}
	graphData, _ := json.Marshal(graph)
	return graphData
}
