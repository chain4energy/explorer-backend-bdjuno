package main

import (
	cfeevapp "github.com/chain4energy/cfeev/app"
	"github.com/chain4energy/juno/v4/cmd"
	initcmd "github.com/chain4energy/juno/v4/cmd/init"
	parsetypes "github.com/chain4energy/juno/v4/cmd/parse/types"
	startcmd "github.com/chain4energy/juno/v4/cmd/start"
	"github.com/chain4energy/juno/v4/modules/messages"
	"github.com/cosmos/cosmos-sdk/types/module"
	migratecmd "github.com/forbole/bdjuno/v4/cmd/migrate"
	parsecmd "github.com/forbole/bdjuno/v4/cmd/parse"

	"github.com/forbole/bdjuno/v4/types/config"

	"github.com/forbole/bdjuno/v4/database"
	"github.com/forbole/bdjuno/v4/modules"
)

func main() {
	initCfg := initcmd.NewConfig().
		WithConfigCreator(config.Creator)

	parseCfg := parsetypes.NewConfig().
		WithDBBuilder(database.Builder).
		WithEncodingConfigBuilder(config.MakeEncodingConfig(getBasicManagers())).
		WithRegistrar(modules.NewRegistrar(getAddressesParser()))

	cfg := cmd.NewConfig("bdjuno").
		WithInitConfig(initCfg).
		WithParseConfig(parseCfg)

	// Run the command
	rootCmd := cmd.RootCmd(cfg.GetName())

	rootCmd.AddCommand(
		cmd.VersionCmd(),
		initcmd.NewInitCmd(cfg.GetInitConfig()),
		parsecmd.NewParseCmd(cfg.GetParseConfig()),
		migratecmd.NewMigrateCmd(cfg.GetName(), cfg.GetParseConfig()),
		startcmd.NewStartCmd(cfg.GetParseConfig()),
	)

	executor := cmd.PrepareRootCmd(cfg.GetName(), rootCmd)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

// getBasicManagers returns the various basic managers that are used to register the encoding to
// support custom messages.
// This should be edited by custom implementations if needed.
func getBasicManagers() []module.BasicManager {
	return []module.BasicManager{
		cfeevapp.ModuleBasics,
	}
}

// getAddressesParser returns the messages parser that should be used to get the users involved in
// a specific message.
// This should be edited by custom implementations if needed.
func getAddressesParser() messages.MessageAddressesParser {
	return messages.JoinMessageParsers(
		messages.CosmosMessageAddressesParser,
	)
}
