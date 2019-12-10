package main

type Graph struct {
	Nodes []GraphNode `json:"nodes"`
	Links []GraphLink `json:"links"`
}

type GraphNode struct {
	Name string `json:"name"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
}

type GraphLink struct {
	Source string `json:"source"`
	Target string `json:"target"`
}
