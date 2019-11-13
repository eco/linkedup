package crypto_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	tmCrypto "github.com/tendermint/tendermint/crypto"
	secp "github.com/tendermint/tendermint/crypto/secp256k1"
)

var _ = Describe("Signer Keeper Tests", func() {
	var key tmCrypto.PrivKey
	BeforeEach(func() {
		BeforeTestRun()

		key = secp.GenPrivKeySecp256k1([]byte("master"))
	})

	It("should fail to set an empty AccAddress", func() {
		addr := sdk.AccAddress(key.PubKey().Address())
		signer := crypto.NewSigner(addr, key)

		msg := []byte("asdfasdfas")
		sig, err := signer.PrivKey.Sign(msg)
		Expect(err).To(BeNil())
		same := signer.PrivKey.PubKey().VerifyBytes(msg, sig)
		Expect(same).To(BeTrue())

	})
})
