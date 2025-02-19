package libwallet

import (
	"crypto/sha256"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
)

func createNonNativeSegwitRedeemScript(witnessScript []byte) ([]byte, error) {
	witnessScriptHash := sha256.Sum256(witnessScript)

	builder := txscript.NewScriptBuilder()
	builder.AddInt64(0)
	builder.AddData(witnessScriptHash[:])

	return builder.Script()
}

func signNonNativeSegwitInput(input Input, index int, tx *wire.MsgTx, privateKey *HDPrivateKey,
	redeemScript, witnessScript []byte) ([]byte, error) {

	txInput := tx.TxIn[index]

	builder := txscript.NewScriptBuilder()
	builder.AddData(redeemScript)
	script, err := builder.Script()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to generate signing script")
	}
	txInput.SignatureScript = script

	privKey, err := privateKey.key.ECPrivKey()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to produce EC priv key for signing")
	}

	sigHashes := txscript.NewTxSigHashes(tx)
	sig, err := txscript.RawTxInWitnessSignature(tx, sigHashes, index, input.OutPoint().Amount(), witnessScript, txscript.SigHashAll, privKey)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to sign V3 input")
	}

	return sig, nil
}
