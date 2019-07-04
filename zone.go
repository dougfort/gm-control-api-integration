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

const zoneName = "workregion"

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

func queryZoneByName(client *clientStruct) (api.Zones, error) {
	var buffer bytes.Buffer
	var request http.Request

	nameFilter := service.ZoneFilter{Name: zoneName}
	zoneFilters := []service.ZoneFilter{nameFilter}

	if err := json.NewEncoder(&buffer).Encode(zoneFilters); err != nil {
		return nil, errors.Wrap(err, "Encode filters")
	}

	request.Method = "GET"
	request.URL = &url.URL{
		Scheme: "http",
		Host:   client.oldtownAddress,
		Path:   "/v1.0/zone",
	}

	values := url.Values{}
	values.Add("filters", buffer.String())
	request.URL.RawQuery = values.Encode()

	rawMessage, err := client.doHTTP(&request)
	if err != nil {
		return nil, errors.Wrap(err, "doHTTP")
	}

	var zones []api.Zone

	if err = json.Unmarshal(rawMessage, &zones); err != nil {
		return nil, errors.Wrap(err, "Unmarshal")
	}

	return zones, nil
}

func getZoneByKey(client *clientStruct, zoneKey api.ZoneKey) (api.Zone, error) {
	var respZone api.Zone
	var request http.Request

	request.Method = "GET"
	request.URL = &url.URL{
		Scheme: "http",
		Host:   client.oldtownAddress,
		Path:   fmt.Sprintf("/v1.0/zone/%s", url.PathEscape(string(zoneKey))),
	}

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
