package stewpot

import (
	"github.com/heeeeeng/node_stewpot/types"
	"math/rand"
)

type NodeSet struct {
	keys  []string
	nodes map[string]types.Node
}

func NewNodeSet() *NodeSet {
	return &NodeSet{}
}

func (ns *NodeSet) Len() int {
	return len(ns.keys)
}

func (ns *NodeSet) Keys() []string {
	return ns.keys
}

func (ns *NodeSet) NodesMap() map[string]types.Node {
	return ns.nodes
}

func (ns *NodeSet) NodesList() []types.Node {
	var nodes []types.Node
	for _, node := range ns.nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

func (ns *NodeSet) RandomNode() types.Node {
	key := ns.keys[rand.Intn(len(ns.keys))]
	return ns.Get(key)
}

func (ns *NodeSet) Get(key string) types.Node {
	return ns.nodes[key]
}

func (ns *NodeSet) Add(node types.Node) {
	_, exists := ns.nodes[node.IP()]
	if !exists {
		ns.keys = append(ns.keys, node.IP())
	}
	ns.nodes[node.IP()] = node
}

func (ns *NodeSet) Update(node *Node) {
	// no need to implement yet.
}
