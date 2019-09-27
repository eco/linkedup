package longy

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/errors"
	"github.com/eco/longy/x/longy/internal/types"
)

//ClaimIDHandler is a struct that wraps the handler params in order to process identity publishes
type ClaimIDHandler struct {
	*BaseHandler
	msg *types.MsgClaimID
}

//NewClaimIDHandler initializes a new handler and returns a pointer to it
func NewClaimIDHandler(ctx *sdk.Context, keeper *Keeper,
	msg *types.MsgClaimID) *ClaimIDHandler {
	return &ClaimIDHandler{
		BaseHandler: NewBaseHandler(ctx, keeper),
		msg:         msg,
	}
}

// HandleClaimIDMsg processes MsgClaimID in order to associate an address with an id
// nolint: unparam
func HandleClaimIDMsg(ctx *sdk.Context, keeper *Keeper, msg *types.MsgClaimID) sdk.Result {
	return NewClaimIDHandler(ctx, keeper, msg).handleMsgClaimID()
}

func (h *ClaimIDHandler) handleMsgClaimID() sdk.Result {
	if !h.isSuperUser(h.msg.Sender) {
		return errors.ErrInsufficientPrivileges("Insufficient privilege to make this call").Result()
	}

	attendee, err := h.getAttendee(h.msg.ID)
	if err != nil {
		return err.Result()
	}

	attendee.Address = h.msg.Address
	h.setAttendee(attendee)

	return sdk.Result{}
}
