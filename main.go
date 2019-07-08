package main

import (
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
		logger:         logger,
		oldtownAddress: viper.GetString("oldtown_address"),
	}

	if err = model.loadZone(logger, &client); err != nil {
		logger.Fatal().AnErr("model.loadZone", err).Msg("main")
	}

	if err = model.loadCluster(logger, &client); err != nil {
		logger.Fatal().AnErr("model.loadCluster", err).Msg("main")
	}

	if err = model.loadDomain(logger, &client); err != nil {
		logger.Fatal().AnErr("model.loadDomain", err).Msg("main")
	}

	if err = model.loadListener(logger, &client); err != nil {
		logger.Fatal().AnErr("model.loadListener", err).Msg("main")
	}

	if err = model.loadSharedRules(logger, &client); err != nil {
		logger.Fatal().AnErr("model.loadSharedRules", err).Msg("main")
	}

	if err = model.loadRoute(logger, &client); err != nil {
		logger.Fatal().AnErr("model.loadRoute", err).Msg("main")
	}

	if err = model.getZone(logger, &client); err != nil {
		logger.Fatal().AnErr("model.getZone", err).Msg("main")
	}

	if err = model.modifyCluster(logger, &client); err != nil {
		logger.Fatal().AnErr("model.modifyCluster", err).Msg("main")
	}

	if err = model.modifyDomain(logger, &client); err != nil {
		logger.Fatal().AnErr("model.modifyDomain", err).Msg("main")
	}

	if err = model.modifyListener(logger, &client); err != nil {
		logger.Fatal().AnErr("model.modifyListener", err).Msg("main")
	}

	if err = model.modifySharedRules(logger, &client); err != nil {
		logger.Fatal().AnErr("model.modifySharedRules", err).Msg("main")
	}

	if err = model.modifyRoute(logger, &client); err != nil {
		logger.Fatal().AnErr("model.modifySharedRules", err).Msg("main")
	}

	if err = model.deleteSharedRules(logger, &client); err != nil {
		logger.Fatal().AnErr("model.deleteSharedRules", err).Msg("main")
	}

	if err = model.deleteRoute(logger, &client); err != nil {
		logger.Fatal().AnErr("model.deleteRoute", err).Msg("main")
	}

	if err = model.deleteListener(logger, &client); err != nil {
		logger.Fatal().AnErr("model.deleteListener", err).Msg("main")
	}

	if err = model.deleteDomain(logger, &client); err != nil {
		logger.Fatal().AnErr("model.deletDomain", err).Msg("main")
	}

	if err = model.deleteCluster(logger, &client); err != nil {
		logger.Fatal().AnErr("model.deleteCluster", err).Msg("main")
	}

	if err = model.deleteZone(logger, &client); err != nil {
		logger.Fatal().AnErr("model.deleteZone", err).Msg("main")
	}

}

func setEnvironmentDefaults() {
	viper.SetDefault("oldtown_address", "localhost:5555")
	viper.SetDefault("oldtown_org_key", "deciphernow")
	viper.SetDefault("log_level", "debug")
}
