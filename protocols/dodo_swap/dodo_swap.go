// Package dodo_swap provides swap event parsing for DODOSwap protocols.
package dodo_swap

import (
	"math/big"

	"github.com/48Club/bscexorcist/types"
	"github.com/48Club/bscexorcist/utils"
	"github.com/ethereum/go-ethereum/common"
	eth "github.com/ethereum/go-ethereum/core/types"
)

// SwapEventSignature for DodoSwap
var SwapEventSignature = common.HexToHash("0xc2c0245e056d5fb095f04cd6373bc770802ebd1e6c918eb78fdef843cdb37b0f")

// DODOSwap implements SwapEvent for DODOSwap protocol.
type DODOSwap struct {
	poolID     common.Address
	tokenFrom  common.Address
	tokenTo    common.Address
	amountFrom *big.Int
	amountTo   *big.Int
}

// isTokenAFirst checks if token A is less than token B.
func isTokenAFirst(tkA, tkB common.Address) bool {
	tokenARep := utils.BigIntFromBytes(tkA.Bytes())
	tokenBRep := utils.BigIntFromBytes(tkB.Bytes())
	return tokenARep.Cmp(tokenBRep) < 0
}

// calcPoolID returns a pseudo-address derived from the first 10 bytes of each sorted tokenFrom/tokenTo.
func calcPoolID(tkA, tkB common.Address) (pool common.Address) {
	tokenARep := utils.BigIntFromBytes(tkA.Bytes())
	tokenBRep := utils.BigIntFromBytes(tkB.Bytes())
	if tokenARep.Cmp(tokenBRep) < 0 {
		copy(pool[:], tkA.Bytes()[:10])
		copy(pool[10:], tkB.Bytes()[:10])
	} else {
		copy(pool[:], tkB.Bytes()[:10])
		copy(pool[10:], tkA.Bytes()[:10])
	}
	return pool
}

// PairID returns a pseudo-address derived from the first 10 bytes of each token in the pair.
func (s *DODOSwap) PairID() types.Addresses {
	return types.AddressesB20(s.poolID)
}

// IsToken0To1 returns true if the swap direction is token0 -> token1.
func (s *DODOSwap) IsToken0To1() bool {
	return isTokenAFirst(s.tokenFrom, s.tokenTo)
}

// AmountIn returns the input amount for the swap.
func (s *DODOSwap) AmountIn() *big.Int {
	if isTokenAFirst(s.tokenFrom, s.tokenTo) {
		// tokenFrom is token0, so amountFrom is the input amount
		return utils.BigIntFromPointer(s.amountFrom)
	}
	return utils.BigIntFromPointer(s.amountTo)
}

// AmountOut returns the output amount for the swap.
func (s *DODOSwap) AmountOut() *big.Int {
	if isTokenAFirst(s.tokenFrom, s.tokenTo) {
		// tokenFrom is token0, so amountFrom is the input amount
		return utils.BigIntFromPointer(s.amountTo)
	}
	return utils.BigIntFromPointer(s.amountFrom)
}

// ParseSwap parses a DODOSwap log into a DODOSwap struct.
// Returns nil if the log is not a valid swap event.
func ParseSwap(log *eth.Log) *DODOSwap {
	if len(log.Topics) != 1 || len(log.Data) < 128 {
		return nil
	}

	fromToken := common.BytesToAddress(log.Data[:32])
	toToken := common.BytesToAddress(log.Data[32:64])

	return &DODOSwap{
		poolID:     calcPoolID(fromToken, toToken),
		tokenFrom:  fromToken,
		tokenTo:    toToken,
		amountFrom: utils.BigIntFromBytes(log.Data[64:96]),
		amountTo:   utils.BigIntFromBytes(log.Data[96:128]),
	}
}
