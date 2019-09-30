package longy

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/internal/types"
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

//nolint: unused
func (h *BaseHandler) getAttendee(id string) (attendee *types.Attendee, err sdk.Error) {
	store := h.keeper.GetAttendeeStore()
	attendee, err = store.GetAttendee(h.ctx, []byte(id))
	return
}

func (h *BaseHandler) setAttendee(attendee *types.Attendee) {
	store := h.keeper.GetAttendeeStore()
	store.SetAttendee(h.ctx, attendee)
}
