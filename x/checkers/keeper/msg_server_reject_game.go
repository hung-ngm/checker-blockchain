package keeper

import (
	"context"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/minhhung123/checkers/x/checkers/types"
)

func (k msgServer) RejectGame(goCtx context.Context, msg *types.MsgRejectGame) (*types.MsgRejectGameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	storedGame, found := k.Keeper.GetStoredGame(ctx, msg.IdValue)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrGameNotFound, "game not found: %s", msg.IdValue)
	}

	// CHeck if the player is expected and did they already play
	if strings.Compare(storedGame.Red, msg.Creator) == 0 {
		if 1 < storedGame.MoveCount {
			return nil, types.ErrRedAlreadyPlayed
		}
	} else if strings.Compare(storedGame.Black, msg.Creator) == 0 {
		if 0 < storedGame.MoveCount { // Black plays first
			return nil, types.ErrBlackAlreadyPlayed
		}
	} else {
		return nil, types.ErrCreatorNotPlayer
	}

	// Remove the game from the FIFO
	nextGame, found := k.Keeper.GetNextGame(ctx)
	if !found {
		panic("NextGame not found")
	}
	k.Keeper.RemoveFromFifo(ctx, &storedGame, &nextGame)

	// Remove the game completely as it is not intersting to keep it
	// Remove the game
	k.Keeper.RemoveStoredGame(ctx, msg.IdValue)
	k.Keeper.SetNextGame(ctx, nextGame)

	// Emit the Event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage, 
			sdk.NewAttribute(sdk.AttributeKeyModule, "checkers"),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.RejectGameEventKey),
			sdk.NewAttribute(types.RejectGameEventCreator, msg.Creator),
			sdk.NewAttribute(types.RejectGameEventIdValue, msg.IdValue),
		),
	)

	return &types.MsgRejectGameResponse{}, nil
}
