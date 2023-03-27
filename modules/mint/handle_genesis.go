package mint

import (
	"encoding/json"
	"fmt"
	cfemintertypes "github.com/chain4energy/c4e-chain/x/cfeminter/types"

	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/bdjuno/v4/types"

	"github.com/rs/zerolog/log"
)

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "mint").Msg("parsing genesis")

	// Read the genesis state
	var genState cfemintertypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[cfemintertypes.ModuleName], &genState)
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
