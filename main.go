package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

func main() {
	logger := zerolog.New(os.Stdout).
		With().Timestamp().Str("program", "integration").Logger()
	logger.Info().Msg("program starts")

	viper.AutomaticEnv()
	setEnvironmentDefaults()

	if viper.GetString("log_level") == "debug" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		logger.Debug().Msg("log level set to debug")
	}

	client := clientStruct{
		logger:         logger,
		oldtownAddress: viper.GetString("oldtown_address"),
	}

	logger.Debug().Msg("verifying that zone does not exist before test")
	zones, err := queryZoneByName(&client)
	if err != nil {
		logger.Fatal().AnErr("queryZoneByName", err).Msg("main")
	}
	if len(zones) != 0 {
		logger.Fatal().Str("zones", fmt.Sprintf("%+v", zones)).Msg("zone found before test")
	}
	logger.Debug().Msg("creating zone")
	zone, err := createZone(&client)
	if err != nil {
		logger.Fatal().AnErr("createZone", err).Msg("main")
	}
	logger.Debug().Msg("verifying that zone exists")
	zones, err = queryZoneByName(&client)
	if err != nil {
		logger.Fatal().AnErr("queryZoneByName", err).Msg("main")
	}
	if len(zones) != 1 {
		logger.Fatal().Str("zones", fmt.Sprintf("%+v", zones)).Msg("wrong number of zones found")
	}
	logger.Debug().Msg("getting the zone object")
	zone2, err := getZoneByKey(&client, zone.ZoneKey)
	if err != nil {
		logger.Fatal().AnErr("getZoneByKey", err).Msg("main")
	}
	if !zone2.Equals(zone) {
		logger.Fatal().
			Str("zone", fmt.Sprintf("%+v", zone)).
			Str("zone2", fmt.Sprintf("%+v", zone2)).
			Msg("zone object mismatch")
	}
	logger.Debug().Msg("deleting zone")
	err = deleteZone(&client, zone)
	if err != nil {
		logger.Fatal().AnErr("deleteZone", err).Msg("main")
	}
	logger.Debug().Msg("verifying that zone does not exist after test")
	zones, err = queryZoneByName(&client)
	if err != nil {
		logger.Fatal().AnErr("queryZoneByName", err).Msg("main")
	}
	if len(zones) != 0 {
		logger.Fatal().Str("zones", fmt.Sprintf("%+v", zones)).Msg("zone found after test")
	}

	logger.Debug().Str("zone", fmt.Sprintf("%+v", zone)).Msg("main")
}

func setEnvironmentDefaults() {
	viper.SetDefault("oldtown_address", "localhost:5555")
	viper.SetDefault("oldtown_org_key", "deciphernow")
	viper.SetDefault("log_level", "debug")
}
