package handler_test

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
	RunSpecs(t, "Handler Test Suite")
}

var simApp *sim.LongyApp
var ctx sdk.Context
var keeper longy.Keeper
var handler sdk.Handler

// BeforeTestRun sets up common data needed by all handler tests.
func BeforeTestRun() {
	simApp, ctx = sim.CreateTestApp(true)
	keeper = simApp.LongyKeeper
	handler = longy.NewHandler(keeper)
}
