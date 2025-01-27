package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	merlion "github.com/merlion-zone/merlion/types"
	"github.com/merlion-zone/merlion/x/maker/types"
	oracletypes "github.com/merlion-zone/merlion/x/oracle/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (m msgServer) MintBySwap(c context.Context, msg *types.MsgMintBySwap) (*types.MsgMintBySwapResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	backingDenom := msg.BackingInMax.Denom

	sender, receiver, err := getSenderReceiver(msg.Sender, msg.To)
	if err != nil {
		return nil, err
	}

	// get prices in usd
	backingPrice, err := m.Keeper.oracleKeeper.GetExchangeRate(ctx, backingDenom)
	if err != nil {
		return nil, err
	}
	lionPrice, err := m.Keeper.oracleKeeper.GetExchangeRate(ctx, merlion.AttoLionDenom)
	if err != nil {
		return nil, err
	}
	merPrice, err := m.Keeper.oracleKeeper.GetExchangeRate(ctx, merlion.MicroUSDDenom)
	if err != nil {
		return nil, err
	}

	// check price lower bound
	merPriceLowerBound := merlion.MicroUSDTarget.Mul(sdk.OneDec().Sub(m.Keeper.MintPriceBias(ctx)))
	if merPrice.LT(merPriceLowerBound) {
		return nil, sdkerrors.Wrapf(types.ErrMerPriceTooLow, "%s price too low: %s", merlion.MicroUSDDenom, merPrice)
	}

	collateralRatio := m.Keeper.GetCollateralRatio(ctx)

	backingParams, found := m.Keeper.GetBackingRiskParams(ctx, backingDenom)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrBackingCoinNotFound, "backing coin denomination not found: %s", backingDenom)
	}
	if !backingParams.Enabled {
		return nil, sdkerrors.Wrapf(types.ErrBackingCoinDisabled, "backing coin disabled: %s", backingDenom)
	}

	totalBacking, poolBacking, err := m.Keeper.getBacking(ctx, backingDenom)
	if err != nil {
		return nil, err
	}

	mintOut := msg.MintOut
	mintFee := computeFee(mintOut, backingParams.MintFee)
	mintTotal := mintOut.AddAmount(mintFee.Amount)
	mintTotalInUSD := mintTotal.Amount.ToDec().Mul(merlion.MicroUSDTarget)

	poolBacking.MerMinted = poolBacking.MerMinted.Add(mintTotal)
	totalBacking.MerMinted = totalBacking.MerMinted.Add(mintTotal)
	if backingParams.MaxMerMint != nil {
		if poolBacking.MerMinted.Amount.GT(*backingParams.MaxMerMint) {
			return nil, sdkerrors.Wrapf(types.ErrMerCeiling, "mer over ceiling")
		}
	}

	backingIn := sdk.NewCoin(backingDenom, sdk.ZeroInt())
	lionIn := sdk.NewCoin(merlion.AttoLionDenom, sdk.ZeroInt())
	if collateralRatio.GTE(sdk.OneDec()) || msg.LionInMax.IsZero() {
		// full/over collateralized, or user selects full collateralization
		backingIn.Amount = mintTotalInUSD.QuoRoundUp(backingPrice).RoundInt()
	} else if collateralRatio.IsZero() {
		// algorithmic
		lionIn.Amount = mintTotalInUSD.QuoRoundUp(lionPrice).RoundInt()
	} else {
		// fractional
		backingIn.Amount = mintTotalInUSD.Mul(collateralRatio).QuoRoundUp(backingPrice).RoundInt()
		lionIn.Amount = mintTotalInUSD.Mul(sdk.OneDec().Sub(collateralRatio)).QuoRoundUp(lionPrice).RoundInt()
	}

	if msg.BackingInMax.IsLT(backingIn) {
		return nil, sdkerrors.Wrapf(types.ErrBackingCoinSlippage, "backing coin needed: %s", backingIn)
	}
	if msg.LionInMax.IsLT(lionIn) {
		return nil, sdkerrors.Wrapf(types.ErrLionCoinSlippage, "lion coin needed: %s", lionIn)
	}

	poolBacking.Backing = poolBacking.Backing.Add(backingIn)
	if backingParams.MaxBacking != nil {
		if poolBacking.Backing.Amount.GT(*backingParams.MaxBacking) {
			return nil, sdkerrors.Wrapf(types.ErrBackingCeiling, "backing over ceiling")
		}
	}

	poolBacking.LionBurned = poolBacking.LionBurned.Add(lionIn)
	totalBacking.LionBurned = totalBacking.LionBurned.Add(lionIn)

	m.Keeper.SetPoolBacking(ctx, poolBacking)
	m.Keeper.SetTotalBacking(ctx, totalBacking)

	// take backing and lion coin
	err = m.Keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(backingIn, lionIn))
	if err != nil {
		return nil, err
	}
	// burn lion
	if lionIn.IsPositive() {
		err = m.Keeper.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(lionIn))
		if err != nil {
			return nil, err
		}
	}

	// mint mer stablecoin
	err = m.Keeper.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(mintTotal))
	if err != nil {
		return nil, err
	}
	// send mer to receiver
	err = m.Keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, sdk.NewCoins(mintOut))
	if err != nil {
		return nil, err
	}
	// send mer fee to oracle
	if mintFee.IsPositive() {
		err = m.Keeper.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, oracletypes.ModuleName, sdk.NewCoins(mintFee))
		if err != nil {
			return nil, err
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventTypeMintBySwap,
			sdk.NewAttribute(types.AttributeKeyCoinIn, sdk.NewCoins(backingIn, lionIn).String()),
			sdk.NewAttribute(types.AttributeKeyCoinOut, mintOut.String()),
			sdk.NewAttribute(types.AttributeKeyFee, mintFee.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &types.MsgMintBySwapResponse{
		MintFee:   mintFee,
		BackingIn: backingIn,
		LionIn:    lionIn,
	}, nil
}

func (m msgServer) BurnBySwap(c context.Context, msg *types.MsgBurnBySwap) (*types.MsgBurnBySwapResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	backingDenom := msg.BackingOutMin.Denom

	sender, receiver, err := getSenderReceiver(msg.Sender, msg.To)
	if err != nil {
		return nil, err
	}

	// get prices in usd
	backingPrice, err := m.Keeper.oracleKeeper.GetExchangeRate(ctx, backingDenom)
	if err != nil {
		return nil, err
	}
	lionPrice, err := m.Keeper.oracleKeeper.GetExchangeRate(ctx, merlion.AttoLionDenom)
	if err != nil {
		return nil, err
	}
	merPrice, err := m.Keeper.oracleKeeper.GetExchangeRate(ctx, merlion.MicroUSDDenom)
	if err != nil {
		return nil, err
	}

	// check price upper bound
	merPriceUpperBound := merlion.MicroUSDTarget.Mul(sdk.OneDec().Add(m.Keeper.BurnPriceBias(ctx)))
	if merPrice.GT(merPriceUpperBound) {
		return nil, sdkerrors.Wrapf(types.ErrMerPriceTooHigh, "%s price too high: %s", merlion.MicroUSDDenom, merPrice)
	}

	collateralRatio := m.Keeper.GetCollateralRatio(ctx)

	backingParams, found := m.Keeper.GetBackingRiskParams(ctx, backingDenom)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrBackingCoinNotFound, "backing coin denomination not found: %s", backingDenom)
	}
	if !backingParams.Enabled {
		return nil, sdkerrors.Wrapf(types.ErrBackingCoinDisabled, "backing coin disabled: %s", backingDenom)
	}

	totalBacking, poolBacking, err := m.Keeper.getBacking(ctx, backingDenom)
	if err != nil {
		return nil, err
	}

	burnIn := msg.BurnIn
	burnFee := computeFee(burnIn, backingParams.BurnFee)
	burn := burnIn.SubAmount(burnFee.Amount)
	burnInUSD := burn.Amount.ToDec().Mul(merlion.MicroUSDTarget)

	backingOut := sdk.NewCoin(backingDenom, sdk.ZeroInt())
	lionOut := sdk.NewCoin(merlion.AttoLionDenom, sdk.ZeroInt())
	if collateralRatio.GTE(sdk.OneDec()) {
		// full/over collateralized
		backingOut.Amount = burnInUSD.QuoRoundUp(backingPrice).RoundInt()
	} else if collateralRatio.IsZero() {
		// algorithmic
		lionOut.Amount = burnInUSD.QuoRoundUp(lionPrice).RoundInt()
	} else {
		// fractional
		backingOut.Amount = burnInUSD.Mul(collateralRatio).QuoRoundUp(backingPrice).RoundInt()
		lionOut.Amount = burnInUSD.Mul(sdk.OneDec().Sub(collateralRatio)).QuoRoundUp(lionPrice).RoundInt()
	}

	if backingOut.IsLT(msg.BackingOutMin) {
		return nil, sdkerrors.Wrapf(types.ErrBackingCoinSlippage, "backing coin out: %s", backingOut)
	}
	if lionOut.IsLT(msg.LionOutMin) {
		return nil, sdkerrors.Wrapf(types.ErrLionCoinSlippage, "lion coin out: %s", lionOut)
	}

	moduleOwnedBacking := m.Keeper.bankKeeper.GetBalance(ctx, m.Keeper.accountKeeper.GetModuleAddress(types.ModuleName), backingDenom)
	if moduleOwnedBacking.IsLT(backingOut) {
		return nil, sdkerrors.Wrapf(types.ErrBackingCoinInsufficient, "backing coin out(%s) < balance(%s)", backingOut, moduleOwnedBacking)
	}

	poolBacking.Backing = poolBacking.Backing.Sub(backingOut)
	// allow LionBurned to be negative
	// here use AddAmount(Neg()) to bypass Sub negativeness check
	poolBacking.LionBurned = poolBacking.LionBurned.AddAmount(lionOut.Amount.Neg())
	totalBacking.LionBurned = totalBacking.LionBurned.AddAmount(lionOut.Amount.Neg())
	// allow MerMinted to be negative
	poolBacking.MerMinted = poolBacking.MerMinted.AddAmount(burn.Amount.Neg())
	poolBacking.MerMinted = poolBacking.MerMinted.AddAmount(burn.Amount.Neg())

	m.Keeper.SetPoolBacking(ctx, poolBacking)
	m.Keeper.SetTotalBacking(ctx, totalBacking)

	// take mer stablecoin
	err = m.Keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(burnIn))
	if err != nil {
		return nil, err
	}
	// burn mer
	err = m.Keeper.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(burn))
	if err != nil {
		return nil, err
	}
	// send mer fee to oracle
	err = m.Keeper.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, oracletypes.ModuleName, sdk.NewCoins(burnFee))
	if err != nil {
		return nil, err
	}

	// mint lion
	err = m.Keeper.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(lionOut))
	if err != nil {
		return nil, err
	}
	// send backing and lion to receiver
	err = m.Keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, sdk.NewCoins(backingOut, lionOut))
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventTypeBurnBySwap,
			sdk.NewAttribute(types.AttributeKeyCoinIn, burn.String()),
			sdk.NewAttribute(types.AttributeKeyCoinOut, sdk.NewCoins(backingOut, lionOut).String()),
			sdk.NewAttribute(types.AttributeKeyFee, burnFee.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &types.MsgBurnBySwapResponse{
		BurnFee:    burnFee,
		BackingOut: backingOut,
		LionOut:    lionOut,
	}, nil
}

func (m msgServer) BuyBacking(c context.Context, msg *types.MsgBuyBacking) (*types.MsgBuyBackingResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	backingDenom := msg.BackingOutMin.Denom

	sender, receiver, err := getSenderReceiver(msg.Sender, msg.To)
	if err != nil {
		return nil, err
	}

	// get prices in usd
	backingPrice, err := m.Keeper.oracleKeeper.GetExchangeRate(ctx, backingDenom)
	if err != nil {
		return nil, err
	}
	lionPrice, err := m.Keeper.oracleKeeper.GetExchangeRate(ctx, merlion.AttoLionDenom)
	if err != nil {
		return nil, err
	}

	collateralRatio := m.Keeper.GetCollateralRatio(ctx)

	backingParams, found := m.Keeper.GetBackingRiskParams(ctx, backingDenom)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrBackingCoinNotFound, "backing coin denomination not found: %s", backingDenom)
	}
	if !backingParams.Enabled {
		return nil, sdkerrors.Wrapf(types.ErrBackingCoinDisabled, "backing coin disabled: %s", backingDenom)
	}

	totalBacking, poolBacking, err := m.Keeper.getBacking(ctx, backingDenom)
	if err != nil {
		return nil, err
	}

	totalBackingValue, err := m.Keeper.totalBackingInUSD(ctx)
	if err != nil {
		return nil, err
	}

	if !totalBacking.MerMinted.IsPositive() {
		return nil, sdkerrors.Wrap(types.ErrBackingCoinInsufficient, "insufficient available backing coin")
	}
	requiredBackingValue := totalBacking.MerMinted.Amount.ToDec().Mul(collateralRatio).TruncateInt()

	availableExcessBackingValue := sdk.NewCoin(merlion.MicroUSDDenom, sdk.ZeroInt())
	if requiredBackingValue.LT(totalBackingValue.Amount) {
		availableExcessBackingValue.Amount = totalBackingValue.Amount.Sub(requiredBackingValue)
	}

	lionInValue := msg.LionIn.Amount.ToDec().Mul(lionPrice)
	if lionInValue.TruncateInt().GT(availableExcessBackingValue.Amount) {
		return nil, sdkerrors.Wrap(types.ErrBackingCoinInsufficient, "insufficient available backing coin")
	}

	backingOut := sdk.NewCoin(backingDenom, lionInValue.Quo(backingPrice).TruncateInt())
	if poolBacking.Backing.IsLT(backingOut) {
		return nil, sdkerrors.Wrap(types.ErrBackingCoinInsufficient, "insufficient available backing coin")
	}
	poolBacking.Backing = poolBacking.Backing.Sub(backingOut)

	fee := computeFee(backingOut, backingParams.BuybackFee)
	backingOut = backingOut.Sub(fee)

	if backingOut.IsLT(msg.BackingOutMin) {
		return nil, sdkerrors.Wrap(types.ErrBackingCoinSlippage, "backing coin over slippage")
	}

	poolBacking.LionBurned = poolBacking.LionBurned.Add(msg.LionIn)
	totalBacking.LionBurned = totalBacking.LionBurned.Add(msg.LionIn)

	m.Keeper.SetPoolBacking(ctx, poolBacking)
	m.Keeper.SetTotalBacking(ctx, totalBacking)

	// take lion-in
	err = m.Keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(msg.LionIn))
	if err != nil {
		return nil, err
	}
	// burn lion
	err = m.Keeper.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(msg.LionIn))
	if err != nil {
		return nil, err
	}

	// send backing to receiver
	err = m.Keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, sdk.NewCoins(backingOut))
	if err != nil {
		return nil, err
	}
	// send fee to oracle
	err = m.Keeper.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, oracletypes.ModuleName, sdk.NewCoins(fee))
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventTypeBuyBacking,
			sdk.NewAttribute(types.AttributeKeyCoinIn, msg.LionIn.String()),
			sdk.NewAttribute(types.AttributeKeyCoinOut, backingOut.String()),
			sdk.NewAttribute(types.AttributeKeyFee, fee.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &types.MsgBuyBackingResponse{
		BackingOut: backingOut,
	}, nil
}

func (m msgServer) SellBacking(c context.Context, msg *types.MsgSellBacking) (*types.MsgSellBackingResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	backingDenom := msg.BackingIn.Denom

	sender, receiver, err := getSenderReceiver(msg.Sender, msg.To)
	if err != nil {
		return nil, err
	}

	// get prices in usd
	backingPrice, err := m.Keeper.oracleKeeper.GetExchangeRate(ctx, backingDenom)
	if err != nil {
		return nil, err
	}
	lionPrice, err := m.Keeper.oracleKeeper.GetExchangeRate(ctx, merlion.AttoLionDenom)
	if err != nil {
		return nil, err
	}

	collateralRatio := m.Keeper.GetCollateralRatio(ctx)

	backingParams, found := m.Keeper.GetBackingRiskParams(ctx, backingDenom)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrBackingCoinNotFound, "backing coin denomination not found: %s", backingDenom)
	}
	if !backingParams.Enabled {
		return nil, sdkerrors.Wrapf(types.ErrBackingCoinDisabled, "backing coin disabled: %s", backingDenom)
	}

	totalBacking, poolBacking, err := m.Keeper.getBacking(ctx, backingDenom)
	if err != nil {
		return nil, err
	}

	totalBackingValue, err := m.Keeper.totalBackingInUSD(ctx)
	if err != nil {
		return nil, err
	}

	requiredBackingValue := totalBacking.MerMinted.Amount.ToDec().Mul(collateralRatio).TruncateInt()

	availableMissingBackingValue := sdk.NewCoin(merlion.MicroUSDDenom, sdk.ZeroInt())
	if requiredBackingValue.GT(totalBackingValue.Amount) {
		availableMissingBackingValue.Amount = requiredBackingValue.Sub(totalBackingValue.Amount)
	}
	availableLionOut := availableMissingBackingValue.Amount.ToDec().Quo(lionPrice)

	bonusRatio := m.Keeper.RecollateralizeBonus(ctx)
	lionMint := sdk.NewCoin(merlion.AttoLionDenom, msg.BackingIn.Amount.ToDec().Mul(backingPrice).Quo(lionPrice).TruncateInt())
	bonus := computeFee(lionMint, &bonusRatio)
	fee := computeFee(lionMint, backingParams.RecollateralizeFee)

	lionMint = lionMint.Add(bonus)
	if lionMint.Amount.ToDec().GT(availableLionOut) {
		return nil, sdkerrors.Wrap(types.ErrLionCoinInsufficient, "insufficient available lion coin")
	}
	lionOut := lionMint.Sub(fee)

	if lionOut.IsLT(msg.LionOutMin) {
		return nil, sdkerrors.Wrap(types.ErrLionCoinSlippage, "lion coin over slippage")
	}

	poolBacking.Backing = poolBacking.Backing.Add(msg.BackingIn)
	if backingParams.MaxBacking != nil {
		if poolBacking.Backing.Amount.GT(*backingParams.MaxBacking) {
			return nil, sdkerrors.Wrap(types.ErrBackingCeiling, "over ceiling")
		}
	}

	// allow LionBurned to be negative
	// here use AddAmount(Neg()) to bypass Sub negativeness check
	poolBacking.LionBurned = poolBacking.LionBurned.AddAmount(lionMint.Amount.Neg())
	totalBacking.LionBurned = totalBacking.LionBurned.AddAmount(lionMint.Amount.Neg())

	m.Keeper.SetPoolBacking(ctx, poolBacking)
	m.Keeper.SetTotalBacking(ctx, totalBacking)

	// take backing-in
	err = m.Keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(msg.BackingIn))
	if err != nil {
		return nil, err
	}

	// mint lion
	err = m.Keeper.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(lionMint))
	if err != nil {
		return nil, err
	}
	// send lion to receiver
	err = m.Keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, sdk.NewCoins(lionOut))
	if err != nil {
		return nil, err
	}
	// send fee to oracle
	err = m.Keeper.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, oracletypes.ModuleName, sdk.NewCoins(fee))
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventTypeSellBacking,
			sdk.NewAttribute(types.AttributeKeyCoinIn, msg.BackingIn.String()),
			sdk.NewAttribute(types.AttributeKeyCoinOut, lionOut.String()),
			sdk.NewAttribute(types.AttributeKeyFee, fee.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &types.MsgSellBackingResponse{
		LionOut: lionOut,
	}, nil
}

func (m msgServer) MintByCollateral(c context.Context, msg *types.MsgMintByCollateral) (*types.MsgMintByCollateralResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	collateralDenom := msg.CollateralDenom

	sender, receiver, err := getSenderReceiver(msg.Sender, msg.To)
	if err != nil {
		return nil, err
	}

	// get prices in usd
	collateralPrice, err := m.Keeper.oracleKeeper.GetExchangeRate(ctx, collateralDenom)
	if err != nil {
		return nil, err
	}
	lionPrice, err := m.Keeper.oracleKeeper.GetExchangeRate(ctx, merlion.AttoLionDenom)
	if err != nil {
		return nil, err
	}
	merPrice, err := m.Keeper.oracleKeeper.GetExchangeRate(ctx, merlion.MicroUSDDenom)
	if err != nil {
		return nil, err
	}

	// check price lower bound
	merPriceLowerBound := merlion.MicroUSDTarget.Mul(sdk.OneDec().Sub(m.Keeper.MintPriceBias(ctx)))
	if merPrice.LT(merPriceLowerBound) {
		return nil, sdkerrors.Wrapf(types.ErrMerPriceTooLow, "%s price too low: %s", merlion.MicroUSDDenom, merPrice)
	}

	collateralParams, found := m.Keeper.GetCollateralRiskParams(ctx, collateralDenom)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrCollateralCoinNotFound, "collateral coin denomination not found: %s", collateralDenom)
	}
	if !collateralParams.Enabled {
		return nil, sdkerrors.Wrapf(types.ErrCollateralCoinDisabled, "collateral coin disabled: %s", collateralDenom)
	}

	totalColl, poolColl, accColl, err := m.Keeper.getCollateral(ctx, sender, collateralDenom)
	if err != nil {
		return nil, err
	}

	// settle interestFee fee
	settleInterestFee(ctx, &accColl, &poolColl, &totalColl, collateralParams.InterestFee)

	// compute mint amount
	mintFee := computeFee(msg.MintOut, collateralParams.MintFee)
	mint := msg.MintOut.Add(mintFee)

	// update debt
	accColl.MerDebt = accColl.MerDebt.Add(mint)
	poolColl.MerDebt = poolColl.MerDebt.Add(mint)
	totalColl.MerDebt = totalColl.MerDebt.Add(mint)

	if collateralParams.MaxMerMint != nil {
		if poolColl.MerDebt.Amount.GT(*collateralParams.MaxMerMint) {
			return nil, sdkerrors.Wrapf(types.ErrMerCeiling, "mer over ceiling")
		}
	}

	// compute actual catalytic lion
	merDue := accColl.MerDebt.Add(accColl.MerByLion)
	bestCatalyticLionInUSD := merDue.Amount.ToDec().Mul(*collateralParams.CatalyticLionRatio)
	lionInMaxInUSD := msg.LionInMax.Amount.ToDec().Mul(lionPrice).TruncateInt()
	catalyticLionInUSD := sdk.MinDec(bestCatalyticLionInUSD, accColl.MerByLion.Amount.Add(lionInMaxInUSD).ToDec()).TruncateInt()

	// compute actual lion-in
	lionInInUSD := catalyticLionInUSD.Sub(accColl.MerByLion.Amount)
	if !lionInInUSD.IsPositive() {
		lionInInUSD = sdk.ZeroInt()
	} else {
		accColl.MerByLion = accColl.MerByLion.AddAmount(lionInInUSD)
		poolColl.MerByLion = poolColl.MerByLion.AddAmount(lionInInUSD)
		totalColl.MerByLion = totalColl.MerByLion.AddAmount(lionInInUSD)
		accColl.MerDebt = accColl.MerDebt.SubAmount(lionInInUSD)
		poolColl.MerDebt = poolColl.MerDebt.SubAmount(lionInInUSD)
		totalColl.MerDebt = totalColl.MerDebt.SubAmount(lionInInUSD)
	}
	lionIn := sdk.NewCoin(merlion.AttoLionDenom, lionInInUSD.ToDec().Quo(lionPrice).TruncateInt())

	accColl.LionBurned = accColl.LionBurned.Add(lionIn)
	poolColl.LionBurned = poolColl.LionBurned.Add(lionIn)
	totalColl.LionBurned = totalColl.LionBurned.Add(lionIn)

	// compute actual catalytic ratio and max loan-to-value
	maxLoanToValue := maxLoanToValueForAccount(&accColl, &collateralParams)

	// check max mintable mer
	collateralInUSD := accColl.Collateral.Amount.ToDec().Mul(collateralPrice)
	maxMerMintable := collateralInUSD.Mul(maxLoanToValue)
	if maxMerMintable.LT(accColl.MerDebt.Amount.ToDec()) {
		return nil, sdkerrors.Wrapf(types.ErrAccountInsufficientCollateral, "account has insufficient collateral: %s", collateralDenom)
	}

	// eventually update collateral
	m.Keeper.SetAccountCollateral(ctx, sender, accColl)
	m.Keeper.SetPoolCollateral(ctx, poolColl)
	m.Keeper.SetTotalCollateral(ctx, totalColl)

	// take lion and burn it
	if lionIn.IsPositive() {
		err = m.Keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(lionIn))
		if err != nil {
			return nil, err
		}
		err = m.Keeper.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(lionIn))
		if err != nil {
			return nil, err
		}
	}

	// mint mer
	err = m.Keeper.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(mint))
	if err != nil {
		return nil, err
	}
	// send mer to receiver
	err = m.Keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, sdk.NewCoins(msg.MintOut))
	if err != nil {
		return nil, err
	}
	// send mint fee to oracle
	if mintFee.IsPositive() {
		err = m.Keeper.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, oracletypes.ModuleName, sdk.NewCoins(mintFee))
		if err != nil {
			return nil, err
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventTypeMintByCollateral,
			sdk.NewAttribute(types.AttributeKeyCoinIn, lionIn.String()),
			sdk.NewAttribute(types.AttributeKeyCoinOut, msg.MintOut.String()),
			sdk.NewAttribute(types.AttributeKeyFee, mintFee.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &types.MsgMintByCollateralResponse{
		MintFee: mintFee,
		LionIn:  lionIn,
	}, nil
}

func (m msgServer) BurnByCollateral(c context.Context, msg *types.MsgBurnByCollateral) (*types.MsgBurnByCollateralResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	collateralDenom := msg.CollateralDenom

	sender, _, err := getSenderReceiver(msg.Sender, "")
	if err != nil {
		return nil, err
	}

	// get prices in usd
	merPrice, err := m.Keeper.oracleKeeper.GetExchangeRate(ctx, merlion.MicroUSDDenom)
	if err != nil {
		return nil, err
	}

	// check price upper bound
	merPriceUpperBound := merlion.MicroUSDTarget.Mul(sdk.OneDec().Add(m.Keeper.BurnPriceBias(ctx)))
	if merPrice.GT(merPriceUpperBound) {
		return nil, sdkerrors.Wrapf(types.ErrMerPriceTooHigh, "%s price too high: %s", merlion.MicroUSDDenom, merPrice)
	}

	collateralParams, found := m.Keeper.GetCollateralRiskParams(ctx, collateralDenom)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrCollateralCoinNotFound, "collateral coin denomination not found: %s", collateralDenom)
	}
	if !collateralParams.Enabled {
		return nil, sdkerrors.Wrapf(types.ErrCollateralCoinDisabled, "collateral coin disabled: %s", collateralDenom)
	}

	totalColl, poolColl, accColl, err := m.Keeper.getCollateral(ctx, sender, collateralDenom)
	if err != nil {
		return nil, err
	}

	// settle interestFee fee
	settleInterestFee(ctx, &accColl, &poolColl, &totalColl, collateralParams.InterestFee)

	// compute burn-in
	if !accColl.MerDebt.IsPositive() {
		return nil, sdkerrors.Wrapf(types.ErrAccountNoDebt, "account has no debt for %s collateral", collateralDenom)
	}
	repayIn := sdk.NewCoin(msg.RepayInMax.Denom, sdk.MinInt(accColl.MerDebt.Amount, msg.RepayInMax.Amount))
	interestFee := sdk.NewCoin(msg.RepayInMax.Denom, sdk.MinInt(accColl.LastInterest.Amount, repayIn.Amount))
	burn := repayIn.Sub(interestFee)

	// update debt
	accColl.LastInterest = accColl.LastInterest.Sub(interestFee)
	accColl.MerDebt = accColl.MerDebt.Sub(repayIn)
	poolColl.MerDebt = poolColl.MerDebt.Sub(repayIn)
	totalColl.MerDebt = totalColl.MerDebt.Sub(repayIn)

	// eventually update collateral
	m.Keeper.SetAccountCollateral(ctx, sender, accColl)
	m.Keeper.SetPoolCollateral(ctx, poolColl)
	m.Keeper.SetTotalCollateral(ctx, totalColl)

	// take mer
	err = m.Keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(repayIn))
	if err != nil {
		return nil, err
	}
	// burn mer
	if burn.IsPositive() {
		err = m.Keeper.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(burn))
		if err != nil {
			return nil, err
		}
	}
	// send fee to oracle
	if interestFee.IsPositive() {
		err = m.Keeper.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, oracletypes.ModuleName, sdk.NewCoins(interestFee))
		if err != nil {
			return nil, err
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventTypeBurnByCollateral,
			sdk.NewAttribute(types.AttributeKeyCoinIn, repayIn.String()),
			sdk.NewAttribute(types.AttributeKeyFee, interestFee.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &types.MsgBurnByCollateralResponse{
		RepayIn: repayIn,
	}, nil
}

func (m msgServer) DepositCollateral(c context.Context, msg *types.MsgDepositCollateral) (*types.MsgDepositCollateralResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	collateralDenom := msg.Collateral.Denom

	sender, receiver, err := getSenderReceiver(msg.Sender, msg.To)
	if err != nil {
		return nil, err
	}

	collateralParams, found := m.Keeper.GetCollateralRiskParams(ctx, collateralDenom)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrCollateralCoinNotFound, "collateral coin denomination not found: %s", collateralDenom)
	}
	if !collateralParams.Enabled {
		return nil, sdkerrors.Wrapf(types.ErrCollateralCoinDisabled, "collateral coin disabled: %s", collateralDenom)
	}

	totalColl, poolColl, accColl, err := m.Keeper.getCollateral(ctx, receiver, collateralDenom)
	if err != nil {
		return nil, err
	}

	accColl.Collateral = accColl.Collateral.Add(msg.Collateral)
	poolColl.Collateral = poolColl.Collateral.Add(msg.Collateral)

	if collateralParams.MaxCollateral != nil {
		if poolColl.Collateral.Amount.GT(*collateralParams.MaxCollateral) {
			return nil, sdkerrors.Wrap(types.ErrCollateralCeiling, "collateral over ceiling")
		}
	}

	settleInterestFee(ctx, &accColl, &poolColl, &totalColl, collateralParams.InterestFee)

	m.Keeper.SetAccountCollateral(ctx, receiver, accColl)
	m.Keeper.SetPoolCollateral(ctx, poolColl)
	m.Keeper.SetTotalCollateral(ctx, totalColl)

	// take collateral from sender
	err = m.Keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(msg.Collateral))
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventTypeDepositCollateral,
			sdk.NewAttribute(types.AttributeKeyCoinIn, msg.Collateral.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &types.MsgDepositCollateralResponse{}, nil
}

func (m msgServer) RedeemCollateral(c context.Context, msg *types.MsgRedeemCollateral) (*types.MsgRedeemCollateralResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	collateralDenom := msg.Collateral.Denom

	sender, receiver, err := getSenderReceiver(msg.Sender, msg.To)
	if err != nil {
		return nil, err
	}

	collateralParams, found := m.Keeper.GetCollateralRiskParams(ctx, collateralDenom)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrCollateralCoinNotFound, "collateral coin denomination not found: %s", collateralDenom)
	}
	if !collateralParams.Enabled {
		return nil, sdkerrors.Wrapf(types.ErrCollateralCoinDisabled, "collateral coin disabled: %s", collateralDenom)
	}

	totalColl, poolColl, accColl, err := m.Keeper.getCollateral(ctx, receiver, collateralDenom)
	if err != nil {
		return nil, err
	}

	// update collateral
	poolColl.Collateral = poolColl.Collateral.Sub(msg.Collateral)
	accColl.Collateral = accColl.Collateral.Sub(msg.Collateral)

	settleInterestFee(ctx, &accColl, &poolColl, &totalColl, collateralParams.InterestFee)

	// get prices in usd
	collateralPrice, err := m.Keeper.oracleKeeper.GetExchangeRate(ctx, collateralDenom)
	if err != nil {
		return nil, err
	}

	maxLoanToValue := maxLoanToValueForAccount(&accColl, &collateralParams)

	collateralInUSD := accColl.Collateral.Amount.ToDec().Mul(collateralPrice)
	if accColl.MerDebt.Amount.ToDec().LT(collateralInUSD.Mul(maxLoanToValue)) {
		return nil, sdkerrors.Wrapf(types.ErrAccountInsufficientCollateral, "account has insufficient collateral: %s", collateralDenom)
	}

	// eventually persist collateral
	m.Keeper.SetAccountCollateral(ctx, sender, accColl)
	m.Keeper.SetPoolCollateral(ctx, poolColl)
	m.Keeper.SetTotalCollateral(ctx, totalColl)

	// send collateral to receiver
	err = m.Keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, sdk.NewCoins(msg.Collateral))
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventTypeRedeemCollateral,
			sdk.NewAttribute(types.AttributeKeyCoinOut, msg.Collateral.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &types.MsgRedeemCollateralResponse{}, nil
}

func (m msgServer) LiquidateCollateral(c context.Context, msg *types.MsgLiquidateCollateral) (*types.MsgLiquidateCollateralResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	collateralDenom := msg.Collateral.Denom

	sender, receiver, err := getSenderReceiver(msg.Sender, msg.To)
	if err != nil {
		return nil, err
	}
	debtor, err := sdk.AccAddressFromBech32(msg.Debtor)
	if err != nil {
		return nil, err
	}

	collateralParams, found := m.Keeper.GetCollateralRiskParams(ctx, collateralDenom)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrCollateralCoinNotFound, "collateral coin denomination not found: %s", collateralDenom)
	}

	if !collateralParams.Enabled {
		return nil, sdkerrors.Wrapf(types.ErrCollateralCoinDisabled, "collateral coin disabled: %s", collateralDenom)
	}

	totalColl, poolColl, accColl, err := m.Keeper.getCollateral(ctx, receiver, collateralDenom)
	if err != nil {
		return nil, err
	}

	settleInterestFee(ctx, &accColl, &poolColl, &totalColl, collateralParams.InterestFee)

	// get prices in usd
	collateralPrice, err := m.Keeper.oracleKeeper.GetExchangeRate(ctx, collateralDenom)
	if err != nil {
		return nil, err
	}

	// check whether undercollateralized
	liquidationValue := accColl.Collateral.Amount.ToDec().Mul(collateralPrice).Mul(*collateralParams.LiquidationThreshold)
	if accColl.MerDebt.Amount.ToDec().LTE(liquidationValue) {
		return nil, sdkerrors.Wrap(types.ErrNotUndercollateralized, "not undercollateralized")
	}

	if msg.Collateral.Amount.GT(accColl.Collateral.Amount) {
		return nil, sdkerrors.Wrap(types.ErrCollateralCoinInsufficient, "collateral coin balance insufficient")
	}

	liquidationFee := msg.Collateral.Amount.ToDec().Mul(*collateralParams.LiquidationFee)
	repayIn := sdk.NewCoin(collateralDenom, msg.Collateral.Amount.ToDec().Sub(liquidationFee).Mul(collateralPrice).TruncateInt())
	commissionFee := sdk.NewCoin(collateralDenom, liquidationFee.Mul(m.Keeper.LiquidationCommissionFee(ctx)).TruncateInt())
	collateralOut := msg.Collateral.Sub(commissionFee)

	// repay for debtor as much as possible
	repayDebt := sdk.NewCoin(merlion.MicroUSDDenom, sdk.MinInt(accColl.MerDebt.Amount, repayIn.Amount))
	merRefund := repayIn.Sub(repayDebt)
	accColl.MerDebt = accColl.MerDebt.Sub(repayDebt)
	poolColl.MerDebt = poolColl.MerDebt.Sub(repayDebt)
	totalColl.MerDebt = totalColl.MerDebt.Sub(repayDebt)

	// eventually persist collateral
	m.Keeper.SetAccountCollateral(ctx, sender, accColl)
	m.Keeper.SetPoolCollateral(ctx, poolColl)
	m.Keeper.SetTotalCollateral(ctx, totalColl)

	// take mer from sender
	err = m.Keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(repayIn))
	if err != nil {
		return nil, err
	}
	// burn mer debt
	err = m.Keeper.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(repayDebt))
	if err != nil {
		return nil, err
	}
	// send excess mer to debtor
	if merRefund.IsPositive() {
		err = m.Keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, debtor, sdk.NewCoins(merRefund))
		if err != nil {
			return nil, err
		}
	}

	// send collateral to receiver
	err = m.Keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, sdk.NewCoins(collateralOut))
	if err != nil {
		return nil, err
	}
	// send liquidation commission fee to oracle
	err = m.Keeper.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, oracletypes.ModuleName, sdk.NewCoins(commissionFee))
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventTypeLiquidateCollateral,
			sdk.NewAttribute(types.AttributeKeyCoinIn, repayIn.String()),
			sdk.NewAttribute(types.AttributeKeyCoinOut, collateralOut.String()),
			sdk.NewAttribute(types.AttributeKeyFee, commissionFee.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &types.MsgLiquidateCollateralResponse{
		RepayIn:       repayIn,
		CollateralOut: collateralOut,
	}, nil
}

func (k Keeper) getBacking(ctx sdk.Context, denom string) (total types.TotalBacking, pool types.PoolBacking, err error) {
	total, found := k.GetTotalBacking(ctx)
	if !found {
		err = sdkerrors.Wrapf(types.ErrBackingCoinNotFound, "backing coin denomination not found: %s", denom)
		return
	}
	pool, found = k.GetPoolBacking(ctx, denom)
	if !found {
		err = sdkerrors.Wrapf(types.ErrBackingCoinNotFound, "backing coin denomination not found: %s", denom)
		return
	}
	return
}

func (k Keeper) getCollateral(ctx sdk.Context, account sdk.AccAddress, denom string, allowNewAccount ...bool) (total types.TotalCollateral, pool types.PoolCollateral, acc types.AccountCollateral, err error) {
	total, found := k.GetTotalCollateral(ctx)
	if !found {
		err = sdkerrors.Wrapf(types.ErrCollateralCoinNotFound, "collateral coin denomination not found: %s", denom)
		return
	}
	pool, found = k.GetPoolCollateral(ctx, denom)
	if !found {
		err = sdkerrors.Wrapf(types.ErrCollateralCoinNotFound, "collateral coin denomination not found: %s", denom)
		return
	}
	acc, found = k.GetAccountCollateral(ctx, account, denom)
	if !found {
		if len(allowNewAccount) > 0 && allowNewAccount[0] {
			acc = types.AccountCollateral{
				Account:             account.String(),
				Collateral:          sdk.NewCoin(denom, sdk.ZeroInt()),
				MerDebt:             sdk.NewCoin(merlion.MicroUSDDenom, sdk.ZeroInt()),
				MerByLion:           sdk.NewCoin(merlion.MicroUSDDenom, sdk.ZeroInt()),
				LionBurned:          sdk.NewCoin(merlion.AttoLionDenom, sdk.ZeroInt()),
				LastInterest:        sdk.NewCoin(merlion.MicroUSDDenom, sdk.ZeroInt()),
				LastSettlementBlock: ctx.BlockHeight(),
			}
		} else {
			err = sdkerrors.Wrapf(types.ErrAccountNoCollateral, "account has no collateral: %s", denom)
			return
		}
	}
	return
}

func (k Keeper) totalBackingInUSD(ctx sdk.Context) (sdk.Coin, error) {
	totalBackingValue := sdk.ZeroDec()
	for _, pool := range k.GetAllPoolBacking(ctx) {
		// get price in usd
		backingPrice, err := k.oracleKeeper.GetExchangeRate(ctx, pool.Backing.Denom)
		if err != nil {
			return sdk.Coin{}, err
		}
		totalBackingValue = totalBackingValue.Add(pool.Backing.Amount.ToDec().Mul(backingPrice))
	}
	return sdk.NewCoin(merlion.MicroUSDDenom, totalBackingValue.TruncateInt()), nil
}

func settleInterestFee(ctx sdk.Context, acc *types.AccountCollateral, pool *types.PoolCollateral, total *types.TotalCollateral, apr *sdk.Dec) {
	if apr != nil {
		period := ctx.BlockHeight() - acc.LastSettlementBlock
		if period == 0 {
			// short circuit
			return
		}
		// principal debt, excluding interest debt
		principalDebt := acc.MerDebt.Sub(acc.LastInterest)
		interestOfPeriod := principalDebt.Amount.ToDec().Mul(*apr).MulInt64(period).QuoInt64(int64(merlion.BlocksPerYear)).RoundInt()
		// update remaining interest accumulation
		acc.LastInterest = acc.LastInterest.AddAmount(interestOfPeriod)
		// update debt
		acc.MerDebt = acc.MerDebt.AddAmount(interestOfPeriod)
		pool.MerDebt = pool.MerDebt.AddAmount(interestOfPeriod)
		total.MerDebt = total.MerDebt.AddAmount(interestOfPeriod)
	}
	// update settlement block
	acc.LastSettlementBlock = ctx.BlockHeight()
}

func computeFee(coin sdk.Coin, rate *sdk.Dec) sdk.Coin {
	amt := sdk.ZeroInt()
	if rate != nil {
		amt = coin.Amount.ToDec().Mul(*rate).TruncateInt()
	}
	return sdk.NewCoin(coin.Denom, amt)
}

func maxLoanToValueForAccount(acc *types.AccountCollateral, collateralParams *types.CollateralRiskParams) sdk.Dec {
	merDue := acc.MerDebt.Add(acc.MerByLion)
	catalyticRatio := acc.MerByLion.Amount.ToDec().QuoInt(merDue.Amount).Quo(*collateralParams.CatalyticLionRatio)
	if catalyticRatio.GT(sdk.OneDec()) {
		catalyticRatio = sdk.OneDec()
	}
	return collateralParams.BasicLoanToValue.Add(collateralParams.LoanToValue.Sub(*collateralParams.BasicLoanToValue).Mul(catalyticRatio))
}

func getSenderReceiver(senderStr, toStr string) (sender sdk.AccAddress, receiver sdk.AccAddress, err error) {
	sender, err = sdk.AccAddressFromBech32(senderStr)
	if err != nil {
		return
	}
	receiver = sender
	if len(toStr) > 0 {
		// user specifies receiver
		receiver, err = sdk.AccAddressFromBech32(toStr)
		if err != nil {
			return
		}
	}
	return
}
