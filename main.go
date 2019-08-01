package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

func main() {
	var err error

	logger := zerolog.New(os.Stdout).
		With().Timestamp().Str("program", "integration").Logger()
	logger.Info().Msg("program starts")

	viper.AutomaticEnv()
	setEnvironmentDefaults()

	if viper.GetString("log_level") == "debug" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		logger.Debug().Msg("log level set to debug")
	}

	model := Model{}

	client := clientStruct{
		logger:        logger,
		serverAddress: viper.GetString("gm_control_api_address"),
	}

	for i, f := range []func(zerolog.Logger, *clientStruct) error{
		model.loadZone,
		model.loadCluster,
		model.loadDomain,
		model.loadListener,
		model.loadSharedRules,
		model.loadRoute,
		model.loadProxy,
		model.getZone,
		model.modifyCluster,
		model.modifyDomain,
		model.modifyListener,
		model.modifySharedRules,
		model.modifyRoute,
		model.modifyProxy,
		model.deleteProxy,
		model.deleteSharedRules,
		model.deleteRoute,
		model.deleteListener,
		model.deleteDomain,
		model.deleteCluster,
		model.deleteZone,
	} {
		if err =f(logger, &client); err != nil {
			logger.Fatal().AnErr(fmt.Sprintf("%d", i), err).Msg("main")
		}		
	}


}

func setEnvironmentDefaults() {
	viper.SetDefault("gm_control_api_address", "localhost:5555")
	viper.SetDefault("gm_control_api_org_key", "deciphernow")
	viper.SetDefault("log_level", "debug")
}
