package smf_context

type UEPathGraph struct {
	SUPI  string
	Graph []*UEPathNode
}

type UEPathNode struct {
	UPFName             string
	Parent              string
	Neighbors           map[string]*UEPathNode
	IsBranchingPoint    bool
	EndPointOfEachChild map[string]*UEPathEndPoint
}

type UEPathEndPoint struct {
	EndPointIP   string
	EndPointPort string
}

func (node *UEPathNode) AddNeighbor(neighbor *UEPathNode) {
	//check if neighbor exist first

	if _, exist := node.Neighbors[neighbor.UPFName]; !exist {
		node.Neighbors[neighbor.UPFName] = neighbor
	}
}

//Add End Point Info to of child node to the map "EndPointOfEachChild"
//If the node is leaf node, it will add the end point info for itself name.
func (node *UEPathNode) AddEndPointOfChild(neighbor *UEPathNode, EndPoint *UEPathEndPoint) {

	if _, exist := node.EndPointOfEachChild[neighbor.UPFName]; !exist {
		node.EndPointOfEachChild[neighbor.UPFName] = EndPoint
	}
}

func (node *UEPathNode) RmbParent(parent string) {

	node.Parent = parent
}

func (node *UEPathNode) GetChild() []*UEPathNode {

	child := make([]*UEPathNode, 0)
	for upfName, upfNode := range node.Neighbors {
		if upfName != node.Parent {
			child = append(child, upfNode)
		}
	}

	return child
}

func (node *UEPathNode) IsLeafNode() bool {

	if len(node.Neighbors) == 1 {

		if _, exist := node.Neighbors[node.Parent]; exist {
			return true
		}
	}

	return false
}

func NewUEPathNode(name string) (node *UEPathNode) {
	node = &UEPathNode{
		UPFName:             name,
		Neighbors:           make(map[string]*UEPathNode),
		EndPointOfEachChild: make(map[string]*UEPathEndPoint),
		IsBranchingPoint:    false,
	}
	return
}

//check a given upf name is a branching point or not
func (uepg *UEPathGraph) IsBranchingPoint(name string) bool {

	for _, upfNode := range uepg.Graph {

		if name == upfNode.UPFName {
			return upfNode.IsBranchingPoint
		}
	}

	return false
}

func NewUEPathGraph(SUPI string) (UEPGraph *UEPathGraph) {

	UEPGraph = new(UEPathGraph)
	UEPGraph.Graph = make([]*UEPathNode, 0)
	UEPGraph.SUPI = SUPI

	paths := smfContext.UERoutingPaths[SUPI]
	lowerBound := 0

	NodeCreated := make(map[string]*UEPathNode)

	for _, path := range paths {
		upperBound := len(path.UPF) - 1

		pathEndPoint := &UEPathEndPoint{
			EndPointIP:   path.DestinationIP,
			EndPointPort: path.DestinationPort,
		}
		for idx, node_name := range path.UPF {

			var ue_node *UEPathNode
			var child_node *UEPathNode
			var parent_node *UEPathNode
			var exist bool

			if ue_node, exist = NodeCreated[node_name]; !exist {
				ue_node = NewUEPathNode(node_name)
				NodeCreated[node_name] = ue_node
				UEPGraph.Graph = append(UEPGraph.Graph, ue_node)
			}

			switch idx {
			case lowerBound:
				child_name := path.UPF[idx+1]

				if child_node, exist = NodeCreated[child_name]; !exist {
					child_node = NewUEPathNode(child_name)
					NodeCreated[child_name] = child_node
					UEPGraph.Graph = append(UEPGraph.Graph, child_node)
				}

				ue_node.AddNeighbor(child_node)
				ue_node.AddEndPointOfChild(child_node, pathEndPoint)

			case upperBound:
				parent_name := path.UPF[idx-1]

				if parent_node, exist = NodeCreated[parent_name]; !exist {
					parent_node = NewUEPathNode(parent_name)
					NodeCreated[parent_name] = parent_node
					UEPGraph.Graph = append(UEPGraph.Graph, parent_node)
				}

				ue_node.AddNeighbor(parent_node)
				ue_node.AddEndPointOfChild(ue_node, pathEndPoint)
				ue_node.RmbParent(parent_name)
			default:
				child_name := path.UPF[idx+1]

				if child_node, exist = NodeCreated[child_name]; !exist {
					child_node = NewUEPathNode(child_name)
					NodeCreated[child_name] = child_node
					UEPGraph.Graph = append(UEPGraph.Graph, child_node)
				}

				parent_name := path.UPF[idx-1]

				if parent_node, exist = NodeCreated[parent_name]; !exist {
					parent_node = NewUEPathNode(parent_name)
					NodeCreated[parent_name] = parent_node
					UEPGraph.Graph = append(UEPGraph.Graph, parent_node)
				}

				ue_node.AddNeighbor(child_node)
				ue_node.AddEndPointOfChild(child_node, pathEndPoint)
				ue_node.AddNeighbor(parent_node)
				ue_node.RmbParent(parent_name)
			}

		}
	}

	return
}

func (uepg *UEPathGraph) FindBranchingPoints() {
	//BFS algo implementation
	const (
		WHITE int = 0
		GREY  int = 1
		BLACK int = 2
	)

	num_of_nodes := len(uepg.Graph)

	color := make(map[string]int)
	distance := make(map[string]int)
	queue := make(chan *UEPathNode, num_of_nodes)

	for _, node := range uepg.Graph {

		color[node.UPFName] = WHITE
		distance[node.UPFName] = num_of_nodes + 1
	}

	cur_idx := 0 // start point
	for j := 0; j < num_of_nodes; j++ {

		cur_name := uepg.Graph[cur_idx].UPFName
		if color[cur_name] == WHITE {
			color[cur_name] = GREY
			distance[cur_name] = 0

			queue <- uepg.Graph[cur_idx]
			for len(queue) > 0 {
				node := <-queue
				branchingCount := 0
				for neighbor_name, neighbor_node := range node.Neighbors {

					if color[neighbor_name] == WHITE {
						color[neighbor_name] = GREY
						distance[neighbor_name] = distance[cur_name] + 1
						queue <- neighbor_node
					}

					if color[neighbor_name] == WHITE || color[neighbor_name] == GREY {
						branchingCount += 1
					}
				}

				if branchingCount >= 2 {
					node.IsBranchingPoint = true
				}
				color[node.UPFName] = BLACK
			}
		}

		//Keep finding other connected components
		cur_idx = j
	}

}
