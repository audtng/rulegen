package main

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/codec/unknownproto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/mempool"

	bApp       *baseapp.BaseApp

	// used for building and parsing the injected tx
	txConfig          client.TxConfig
	interfaceRegistry codectypes.InterfaceRegistry
	mp                mempool.Mempool

	defaultPrepareProposalHandler sdk.PrepareProposalHandler
	defaultProcessProposalHandler sdk.ProcessProposalHandler
		ckptKeeper:                    ckptKeeper,
		bApp:                          bApp,
		txConfig:                      encCfg.TxConfig,
		interfaceRegistry:             encCfg.InterfaceRegistry,
		defaultPrepareProposalHandler: defaultHandler.PrepareProposalHandler(),
		defaultProcessProposalHandler: defaultHandler.ProcessProposalHandler(),
	}
	}

	var ve ckpttypes.VoteExtension
	if err := unknownproto.RejectUnknownFieldsStrict(veBytes, &ve, h.interfaceRegistry); err != nil {
		return nil, fmt.Errorf("vote extension contains unknown or extra bytes: %w", err)
	}

	if err := ve.Unmarshal(veBytes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal vote extension: %w", err)
	}
