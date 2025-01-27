package maker

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/merlion-zone/merlion/x/maker/keeper"
	"github.com/merlion-zone/merlion/x/maker/types"
)

// NewHandler ...
func NewHandler(k keeper.Keeper) sdk.Handler {
	// this line is used by starport scaffolding # handler/msgServer

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		// this line is used by starport scaffolding # 1
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func NewMakerProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.RegisterBackingProposal:
			return keeper.HandleRegisterBackingProposal(ctx, k, c)
		case *types.RegisterCollateralProposal:
			return keeper.HandleRegisterCollateralProposal(ctx, k, c)
		case *types.SetBackingRiskParamsProposal:
			return keeper.HandleSetBackingRiskParamsProposal(ctx, k, c)
		case *types.SetCollateralRiskParamsProposal:
			return keeper.HandleSetCollateralRiskParamsProposal(ctx, k, c)
		case *types.BatchSetBackingRiskParamsProposal:
			return keeper.HandleBatchSetBackingRiskParamsProposal(ctx, k, c)
		case *types.BatchSetCollateralRiskParamsProposal:
			return keeper.HandleBatchSetCollateralRiskParamsProposal(ctx, k, c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s proposal content type: %T", types.ModuleName, c)
		}
	}
}
