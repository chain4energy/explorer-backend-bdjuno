package mint

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

// UpdateParams gets the updated params and stores them inside the database
func (m *Module) UpdateParams(height int64) error {
	log.Debug().Str("module", "mint").Int64("height", height).
		Msg("updating params")

	_, err := m.source.Params(height)
	if err != nil {
		return fmt.Errorf("error while getting params: %s", err)
	}
	return nil
	//return m.db.SaveMintParams(types.NewMintParams(params, height))

}
