package smartaccount_test

import (
	"fmt"
	"os"

	tests "github.com/aura-nw/aura/tests"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

const (
	mockChainID = "aura-testnet"
	signMode    = signing.SignMode_SIGN_MODE_DIRECT
)

func makeMockAccount(keybase keyring.Keyring, uid string) (authtypes.AccountI, error) {
	record, _, err := keybase.NewMnemonic(
		uid,
		keyring.English,
		sdk.FullFundraiserPath,
		keyring.DefaultBIP39Passphrase,
		hd.Secp256k1,
	)
	if err != nil {
		return nil, err
	}

	pk := record.GetPubKey()
	if pk == nil {
		return nil, fmt.Errorf("pubkey error")
	}

	return authtypes.NewBaseAccount(pk.Address().Bytes(), pk, 0, 0), nil
}

type Signer struct {
	keyName        string             // the name of the key in the keyring
	acc            authtypes.AccountI // the account corresponding to the address
	overrideAccNum *uint64            // if not nil, will override the account number in the AccountI
	overrideSeq    *uint64            // if not nil, will override the sequence in the AccountI
}

func (s *Signer) AccountNumber() uint64 {
	if s.overrideAccNum != nil {
		return *s.overrideAccNum
	}

	return s.acc.GetAccountNumber()
}

func (s *Signer) Sequence() uint64 {
	if s.overrideSeq != nil {
		return *s.overrideSeq
	}

	return s.acc.GetSequence()
}

// Logics in this function is mostly copied from:
// cosmos/cosmos-sdk/x/auth/ante/testutil_test.go/CreateTestTx
func prepareTx(
	ctx sdk.Context, keybase keyring.Keyring,
	msgs []sdk.Msg, signers []Signer, chainID string,
	sign bool,
) (authsigning.Tx, error) {

	encoding := tests.MakeTestEncodingConfig()
	txBuilder := encoding.TxConfig.NewTxBuilder()

	if err := txBuilder.SetMsgs(msgs...); err != nil {
		return nil, err
	}

	// if the tx doesn't need to be signed, we can return here
	if !sign {
		return txBuilder.GetTx(), nil
	}

	// round 1: set empty signature
	sigs := []signing.SignatureV2{}

	for _, signer := range signers {
		sig := signing.SignatureV2{
			PubKey: signer.acc.GetPubKey(),
			Data: &signing.SingleSignatureData{
				SignMode:  signMode,
				Signature: nil, // empty
			},
			Sequence: signer.acc.GetSequence(),
		}

		sigs = append(sigs, sig)
	}

	fmt.Fprintln(os.Stdout, sigs)

	if err := txBuilder.SetSignatures(sigs...); err != nil {
		return nil, err
	}

	// round 2: sign the tx
	sigs = []signing.SignatureV2{}

	for _, signer := range signers {
		signerData := authsigning.SignerData{
			ChainID:       chainID,
			AccountNumber: signer.AccountNumber(),
			Sequence:      signer.Sequence(),
		}

		signBytes, err := encoding.TxConfig.SignModeHandler().GetSignBytes(signMode, signerData, txBuilder.GetTx())
		if err != nil {
			return nil, err
		}

		sigBytes, _, err := keybase.Sign(signer.keyName, signBytes)
		if err != nil {
			return nil, err
		}

		sig := signing.SignatureV2{
			PubKey: signer.acc.GetPubKey(),
			Data: &signing.SingleSignatureData{
				SignMode:  signMode,
				Signature: sigBytes,
			},
			Sequence: signer.Sequence(),
		}

		sigs = append(sigs, sig)
	}

	if err := txBuilder.SetSignatures(sigs...); err != nil {
		return nil, err
	}

	return txBuilder.GetTx(), nil
}
