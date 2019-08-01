package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"

	api "github.com/deciphernow/gm-control-api/api"
	service "github.com/deciphernow/gm-control-api/api/service"
)

const routePath = "/path/metrics"

func createRoute(
	client *clientStruct,
	zone api.Zone,
	domain api.Domain,
	sharedRules api.SharedRules,
) (api.Route, error) {
	var reqRoute api.Route
	var respRoute api.Route
	var buffer bytes.Buffer
	var request http.Request

	reqRoute.Path = routePath
	reqRoute.ZoneKey = zone.ZoneKey
	reqRoute.DomainKey = domain.DomainKey
	reqRoute.SharedRulesKey = sharedRules.SharedRulesKey

	if err := json.NewEncoder(&buffer).Encode(&reqRoute); err != nil {
		return api.Route{}, errors.Wrap(err, "Encode")
	}

	request.Method = "POST"
	request.URL = &url.URL{
		Scheme: "http",
		Host:   client.serverAddress,
		Path:   "/v1.0/route",
	}
	request.Body = ioutil.NopCloser(&buffer)

	rawMessage, err := client.doHTTP(&request)
	if err != nil {
		return api.Route{}, errors.Wrap(err, "doHTTP")
	}

	if err = json.Unmarshal(rawMessage, &respRoute); err != nil {
		return api.Route{}, errors.Wrap(err, "Unmarshal")
	}

	return respRoute, nil
}

func queryRouteByPath(client *clientStruct) (api.Routes, error) {
	var buffer bytes.Buffer
	var request http.Request

	pathFilter := service.RouteFilter{Path: routePath}
	pathFilters := []service.RouteFilter{pathFilter}

	if err := json.NewEncoder(&buffer).Encode(pathFilters); err != nil {
		return nil, errors.Wrap(err, "Encode filters")
	}

	request.Method = "GET"
	request.URL = &url.URL{
		Scheme: "http",
		Host:   client.serverAddress,
		Path:   "/v1.0/route",
	}

	values := url.Values{}
	values.Add("filters", buffer.String())
	request.URL.RawQuery = values.Encode()

	rawMessage, err := client.doHTTP(&request)
	if err != nil {
		return nil, errors.Wrap(err, "doHTTP")
	}

	var routes []api.Route

	if err = json.Unmarshal(rawMessage, &routes); err != nil {
		return nil, errors.Wrap(err, "Unmarshal")
	}

	return routes, nil
}

func getRouteByKey(client *clientStruct, routeKey api.RouteKey) (api.Route, error) {
	var respRoute api.Route
	var request http.Request

	request.Method = "GET"
	request.URL = &url.URL{
		Scheme: "http",
		Host:   client.serverAddress,
		Path:   fmt.Sprintf("/v1.0/route/%s", url.PathEscape(string(routeKey))),
	}

	rawMessage, err := client.doHTTP(&request)
	if err != nil {
		return api.Route{}, errors.Wrap(err, "doHTTP")
	}

	if err = json.Unmarshal(rawMessage, &respRoute); err != nil {
		return api.Route{}, errors.Wrap(err, "Unmarshal")
	}

	return respRoute, nil
}

func editRoute(client *clientStruct, route api.Route) (api.Route, error) {
	var respRoute api.Route
	var buffer bytes.Buffer
	var request http.Request

	if err := json.NewEncoder(&buffer).Encode(&route); err != nil {
		return api.Route{}, errors.Wrap(err, "Encode")
	}

	request.Method = "PUT"
	request.URL = &url.URL{
		Scheme: "http",
		Host:   client.serverAddress,
		Path:   routeKeyPath(route),
	}
	request.Body = ioutil.NopCloser(&buffer)

	rawMessage, err := client.doHTTP(&request)
	if err != nil {
		return api.Route{}, errors.Wrap(err, "doHTTP")
	}

	if err = json.Unmarshal(rawMessage, &respRoute); err != nil {
		return api.Route{}, errors.Wrap(err, "Unmarshal")
	}

	return respRoute, nil
}

func deleteRoute(client *clientStruct, route api.Route) error {
	var request http.Request

	request.Method = "DELETE"
	request.URL = &url.URL{
		Scheme: "http",
		Host:   client.serverAddress,
		Path:   routeKeyPath(route),
	}

	values := url.Values{}
	values.Add("checksum", route.Checksum.Checksum)
	request.URL.RawQuery = values.Encode()

	_, err := client.doHTTP(&request)
	if err != nil {
		return errors.Wrap(err, "doHTTP")
	}

	return nil
}

func routeKeyPath(route api.Route) string {
	return fmt.Sprintf("/v1.0/route/%s", url.PathEscape(string(route.RouteKey)))
}
