package main

package v4_2

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/math"
	store "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/babylonlabs-io/babylon/v4/app/keepers"
	"github.com/babylonlabs-io/babylon/v4/app/upgrades"
	bbn "github.com/babylonlabs-io/babylon/v4/types"
	btcstkkeeper "github.com/babylonlabs-io/babylon/v4/x/btcstaking/keeper"
	btcstktypes "github.com/babylonlabs-io/babylon/v4/x/btcstaking/types"
	costkkeeper "github.com/babylonlabs-io/babylon/v4/x/costaking/keeper"
	costktypes "github.com/babylonlabs-io/babylon/v4/x/costaking/types"
	fkeeper "github.com/babylonlabs-io/babylon/v4/x/finality/keeper"
	ftypes "github.com/babylonlabs-io/babylon/v4/x/finality/types"
)

const UpgradeName = "v4.2"

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
		// Run migrations before applying any other state changes.
		migrations, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return nil, err
		}

		costkStoreKey := keepers.GetKey(costktypes.StoreKey)
		if costkStoreKey == nil {
			return nil, errors.New("invalid costaking types store key")
		}
		coStkStoreService := runtime.NewKVStoreService(costkStoreKey)

		// Reset co-staker rewards tracker
		if err := ResetCoStakerRwdsTracker(
			ctx,
			keepers.EncCfg.Codec,
			coStkStoreService,
			keepers.BTCStakingKeeper,
			keepers.CostakingKeeper,
			keepers.FinalityKeeper,
		); err != nil {
			return nil, err
		}

		return migrations, nil
	}
}

// ResetCoStakerRwdsTracker resets the costaker rewards tracker
// It resets tracked ActiveSats for all BTC stakers
func ResetCoStakerRwdsTracker(
	ctx context.Context,
	cdc codec.BinaryCodec,
	costkStoreService corestoretypes.KVStoreService,
	btcStkKeeper btcstkkeeper.Keeper,
	coStkKeeper costkkeeper.Keeper,
	fKeeper fkeeper.Keeper,
) error {
	sb := collections.NewSchemaBuilder(costkStoreService)
	rwdTrackers := collections.NewMap(
		sb,
		costktypes.CostakerRewardsTrackerKeyPrefix,
		"costaker_rewards_tracker",
		collections.BytesKey,
		codec.CollValue[costktypes.CostakerRewardsTracker](cdc),
	)

	// Zero out tracked amounts in existing rewards trackers
	accsWithActiveSats, err := zeroOutCoStakerRwdsActiveSats(ctx, rwdTrackers)
	if err != nil {
		return err
	}

	params := coStkKeeper.GetParams(ctx)
	endedPeriod, err := coStkKeeper.IncrementRewardsPeriod(ctx)
	if err != nil {
		return err
	}

	// Update co-staker rwd tracker for all BTC stakers
	if err := updateBTCStakersRwdTracker(ctx, endedPeriod, rwdTrackers, accsWithActiveSats, btcStkKeeper, fKeeper, coStkKeeper, params); err != nil {
		return err
	}

	totalScore, err := getTotalScore(ctx, rwdTrackers)
	if err != nil {
		return err
	}

	currentRwd, err := coStkKeeper.GetCurrentRewards(ctx)
	if err != nil {
		return err
	}

	currentRwd.TotalScore = totalScore
	if err := currentRwd.Validate(); err != nil {
		return err
	}

	return coStkKeeper.SetCurrentRewards(ctx, *currentRwd)
}

type ActiveSatsTracked struct {
	PreviousActiveSats math.Int
	CurrentActiveSats  math.Int
}

// updateBTCStakersRwdTracker retrieves all active BTC stakers with pagination and updates their costaker rewards trackers (ActiveSatoshis only)
func updateBTCStakersRwdTracker(
	ctx context.Context,
	period uint64,
	rwdTrackers collections.Map[[]byte, costktypes.CostakerRewardsTracker],
	accsWithActiveSats map[string]ActiveSatsTracked,
	btcStkKeeper btcstkkeeper.Keeper,
	fKeeper fkeeper.Keeper,
	coStkKeeper costkkeeper.Keeper,
	params costktypes.Params,
) error {
	// To count as btc staker for the co-staking rewards
	// need to be delegating to a FP within the current active set
	// This runs on preblocker (before BeginBlock), so the active set to consider should be from previous height
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	height := uint64(sdkCtx.HeaderInfo().Height)
	vp := fKeeper.GetVotingPowerDistCache(ctx, height-1)
	if vp == nil {
		vp = ftypes.NewVotingPowerDistCache()
	}
	activeFps := vp.GetActiveFinalityProviderSet()

	var nextKey []byte

	for {
		req := &btcstktypes.QueryBTCDelegationsRequest{
			Status: btcstktypes.BTCDelegationStatus_ACTIVE,
			Pagination: &query.PageRequest{
				Key: nextKey,
			},
		}

		btcDelRes, err := btcStkKeeper.BTCDelegations(ctx, req)
		if err != nil {
			return err
		}

		for _, del := range btcDelRes.BtcDelegations {
			// check if delegating to an active FP
			if !delegatingToActiveFP(del.FpBtcPkList, activeFps) {
				continue
			}
			// add all current active sats to the memory cache for later checking which have a diff
			delSat := math.NewIntFromUint64(del.TotalSat)
			data, found := accsWithActiveSats[del.StakerAddr]
			if found {
				data.CurrentActiveSats = data.CurrentActiveSats.Add(delSat)
				accsWithActiveSats[del.StakerAddr] = data
			} else {
				accsWithActiveSats[del.StakerAddr] = ActiveSatsTracked{
					PreviousActiveSats: math.ZeroInt(),
					CurrentActiveSats:  delSat,
				}
			}
		}

		if btcDelRes.Pagination == nil || len(btcDelRes.Pagination.NextKey) == 0 {
			break
		}
		nextKey = btcDelRes.Pagination.NextKey
	}

	// now update the costaker rewards trackers (all of them because we zeroed them out before)
	for accAddrStr, satsTracked := range accsWithActiveSats {
		diff := satsTracked.CurrentActiveSats.Sub(satsTracked.PreviousActiveSats)
		needsCorrection := !diff.IsZero()
		if err := updateCostakerActiveSatsRewardsTracker(ctx, coStkKeeper, period, rwdTrackers, sdk.MustAccAddressFromBech32(accAddrStr), satsTracked.CurrentActiveSats, params, needsCorrection); err != nil {
			return err
		}
	}

	return nil
}

// updateCostakerActiveSatsRewardsTracker creates or updates a costaker rewards tracker
func updateCostakerActiveSatsRewardsTracker(
	ctx context.Context,
	coStkKeeper costkkeeper.Keeper,
	endedPeriod uint64,
	rwdTrackers collections.Map[[]byte, costktypes.CostakerRewardsTracker],
	stakerAddr sdk.AccAddress,
	btcAmount math.Int,
	params costktypes.Params,
	needsCorrection bool,
) error {
	addrKey := []byte(stakerAddr)

	// Try to get existing tracker
	rt, err := rwdTrackers.Get(ctx, addrKey)
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		return err
	}

	if errors.Is(err, collections.ErrNotFound) {
		// this should not happen as we're updating existing trackers only
		return nil
	}
	// Update existing tracker (need to set the ActiveSatoshis because these were zeroed out before)
	// Update the StartPeriodCumulativeReward only if the ActiveSatoshis is changing
	rt.ActiveSatoshis = rt.ActiveSatoshis.Add(btcAmount)
	if needsCorrection {
		rt.StartPeriodCumulativeReward = endedPeriod
		if err := coStkKeeper.CalculateCostakerRewardsAndSendToGauge(ctx, stakerAddr, endedPeriod); err != nil {
			return err
		}
	}

	// Update score
	rt.UpdateScore(params.ScoreRatioBtcByBaby)

	// Save tracker
	if err := rwdTrackers.Set(ctx, addrKey, rt); err != nil {
		return err
	}

	return nil
}

// zeroOutCoStakerRwdsActiveSats zeros out ActiveSatoshis in all costaker rewards trackers
func zeroOutCoStakerRwdsActiveSats(
	ctx context.Context,
	rwdTrackers collections.Map[[]byte, costktypes.CostakerRewardsTracker],
) (map[string]ActiveSatsTracked, error) {
	accsWithActiveSats := make(map[string]ActiveSatsTracked)
	iter, err := rwdTrackers.Iterate(ctx, nil)
	if err != nil {
		return accsWithActiveSats, err
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		costakerAddr, err := iter.Key()
		if err != nil {
			return accsWithActiveSats, err
		}

		tracker, err := iter.Value()
		if err != nil {
			return accsWithActiveSats, err
		}
		if tracker.ActiveSatoshis.IsZero() {
			continue
		}

		sdkAddr := sdk.AccAddress(costakerAddr)
		accsWithActiveSats[sdkAddr.String()] = ActiveSatsTracked{
			PreviousActiveSats: tracker.ActiveSatoshis,
			CurrentActiveSats:  math.ZeroInt(),
		}

		// Zero out ActiveSatoshis
		tracker.ActiveSatoshis = math.ZeroInt()
		if err := rwdTrackers.Set(ctx, costakerAddr, tracker); err != nil {
			return accsWithActiveSats, err
		}
	}

	return accsWithActiveSats, nil
}

func getTotalScore(
	ctx context.Context,
	rwdTrackers collections.Map[[]byte, costktypes.CostakerRewardsTracker],
) (math.Int, error) {
	totalScore := math.ZeroInt()
	iter, err := rwdTrackers.Iterate(ctx, nil)
	if err != nil {
		return totalScore, err
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		tracker, err := iter.Value()
		if err != nil {
			return totalScore, err
		}

		totalScore = totalScore.Add(tracker.TotalScore)
	}

	return totalScore, nil
}

func delegatingToActiveFP(fpBtcPks []bbn.BIP340PubKey, activeFps map[string]*ftypes.FinalityProviderDistInfo) bool {
	// check if delegating to an active FP
	isActiveDel := false
	for _, fpBtcPk := range fpBtcPks {
		if _, ok := activeFps[fpBtcPk.MarshalHex()]; ok {
			isActiveDel = true
			break
		}
	}

	return isActiveDel
}
	return block, bbnBlock
}

// BatchUnbondDelegations processes multiple unbondings in a single Babylon block.
// All unbonding BTC txs are included in the same BTC block, and all MsgBTCUndelegate
// messages are submitted before generating the Babylon block.
func (d *BabylonAppDriver) BatchUnbondDelegations(unbondings []*UnbondingInfo) {
	if len(unbondings) == 0 {
		return
	}

	// Collect all unbonding txs
	var unbondingTxs []*wire.MsgTx
	for _, ub := range unbondings {
		unbondingTxs = append(unbondingTxs, ub.UnbondingTx)
	}

	// Include all unbonding txs in a single BTC block and confirm
	btcCheckpointParams := d.GetBTCCkptParams(d.t)
	tip, _ := d.GetBTCLCTip()

	block := datagen.GenRandomBtcdBlockWithTransactions(d.r, unbondingTxs, tip)
	headers := BlocksWithProofsToHeaderBytes([]*datagen.BlockWithProofs{block})

	confirmationBlocks := datagen.GenNEmptyBlocks(
		d.r,
		uint64(btcCheckpointParams.BtcConfirmationDepth),
		&block.Block.Header,
	)
	confirmationHeaders := BlocksWithProofsToHeaderBytes(confirmationBlocks)
	headers = append(headers, confirmationHeaders...)

	// Insert BTC headers (this creates a Babylon block)
	d.SendTxWithMsgsFromDriverAccount(d.t, &btclighttypes.MsgInsertHeaders{
		Signer:  d.GetDriverAccountAddress().String(),
		Headers: headers,
	})

	// Now send all MsgBTCUndelegate messages (one per staker)
	// Each proof index is: 0 = coinbase, 1 = first unbonding tx, 2 = second, etc.
	for i, ub := range unbondings {
		stakingTxBz, err := bbn.SerializeBTCTx(ub.StakingTx)
		require.NoError(d.t, err)

		unbondingTxBytes, err := bbn.SerializeBTCTx(ub.UnbondingTx)
		require.NoError(d.t, err)

		// Proof index is i+1 because index 0 is the coinbase tx
		msg := &bstypes.MsgBTCUndelegate{
			Signer:                        ub.Staker.AddressString(),
			StakingTxHash:                 ub.StakingTxHash.String(),
			StakeSpendingTx:               unbondingTxBytes,
			StakeSpendingTxInclusionProof: bstypes.NewInclusionProofFromSpvProof(block.Proofs[i+1]),
			FundingTransactions:           [][]byte{stakingTxBz},
		}

		ub.Staker.SendMessage(msg)
	}

	// Generate the Babylon block that includes all unbonding messages
	d.GenerateNewBlockAssertExecutionSuccess()
}

func (d *BabylonAppDriver) IncludeTxsInBTC(txs []*wire.MsgTx) *datagen.BlockWithProofs {
	tip, _ := d.GetBTCLCTip()

	BTCPrivateKey *btcec.PrivateKey
}

// UnbondingInfo holds the data needed to submit an unbonding transaction
type UnbondingInfo struct {
	Staker        *Staker
	StakingTx     *wire.MsgTx
	StakingTxHash *chainhash.Hash
	UnbondingTx   *wire.MsgTx
}

func (s *Staker) BTCPublicKey() *bbn.BIP340PubKey {
	pk := bbn.NewBIP340PubKeyFromBTCPK(s.BTCPrivateKey.PubKey())
	return pk

	s.SendMessage(msg)
}

// PrepareUnbonding prepares unbonding data without including it in BTC or sending the message.
// Returns the UnbondingInfo which can be used with BatchUnbondDelegations.
func (s *Staker) PrepareUnbonding(
	stakingTxHash *chainhash.Hash,
	stakingTx *wire.MsgTx,
	covSender *CovenantSender,
) *UnbondingInfo {
	params := s.d.GetBTCStakingParams(s.t)

	delegation := s.d.GetBTCDelegation(s.t, stakingTxHash.String())
	require.NotNil(s.t, delegation, "delegation should exist")
	infos := parseInfos(s.t, delegation, params)

	unbondingPathSpendInfo, err := infos.StakingSlashingInfo.StakingInfo.UnbondingPathSpendInfo()
	require.NoError(s.t, err)

	stakingOutput := stakingTx.TxOut[delegation.StakingOutputIdx]

	covenantSKs := covSender.CovenantPrivateKeys()
	covenantSigs := datagen.GenerateSignatures(
		s.t,
		covenantSKs,
		infos.UnbondingSlashingInfo.UnbondingTx,
		stakingOutput,
		unbondingPathSpendInfo.RevealedLeaf,
	)

	stakerSig, err := btcstaking.SignTxWithOneScriptSpendInputFromTapLeaf(
		infos.UnbondingSlashingInfo.UnbondingTx,
		stakingOutput,
		s.BTCPrivateKey,
		unbondingPathSpendInfo.RevealedLeaf,
	)
	require.NoError(s.t, err)

	witness, err := unbondingPathSpendInfo.CreateUnbondingPathWitness(covenantSigs, stakerSig)
	require.NoError(s.t, err)

	unbondingTxMsg := infos.UnbondingSlashingInfo.UnbondingTx
	unbondingTxMsg.TxIn[0].Witness = witness

	return &UnbondingInfo{
		Staker:        s,
		StakingTx:     stakingTx,
		StakingTxHash: stakingTxHash,
		UnbondingTx:   unbondingTxMsg,
	}
}
