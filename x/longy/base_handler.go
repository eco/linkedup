package longy

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//BaseHandler is a struct that supplies data and methods common across all message handlers.
type BaseHandler struct {
	ctx    sdk.Context
	keeper *Keeper
}

//NewBaseHandler initializes a new handler and returns a pointer to it
// nolint: gocritic
func NewBaseHandler(ctx sdk.Context, keeper *Keeper) *BaseHandler {
	return &BaseHandler{
		ctx:    ctx,
		keeper: keeper,
	}
}

func (h *ClaimIDHandler) isSuperUser(acc sdk.AccAddress) bool {
	//todo make superuser and update func
	return true
}
