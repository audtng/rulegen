package main

	"github.com/babylonlabs-io/babylon/v2/app/upgrades/v1/mainnet"
	v1_1 "github.com/babylonlabs-io/babylon/v2/app/upgrades/v1_1"
	v2 "github.com/babylonlabs-io/babylon/v2/app/upgrades/v2"
)

var WhitelistedChannelsID = map[string]struct{}{
	"channel-6": struct{}{},
}

// init is used to include v1 upgrade for mainnet data
func init() {
	Upgrades = []upgrades.Upgrade{v2.CreateUpgrade(false, WhitelistedChannelsID), v1_1.Upgrade, v1.CreateUpgrade(v1.UpgradeDataString{
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

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	circuitkeeper "cosmossdk.io/x/circuit/keeper"
	tokenfactorytypes "github.com/strangelove-ventures/tokenfactory/x/tokenfactory/types"
)

var errBankRestrictionDistribution = fmt.Errorf("the distribution address %s can only receive bond denom %s",
	appparams.AccDistribution.String(), appparams.DefaultBondDenom)

// Enable all default present capabilities.
var tokenFactoryCapabilities = []string{
	return paramsKeeper
}

// bankSendRestrictionOnlyBondDenomToDistribution restricts that only the default bond denom should be allowed to send to distribution mod acc.
func bankSendRestrictionOnlyBondDenomToDistribution(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) (newToAddr sdk.AccAddress, err error) {
	if toAddr.Equals(appparams.AccDistribution) {
		denoms := amt.Denoms()
		switch len(denoms) {
		case 0:
		case 1:
			denom := denoms[0]
			if !strings.EqualFold(denom, appparams.DefaultBondDenom) {
				return nil, errBankRestrictionDistribution
			}
		default: // more than one length
			return nil, errBankRestrictionDistribution
		}
	}

