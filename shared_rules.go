package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"

	api "github.com/deciphernow/gm-control-api"
	service "github.com/deciphernow/gm-control-api/service"
)

const sharedRulesName = "sharedRules1"

func createSharedRules(
	client *clientStruct,
	zone api.Zone,
) (api.SharedRules, error) {
	var reqSharedRules api.SharedRules
	var respSharedRules api.SharedRules
	var buffer bytes.Buffer
	var request http.Request

	reqSharedRules.ZoneKey = zone.ZoneKey
	reqSharedRules.Name = sharedRulesName

	if err := json.NewEncoder(&buffer).Encode(&reqSharedRules); err != nil {
		return api.SharedRules{}, errors.Wrap(err, "Encode")
	}

	request.Method = "POST"
	request.URL = &url.URL{
		Scheme: "http",
		Host:   client.oldtownAddress,
		Path:   "/v1.0/shared_rules",
	}
	request.Body = ioutil.NopCloser(&buffer)

	rawMessage, err := client.doHTTP(&request)
	if err != nil {
		return api.SharedRules{}, errors.Wrap(err, "doHTTP")
	}

	if err = json.Unmarshal(rawMessage, &respSharedRules); err != nil {
		return api.SharedRules{}, errors.Wrap(err, "Unmarshal")
	}

	return respSharedRules, nil
}

func querySharedRulesByName(client *clientStruct) (api.SharedRulesSlice, error) {
	var buffer bytes.Buffer
	var request http.Request

	nameFilter := service.SharedRulesFilter{Name: sharedRulesName}
	sharedRulesFilters := []service.SharedRulesFilter{nameFilter}

	if err := json.NewEncoder(&buffer).Encode(sharedRulesFilters); err != nil {
		return nil, errors.Wrap(err, "Encode filters")
	}

	request.Method = "GET"
	request.URL = &url.URL{
		Scheme: "http",
		Host:   client.oldtownAddress,
		Path:   "/v1.0/shared_rules",
	}

	values := url.Values{}
	values.Add("filters", buffer.String())
	request.URL.RawQuery = values.Encode()

	rawMessage, err := client.doHTTP(&request)
	if err != nil {
		return nil, errors.Wrap(err, "doHTTP")
	}

	var sharedRulesSlice api.SharedRulesSlice

	if err = json.Unmarshal(rawMessage, &sharedRulesSlice); err != nil {
		return nil, errors.Wrap(err, "Unmarshal")
	}

	return sharedRulesSlice, nil
}

func getSharedRulesByKey(client *clientStruct, sharedRulesKey api.SharedRulesKey) (api.SharedRules, error) {
	var respSharedRules api.SharedRules
	var request http.Request

	request.Method = "GET"
	request.URL = &url.URL{
		Scheme: "http",
		Host:   client.oldtownAddress,
		Path:   fmt.Sprintf("/v1.0/shared_rules/%s", url.PathEscape(string(sharedRulesKey))),
	}

	rawMessage, err := client.doHTTP(&request)
	if err != nil {
		return api.SharedRules{}, errors.Wrap(err, "doHTTP")
	}

	if err = json.Unmarshal(rawMessage, &respSharedRules); err != nil {
		return api.SharedRules{}, errors.Wrap(err, "Unmarshal")
	}

	return respSharedRules, nil
}

func editSharedRules(client *clientStruct, sharedRules api.SharedRules) (api.SharedRules, error) {
	var respSharedRules api.SharedRules
	var buffer bytes.Buffer
	var request http.Request

	if err := json.NewEncoder(&buffer).Encode(&sharedRules); err != nil {
		return api.SharedRules{}, errors.Wrap(err, "Encode")
	}

	request.Method = "PUT"
	request.URL = &url.URL{
		Scheme: "http",
		Host:   client.oldtownAddress,
		Path:   sharedRulesKeyPath(sharedRules),
	}
	request.Body = ioutil.NopCloser(&buffer)

	rawMessage, err := client.doHTTP(&request)
	if err != nil {
		return api.SharedRules{}, errors.Wrap(err, "doHTTP")
	}

	if err = json.Unmarshal(rawMessage, &respSharedRules); err != nil {
		return api.SharedRules{}, errors.Wrap(err, "Unmarshal")
	}

	return respSharedRules, nil
}

func deleteSharedRules(client *clientStruct, sharedRules api.SharedRules) error {
	var request http.Request

	request.Method = "DELETE"
	request.URL = &url.URL{
		Scheme: "http",
		Host:   client.oldtownAddress,
		Path:   sharedRulesKeyPath(sharedRules),
	}

	values := url.Values{}
	values.Add("checksum", sharedRules.Checksum.Checksum)
	request.URL.RawQuery = values.Encode()

	_, err := client.doHTTP(&request)
	if err != nil {
		return errors.Wrap(err, "doHTTP")
	}

	return nil
}

func sharedRulesKeyPath(sharedRules api.SharedRules) string {
	return fmt.Sprintf("/v1.0/shared_rules/%s", url.PathEscape(string(sharedRules.SharedRulesKey)))
}
