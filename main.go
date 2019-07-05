package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"

	api "github.com/deciphernow/gm-control-api"
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

	var model Model

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
	model.Zone, err = createZone(&client)
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

	logger.Debug().Msg("verifying that cluster does not exist before test")
	clusters, err := queryClusterByName(&client)
	if err != nil {
		logger.Fatal().AnErr("queryClustersByName", err).Msg("main")
	}
	if len(clusters) != 0 {
		logger.Fatal().Str("clusters", fmt.Sprintf("%+v", clusters)).Msg("cluster found before test")
	}
	logger.Debug().Msg("creating cluster")
	model.Cluster1, err = createCluster(&client, model.Zone)
	if err != nil {
		logger.Fatal().AnErr("createCluster", err).Msg("main")
	}
	logger.Debug().Msg("verifying that cluster exists")
	clusters, err = queryClusterByName(&client)
	if err != nil {
		logger.Fatal().AnErr("queryClusterByName", err).Msg("main")
	}
	if len(clusters) != 1 {
		logger.Fatal().Str("clusters", fmt.Sprintf("%+v", clusters)).Msg("wrong number of clusters found")
	}

	logger.Debug().Msg("getting the zone object")
	zone2, err := getZoneByKey(&client, model.Zone.ZoneKey)
	if err != nil {
		logger.Fatal().AnErr("getZoneByKey", err).Msg("main")
	}
	if !zone2.Equals(model.Zone) {
		logger.Fatal().
			Str("zone", fmt.Sprintf("%+v", model.Zone)).
			Str("zone2", fmt.Sprintf("%+v", zone2)).
			Msg("zone object mismatch")
	}

	logger.Debug().Msg("getting the cluster object")
	cluster2, err := getClusterByKey(&client, model.Cluster1.ClusterKey)
	if err != nil {
		logger.Fatal().AnErr("getClusterByKey", err).Msg("main")
	}
	if !cluster2.Equals(model.Cluster1) {
		logger.Fatal().
			Str("cluster", fmt.Sprintf("%+v", model.Cluster1)).
			Str("cluster2", fmt.Sprintf("%+v", cluster2)).
			Msg("cluster object mismatch")
	}

	logger.Debug().Msg("deleting cluster")
	err = deleteCluster(&client, model.Cluster1)
	if err != nil {
		logger.Fatal().AnErr("deleteCluster", err).Msg("main")
	}
	logger.Debug().Msg("verifying that cluster does not exist after test")
	clusters, err = queryClusterByName(&client)
	if err != nil {
		logger.Fatal().AnErr("queryClusterByName", err).Msg("main")
	}
	if len(clusters) != 0 {
		logger.Fatal().Str("clusters", fmt.Sprintf("%+v", clusters)).Msg("zone found after test")
	}

	model.Cluster1 = api.Cluster{}

	logger.Debug().Msg("deleting zone")
	err = deleteZone(&client, model.Zone)
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
}

func setEnvironmentDefaults() {
	viper.SetDefault("oldtown_address", "localhost:5555")
	viper.SetDefault("oldtown_org_key", "deciphernow")
	viper.SetDefault("log_level", "debug")
}
