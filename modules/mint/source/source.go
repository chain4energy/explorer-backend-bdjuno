package source

import (
	cfeminter "github.com/chain4energy/c4e-chain/v2/x/cfeminter/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Source interface {
	GetInflation() (sdk.Dec, error)
	Params(height int64) (cfeminter.Params, error)
}
