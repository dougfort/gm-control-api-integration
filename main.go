package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"

	api "github.com/deciphernow/gm-control-api"
)

const zoneName = "workregion"

type clientStruct struct {
	logger zerolog.Logger
	oldtownAddress string
	httpClient http.Client
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

	client := clientStruct {
		logger: logger,
		oldtownAddress: viper.GetString("oldtown_address"),		
	}

	zone, err := createZone(&client)
	if err != nil {
		logger.Fatal().AnErr("createZone", err).Msg("main")
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

	reqZone.Name = zoneName

	if err := json.NewEncoder(&buffer).Encode(&reqZone); err != nil {
		return api.Zone{}, errors.Wrap(err, "Encode")
	}
	
	rawMessage, err := client.doHTTP("POST", "/v1.0/zone", &buffer)
	if err != nil {
		return api.Zone{}, errors.Wrap(err, "doHTTP")
	}

	if err = json.Unmarshal(rawMessage, &respZone); err != nil {
		return api.Zone{}, errors.Wrap(err, "Unmarshal")
	}

	return respZone, nil
}

func (client *clientStruct) doHTTP(method string, path string, body io.Reader) (json.RawMessage, error) {

	uri := fmt.Sprintf("http://%s%s", client.oldtownAddress, path)

	request, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil, errors.Wrap(err, "NewRequest")
	}

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
		rawMessage, _ := bodyMap["error"]
		var errorMap map[string]string
		if err = json.Unmarshal(rawMessage, &errorMap); err != nil {
			client.logger.Error().AnErr("Unmarshal", err).Msg("error in error handling")
		}
		return nil, errors.Errorf("HTTP request failed: (%d) %s: %+v", 
			response.StatusCode, response.Status, errorMap)
	}

	rawMessage, ok := bodyMap["result"]
	if !ok {
		return nil, errors.Errorf("unexpected result %+v", 
			bodyMap)
	}

	return rawMessage, nil
}