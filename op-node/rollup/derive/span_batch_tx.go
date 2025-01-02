package derive

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/holiman/uint256"
)

type spanBatchTxData interface {
	txType() byte // returns the type ID
}

type spanBatchTx struct {
	inner spanBatchTxData
}

type spanBatchLegacyTxData struct {
	Value    *big.Int // wei amount
	GasPrice *big.Int // wei per gas
	Data     []byte
}

func (txData *spanBatchLegacyTxData) txType() byte { return types.LegacyTxType }

type spanBatchAccessListTxData struct {
	Value      *big.Int // wei amount
	GasPrice   *big.Int // wei per gas
	Data       []byte
	AccessList types.AccessList // EIP-2930 access list
}

func (txData *spanBatchAccessListTxData) txType() byte { return types.AccessListTxType }

type spanBatchDynamicFeeTxData struct {
	Value      *big.Int
	GasTipCap  *big.Int // a.k.a. maxPriorityFeePerGas
	GasFeeCap  *big.Int // a.k.a. maxFeePerGas
	Data       []byte
	AccessList types.AccessList
}

func (txData *spanBatchDynamicFeeTxData) txType() byte { return types.DynamicFeeTxType }

type spanBatchBlobTxData struct {
	Value      *uint256.Int
	GasTipCap  *uint256.Int // a.k.a. maxPriorityFeePerGas
	GasFeeCap  *uint256.Int // a.k.a. maxFeePerGas
	Data       []byte
	AccessList types.AccessList
	BlobFeeCap *uint256.Int // a.k.a. maxFeePerBlobGas
	BlobHashes []common.Hash
}

func (txData *spanBatchBlobTxData) txType() byte { return types.BlobTxType }

// Type returns the transaction type.
func (tx *spanBatchTx) Type() uint8 {
	return tx.inner.txType()
}

// encodeTyped writes the canonical encoding of a typed transaction to w.
func (tx *spanBatchTx) encodeTyped(w *bytes.Buffer) error {
	w.WriteByte(tx.Type())
	return rlp.Encode(w, tx.inner)
}

// MarshalBinary returns the canonical encoding of the transaction.
// For legacy transactions, it returns the RLP encoding. For EIP-2718 typed
// transactions, it returns the type and payload.
func (tx *spanBatchTx) MarshalBinary() ([]byte, error) {
	if tx.Type() == types.LegacyTxType {
		return rlp.EncodeToBytes(tx.inner)
	}
	var buf bytes.Buffer
	err := tx.encodeTyped(&buf)
	return buf.Bytes(), err
}

// setDecoded sets the inner transaction after decoding.
func (tx *spanBatchTx) setDecoded(inner spanBatchTxData, size uint64) {
	tx.inner = inner
}

// decodeTyped decodes a typed transaction from the canonical format.
func (tx *spanBatchTx) decodeTyped(b []byte) (spanBatchTxData, error) {
	if len(b) <= 1 {
		return nil, fmt.Errorf("failed to decode span batch: %w", ErrTypedTxTooShort)
	}
	switch b[0] {
	case types.AccessListTxType:
		var inner spanBatchAccessListTxData
		err := rlp.DecodeBytes(b[1:], &inner)
		if err != nil {
			return nil, fmt.Errorf("failed to decode spanBatchAccessListTxData: %w", err)
		}
		return &inner, nil
	case types.DynamicFeeTxType:
		var inner spanBatchDynamicFeeTxData
		err := rlp.DecodeBytes(b[1:], &inner)
		if err != nil {
			return nil, fmt.Errorf("failed to decode spanBatchDynamicFeeTxData: %w", err)
		}
		return &inner, nil
	case types.BlobTxType:
		var inner spanBatchBlobTxData
		err := rlp.DecodeBytes(b[1:], &inner)
		if err != nil {
			return nil, fmt.Errorf("failed to decode spanBatchBlobTxData: %w", err)
		}
		return &inner, nil
	default:
		return nil, types.ErrTxTypeNotSupported
	}
}

// UnmarshalBinary decodes the canonical encoding of transactions.
// It supports legacy RLP transactions and EIP2718 typed transactions.
func (tx *spanBatchTx) UnmarshalBinary(b []byte) error {
	if len(b) > 0 && b[0] > 0x7f {
		// It's a legacy transaction.
		var data spanBatchLegacyTxData
		err := rlp.DecodeBytes(b, &data)
		if err != nil {
			return fmt.Errorf("failed to decode spanBatchLegacyTxData: %w", err)
		}
		tx.setDecoded(&data, uint64(len(b)))
		return nil
	}
	// It's an EIP2718 typed transaction envelope.
	inner, err := tx.decodeTyped(b)
	if err != nil {
		return err
	}
	tx.setDecoded(inner, uint64(len(b)))
	return nil
}

// convertToFullTx takes values and convert spanBatchTx to types.Transaction
func (tx *spanBatchTx) convertToFullTx(nonce, gas uint64, to *common.Address, chainID, V, R, S *big.Int) (*types.Transaction, error) {
	var inner types.TxData
	switch tx.Type() {
	case types.LegacyTxType:
		batchTxInner := tx.inner.(*spanBatchLegacyTxData)
		inner = &types.LegacyTx{
			Nonce:    nonce,
			GasPrice: batchTxInner.GasPrice,
			Gas:      gas,
			To:       to,
			Value:    batchTxInner.Value,
			Data:     batchTxInner.Data,
			V:        V,
			R:        R,
			S:        S,
		}
	case types.AccessListTxType:
		batchTxInner := tx.inner.(*spanBatchAccessListTxData)
		inner = &types.AccessListTx{
			ChainID:    chainID,
			Nonce:      nonce,
			GasPrice:   batchTxInner.GasPrice,
			Gas:        gas,
			To:         to,
			Value:      batchTxInner.Value,
			Data:       batchTxInner.Data,
			AccessList: batchTxInner.AccessList,
			V:          V,
			R:          R,
			S:          S,
		}
	case types.DynamicFeeTxType:
		batchTxInner := tx.inner.(*spanBatchDynamicFeeTxData)
		inner = &types.DynamicFeeTx{
			ChainID:    chainID,
			Nonce:      nonce,
			GasTipCap:  batchTxInner.GasTipCap,
			GasFeeCap:  batchTxInner.GasFeeCap,
			Gas:        gas,
			To:         to,
			Value:      batchTxInner.Value,
			Data:       batchTxInner.Data,
			AccessList: batchTxInner.AccessList,
			V:          V,
			R:          R,
			S:          S,
		}
	case types.BlobTxType:
		if to == nil {
			return nil, fmt.Errorf("invalid blob tx: to can't be nil")
		}
		VU256, overflow := uint256.FromBig(V)
		if overflow {
			return nil, fmt.Errorf("invalid blob tx: V overflow:%v", V)
		}
		RU256, overflow := uint256.FromBig(R)
		if overflow {
			return nil, fmt.Errorf("invalid blob tx: R overflow:%v", R)
		}
		SU256, overflow := uint256.FromBig(S)
		if overflow {
			return nil, fmt.Errorf("invalid blob tx: S overflow:%v", R)
		}
		batchTxInner := tx.inner.(*spanBatchBlobTxData)
		inner = &types.BlobTx{
			ChainID:    uint256.MustFromBig(chainID),
			Nonce:      nonce,
			GasTipCap:  batchTxInner.GasTipCap,
			GasFeeCap:  batchTxInner.GasFeeCap,
			Gas:        gas,
			To:         *to,
			Value:      batchTxInner.Value,
			Data:       batchTxInner.Data,
			AccessList: batchTxInner.AccessList,
			BlobFeeCap: batchTxInner.BlobFeeCap,
			BlobHashes: batchTxInner.BlobHashes,
			V:          VU256,
			R:          RU256,
			S:          SU256,
		}

	default:
		return nil, fmt.Errorf("invalid tx type: %d", tx.Type())
	}
	return types.NewTx(inner), nil
}

// newSpanBatchTx converts types.Transaction to spanBatchTx
func newSpanBatchTx(tx *types.Transaction) (*spanBatchTx, error) {
	var inner spanBatchTxData
	switch tx.Type() {
	case types.LegacyTxType:
		inner = &spanBatchLegacyTxData{
			GasPrice: tx.GasPrice(),
			Value:    tx.Value(),
			Data:     tx.Data(),
		}
	case types.AccessListTxType:
		inner = &spanBatchAccessListTxData{
			GasPrice:   tx.GasPrice(),
			Value:      tx.Value(),
			Data:       tx.Data(),
			AccessList: tx.AccessList(),
		}
	case types.DynamicFeeTxType:
		inner = &spanBatchDynamicFeeTxData{
			GasTipCap:  tx.GasTipCap(),
			GasFeeCap:  tx.GasFeeCap(),
			Value:      tx.Value(),
			Data:       tx.Data(),
			AccessList: tx.AccessList(),
		}
	case types.BlobTxType:
		gasTipCap, overflow := uint256.FromBig(tx.GasTipCap())
		if overflow {
			return nil, fmt.Errorf("tx.GasTipCap() overflow: %v", tx.GasTipCap())
		}
		gasFeeCap, overflow := uint256.FromBig(tx.GasFeeCap())
		if overflow {
			return nil, fmt.Errorf("tx.GasFeeCap() overflow: %v", tx.GasFeeCap())
		}
		value, overflow := uint256.FromBig(tx.Value())
		if overflow {
			return nil, fmt.Errorf("tx.Value() overflow: %v", tx.Value())
		}
		blobFeeCap, overflow := uint256.FromBig(tx.BlobGasFeeCap())
		if overflow {
			return nil, fmt.Errorf("tx.BlobGasFeeCap() overflow: %v", tx.BlobGasFeeCap())
		}
		inner = &spanBatchBlobTxData{
			Value:      value,
			GasTipCap:  gasTipCap,
			GasFeeCap:  gasFeeCap,
			Data:       tx.Data(),
			AccessList: tx.AccessList(),
			BlobFeeCap: blobFeeCap,
			BlobHashes: tx.BlobHashes(),
		}
	default:
		return nil, fmt.Errorf("invalid tx type: %d", tx.Type())
	}
	return &spanBatchTx{inner: inner}, nil
}
