package main

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type clientStruct struct {
	logger         zerolog.Logger
	oldtownAddress string
	httpClient     http.Client
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
