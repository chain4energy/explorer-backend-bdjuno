package v3

import (
	v3 "github.com/chain4energy/juno/v4/cmd/migrate/v3"

	"github.com/forbole/bdjuno/v4/modules/actions"
)

type Config struct {
	v3.Config `yaml:"-,inline"`

	// The following are there to support modules which config are present if they are enabled

	Actions *actions.Config `yaml:"actions"`
}
