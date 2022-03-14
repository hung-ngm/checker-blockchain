package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	rules "github.com/minhhung123/checkers/x/checkers/rules"
	"github.com/minhhung123/checkers/x/checkers/types"
)

// Return error if the player does not have enough funds
func (k *Keeper) CollectWager(ctx sdk.Context, storedGame *types.StoredGame) error {
	// Make the player pay the wager at the beginning
	if storedGame.MoveCount == 0 {
		// Black plays first
		black, err := storedGame.GetBlackAddress()
		if err != nil {
			panic(err.Error())

		}
		err = k.bank.SendCoinsFromAccountToModule(ctx, black, types.ModuleName, sdk.NewCoins(storedGame.GetWagerCoin()))
		if err != nil {
			return sdkerrors.Wrapf(err, types.ErrBlackCannotPay.Error())
		}
	} else if storedGame.MoveCount == 1 {
		// Red plays second
		red, err := storedGame.GetRedAddress()
		if err != nil {
			panic(err.Error())
		}
		err = k.bank.SendCoinsFromAccountToModule(ctx, red, types.ModuleName, sdk.NewCoins(storedGame.GetWagerCoin()))
		if err != nil {
			return sdkerrors.Wrapf(err, types.ErrRedCannotPay.Error())
		}
	}
	return nil
}

func (k *Keeper) MustPayWinnings(ctx sdk.Context, storedGame *types.StoredGame) {
	winnerAddress, found, err := storedGame.GetWinnerAddress()
	if err != nil {
		panic(err.Error())
	}
	if !found {
		panic(fmt.Sprintf(types.ErrCannotFindWinnerByColor.Error(), storedGame.Winner))
	}
	winnings := storedGame.GetWagerCoin()
	if storedGame.MoveCount == 0 {
		panic(types.ErrNothingToPay.Error())
	} else if 1 < storedGame.MoveCount {
		winnings = winnings.Add(winnings)
	}
	err = k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, winnerAddress, sdk.NewCoins(winnings))
	if err != nil {
		panic(types.ErrCannotPayWinnings.Error())
	}

}

// Game must be in a state where it can be refunded
func (k *Keeper) MustRefundWager(ctx sdk.Context, storedGame *types.StoredGame) {
	// Refund wager to black player if red rejects after black has played
	if storedGame.MoveCount == 1 {
		// Refund
		black, err := storedGame.GetBlackAddress()
		if err != nil {
			panic(err.Error())
		}
		err = k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, black, sdk.NewCoins(storedGame.GetWagerCoin()))
		if err != nil {
			panic(fmt.Sprintf(types.ErrCannotRefundWager.Error(), rules.BLACK_PLAYER.Color))
		}
	} else if storedGame.MoveCount == 0 {
		// Do nothing
	} else {
		// Draw Mechanism
		panic(fmt.Sprintf(types.ErrNotInRefundState.Error(), storedGame.MoveCount))
	}
}
