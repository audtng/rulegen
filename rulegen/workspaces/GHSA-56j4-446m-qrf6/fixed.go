package main

package v2_2

import (
	"context"

	store "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/babylonlabs-io/babylon/v2/app/keepers"
	"github.com/babylonlabs-io/babylon/v2/app/upgrades"
)

// UpgradeName defines the on-chain upgrade name for the Babylon v2.2 upgrade
const UpgradeName = "v2.2"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added:   []string{},
		Deleted: []string{},
	},
}

func CreateUpgradeHandler(mm *module.Manager, configurator module.Configurator, keepers *keepers.AppKeepers) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
	"github.com/babylonlabs-io/babylon/v2/app/upgrades/v1/mainnet"
	v1_1 "github.com/babylonlabs-io/babylon/v2/app/upgrades/v1_1"
	v2 "github.com/babylonlabs-io/babylon/v2/app/upgrades/v2"
	v22 "github.com/babylonlabs-io/babylon/v2/app/upgrades/v2_2"
)

var WhitelistedChannelsID = map[string]struct{}{
	"channel-6": struct{}{},
}

// init is used to include v2.2 upgrade for mainnet data
func init() {
	Upgrades = []upgrades.Upgrade{
		v22.Upgrade,
		v2.CreateUpgrade(false, WhitelistedChannelsID),
		v1_1.Upgrade,
		v1.CreateUpgrade(v1.UpgradeDataString{
			BtcStakingParamsStr:       mainnet.BtcStakingParamsStr,
			FinalityParamStr:          mainnet.FinalityParamStr,
			IncentiveParamStr:         mainnet.IncentiveParamStr,
			CosmWasmParamStr:          mainnet.CosmWasmParamStr,
			NewBtcHeadersStr:          mainnet.NewBtcHeadersStr,
			TokensDistributionStr:     mainnet.TokensDistributionStr,
			AllowedStakingTxHashesStr: mainnet.AllowedStakingTxHashesStr,
		}, mainnet.ParamUpgrade)}
}
	"path/filepath"
	"strings"

	"cosmossdk.io/errors"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	circuitkeeper "cosmossdk.io/x/circuit/keeper"
	tokenfactorytypes "github.com/strangelove-ventures/tokenfactory/x/tokenfactory/types"
)

var errBankRestriction = fmt.Errorf("can only receive bond denom %s", appparams.DefaultBondDenom)

// Enable all default present capabilities.
var tokenFactoryCapabilities = []string{
	return paramsKeeper
}

// bankSendRestrictionOnlyBondDenomToDistribution restricts that only the default bond denom should be allowed to send to distribution and fee collector mod accs.
func bankSendRestrictionOnlyBondDenomToDistribution(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) (newToAddr sdk.AccAddress, err error) {
	if toAddr.Equals(appparams.AccDistribution) || toAddr.Equals(appparams.AccFeeCollector) {
		denoms := amt.Denoms()
		switch len(denoms) {
		case 0:
		case 1:
			denom := denoms[0]
			if !strings.EqualFold(denom, appparams.DefaultBondDenom) {
				return nil, errors.Wrapf(errBankRestriction, "address %s", toAddr)
			}
		default: // more than one length
			return nil, errors.Wrapf(errBankRestriction, "address %s", toAddr)
		}
	}

