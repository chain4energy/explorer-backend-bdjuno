package remote

import (
	cfeminter "github.com/chain4energy/c4e-chain/x/cfeminter/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/juno/v3/node/remote"

	mintsource "github.com/forbole/bdjuno/v3/modules/mint/source"
)

var (
	_ mintsource.Source = &Source{}
)

// Source implements mintsource.Source using a remote node
type Source struct {
	*remote.Source
	querier cfeminter.QueryClient
}

// NewSource returns a new Source instance
func NewSource(source *remote.Source, querier cfeminter.QueryClient) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// GetInflation implements mintsource.Source
func (s Source) GetInflation(height int64) (sdk.Dec, error) {
	res, err := s.querier.Inflation(remote.GetHeightRequestContext(s.Ctx, height), &cfeminter.QueryInflationRequest{})
	if err != nil {
		return sdk.Dec{}, err
	}

	return res.Inflation, nil
}

// Params implements mintsource.Source
func (s Source) Params(height int64) (cfeminter.Params, error) {
	res, err := s.querier.Params(remote.GetHeightRequestContext(s.Ctx, height), &cfeminter.QueryParamsRequest{})
	if err != nil {
		return cfeminter.Params{}, nil
	}

	return res.Params, nil
}
