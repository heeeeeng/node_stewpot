package stewpot

type GraphStatus struct {
	links []graphLink
	nodes []graphNode
}

type graphLink struct {
	index int
	color string
	width int
}

type graphNode struct {
	index int
	color string
}
