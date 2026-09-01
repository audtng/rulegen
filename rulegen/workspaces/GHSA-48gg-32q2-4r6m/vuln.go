package main

		discprotos.RegisterDiscoveryServer(peerServer.Server(), discoveryService)
	}

	if coreConfig.GatewayOptions.Enabled {
		if coreConfig.DiscoveryEnabled {
			logger.Info("Starting peer with Gateway enabled")

			gatewayServer := gateway.CreateServer(
				serverEndorser,
				discoveryService,
				peerInstance,
				&serverConfig.SecOpts,
				aclProvider,
				coreConfig.LocalMSPID,
				coreConfig.GatewayOptions,
				builtinSCCs,
			)
			gatewayprotos.RegisterGatewayServer(peerServer.Server(), gatewayServer)
		} else {
			logger.Warning("Discovery service must be enabled for embedded gateway")
		}
	}

	logger.Infof("Starting peer with ID=[%s], network ID=[%s], address=[%s]", coreConfig.PeerID, coreConfig.NetworkID, coreConfig.PeerAddress)

	// Get configuration before starting go routines to avoid
	// Register the Endorser server
	pb.RegisterEndorserServer(peerServer.Server(), auth)

	// register the snapshot server
	snapshotSvc := &snapshotgrpc.SnapshotService{LedgerGetter: peerInstance, ACLProvider: aclProvider}
	pb.RegisterSnapshotServer(peerServer.Server(), snapshotSvc)
package library

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)
// Config configures the factory methods
// and plugins for the registry
type Config struct {
	AuthFilters []*HandlerConfig `yaml:"authFilters"`
	Decorators  []*HandlerConfig `yaml:"decorators"`
	Endorsers   PluginMapping    `yaml:"endorsers"`
	Validators  PluginMapping    `yaml:"validators"`
}

// PluginMapping stores a map between chaincode id to plugin config
		validators[k] = &HandlerConfig{Name: name, Library: library}
	}

	return Config{
		AuthFilters: authFilters,
		Decorators:  decorators,
		Endorsers:   endorsers,
		Validators:  validators,
	}, nil
}
