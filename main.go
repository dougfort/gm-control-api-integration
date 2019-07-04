package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"

	api "github.com/deciphernow/gm-control-api"
)

const zoneName = "workregion"

type clientStruct struct {
	logger         zerolog.Logger
	oldtownAddress string
	httpClient     http.Client
}

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

	logger.Debug().Msg("creating zone")
	zone, err := createZone(&client)
	if err != nil {
		logger.Fatal().AnErr("createZone", err).Msg("main")
	}
	logger.Debug().Msg("deleting zone")
	err = deleteZone(&client, zone)
	if err != nil {
		logger.Fatal().AnErr("deleteZone", err).Msg("main")
	}

	logger.Debug().Str("zone", fmt.Sprintf("%+v", zone)).Msg("main")
}

func setEnvironmentDefaults() {
	viper.SetDefault("oldtown_address", "localhost:5555")
	viper.SetDefault("oldtown_org_key", "deciphernow")
	viper.SetDefault("log_level", "debug")
}

func createZone(client *clientStruct) (api.Zone, error) {
	var reqZone api.Zone
	var respZone api.Zone
	var buffer bytes.Buffer
	var request http.Request

	reqZone.Name = zoneName

	if err := json.NewEncoder(&buffer).Encode(&reqZone); err != nil {
		return api.Zone{}, errors.Wrap(err, "Encode")
	}

	request.Method = "POST"
	request.URL = &url.URL{
		Scheme: "http",
		Host:   client.oldtownAddress,
		Path:   "/v1.0/zone",
	}
	request.Body = ioutil.NopCloser(&buffer)

	rawMessage, err := client.doHTTP(&request)
	if err != nil {
		return api.Zone{}, errors.Wrap(err, "doHTTP")
	}

	if err = json.Unmarshal(rawMessage, &respZone); err != nil {
		return api.Zone{}, errors.Wrap(err, "Unmarshal")
	}

	return respZone, nil
}

func deleteZone(client *clientStruct, zone api.Zone) error {
	var request http.Request

	request.Method = "DELETE"
	request.URL = &url.URL{
		Scheme: "http",
		Host:   client.oldtownAddress,
		Path:   fmt.Sprintf("/v1.0/zone/%s", url.PathEscape(string(zone.ZoneKey))),
	}

	values := url.Values{}
	values.Add("checksum", zone.Checksum.Checksum)
	request.URL.RawQuery = values.Encode()

	_, err := client.doHTTP(&request)
	if err != nil {
		return errors.Wrap(err, "doHTTP")
	}

	return nil
}

func (client *clientStruct) doHTTP(request *http.Request) (json.RawMessage, error) {
	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "client.Do")
	}
	defer response.Body.Close()

	var bodyMap map[string]json.RawMessage
	if err = json.NewDecoder(response.Body).Decode(&bodyMap); err != nil {
		return nil, errors.Wrap(err, "Decode")
	}

	if response.StatusCode != http.StatusOK {
		rawMessage := bodyMap["error"]
		var errorMap map[string]string
		if err = json.Unmarshal(rawMessage, &errorMap); err != nil {
			client.logger.Error().AnErr("Unmarshal", err).Msg("error in error handling")
		}
		return nil, errors.Errorf("HTTP request failed: (%d) %s: %+v",
			response.StatusCode, response.Status, errorMap)
	}

	return bodyMap["result"], nil
}
