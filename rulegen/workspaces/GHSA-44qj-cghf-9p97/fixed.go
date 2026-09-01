package main

}

func (upi *UserPlaneInformation) UpNodesFromConfiguration(upTopology *factory.UserPlaneInformation) error {
	candidateUPNodes := make(map[string]*UPNode, len(upi.UPNodes))
	for name, node := range upi.UPNodes {
		candidateUPNodes[name] = node
	}

	candidateUPFs := make(map[string]*UPNode, len(upi.UPFs))
	for name, node := range upi.UPFs {
		candidateUPFs[name] = node
	}

	candidateAccessNetwork := make(map[string]*UPNode, len(upi.AccessNetwork))
	for name, node := range upi.AccessNetwork {
		candidateAccessNetwork[name] = node
	}

	candidateUPFsID := make(map[string]string, len(upi.UPFsID))
	for name, id := range upi.UPFsID {
		candidateUPFsID[name] = id
	}

	candidateUPFsIPtoID := make(map[string]string, len(upi.UPFsIPtoID))
	for ip, id := range upi.UPFsIPtoID {
		candidateUPFsIPtoID[ip] = id
	}

	candidateUPFIPToName := make(map[string]string, len(upi.UPFIPToName))
	for ip, name := range upi.UPFIPToName {
		candidateUPFIPToName[ip] = name
	}

	createdUPFs := make([]*UPF, 0)
	cleanupCreatedUPFs := func() {
		for _, upf := range createdUPFs {
			upfPool.Delete(upf.UUID())
		}
	}

	for name, node := range upTopology.UPNodes {
		if _, ok := candidateUPNodes[name]; ok {
			logger.InitLog.Warningf("Node [%s] already exists in SMF.\n", name)
			continue
		}
		upNode := new(UPNode)
		upNode.Name = name
		upNode.Type = UPNodeType(node.Type)
		switch upNode.Type {
		case UPNODE_UPF:
			}

			upNode.UPF = NewUPF(&upNode.NodeID, node.InterfaceUpfInfoList)
			createdUPFs = append(createdUPFs, upNode.UPF)
			snssaiInfos := make([]*SnssaiUPFInfo, 0)
			for _, snssaiInfoConfig := range node.SNssaiInfos {
				snssaiInfo := &SnssaiUPFInfo{
					for _, pool := range dnnInfoConfig.Pools {
						ueIPPool := NewUEIPPool(pool)
						if ueIPPool == nil {
							cleanupCreatedUPFs()
							return fmt.Errorf("invalid pools value: %+v", pool)
						} else {
							ueIPPools = append(ueIPPools, ueIPPool)
					for _, pool := range dnnInfoConfig.StaticPools {
						ueIPPool := NewUEIPPool(pool)
						if ueIPPool == nil {
							cleanupCreatedUPFs()
							return fmt.Errorf("invalid pools value: %+v", pool)
						} else {
							staticUeIPPools = append(staticUeIPPools, ueIPPool)
							for _, dynamicUePool := range ueIPPools {
								if dynamicUePool.ueSubNet.Contains(ueIPPool.ueSubNet.IP) {
									if err := dynamicUePool.Exclude(ueIPPool); err != nil {
										cleanupCreatedUPFs()
										return fmt.Errorf("exclude static Pool[%s] failed: %v",
											ueIPPool.ueSubNet, err)
									}
				snssaiInfos = append(snssaiInfos, snssaiInfo)
			}
			upNode.UPF.SNssaiInfos = snssaiInfos
			candidateUPFs[name] = upNode

			// AllocateUPFID
			upfid := upNode.UPF.UUID()
			upfip := upNode.NodeID.ResolveNodeIdToIp().String()
			candidateUPFsID[name] = upfid
			candidateUPFsIPtoID[upfip] = upfid

		case UPNODE_AN:
			upNode.ANIP = net.ParseIP(node.ANIP)
			candidateAccessNetwork[name] = upNode
		default:
			logger.InitLog.Warningf("invalid UPNodeType: %s\n", upNode.Type)
		}

		candidateUPNodes[name] = upNode

		ipStr := upNode.NodeID.ResolveNodeIdToIp().String()
		candidateUPFIPToName[ipStr] = name
	}

	// overlap UE IP pool validation
	allUEIPPools := []*UeIPPool{}
	for _, upf := range candidateUPFs {
		for _, snssaiInfo := range upf.UPF.SNssaiInfos {
			for _, dnn := range snssaiInfo.DnnList {
				allUEIPPools = append(allUEIPPools, dnn.UeIPPools...)
		}
	}
	if isOverlap(allUEIPPools) {
		cleanupCreatedUPFs()
		return fmt.Errorf("overlap cidr value between UPFs")
	}

	upi.UPNodes = candidateUPNodes
	upi.UPFs = candidateUPFs
	upi.AccessNetwork = candidateAccessNetwork
	upi.UPFsID = candidateUPFsID
	upi.UPFsIPtoID = candidateUPFsIPtoID
	upi.UPFIPToName = candidateUPFIPToName

	return nil
}


import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"reflect"
}

// NewUserPlaneInformation process the configuration then returns a new instance of UserPlaneInformation
func NewUserPlaneInformation(upTopology *factory.UserPlaneInformation) (*UserPlaneInformation, error) {
	nodePool := make(map[string]*UPNode)
	upfPool := make(map[string]*UPNode)
	anPool := make(map[string]*UPNode)
					for _, pool := range dnnInfoConfig.Pools {
						ueIPPool := NewUEIPPool(pool)
						if ueIPPool == nil {
							return nil, fmt.Errorf("invalid pools value: %+v", pool)
						} else {
							ueIPPools = append(ueIPPools, ueIPPool)
							allUEIPPools = append(allUEIPPools, ueIPPool)
					for _, staticPool := range dnnInfoConfig.StaticPools {
						staticUeIPPool := NewUEIPPool(staticPool)
						if staticUeIPPool == nil {
							return nil, fmt.Errorf("invalid pools value: %+v", staticPool)
						} else {
							staticUeIPPools = append(staticUeIPPools, staticUeIPPool)
							for _, dynamicUePool := range ueIPPools {
								if dynamicUePool.ueSubNet.Contains(staticUeIPPool.ueSubNet.IP) {
									if err := dynamicUePool.Exclude(staticUeIPPool); err != nil {
										return nil, fmt.Errorf("exclude static Pool[%s] failed: %v",
											staticUeIPPool.ueSubNet, err)
									}
								}
	}

	if isOverlap(allUEIPPools) {
		return nil, fmt.Errorf("overlap cidr value between UPFs")
	}

	for _, link := range upTopology.Links {
		DefaultUserPlanePathToUPF: make(map[string]map[string][]*UPNode),
	}

	return userplaneInformation, nil
}

func (upi *UserPlaneInformation) UpNodesToConfiguration() map[string]*factory.UPNode {
	return links
}

func (upi *UserPlaneInformation) UpNodesFromConfiguration(upTopology *factory.UserPlaneInformation) error {
	for name, node := range upTopology.UPNodes {
		if _, ok := upi.UPNodes[name]; ok {
			logger.InitLog.Warningf("Node [%s] already exists in SMF.\n", name)
					for _, pool := range dnnInfoConfig.Pools {
						ueIPPool := NewUEIPPool(pool)
						if ueIPPool == nil {
							return fmt.Errorf("invalid pools value: %+v", pool)
						} else {
							ueIPPools = append(ueIPPools, ueIPPool)
						}
					for _, pool := range dnnInfoConfig.StaticPools {
						ueIPPool := NewUEIPPool(pool)
						if ueIPPool == nil {
							return fmt.Errorf("invalid pools value: %+v", pool)
						} else {
							staticUeIPPools = append(staticUeIPPools, ueIPPool)
							for _, dynamicUePool := range ueIPPools {
								if dynamicUePool.ueSubNet.Contains(ueIPPool.ueSubNet.IP) {
									if err := dynamicUePool.Exclude(ueIPPool); err != nil {
										return fmt.Errorf("exclude static Pool[%s] failed: %v",
											ueIPPool.ueSubNet, err)
									}
								}
		}
	}
	if isOverlap(allUEIPPools) {
		return fmt.Errorf("overlap cidr value between UPFs")
	}

	return nil
}

func (upi *UserPlaneInformation) LinksFromConfiguration(upTopology *factory.UserPlaneInformation) {
	return atomic.AddUint64(&smfContext.LocalSEIDCount, 1)
}

func InitSmfContext(config *factory.Config) error {
	if config == nil {
		return fmt.Errorf("config is nil")
	}

	logger.CtxLog.Infof("smfconfig Info: Version[%s] Description[%s]", config.Info.Version, config.Info.Description)

	sbi := configuration.Sbi
	if sbi == nil {
		return fmt.Errorf("configuration needs \"sbi\" value")
	} else {
		smfContext.URIScheme = models.UriScheme(sbi.Scheme)
		smfContext.RegisterIPv4 = factory.SmfSbiDefaultIPv4 // default localhost

	smfContext.SupportedPDUSessionType = "IPv4"

	userPlaneInformation, err := NewUserPlaneInformation(&configuration.UserPlaneInformation)
	if err != nil {
		return fmt.Errorf("initialize user plane information failed: %w", err)
	}
	smfContext.UserPlaneInformation = userPlaneInformation

	smfContext.ChargingIDGenerator = idgenerator.NewGenerator(1, math.MaxUint32)

	TeidGenerator = idgenerator.NewGenerator(1, math.MaxUint32)

	smfContext.Ues = InitSmfUeData()

	return nil
}

func InitSMFUERouting(routingConfig *factory.RoutingConfig) {
internal/context/sm_context_policy_test.go | 4 +++-
internal/pfcp/message/build_test.go        | 4 +++-
internal/sbi/consumer/pcf_service_test.go  | 2 +-
internal/sbi/processor/pdu_session_test.go | 4 +++-
internal/sbi/server.go                     | 4 +++-
5 files changed, 13 insertions(+), 5 deletions(-)
