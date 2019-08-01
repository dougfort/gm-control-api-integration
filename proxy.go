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

const proxyName = "proxy-name"

func createProxy(
	client *clientStruct,
	zone api.Zone,
	domain api.Domain,
	listener api.Listener,
) (api.Proxy, error) {
	var reqProxy api.Proxy
	var respProxy api.Proxy
	var buffer bytes.Buffer
	var request http.Request

	reqProxy.Name = proxyName
	reqProxy.ZoneKey = zone.ZoneKey
	reqProxy.DomainKeys = []api.DomainKey{domain.DomainKey}
	reqProxy.ListenerKeys = []api.ListenerKey{listener.ListenerKey}

	if err := json.NewEncoder(&buffer).Encode(&reqProxy); err != nil {
		return api.Proxy{}, errors.Wrap(err, "Encode")
	}

	request.Method = "POST"
	request.URL = &url.URL{
		Scheme: "http",
		Host:   client.serverAddress,
		Path:   "/v1.0/proxy",
	}
	request.Body = ioutil.NopCloser(&buffer)

	rawMessage, err := client.doHTTP(&request)
	if err != nil {
		return api.Proxy{}, errors.Wrap(err, "doHTTP")
	}

	if err = json.Unmarshal(rawMessage, &respProxy); err != nil {
		return api.Proxy{}, errors.Wrap(err, "Unmarshal")
	}

	return respProxy, nil
}

func queryProxyByName(client *clientStruct) (api.Proxies, error) {
	var buffer bytes.Buffer
	var request http.Request

	nameFilter := service.ProxyFilter{Name: proxyName}
	nameFilters := []service.ProxyFilter{nameFilter}

	if err := json.NewEncoder(&buffer).Encode(nameFilters); err != nil {
		return nil, errors.Wrap(err, "Encode filters")
	}

	request.Method = "GET"
	request.URL = &url.URL{
		Scheme: "http",
		Host:   client.serverAddress,
		Path:   "/v1.0/proxy",
	}

	values := url.Values{}
	values.Add("filters", buffer.String())
	request.URL.RawQuery = values.Encode()

	rawMessage, err := client.doHTTP(&request)
	if err != nil {
		return nil, errors.Wrap(err, "doHTTP")
	}

	var proxys []api.Proxy

	if err = json.Unmarshal(rawMessage, &proxys); err != nil {
		return nil, errors.Wrap(err, "Unmarshal")
	}

	return proxys, nil
}

func getProxyByKey(client *clientStruct, proxyKey api.ProxyKey) (api.Proxy, error) {
	var respProxy api.Proxy
	var request http.Request

	request.Method = "GET"
	request.URL = &url.URL{
		Scheme: "http",
		Host:   client.serverAddress,
		Path:   fmt.Sprintf("/v1.0/proxy/%s", url.PathEscape(string(proxyKey))),
	}

	rawMessage, err := client.doHTTP(&request)
	if err != nil {
		return api.Proxy{}, errors.Wrap(err, "doHTTP")
	}

	if err = json.Unmarshal(rawMessage, &respProxy); err != nil {
		return api.Proxy{}, errors.Wrap(err, "Unmarshal")
	}

	return respProxy, nil
}

func editProxy(client *clientStruct, proxy api.Proxy) (api.Proxy, error) {
	var respProxy api.Proxy
	var buffer bytes.Buffer
	var request http.Request

	if err := json.NewEncoder(&buffer).Encode(&proxy); err != nil {
		return api.Proxy{}, errors.Wrap(err, "Encode")
	}

	request.Method = "PUT"
	request.URL = &url.URL{
		Scheme: "http",
		Host:   client.serverAddress,
		Path:   proxyKeyPath(proxy),
	}
	request.Body = ioutil.NopCloser(&buffer)

	rawMessage, err := client.doHTTP(&request)
	if err != nil {
		return api.Proxy{}, errors.Wrap(err, "doHTTP")
	}

	if err = json.Unmarshal(rawMessage, &respProxy); err != nil {
		return api.Proxy{}, errors.Wrap(err, "Unmarshal")
	}

	return respProxy, nil
}

func deleteProxy(client *clientStruct, proxy api.Proxy) error {
	var request http.Request

	request.Method = "DELETE"
	request.URL = &url.URL{
		Scheme: "http",
		Host:   client.serverAddress,
		Path:   proxyKeyPath(proxy),
	}

	values := url.Values{}
	values.Add("checksum", proxy.Checksum.Checksum)
	request.URL.RawQuery = values.Encode()

	_, err := client.doHTTP(&request)
	if err != nil {
		return errors.Wrap(err, "doHTTP")
	}

	return nil
}

func proxyKeyPath(proxy api.Proxy) string {
	return fmt.Sprintf("/v1.0/proxy/%s", url.PathEscape(string(proxy.ProxyKey)))
}
