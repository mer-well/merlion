package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	merlion "github.com/merlion-zone/merlion/types"
	"github.com/merlion-zone/merlion/x/oracle/types"
)

// RewardBallotWinners gives out portion of seigniorage reward to the
// oracle voters that voted faithfully, at the end of every VotePeriod.
func (k Keeper) RewardBallotWinners(
	ctx sdk.Context,
	votePeriod int64,
	rewardDistributionWindow int64,
	voteTargets map[string]struct{},
	ballotWinners map[string]types.Claim,
) {
	rewardDenoms := make([]string, len(voteTargets)+1)
	rewardDenoms[0] = merlion.AttoLionDenom

	i := 1
	for denom := range voteTargets {
		rewardDenoms[i] = denom
		i++
	}

	// Sum weight of the claims
	ballotPowerSum := int64(0)
	for _, winner := range ballotWinners {
		ballotPowerSum += winner.Weight
	}

	// Exit if the ballot is empty
	if ballotPowerSum == 0 {
		return
	}

	distributionRatio := sdk.NewDec(votePeriod).QuoInt64(rewardDistributionWindow)

	var periodRewards sdk.DecCoins
	for _, denom := range rewardDenoms {
		rewardPool := k.GetRewardPool(ctx, denom)

		// Return if there's no rewards to give out
		if rewardPool.IsZero() {
			continue
		}

		periodRewards = periodRewards.Add(sdk.NewDecCoinFromDec(
			denom,
			sdk.NewDecFromInt(rewardPool.Amount).Mul(distributionRatio),
		))
	}

	// Dole out rewards
	var distributedReward sdk.Coins
	for _, winner := range ballotWinners {
		receiverVal := k.StakingKeeper().Validator(ctx, winner.Recipient)

		// Reflects contribution
		rewardCoins, _ := periodRewards.MulDec(sdk.NewDec(winner.Weight).QuoInt64(ballotPowerSum)).TruncateDecimal()

		// In case absence of the validator, we just skip distribution
		if receiverVal != nil && !rewardCoins.IsZero() {
			k.distrKeeper.AllocateTokensToValidator(ctx, receiverVal, sdk.NewDecCoinsFromCoins(rewardCoins...))
			distributedReward = distributedReward.Add(rewardCoins...)
		}
	}

	// Move distributed reward to distribution module
	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.distrName, distributedReward)
	if err != nil {
		panic(fmt.Sprintf("failed to send coins to distribution module %s", err.Error()))
	}
}
