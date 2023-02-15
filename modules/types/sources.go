package types

import (
	"fmt"
	c4eapp "github.com/chain4energy/c4e-chain/app"
	cfemintertypes "github.com/chain4energy/c4e-chain/x/cfeminter/types"
	"github.com/ignite/cli/ignite/pkg/cosmoscmd"
	"os"

	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/forbole/juno/v3/node/remote"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/forbole/juno/v3/node/local"

	nodeconfig "github.com/forbole/juno/v3/node/config"

	banksource "github.com/forbole/bdjuno/v3/modules/bank/source"
	localbanksource "github.com/forbole/bdjuno/v3/modules/bank/source/local"
	remotebanksource "github.com/forbole/bdjuno/v3/modules/bank/source/remote"
	distrsource "github.com/forbole/bdjuno/v3/modules/distribution/source"
	localdistrsource "github.com/forbole/bdjuno/v3/modules/distribution/source/local"
	remotedistrsource "github.com/forbole/bdjuno/v3/modules/distribution/source/remote"
	govsource "github.com/forbole/bdjuno/v3/modules/gov/source"
	localgovsource "github.com/forbole/bdjuno/v3/modules/gov/source/local"
	remotegovsource "github.com/forbole/bdjuno/v3/modules/gov/source/remote"
	mintsource "github.com/forbole/bdjuno/v3/modules/mint/source"
	localmintsource "github.com/forbole/bdjuno/v3/modules/mint/source/local"
	remotemintsource "github.com/forbole/bdjuno/v3/modules/mint/source/remote"
	slashingsource "github.com/forbole/bdjuno/v3/modules/slashing/source"
	localslashingsource "github.com/forbole/bdjuno/v3/modules/slashing/source/local"
	remoteslashingsource "github.com/forbole/bdjuno/v3/modules/slashing/source/remote"
	stakingsource "github.com/forbole/bdjuno/v3/modules/staking/source"
	localstakingsource "github.com/forbole/bdjuno/v3/modules/staking/source/local"
	remotestakingsource "github.com/forbole/bdjuno/v3/modules/staking/source/remote"
)

type Sources struct {
	BankSource     banksource.Source
	DistrSource    distrsource.Source
	GovSource      govsource.Source
	MintSource     mintsource.Source
	SlashingSource slashingsource.Source
	StakingSource  stakingsource.Source
}

func BuildSources(nodeCfg nodeconfig.Config, encodingConfig *params.EncodingConfig) (*Sources, error) {
	switch cfg := nodeCfg.Details.(type) {
	case *remote.Details:
		return buildRemoteSources(cfg)
	case *local.Details:
		return buildLocalSources(cfg, encodingConfig)

	default:
		return nil, fmt.Errorf("invalid configuration type: %T", cfg)
	}
}

func buildLocalSources(cfg *local.Details, encodingConfig *params.EncodingConfig) (*Sources, error) {
	source, err := local.NewSource(cfg.Home, encodingConfig)
	if err != nil {
		return nil, err
	}

	app := c4eapp.New(
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)), source.StoreDB, nil, true, map[int64]bool{},
		cfg.Home, 0, cosmoscmd.EncodingConfig(simapp.MakeTestEncodingConfig()), simapp.EmptyAppOptions{},
	)
	newC4eApp := app.(*c4eapp.App)
	sources := &Sources{
		BankSource:     localbanksource.NewSource(source, banktypes.QueryServer(newC4eApp.BankKeeper)),
		DistrSource:    localdistrsource.NewSource(source, distrtypes.QueryServer(newC4eApp.DistrKeeper)),
		GovSource:      localgovsource.NewSource(source, govtypes.QueryServer(newC4eApp.GovKeeper)),
		MintSource:     localmintsource.NewSource(source, cfemintertypes.QueryServer(newC4eApp.CfeminterKeeper)),
		SlashingSource: localslashingsource.NewSource(source, slashingtypes.QueryServer(newC4eApp.SlashingKeeper)),
		StakingSource:  localstakingsource.NewSource(source, stakingkeeper.Querier{Keeper: newC4eApp.StakingKeeper}),
	}

	// Mount and initialize the stores
	err = source.MountKVStores(newC4eApp, "keys")
	if err != nil {
		return nil, err
	}

	err = source.MountTransientStores(newC4eApp, "tkeys")
	if err != nil {
		return nil, err
	}

	err = source.MountMemoryStores(newC4eApp, "memKeys")
	if err != nil {
		return nil, err
	}

	err = source.InitStores()
	if err != nil {
		return nil, err
	}

	return sources, nil
}

func buildRemoteSources(cfg *remote.Details) (*Sources, error) {
	source, err := remote.NewSource(cfg.GRPC)
	if err != nil {
		return nil, fmt.Errorf("error while creating remote source: %s", err)
	}

	return &Sources{
		BankSource:     remotebanksource.NewSource(source, banktypes.NewQueryClient(source.GrpcConn)),
		DistrSource:    remotedistrsource.NewSource(source, distrtypes.NewQueryClient(source.GrpcConn)),
		GovSource:      remotegovsource.NewSource(source, govtypes.NewQueryClient(source.GrpcConn)),
		MintSource:     remotemintsource.NewSource(source, cfemintertypes.NewQueryClient(source.GrpcConn)),
		SlashingSource: remoteslashingsource.NewSource(source, slashingtypes.NewQueryClient(source.GrpcConn)),
		StakingSource:  remotestakingsource.NewSource(source, stakingtypes.NewQueryClient(source.GrpcConn)),
	}, nil
}
