package types

import (
	mintertypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

// MintParams represents the x/mint parameters
type MintParams struct {
	mintertypes.Params
	Height int64
}

// NewMintParams allows to build a new MintParams instance
func NewMintParams(params mintertypes.Params, height int64) *MintParams {
	return &MintParams{
		Params: params,
		Height: height,
	}
}
