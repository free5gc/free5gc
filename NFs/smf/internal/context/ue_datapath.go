package context

import (
	"fmt"
	"math"

	"github.com/free5gc/smf/internal/logger"
	"github.com/free5gc/smf/pkg/factory"
	"github.com/free5gc/util/idgenerator"
)

type UEPreConfigPaths struct {
	DataPathPool    DataPathPool
	PathIDGenerator *idgenerator.IDGenerator
}

func NewUEDataPathNode(name string) (node *DataPathNode, err error) {
	upNodes := smfContext.UserPlaneInformation.UPNodes

	if _, exist := upNodes[name]; !exist {
		err = fmt.Errorf("UPNode %s isn't exist in smfcfg.yaml, but in UERouting.yaml", name)
		return nil, err
	}

	node = &DataPathNode{
		UPF:            upNodes[name].UPF,
		UpLinkTunnel:   &GTPTunnel{},
		DownLinkTunnel: &GTPTunnel{},
	}
	return
}

func NewUEPreConfigPaths(paths []factory.SpecificPath) (*UEPreConfigPaths, error) {
	var uePreConfigPaths *UEPreConfigPaths
	ueDataPathPool := NewDataPathPool()
	lowerBound := 0
	pathIDGenerator := idgenerator.NewGenerator(1, math.MaxInt32)

	logger.PduSessLog.Infoln("In NewUEPreConfigPaths")

	for _, path := range paths {
		dataPath := NewDataPath()

		var pathID int64
		if allocPathID, err := pathIDGenerator.Allocate(); err != nil {
			logger.CtxLog.Warnf("Allocate pathID error: %+v", err)
			return nil, err
		} else {
			pathID = allocPathID
		}

		dataPath.Destination.DestinationIP = path.DestinationIP
		dataPath.Destination.DestinationPort = path.DestinationPort
		ueDataPathPool[pathID] = dataPath
		var parentNode *DataPathNode = nil
		for idx, nodeName := range path.Path {
			newUeNode, err := NewUEDataPathNode(nodeName)
			if err != nil {
				return nil, err
			}

			if idx == lowerBound {
				dataPath.FirstDPNode = newUeNode
			}
			if parentNode != nil {
				newUeNode.AddPrev(parentNode)
				parentNode.AddNext(newUeNode)
			}
			parentNode = newUeNode
		}

		logger.CtxLog.Traceln("New data path added")
		logger.CtxLog.Traceln("\n" + dataPath.String() + "\n")
	}

	uePreConfigPaths = &UEPreConfigPaths{
		DataPathPool:    ueDataPathPool,
		PathIDGenerator: pathIDGenerator,
	}
	return uePreConfigPaths, nil
}

func GetUEPreConfigPaths(supi string, upfName string) *UEPreConfigPaths {
	groupName := GetULCLGroupNameFromSUPI(supi)
	if groupName == "" {
		return nil
	}
	dataPathPool := NewDataPathPool()
	dataPathPool[1] = smfContext.UEDefaultPathPool[groupName].GetDefaultPath(upfName)
	var i int64 = 2
	for _, dataPath := range smfContext.UEPreConfigPathPool[groupName].DataPathPool {
		firstNode := dataPath.CopyFirstDPNode()
		path := &DataPath{
			Activated:     false,
			IsDefaultPath: false,
			Destination:   dataPath.Destination,
			FirstDPNode:   firstNode,
		}
		dataPathPool[i] = path
		i++
	}
	paths := &UEPreConfigPaths{
		DataPathPool:    dataPathPool,
		PathIDGenerator: smfContext.UEPreConfigPathPool[groupName].PathIDGenerator,
	}
	return paths
}

func CheckUEHasPreConfig(supi string) (exist bool) {
	groupName := GetULCLGroupNameFromSUPI(supi)
	logger.CtxLog.Tracef("UE [%s] belongs to group [%s]", supi, groupName)
	if groupName == "" {
		return false
	}
	_, exist = smfContext.UEPreConfigPathPool[groupName]
	return
}
