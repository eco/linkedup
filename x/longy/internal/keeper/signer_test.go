package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	k "github.com/eco/longy/x/longy/internal/keeper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tendermint/tendermint/crypto"
	tmcrypto "github.com/tendermint/tendermint/crypto/secp256k1"
)

var _ = Describe("Signer Keeper Tests", func() {
	var key crypto.PrivKey
	BeforeEach(func() {
		BeforeTestRun()

		key = tmcrypto.GenPrivKeySecp256k1([]byte("master"))
	})

	It("should fail to set an empty AccAddress", func() {
		addr := sdk.AccAddress(key.PubKey().Address())
		signer := k.NewSigner(addr, key)

		msg := []byte("asdfasdfas")
		sig, err := signer.PrivKey.Sign(msg)
		Expect(err).To(BeNil())
		same := signer.PrivKey.PubKey().VerifyBytes(msg, sig)
		Expect(same).To(BeTrue())

	})
})
