package crypto_test

import (
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/crypto"
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

	It("should fail when address cannot be decoded", func() {
		e := crypto.ValidateSig(key.PubKey(), "notanaddress!", "")
		Expect(e).To(Not(BeNil()))
		Expect(e.Code()).To(Equal(sdk.CodeInvalidAddress))
	})

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
		addrString := addr.String()
		sig, err := signer.PrivKey.Sign([]byte(addrString))
		Expect(err).To(BeNil())

		sigEncoded := hex.EncodeToString(sig)
		err = crypto.ValidateSig(key.PubKey(), addrString, sigEncoded)
		Expect(err).To(BeNil())
	})

	FIt("should succeed when signature is valid", func() {
		var priv secp.PrivKeySecp256k1
		privHex := "6453c9b244aa569dc3b769462c2192e0f6e7532c353fe139e0479986d387acfa"
		tmp := []byte(privHex)
		copy(priv[:], tmp)
		addrString := sdk.AccAddress(priv.PubKey().Address()).String()
		//hash := hex.EncodeToString([]byte(addrString))
		//fmt.Println(hash)
		sig, err := priv.Sign([]byte(addrString))
		Expect(err).To(BeNil())
		sigEncoded := hex.EncodeToString(sig)
		err = crypto.ValidateSig(priv.PubKey(), addrString, sigEncoded)
		Expect(err).To(BeNil())

		fromJs := "050f543a67f78059e3ff022230aa0a9df86b40ce31a77f85b81ccf2d55e0ca7b6455a7ea366a0217aafc7293cafaf7104f2980a52feae325e0155566b19162f8"
		err = crypto.ValidateSig(priv.PubKey(), addrString, fromJs)
		Expect(err).To(BeNil())
	})

})
