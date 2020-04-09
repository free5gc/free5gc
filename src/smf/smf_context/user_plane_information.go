package smf_context

import (
	"fmt"
	"free5gc/lib/pfcp/pfcpType"
	"free5gc/lib/util_3gpp"
	"free5gc/src/smf/factory"
	"free5gc/src/smf/logger"
	"net"
	"reflect"
)

// UserPlaneInformation store userplane topology
type UserPlaneInformation struct {
	UPNodes              map[string]*UPNode
	UPFs                 map[string]*UPNode
	AccessNetwork        map[string]*UPNode
	UPFIPToName          map[string]string
	UPFsID               map[string]string    // name to id
	UPFsIPtoID           map[string]string    // ip->id table, for speed optimization
	DefaultUserPlanePath map[string][]*UPNode // DNN to Default Path
}

type UPNodeType string

const (
	UPNODE_UPF UPNodeType = "UPF"
	UPNODE_AN  UPNodeType = "AN"
)

// UPNode represent the user plane node topology
type UPNode struct {
	Type         UPNodeType
	NodeID       pfcpType.NodeID
	UPResourceIP net.IP
	ANIP         net.IP
	Dnn          string
	Links        []*UPNode
	UPF          *UPF
}

// UPPath represent User Plane Sequence of this path
type UPPath []*UPNode

func AllocateUPFID() {
	UPFsID := smfContext.UserPlaneInformation.UPFsID
	UPFsIPtoID := smfContext.UserPlaneInformation.UPFsIPtoID

	for upfName, upfNode := range smfContext.UserPlaneInformation.UPFs {
		upfid := upfNode.UPF.UUID()
		upfip := upfNode.NodeID.ResolveNodeIdToIp().String()

		UPFsID[upfName] = upfid
		UPFsIPtoID[upfip] = upfid
	}
}

func processUPTopology(upTopology *factory.UserPlaneInformation) {
	nodePool := make(map[string]*UPNode)
	upfPool := make(map[string]*UPNode)
	anPool := make(map[string]*UPNode)
	upfIPMap := make(map[string]string)

	for name, node := range upTopology.UPNodes {
		upNode := new(UPNode)
		upNode.Type = UPNodeType(node.Type)
		switch upNode.Type {
		case UPNODE_AN:
			upNode.ANIP = net.ParseIP(node.ANIP)
			anPool[name] = upNode
		case UPNODE_UPF:
			//ParseIp() always return 16 bytes
			//so we can't use the length of return ip to seperate IPv4 and IPv6
			//This is just a work around
			var ip net.IP
			if net.ParseIP(node.NodeID).To4() == nil {

				ip = net.ParseIP(node.NodeID)
			} else {

				ip = net.ParseIP(node.NodeID).To4()
			}

			switch len(ip) {
			case net.IPv4len:
				upNode.NodeID = pfcpType.NodeID{
					NodeIdType:  pfcpType.NodeIdTypeIpv4Address,
					NodeIdValue: ip,
				}
			case net.IPv6len:
				upNode.NodeID = pfcpType.NodeID{
					NodeIdType:  pfcpType.NodeIdTypeIpv6Address,
					NodeIdValue: ip,
				}
			default:
				upNode.NodeID = pfcpType.NodeID{
					NodeIdType:  pfcpType.NodeIdTypeFqdn,
					NodeIdValue: []byte(node.NodeID),
				}
			}

			upfPool[name] = upNode
		default:
			logger.InitLog.Warningf("invalid UPNodeType: %s\n", upNode.Type)
		}

		nodePool[name] = upNode

		ipStr := upNode.NodeID.ResolveNodeIdToIp().String()
		upfIPMap[ipStr] = name
	}

	for _, link := range upTopology.Links {
		nodeA := nodePool[link.A]
		nodeB := nodePool[link.B]
		if nodeA == nil || nodeB == nil {
			logger.InitLog.Warningf("UPLink [%s] <=> [%s] not establish\n", link.A, link.B)
			continue
		}
		nodeA.Links = append(nodeA.Links, nodeB)
		nodeB.Links = append(nodeB.Links, nodeA)
	}

	//Initialize each UPF
	for _, upfNode := range upfPool {
		upfNode.UPF = NewUPF(&upfNode.NodeID)
	}

	smfContext.UserPlaneInformation.UPNodes = nodePool
	smfContext.UserPlaneInformation.UPFs = upfPool
	smfContext.UserPlaneInformation.AccessNetwork = anPool
	smfContext.UserPlaneInformation.UPFIPToName = upfIPMap
	smfContext.UserPlaneInformation.UPFsID = make(map[string]string)
	smfContext.UserPlaneInformation.UPFsIPtoID = make(map[string]string)
	smfContext.UserPlaneInformation.DefaultUserPlanePath = make(map[string][]*UPNode)
}

func (upi *UserPlaneInformation) GetUPFNameByIp(ip string) string {

	return upi.UPFIPToName[ip]
}

func (upi *UserPlaneInformation) GetUPFNodeIDByName(name string) pfcpType.NodeID {

	return upi.UPFs[name].NodeID
}

func (upi *UserPlaneInformation) GetUPFNodeByIP(ip string) *UPNode {
	upfName := upi.GetUPFNameByIp(ip)
	return upi.UPFs[upfName]
}

func (upi *UserPlaneInformation) GetUPFIDByIP(ip string) string {

	return upi.UPFsIPtoID[ip]
}

func (upi *UserPlaneInformation) GetDefaultUserPlanePathByDNN(dnn string) (path UPPath) {
	path, pathExist := upi.DefaultUserPlanePath[dnn]

	if pathExist {
		return
	} else {
		pathExist = upi.GenerateDefaultPath(dnn)
		if pathExist {
			return upi.DefaultUserPlanePath[dnn]
		}
	}
	return nil
}

func (upi *UserPlaneInformation) ExistDefaultPath(dnn string) bool {

	_, exist := upi.DefaultUserPlanePath[dnn]
	return exist
}

func GenerateDataPath(upPath UPPath, smContext *SMContext) (root *DataPathNode) {
	if len(upPath) < 1 {
		logger.CtxLog.Errorf("invalid path")
	}
	var lowerBound = 0
	var upperBound = len(upPath) - 1
	var curDataPathNode *DataPathNode
	var prevDataPathNode *DataPathNode

	for idx, upNode := range upPath {
		curDataPathNode := NewDataPathNode()
		curDataPathNode.UPF = upNode.UPF
		curDataPathNode.InUse = true
		if idx == lowerBound {
			root = curDataPathNode
			root.DataPathToAN = NewDataPathDownLink()
			root.SetUpLinkSrcNode(nil)
		}
		if idx == upperBound {
			curDataPathNode.SetDownLinkSrcNode(nil)
			// curDataPathNode.AddParent(prevDataPathNode)
			curDataPathNode.DLDataPathLinkForPSA = NewDataPathUpLink()
			// prevDataPathNode.AddChild(curDataPathNode)
		}
		if prevDataPathNode != nil {
			prevDataPathNode.SetDownLinkSrcNode(curDataPathNode)
			curDataPathNode.SetUpLinkSrcNode(prevDataPathNode)

			curDataPathNode.AddParent(prevDataPathNode)
			prevDataPathNode.AddChild(curDataPathNode)
		}
		prevDataPathNode = curDataPathNode
	}

	curDataPathNode = root
	for curDataPathNode != nil {
		fmt.Println("calculate ", curDataPathNode.UPF.PFCPAddr().String())
		curULTunnel := curDataPathNode.UpLinkTunnel
		curDLTunnel := curDataPathNode.DownLinkTunnel

		// Setup UpLink PDR
		if curULTunnel != nil {
			ULPDR := curULTunnel.MatchedPDR
			ULDestUPF := curULTunnel.DestEndPoint.UPF

			ULPDR.Precedence = 32
			ULPDR.PDI = PDI{
				SourceInterface: pfcpType.SourceInterface{InterfaceValue: pfcpType.SourceInterfaceAccess},
				LocalFTeid: &pfcpType.FTEID{
					V4:          true,
					Ipv4Address: ULDestUPF.UPIPInfo.Ipv4Address,
					Teid:        curULTunnel.TEID,
				},
				UEIPAddress: &pfcpType.UEIPAddress{
					V4:          true,
					Ipv4Address: smContext.PDUAddress.To4(),
				},
			}
			ULPDR.OuterHeaderRemoval = &pfcpType.OuterHeaderRemoval{OuterHeaderRemovalDescription: pfcpType.OuterHeaderRemovalGtpUUdpIpv4}

			ULFAR := ULPDR.FAR

			if curDLTunnel != nil {
				if nextULDest := curDLTunnel.SrcEndPoint; nextULDest != nil {
					nextULTunnel := nextULDest.UpLinkTunnel
					ULFAR.ApplyAction = pfcpType.ApplyAction{Buff: false, Drop: false, Dupl: false, Forw: true, Nocp: false}
					ULFAR.ForwardingParameters = &ForwardingParameters{
						DestinationInterface: pfcpType.DestinationInterface{InterfaceValue: pfcpType.DestinationInterfaceCore},
						OuterHeaderCreation: &pfcpType.OuterHeaderCreation{
							OuterHeaderCreationDescription: pfcpType.OuterHeaderCreationGtpUUdpIpv4,
							Ipv4Address:                    nextULTunnel.DestEndPoint.UPF.UPIPInfo.Ipv4Address,
							Teid:                           nextULTunnel.TEID,
						},
					}
				}
			}

		}

		// Setup DownLink
		if curDLTunnel != nil {
			DLPDR := curDLTunnel.MatchedPDR
			DLDestUPF := curDLTunnel.DestEndPoint.UPF

			DLPDR.Precedence = 32
			DLPDR.PDI = PDI{
				SourceInterface: pfcpType.SourceInterface{InterfaceValue: pfcpType.SourceInterfaceCore},
				LocalFTeid: &pfcpType.FTEID{
					V4:          true,
					Ipv4Address: DLDestUPF.UPIPInfo.Ipv4Address,
					Teid:        curDLTunnel.TEID,
				},

				// TODO: Should Uncomment this after FR5GC-1029 is solved
				// UEIPAddress: &pfcpType.UEIPAddress{
				// 	V4:          true,
				// 	Ipv4Address: smContext.PDUAddress.To4(),
				// },
			}

			// TODO: Should delete this after FR5GC-1029 is solved
			if curDataPathNode.IsAnchorUPF() {
				DLPDR.PDI.UEIPAddress = &pfcpType.UEIPAddress{
					V4:          true,
					Ipv4Address: smContext.PDUAddress.To4(),
				}
			}

			fmt.Println("In GenerateDataPath")
			fmt.Println("curDataPathNode IP: ", curDataPathNode.GetNodeIP())
			fmt.Println("Is anchor point: ", curDataPathNode.IsAnchorUPF())

			if !curDataPathNode.IsAnchorUPF() {
				DLPDR.OuterHeaderRemoval = &pfcpType.OuterHeaderRemoval{OuterHeaderRemovalDescription: pfcpType.OuterHeaderRemovalGtpUUdpIpv4}
			}

			DLFAR := DLPDR.FAR

			nextDLTunnel := curDLTunnel.DestEndPoint.DownLinkTunnel
			fmt.Println("DestEndPoint TEID", nextDLTunnel.TEID)
			fmt.Println("SrcEndPoint TEID", curDLTunnel.DestEndPoint.UpLinkTunnel.TEID)
			//fmt.Println("SrcEndPoint IP", curDLTunnel.DestEndPoint.UpLinkTunnel.SrcEndPoint.GetNodeIP())

			if nextDLDest := curULTunnel.SrcEndPoint; nextDLDest != nil {
				fmt.Println("In GenerateDataPath")
				fmt.Println(nextDLDest.GetNodeIP())
				DLFAR.ApplyAction = pfcpType.ApplyAction{Buff: false, Drop: false, Dupl: false, Forw: true, Nocp: false}
				DLFAR.ForwardingParameters = &ForwardingParameters{
					DestinationInterface: pfcpType.DestinationInterface{InterfaceValue: pfcpType.DestinationInterfaceAccess},
					OuterHeaderCreation: &pfcpType.OuterHeaderCreation{
						OuterHeaderCreationDescription: pfcpType.OuterHeaderCreationGtpUUdpIpv4,
						Ipv4Address:                    nextDLDest.UPF.NodeID.ResolveNodeIdToIp(),
						Teid:                           curDLTunnel.DestEndPoint.UpLinkTunnel.TEID,
					},
				}
			}
		}
		if curDataPathNode.DownLinkTunnel != nil {
			if curDataPathNode.DownLinkTunnel.SrcEndPoint == nil {
				DNDLPDR := curDataPathNode.DownLinkTunnel.MatchedPDR
				DNDLPDR.PDI = PDI{
					SourceInterface: pfcpType.SourceInterface{InterfaceValue: pfcpType.SourceInterfaceCore},
					NetworkInstance: util_3gpp.Dnn(smContext.Dnn),
					UEIPAddress: &pfcpType.UEIPAddress{
						V4:          true,
						Ipv4Address: smContext.PDUAddress.To4(),
					},
				}
				break
			}
		}

		if curDataPathNode.DownLinkTunnel == nil {
			break
		}
		curDataPathNode = curDataPathNode.DownLinkTunnel.SrcEndPoint
	}

	return
}

func (upi *UserPlaneInformation) GenerateDefaultPath(dnn string) (pathExist bool) {

	var source *UPNode
	var destination *UPNode

	for _, node := range upi.AccessNetwork {

		if node.Type == UPNODE_AN {
			source = node
			break
		}
	}

	if source == nil {
		logger.CtxLog.Errorf("There is no AN Node in config file!")
		return false
	}

	for _, node := range upi.UPFs {

		if node.UPF.UPIPInfo.NetworkInstance != nil {
			node_dnn := string(node.UPF.UPIPInfo.NetworkInstance)
			if node_dnn == dnn {
				destination = node
				break
			}
		}
	}

	if destination == nil {
		logger.CtxLog.Errorf("Can't find UPF with DNN [%s]\n", dnn)
		return false
	}

	//Run DFS
	var visited map[*UPNode]bool
	visited = make(map[*UPNode]bool)

	for _, upNode := range upi.UPNodes {
		visited[upNode] = false
	}

	var path []*UPNode
	path, pathExist = getPathBetween(source, destination, visited)

	if path[0].Type == UPNODE_AN {
		path = path[1:]
	}
	upi.DefaultUserPlanePath[dnn] = path
	return
}

func getPathBetween(cur *UPNode, dest *UPNode, visited map[*UPNode]bool) (path []*UPNode, pathExist bool) {

	visited[cur] = true

	if reflect.DeepEqual(*cur, *dest) {

		path = make([]*UPNode, 0)
		path = append(path, cur)
		pathExist = true
		return
	}

	for _, nodes := range cur.Links {

		if !visited[nodes] {
			path_tail, path_exist := getPathBetween(nodes, dest, visited)

			if path_exist {
				path = make([]*UPNode, 0)
				path = append(path, cur)

				path = append(path, path_tail...)
				pathExist = true

				return
			}
		}
	}

	return nil, false

}
