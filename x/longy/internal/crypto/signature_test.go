package crypto_test

import (
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/internal/crypto"
	"github.com/eco/longy/x/longy/internal/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	tmCrypto "github.com/tendermint/tendermint/crypto"
	secp "github.com/tendermint/tendermint/crypto/secp256k1"
)

var _ = Describe("Signature Keeper Tests", func() {
	var key tmCrypto.PrivKey
	var addr sdk.AccAddress
	BeforeEach(func() {
		BeforeTestRun()

		key = secp.GenPrivKeySecp256k1([]byte("master"))
		addr = sdk.AccAddress(key.PubKey().Address())
	})

	It("should fail when public key is empty", func() {
		e := crypto.ValidateSig(nil, "", "")
		Expect(e).To(Not(BeNil()))
		Expect(e.Code()).To(Equal(types.InvalidPublicKey))
	})
	//
	//It("should fail when address cannot be decoded", func() {
	//	e := crypto.ValidateSig(key.PubKey(), "notanaddress!", "")
	//	Expect(e).To(Not(BeNil()))
	//	Expect(e.Code()).To(Equal(sdk.CodeInvalidAddress))
	//})

	It("should fail when signature does not match the key", func() {
		e := crypto.ValidateSig(key.PubKey(), addr.String(), "fakesig")
		Expect(e).To(Not(BeNil()))
		Expect(e.Code()).To(Equal(types.SigDecodeError))
	})

	It("should fail when signature does not match the key", func() {
		sigEncoded := hex.EncodeToString([]byte("fakesig"))
		e := crypto.ValidateSig(key.PubKey(), addr.String(), sigEncoded)
		Expect(e).To(Not(BeNil()))
		Expect(e.Code()).To(Equal(types.InvalidSignature))
	})

	It("should succeed when signature is valid", func() {
		signer := crypto.NewSigner(addr, key)
		hash, err := crypto.Hash(addr)
		Expect(err).To(BeNil())

		sig, err := signer.PrivKey.Sign(hash)
		Expect(err).To(BeNil())

		sigEncoded := hex.EncodeToString(sig)
		err = crypto.ValidateSig(key.PubKey(), addr.String(), sigEncoded)
		Expect(err).To(BeNil())
	})
})
