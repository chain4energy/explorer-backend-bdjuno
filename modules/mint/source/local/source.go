package local

import (
	"fmt"
	cfeminter "github.com/chain4energy/c4e-chain/x/cfeminter/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/juno/v3/node/local"

	mintsource "github.com/forbole/bdjuno/v3/modules/mint/source"
)

var (
	_ mintsource.Source = &Source{}
)

// Source implements mintsource.Source using a local node
type Source struct {
	*local.Source
	querier cfeminter.QueryServer
}

// NewSource returns a new Source instace
func NewSource(source *local.Source, querier cfeminter.QueryServer) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// GetInflation implements mintsource.Source
func (s Source) GetInflation(height int64) (sdk.Dec, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return sdk.Dec{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.querier.Inflation(sdk.WrapSDKContext(ctx), &cfeminter.QueryInflationRequest{})
	if err != nil {
		return sdk.Dec{}, err
	}

	return res.Inflation, nil
}

// Params implements mintsource.Source
func (s Source) Params(height int64) (cfeminter.Params, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return cfeminter.Params{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.querier.Params(sdk.WrapSDKContext(ctx), &cfeminter.QueryParamsRequest{})
	if err != nil {
		return cfeminter.Params{}, err
	}

	return res.Params, nil
}
