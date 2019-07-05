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

const clusterName = "cluster1"

func createCluster(client *clientStruct, zone api.Zone) (api.Cluster, error) {
	var reqCluster api.Cluster
	var respCluster api.Cluster
	var buffer bytes.Buffer
	var request http.Request

	reqCluster.ZoneKey = zone.ZoneKey
	reqCluster.Name = clusterName

	if err := json.NewEncoder(&buffer).Encode(&reqCluster); err != nil {
		return api.Cluster{}, errors.Wrap(err, "Encode")
	}

	request.Method = "POST"
	request.URL = &url.URL{
		Scheme: "http",
		Host:   client.oldtownAddress,
		Path:   "/v1.0/cluster",
	}
	request.Body = ioutil.NopCloser(&buffer)

	rawMessage, err := client.doHTTP(&request)
	if err != nil {
		return api.Cluster{}, errors.Wrap(err, "doHTTP")
	}

	if err = json.Unmarshal(rawMessage, &respCluster); err != nil {
		return api.Cluster{}, errors.Wrap(err, "Unmarshal")
	}

	return respCluster, nil
}

func queryClusterByName(client *clientStruct) (api.Clusters, error) {
	var buffer bytes.Buffer
	var request http.Request

	nameFilter := service.ClusterFilter{Name: clusterName}
	clusterFilters := []service.ClusterFilter{nameFilter}

	if err := json.NewEncoder(&buffer).Encode(clusterFilters); err != nil {
		return nil, errors.Wrap(err, "Encode filters")
	}

	request.Method = "GET"
	request.URL = &url.URL{
		Scheme: "http",
		Host:   client.oldtownAddress,
		Path:   "/v1.0/cluster",
	}

	values := url.Values{}
	values.Add("filters", buffer.String())
	request.URL.RawQuery = values.Encode()

	rawMessage, err := client.doHTTP(&request)
	if err != nil {
		return nil, errors.Wrap(err, "doHTTP")
	}

	var clusters []api.Cluster

	if err = json.Unmarshal(rawMessage, &clusters); err != nil {
		return nil, errors.Wrap(err, "Unmarshal")
	}

	return clusters, nil
}

func getClusterByKey(client *clientStruct, clusterKey api.ClusterKey) (api.Cluster, error) {
	var respCluster api.Cluster
	var request http.Request

	request.Method = "GET"
	request.URL = &url.URL{
		Scheme: "http",
		Host:   client.oldtownAddress,
		Path:   fmt.Sprintf("/v1.0/cluster/%s", url.PathEscape(string(clusterKey))),
	}

	rawMessage, err := client.doHTTP(&request)
	if err != nil {
		return api.Cluster{}, errors.Wrap(err, "doHTTP")
	}

	if err = json.Unmarshal(rawMessage, &respCluster); err != nil {
		return api.Cluster{}, errors.Wrap(err, "Unmarshal")
	}

	return respCluster, nil
}

func editCluster(client *clientStruct, cluster api.Cluster) (api.Cluster, error) {
	var respCluster api.Cluster
	var buffer bytes.Buffer
	var request http.Request

	if err := json.NewEncoder(&buffer).Encode(&cluster); err != nil {
		return api.Cluster{}, errors.Wrap(err, "Encode")
	}

	request.Method = "PUT"
	request.URL = &url.URL{
		Scheme: "http",
		Host:   client.oldtownAddress,
		Path:   clusterKeyPath(cluster),
	}
	request.Body = ioutil.NopCloser(&buffer)

	rawMessage, err := client.doHTTP(&request)
	if err != nil {
		return api.Cluster{}, errors.Wrap(err, "doHTTP")
	}

	if err = json.Unmarshal(rawMessage, &respCluster); err != nil {
		return api.Cluster{}, errors.Wrap(err, "Unmarshal")
	}

	return respCluster, nil
}

func deleteCluster(client *clientStruct, cluster api.Cluster) error {
	var request http.Request

	request.Method = "DELETE"
	request.URL = &url.URL{
		Scheme: "http",
		Host:   client.oldtownAddress,
		Path:   clusterKeyPath(cluster),
	}

	values := url.Values{}
	values.Add("checksum", cluster.Checksum.Checksum)
	request.URL.RawQuery = values.Encode()

	_, err := client.doHTTP(&request)
	if err != nil {
		return errors.Wrap(err, "doHTTP")
	}

	return nil
}

func clusterKeyPath(cluster api.Cluster) string {
	return fmt.Sprintf("/v1.0/cluster/%s", url.PathEscape(string(cluster.ClusterKey)))
}
