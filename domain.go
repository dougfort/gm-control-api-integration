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

const domainName = "domain1"

func createDomain(client *clientStruct, zone api.Zone) (api.Domain, error) {
	var reqDomain api.Domain
	var respDomain api.Domain
	var buffer bytes.Buffer
	var request http.Request

	reqDomain.ZoneKey = zone.ZoneKey
	reqDomain.Name = domainName

	if err := json.NewEncoder(&buffer).Encode(&reqDomain); err != nil {
		return api.Domain{}, errors.Wrap(err, "Encode")
	}

	request.Method = "POST"
	request.URL = &url.URL{
		Scheme: "http",
		Host:   client.oldtownAddress,
		Path:   "/v1.0/domain",
	}
	request.Body = ioutil.NopCloser(&buffer)

	rawMessage, err := client.doHTTP(&request)
	if err != nil {
		return api.Domain{}, errors.Wrap(err, "doHTTP")
	}

	if err = json.Unmarshal(rawMessage, &respDomain); err != nil {
		return api.Domain{}, errors.Wrap(err, "Unmarshal")
	}

	return respDomain, nil
}

func queryDomainByName(client *clientStruct) (api.Domains, error) {
	var buffer bytes.Buffer
	var request http.Request

	nameFilter := service.DomainFilter{Name: domainName}
	domainFilters := []service.DomainFilter{nameFilter}

	if err := json.NewEncoder(&buffer).Encode(domainFilters); err != nil {
		return nil, errors.Wrap(err, "Encode filters")
	}

	request.Method = "GET"
	request.URL = &url.URL{
		Scheme: "http",
		Host:   client.oldtownAddress,
		Path:   "/v1.0/domain",
	}

	values := url.Values{}
	values.Add("filters", buffer.String())
	request.URL.RawQuery = values.Encode()

	rawMessage, err := client.doHTTP(&request)
	if err != nil {
		return nil, errors.Wrap(err, "doHTTP")
	}

	var domains []api.Domain

	if err = json.Unmarshal(rawMessage, &domains); err != nil {
		return nil, errors.Wrap(err, "Unmarshal")
	}

	return domains, nil
}

func getDomainByKey(client *clientStruct, domainKey api.DomainKey) (api.Domain, error) {
	var respDomain api.Domain
	var request http.Request

	request.Method = "GET"
	request.URL = &url.URL{
		Scheme: "http",
		Host:   client.oldtownAddress,
		Path:   fmt.Sprintf("/v1.0/domain/%s", url.PathEscape(string(domainKey))),
	}

	rawMessage, err := client.doHTTP(&request)
	if err != nil {
		return api.Domain{}, errors.Wrap(err, "doHTTP")
	}

	if err = json.Unmarshal(rawMessage, &respDomain); err != nil {
		return api.Domain{}, errors.Wrap(err, "Unmarshal")
	}

	return respDomain, nil
}

func editDomain(client *clientStruct, domain api.Domain) (api.Domain, error) {
	var respDomain api.Domain
	var buffer bytes.Buffer
	var request http.Request

	if err := json.NewEncoder(&buffer).Encode(&domain); err != nil {
		return api.Domain{}, errors.Wrap(err, "Encode")
	}

	request.Method = "PUT"
	request.URL = &url.URL{
		Scheme: "http",
		Host:   client.oldtownAddress,
		Path:   domainKeyPath(domain),
	}
	request.Body = ioutil.NopCloser(&buffer)

	rawMessage, err := client.doHTTP(&request)
	if err != nil {
		return api.Domain{}, errors.Wrap(err, "doHTTP")
	}

	if err = json.Unmarshal(rawMessage, &respDomain); err != nil {
		return api.Domain{}, errors.Wrap(err, "Unmarshal")
	}

	return respDomain, nil
}

func deleteDomain(client *clientStruct, domain api.Domain) error {
	var request http.Request

	request.Method = "DELETE"
	request.URL = &url.URL{
		Scheme: "http",
		Host:   client.oldtownAddress,
		Path:   domainKeyPath(domain),
	}

	values := url.Values{}
	values.Add("checksum", domain.Checksum.Checksum)
	request.URL.RawQuery = values.Encode()

	_, err := client.doHTTP(&request)
	if err != nil {
		return errors.Wrap(err, "doHTTP")
	}

	return nil
}

func domainKeyPath(domain api.Domain) string {
	return fmt.Sprintf("/v1.0/domain/%s", url.PathEscape(string(domain.DomainKey)))
}
