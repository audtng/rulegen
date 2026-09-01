package main

/*
Copyright IBM Corp, SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package filter

import (
	"context"
	"time"

	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/core/handlers/auth"
	"github.com/hyperledger/fabric/protoutil"
	"github.com/pkg/errors"
)

// NewTimeWindowCheckFilter creates a new Filter that checks timewindow expiration
func NewTimeWindowCheckFilter(timeWindow time.Duration) auth.Filter {
	return &timewindowCheckFilter{
		timeWindow: timeWindow,
	}
}

type timewindowCheckFilter struct {
	next       peer.EndorserServer
	timeWindow time.Duration
}

// Init initializes the Filter with the next EndorserServer
func (f *timewindowCheckFilter) Init(next peer.EndorserServer) {
	f.next = next
}

func validateTimewindowProposal(signedProp *peer.SignedProposal, timeWindow time.Duration) error {
	prop, err := protoutil.UnmarshalProposal(signedProp.ProposalBytes)
	if err != nil {
		return errors.Wrap(err, "failed parsing proposal")
	}

	hdr, err := protoutil.UnmarshalHeader(prop.Header)
	if err != nil {
		return errors.Wrap(err, "failed parsing header")
	}

	chdr, err := protoutil.UnmarshalChannelHeader(hdr.ChannelHeader)
	if err != nil {
		return errors.Wrap(err, "failed parsing channel header")
	}

	timeProposal := chdr.Timestamp.AsTime().UTC()
	now := time.Now().UTC()

	if timeProposal.Add(timeWindow).Before(now) || timeProposal.Add(-timeWindow).After(now) {
		return errors.Errorf("request unauthorized due to incorrect timestamp %s, peer time %s, peer.authentication.timewindow %s",
			timeProposal.Format(time.RFC3339), now.Format(time.RFC3339), timeWindow.String())
	}

	return nil
}

// ProcessProposal processes a signed proposal
func (f *timewindowCheckFilter) ProcessProposal(ctx context.Context, signedProp *peer.SignedProposal) (*peer.ProposalResponse, error) {
	if err := validateTimewindowProposal(signedProp, f.timeWindow); err != nil {
		return nil, err
	}
	return f.next.ProcessProposal(ctx, signedProp)
}
		discprotos.RegisterDiscoveryServer(peerServer.Server(), discoveryService)
	}

	logger.Infof("Starting peer with ID=[%s], network ID=[%s], address=[%s]", coreConfig.PeerID, coreConfig.NetworkID, coreConfig.PeerAddress)

	// Get configuration before starting go routines to avoid
	// Register the Endorser server
	pb.RegisterEndorserServer(peerServer.Server(), auth)

	if coreConfig.GatewayOptions.Enabled {
		if coreConfig.DiscoveryEnabled {
			logger.Info("Starting peer with Gateway enabled")

			gatewayServer := gateway.CreateServer(
				auth,
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

	// register the snapshot server
	snapshotSvc := &snapshotgrpc.SnapshotService{LedgerGetter: peerInstance, ACLProvider: aclProvider}
	pb.RegisterSnapshotServer(peerServer.Server(), snapshotSvc)
package library

import (
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)
// Config configures the factory methods
// and plugins for the registry
type Config struct {
	AuthFilters              []*HandlerConfig `yaml:"authFilters"`
	Decorators               []*HandlerConfig `yaml:"decorators"`
	Endorsers                PluginMapping    `yaml:"endorsers"`
	Validators               PluginMapping    `yaml:"validators"`
	AuthenticationTimeWindow time.Duration    `yaml:"authenticationTimeWindow"`
}

// PluginMapping stores a map between chaincode id to plugin config
		validators[k] = &HandlerConfig{Name: name, Library: library}
	}

	authenticationTimeWindow := viper.GetDuration("peer.authentication.timewindow")
	if authenticationTimeWindow == 0 {
		defaultTimeWindow := 15 * time.Minute
		logger.Warningf("`peer.authentication.timewindow` not set; defaulting to %s", defaultTimeWindow)
		authenticationTimeWindow = defaultTimeWindow
	}

	return Config{
		AuthFilters:              authFilters,
		Decorators:               decorators,
		Endorsers:                endorsers,
		Validators:               validators,
		AuthenticationTimeWindow: authenticationTimeWindow,
	}, nil
}
