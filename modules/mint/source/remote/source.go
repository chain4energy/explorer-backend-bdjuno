package remote

import (
	"encoding/json"
	cfeminter "github.com/chain4energy/c4e-chain/v2/x/cfeminter/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/juno/v5/node/remote"
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
func (s Source) GetInflation() (sdk.Dec, error) {
	resp, err := http.Get("localhost:1317" + "/c4e/minter/v1beta1/inflation")
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

	var queryInflationResponse cfeminter.QueryInflationResponse
	err = json.Unmarshal(body, &queryInflationResponse)
	if err != nil {
		return sdk.Dec{}, err
	}

	return queryInflationResponse.Inflation, nil
}

// Params implements mintsource.Source
func (s Source) Params(height int64) (cfeminter.Params, error) {
	res, err := s.querier.Params(remote.GetHeightRequestContext(s.Ctx, height), &cfeminter.QueryParamsRequest{})
	if err != nil {
		return cfeminter.Params{}, nil
	}

	return res.Params, nil
}
