package mint

import (
	"encoding/json"
	"fmt"
	mintertypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/bdjuno/v4/types"

	"github.com/rs/zerolog/log"
)

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "mint").Msg("parsing genesis")

	// Read the genesis state
	var genState mintertypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[mintertypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading mint genesis data: %s", err)
	}

	// Save the params
	err = m.db.SaveMintParams(types.NewMintParams(genState.Params, doc.InitialHeight))
	if err != nil {
		return fmt.Errorf("error while storing genesis mint params: %s", err)
	}

	return nil
}
