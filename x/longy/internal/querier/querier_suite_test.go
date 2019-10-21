package querier_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy"
	"github.com/eco/longy/x/longy/sim"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestMonitor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Querier Test Suite")
}

var simApp *sim.LongyApp
var ctx sdk.Context
var keeper longy.Keeper
var querier sdk.Querier

const (
	qr1 = "1234"
)

var sender sdk.AccAddress

// BeforeTestRun sets up common data needed by all querier tests.
func BeforeTestRun() {
	simApp, ctx = sim.CreateTestApp(true)
	keeper = simApp.LongyKeeper
	querier = longy.NewQuerier(keeper)
}
