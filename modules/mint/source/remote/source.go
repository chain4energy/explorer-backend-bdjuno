package remote

import (
	"encoding/json"
	"github.com/chain4energy/juno/v4/node/remote"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mintertypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"io"
	"net/http"

	mintsource "github.com/forbole/bdjuno/v4/modules/mint/source"
)

var (
	_ mintsource.Source = &Source{}
)

// Source implements mintsource.Source using a remote node
type Source struct {
	*remote.Source
	querier mintertypes.QueryClient
}

// NewSource returns a new Source instance
func NewSource(source *remote.Source, querier mintertypes.QueryClient) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// GetInflation implements mintsource.Source
func (s Source) GetInflation() (sdk.Dec, error) {
	resp, err := http.Get(s.Restendpoint + "/c4e/minter/v1beta1/inflation")
	if err != nil {
		return sdk.Dec{}, err
	}
	defer resp.Body.Close()
	if err != nil {
		return sdk.Dec{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return sdk.Dec{}, err
	}

	var queryInflationResponse mintertypes.QueryInflationResponse
	err = json.Unmarshal(body, &queryInflationResponse)
	if err != nil {
		return sdk.Dec{}, err
	}

	return queryInflationResponse.Inflation, nil
}

// Params implements mintsource.Source
func (s Source) Params(height int64) (mintertypes.Params, error) {
	res, err := s.querier.Params(remote.GetHeightRequestContext(s.Ctx, height), &mintertypes.QueryParamsRequest{})
	if err != nil {
		return mintertypes.Params{}, nil
	}

	return res.Params, nil
}
