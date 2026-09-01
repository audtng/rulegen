package main

}

func (upi *UserPlaneInformation) UpNodesFromConfiguration(upTopology *factory.UserPlaneInformation) error {
	for name, node := range upTopology.UPNodes {
		if _, ok := upi.UPNodes[name]; ok {
			logger.InitLog.Warningf("Node [%s] already exists in SMF.\n", name)
			continue
		}
		upNode := new(UPNode)
		upNode.Type = UPNodeType(node.Type)
		switch upNode.Type {
		case UPNODE_UPF:
			}

			upNode.UPF = NewUPF(&upNode.NodeID, node.InterfaceUpfInfoList)
			snssaiInfos := make([]*SnssaiUPFInfo, 0)
			for _, snssaiInfoConfig := range node.SNssaiInfos {
				snssaiInfo := &SnssaiUPFInfo{
					for _, pool := range dnnInfoConfig.Pools {
						ueIPPool := NewUEIPPool(pool)
						if ueIPPool == nil {
							return fmt.Errorf("invalid pools value: %+v", pool)
						} else {
							ueIPPools = append(ueIPPools, ueIPPool)
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
				snssaiInfos = append(snssaiInfos, snssaiInfo)
			}
			upNode.UPF.SNssaiInfos = snssaiInfos
			upi.UPFs[name] = upNode

			// AllocateUPFID
			upfid := upNode.UPF.UUID()
			upfip := upNode.NodeID.ResolveNodeIdToIp().String()
			upi.UPFsID[name] = upfid
			upi.UPFsIPtoID[upfip] = upfid

		case UPNODE_AN:
			upNode.ANIP = net.ParseIP(node.ANIP)
			upi.AccessNetwork[name] = upNode
		default:
			logger.InitLog.Warningf("invalid UPNodeType: %s\n", upNode.Type)
		}

		upi.UPNodes[name] = upNode

		ipStr := upNode.NodeID.ResolveNodeIdToIp().String()
		upi.UPFIPToName[ipStr] = name
	}

	// overlap UE IP pool validation
	allUEIPPools := []*UeIPPool{}
	for _, upf := range upi.UPFs {
		for _, snssaiInfo := range upf.UPF.SNssaiInfos {
			for _, dnn := range snssaiInfo.DnnList {
				allUEIPPools = append(allUEIPPools, dnn.UeIPPools...)
		}
	}
	if isOverlap(allUEIPPools) {
		return fmt.Errorf("overlap cidr value between UPFs")
	}

	return nil
}


import (
	"errors"
	"math/rand"
	"net"
	"reflect"
}

// NewUserPlaneInformation process the configuration then returns a new instance of UserPlaneInformation
func NewUserPlaneInformation(upTopology *factory.UserPlaneInformation) *UserPlaneInformation {
	nodePool := make(map[string]*UPNode)
	upfPool := make(map[string]*UPNode)
	anPool := make(map[string]*UPNode)
					for _, pool := range dnnInfoConfig.Pools {
						ueIPPool := NewUEIPPool(pool)
						if ueIPPool == nil {
							logger.InitLog.Fatalf("invalid pools value: %+v", pool)
						} else {
							ueIPPools = append(ueIPPools, ueIPPool)
							allUEIPPools = append(allUEIPPools, ueIPPool)
					for _, staticPool := range dnnInfoConfig.StaticPools {
						staticUeIPPool := NewUEIPPool(staticPool)
						if staticUeIPPool == nil {
							logger.InitLog.Fatalf("invalid pools value: %+v", staticPool)
						} else {
							staticUeIPPools = append(staticUeIPPools, staticUeIPPool)
							for _, dynamicUePool := range ueIPPools {
								if dynamicUePool.ueSubNet.Contains(staticUeIPPool.ueSubNet.IP) {
									if err := dynamicUePool.Exclude(staticUeIPPool); err != nil {
										logger.InitLog.Fatalf("exclude static Pool[%s] failed: %v",
											staticUeIPPool.ueSubNet, err)
									}
								}
	}

	if isOverlap(allUEIPPools) {
		logger.InitLog.Fatalf("overlap cidr value between UPFs")
	}

	for _, link := range upTopology.Links {
		DefaultUserPlanePathToUPF: make(map[string]map[string][]*UPNode),
	}

	return userplaneInformation
}

func (upi *UserPlaneInformation) UpNodesToConfiguration() map[string]*factory.UPNode {
	return links
}

func (upi *UserPlaneInformation) UpNodesFromConfiguration(upTopology *factory.UserPlaneInformation) {
	for name, node := range upTopology.UPNodes {
		if _, ok := upi.UPNodes[name]; ok {
			logger.InitLog.Warningf("Node [%s] already exists in SMF.\n", name)
					for _, pool := range dnnInfoConfig.Pools {
						ueIPPool := NewUEIPPool(pool)
						if ueIPPool == nil {
							logger.InitLog.Fatalf("invalid pools value: %+v", pool)
						} else {
							ueIPPools = append(ueIPPools, ueIPPool)
						}
					for _, pool := range dnnInfoConfig.StaticPools {
						ueIPPool := NewUEIPPool(pool)
						if ueIPPool == nil {
							logger.InitLog.Fatalf("invalid pools value: %+v", pool)
						} else {
							staticUeIPPools = append(staticUeIPPools, ueIPPool)
							for _, dynamicUePool := range ueIPPools {
								if dynamicUePool.ueSubNet.Contains(ueIPPool.ueSubNet.IP) {
									if err := dynamicUePool.Exclude(ueIPPool); err != nil {
										logger.InitLog.Fatalf("exclude static Pool[%s] failed: %v",
											ueIPPool.ueSubNet, err)
									}
								}
		}
	}
	if isOverlap(allUEIPPools) {
		logger.InitLog.Fatalf("overlap cidr value between UPFs")
	}
}

func (upi *UserPlaneInformation) LinksFromConfiguration(upTopology *factory.UserPlaneInformation) {
	return atomic.AddUint64(&smfContext.LocalSEIDCount, 1)
}

func InitSmfContext(config *factory.Config) {
	if config == nil {
		logger.CtxLog.Error("Config is nil")
		return
	}

	logger.CtxLog.Infof("smfconfig Info: Version[%s] Description[%s]", config.Info.Version, config.Info.Description)

	sbi := configuration.Sbi
	if sbi == nil {
		logger.CtxLog.Errorln("Configuration needs \"sbi\" value")
		return
	} else {
		smfContext.URIScheme = models.UriScheme(sbi.Scheme)
		smfContext.RegisterIPv4 = factory.SmfSbiDefaultIPv4 // default localhost

	smfContext.SupportedPDUSessionType = "IPv4"

	smfContext.UserPlaneInformation = NewUserPlaneInformation(&configuration.UserPlaneInformation)

	smfContext.ChargingIDGenerator = idgenerator.NewGenerator(1, math.MaxUint32)

	TeidGenerator = idgenerator.NewGenerator(1, math.MaxUint32)

	smfContext.Ues = InitSmfUeData()
}

func InitSMFUERouting(routingConfig *factory.RoutingConfig) {
internal/context/sm_context_policy_test.go | 4 +++-
internal/pfcp/message/build_test.go        | 4 +++-
internal/sbi/consumer/pcf_service_test.go  | 2 +-
internal/sbi/processor/pdu_session_test.go | 4 +++-
internal/sbi/server.go                     | 4 +++-
5 files changed, 13 insertions(+), 5 deletions(-)
