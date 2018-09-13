package blockchain

import (
	"encoding/binary"
	"encoding/hex"
	"github.com/aergoio/aergo/account/key"
	"github.com/aergoio/aergo/types"
	"github.com/btcsuite/btcd/btcec"
	"testing"
)

const (
	maxAccount   = 100
	maxRecipient = 100
	maxTx        = 10000
)

var (
	accs      [maxAccount][]byte
	sign      [maxAccount]*btcec.PrivateKey
	recipient [maxRecipient][]byte
	txs       [maxTx]*types.Tx
)

func TestTXs(t *testing.T) {
}

func _itobU32(argv uint32) []byte {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, argv)
	return bs
}

func getAccount(tx *types.Tx) string {
	return hex.EncodeToString(tx.GetBody().GetAccount())
}

func beforeTest() error {
	for i := 0; i < maxAccount; i++ {
		privkey, err := btcec.NewPrivateKey(btcec.S256())
		if err != nil {
			return err
		}
		//gen new address
		accs[i] = key.GenerateAddress(&privkey.PublicKey)
		sign[i] = privkey
		recipient[i] = _itobU32(uint32(i))
	}

	// init Tx
	for i := 0; i < maxTx; i++ {
		txs[i] = nil
	}

	// gen Tx
	accCount := 100
	txCount := 100
	nonce := make([]uint64, txCount)
	for i := 0; i < txCount; i++ {
		nonce[i] = uint64(i + 1)
	}
	for i := 0; i < accCount; i++ {
		for j := 0; j < txCount; j++ {
			tmp := genTx(i, j, nonce[j], uint64(i+1))
			txs[i*txCount+j] = tmp
		}
	}

	return nil
}

func afterTest() {

}

func genTx(acc int, rec int, nonce uint64, amount uint64) *types.Tx {
	tx := types.Tx{
		Body: &types.TxBody{
			Nonce:     nonce,
			Account:   accs[acc],
			Recipient: recipient[rec],
			Amount:    amount,
		},
	}
	//tx.Hash = tx.CalculateTxHash()
	key.SignTx(&tx, sign[acc])
	return &tx
}

func TestInvalidTransactions(t *testing.T) {
	/*
		beforeTest(t)
		defer afterTest()
		tx := genTx(0, 1, 1, 1)

		err := key.SignTx(tx, sign[1])
		if err == nil {
			t.Errorf("put invalid tx should be failed")
		}

		tx.Body.Sign = nil
		tx.Hash = tx.CalculateTxHash()

		if err == nil {
			t.Errorf("put invalid tx should be failed")
		}
	*/
}

// gen sequential transactions
// bench
func TestVerifyValidTxs(t *testing.T) {
	t.Log("TestVerify")
	beforeTest()
	defer afterTest()

	txslice := make([]*types.Tx, 0)
	for _, tx := range txs {
		txslice = append(txslice, tx)
	}

	verifier := NewSignVerifier(DefaultVerifierCnt)

	failed, errors := verifier.VerifyTxs(&types.TxList{Txs: txslice})
	if failed {
		for i, error := range errors {
			if error != nil {
				t.Errorf("failed tx %d:%s", i, error.Error())
			}
		}
	}

	t.Log("block verify succeed")
}

func BenchmarkVerify10000tx(b *testing.B) {
	b.Log("TestVerify")
	beforeTest()
	defer afterTest()

	txslice := make([]*types.Tx, 0)
	for _, tx := range txs {
		txslice = append(txslice, tx)
	}

	verifier := NewSignVerifier(DefaultVerifierCnt)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		failed, errors := verifier.VerifyTxs(&types.TxList{Txs: txslice})
		if failed {
			for i, error := range errors {
				if error != nil {
					b.Errorf("failed tx %d:%s", i, error.Error())
				}
			}
		}
	}

	b.Log("block verify succeed")
}

func BenchmarkVerify10000txSerial(b *testing.B) {
	b.Log("TestVerify")
	beforeTest()
	defer afterTest()

	txslice := make([]*types.Tx, 0)
	for _, tx := range txs {
		txslice = append(txslice, tx)
	}

	verifier := NewSignVerifier(DefaultVerifierCnt)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		failed, errors := verifier.verifyTxsInplace(&types.TxList{Txs: txslice})
		if failed {
			for i, error := range errors {
				if error != nil {
					b.Errorf("failed tx %d:%s", i, error.Error())
				}
			}
		}
	}

	b.Log("block verify succeed")
}
