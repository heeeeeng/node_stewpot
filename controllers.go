package main

import (
	"encoding/json"
	"html/template"
	"net/http"
)

func (s *Stewpot) MainPage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./static/index.html")
	t.Execute(w, "hello!")
}

func (s *Stewpot) StaticController(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir("./"))
	fs.ServeHTTP(w, r)
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

func (s *Stewpot) SendMsg(w http.ResponseWriter, r *http.Request) {

}
