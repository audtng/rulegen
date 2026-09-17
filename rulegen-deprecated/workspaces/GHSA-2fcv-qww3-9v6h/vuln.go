package main

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/mempool"

	bApp       *baseapp.BaseApp

	// used for building and parsing the injected tx
	txConfig client.TxConfig
	mp       mempool.Mempool

	defaultPrepareProposalHandler sdk.PrepareProposalHandler
	defaultProcessProposalHandler sdk.ProcessProposalHandler
		ckptKeeper:                    ckptKeeper,
		bApp:                          bApp,
		txConfig:                      encCfg.TxConfig,
		defaultPrepareProposalHandler: defaultHandler.PrepareProposalHandler(),
		defaultProcessProposalHandler: defaultHandler.ProcessProposalHandler(),
	}
	}

	var ve ckpttypes.VoteExtension
	if err := ve.Unmarshal(veBytes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal vote extension: %w", err)
	}
