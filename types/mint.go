package types

import (
	cfemintertypes "github.com/chain4energy/c4e-chain/x/cfeminter/types"
)

// MintParams represents the x/mint parameters
type MintParams struct {
	cfemintertypes.Params
	Height int64
}

// NewMintParams allows to build a new MintParams instance
func NewMintParams(params cfemintertypes.Params, height int64) *MintParams {
	return &MintParams{
		Params: params,
		Height: height,
	}
}
