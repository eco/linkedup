package longy

var (
//_ module.AppModule      = AppModule{}
//_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic is
type AppModuleBasic struct{}

// AppModule is
// nolint: structcheck, unused
type AppModule struct {
	AppModuleBasic
	attendeeKeeper AttendeeKeeper
}
