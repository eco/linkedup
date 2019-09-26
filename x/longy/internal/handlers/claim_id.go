package handlers

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy"
	"github.com/eco/longy/x/longy/internal/types"
)

//ClaimIDHandler is a struct that wraps the handler params in order to process identity publishes
type ClaimIDHandler struct {
	*BaseHandler
	msg *types.MsgClaimID
}

//NewClaimIDHandler initializes a new handler and returns a pointer to it
func NewClaimIDHandler(ctx *sdk.Context, keeper *longy.Keeper,
	msg *types.MsgClaimID) *ClaimIDHandler {
	return &ClaimIDHandler{
		BaseHandler: NewBaseHandler(ctx, keeper),
		msg:         msg,
	}
}

// handleMsgClaimID processes MsgClaimID in order to associate an address with an id
// nolint: unparam
func handleClaimIDMsg(ctx *sdk.Context, keeper *longy.Keeper, msg *types.MsgClaimID) sdk.Result {
	return NewClaimIDHandler(ctx, keeper, msg).handleMsgClaimID()
}

func (h *ClaimIDHandler) handleMsgClaimID() sdk.Result {

	return sdk.Result{}
}
