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

const listenerName = "listener1"
const listenerIP = "127.0.0.1"
const listenerPort = 3333

var listenerProtocol api.ListenerProtocol = api.HttpListenerProtocol

func createListener(
	client *clientStruct,
	zone api.Zone,
	domain api.Domain,
) (api.Listener, error) {
	var reqListener api.Listener
	var respListener api.Listener
	var buffer bytes.Buffer
	var request http.Request

	reqListener.ZoneKey = zone.ZoneKey
	reqListener.Name = listenerName
	reqListener.IP = listenerIP
	reqListener.Port = listenerPort
	reqListener.Protocol = listenerProtocol
	reqListener.DomainKeys = []api.DomainKey{domain.DomainKey}

	if err := json.NewEncoder(&buffer).Encode(&reqListener); err != nil {
		return api.Listener{}, errors.Wrap(err, "Encode")
	}

	request.Method = "POST"
	request.URL = &url.URL{
		Scheme: "http",
		Host:   client.oldtownAddress,
		Path:   "/v1.0/listener",
	}
	request.Body = ioutil.NopCloser(&buffer)

	rawMessage, err := client.doHTTP(&request)
	if err != nil {
		return api.Listener{}, errors.Wrap(err, "doHTTP")
	}

	if err = json.Unmarshal(rawMessage, &respListener); err != nil {
		return api.Listener{}, errors.Wrap(err, "Unmarshal")
	}

	return respListener, nil
}

func queryListenerByName(client *clientStruct) (api.Listeners, error) {
	var buffer bytes.Buffer
	var request http.Request

	nameFilter := service.ListenerFilter{Name: listenerName}
	listenerFilters := []service.ListenerFilter{nameFilter}

	if err := json.NewEncoder(&buffer).Encode(listenerFilters); err != nil {
		return nil, errors.Wrap(err, "Encode filters")
	}

	request.Method = "GET"
	request.URL = &url.URL{
		Scheme: "http",
		Host:   client.oldtownAddress,
		Path:   "/v1.0/listener",
	}

	values := url.Values{}
	values.Add("filters", buffer.String())
	request.URL.RawQuery = values.Encode()

	rawMessage, err := client.doHTTP(&request)
	if err != nil {
		return nil, errors.Wrap(err, "doHTTP")
	}

	var listeners []api.Listener

	if err = json.Unmarshal(rawMessage, &listeners); err != nil {
		return nil, errors.Wrap(err, "Unmarshal")
	}

	return listeners, nil
}

func getListenerByKey(client *clientStruct, listenerKey api.ListenerKey) (api.Listener, error) {
	var respListener api.Listener
	var request http.Request

	request.Method = "GET"
	request.URL = &url.URL{
		Scheme: "http",
		Host:   client.oldtownAddress,
		Path:   fmt.Sprintf("/v1.0/listener/%s", url.PathEscape(string(listenerKey))),
	}

	rawMessage, err := client.doHTTP(&request)
	if err != nil {
		return api.Listener{}, errors.Wrap(err, "doHTTP")
	}

	if err = json.Unmarshal(rawMessage, &respListener); err != nil {
		return api.Listener{}, errors.Wrap(err, "Unmarshal")
	}

	return respListener, nil
}

func editListener(client *clientStruct, listener api.Listener) (api.Listener, error) {
	var respListener api.Listener
	var buffer bytes.Buffer
	var request http.Request

	if err := json.NewEncoder(&buffer).Encode(&listener); err != nil {
		return api.Listener{}, errors.Wrap(err, "Encode")
	}

	request.Method = "PUT"
	request.URL = &url.URL{
		Scheme: "http",
		Host:   client.oldtownAddress,
		Path:   listenerKeyPath(listener),
	}
	request.Body = ioutil.NopCloser(&buffer)

	rawMessage, err := client.doHTTP(&request)
	if err != nil {
		return api.Listener{}, errors.Wrap(err, "doHTTP")
	}

	if err = json.Unmarshal(rawMessage, &respListener); err != nil {
		return api.Listener{}, errors.Wrap(err, "Unmarshal")
	}

	return respListener, nil
}

func deleteListener(client *clientStruct, listener api.Listener) error {
	var request http.Request

	request.Method = "DELETE"
	request.URL = &url.URL{
		Scheme: "http",
		Host:   client.oldtownAddress,
		Path:   listenerKeyPath(listener),
	}

	values := url.Values{}
	values.Add("checksum", listener.Checksum.Checksum)
	request.URL.RawQuery = values.Encode()

	_, err := client.doHTTP(&request)
	if err != nil {
		return errors.Wrap(err, "doHTTP")
	}

	return nil
}

func listenerKeyPath(listener api.Listener) string {
	return fmt.Sprintf("/v1.0/listener/%s", url.PathEscape(string(listener.ListenerKey)))
}
