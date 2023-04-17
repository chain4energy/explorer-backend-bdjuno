package local

import (
	"encoding/json"
	"fmt"
	"github.com/chain4energy/juno/v4/node/local"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mintertypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"io"
	"net/http"

	mintsource "github.com/forbole/bdjuno/v4/modules/mint/source"
)

var (
	_ mintsource.Source = &Source{}
)

// Source implements mintsource.Source using a local node
type Source struct {
	*local.Source
	querier mintertypes.QueryServer
}

// NewSource returns a new Source instace
func NewSource(source *local.Source, querier mintertypes.QueryServer) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// GetInflation implements mintsource.Source
func (s Source) GetInflation() (sdk.Dec, error) {
	resp, err := http.Get("http://localhost:1317/c4e/minter/v1beta1/inflation")
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
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return mintertypes.Params{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.querier.Params(sdk.WrapSDKContext(ctx), &mintertypes.QueryParamsRequest{})
	if err != nil {
		return mintertypes.Params{}, err
	}

	return res.Params, nil
}
