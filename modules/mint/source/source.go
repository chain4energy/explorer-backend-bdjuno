package source

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	mintertypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

type Source interface {
	GetInflation() (sdk.Dec, error)
	Params(height int64) (mintertypes.Params, error)
}
