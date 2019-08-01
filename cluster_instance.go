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
)

func putClusterInstance(
	client *clientStruct,
	cluster api.Cluster,
	instance api.Instance,
) (api.Cluster, error) {
	var respCluster api.Cluster
	var buffer bytes.Buffer
	var request http.Request

	if err := json.NewEncoder(&buffer).Encode(&instance); err != nil {
		return api.Cluster{}, errors.Wrap(err, "Encode")
	}

	request.Method = "PUT"
	request.URL = &url.URL{
		Scheme: "http",
		Host:   client.serverAddress,
		Path:   clusterInstancesPath(cluster.ClusterKey),
	}
	request.Body = ioutil.NopCloser(&buffer)

	values := url.Values{}
	values.Add("checksum", cluster.Checksum.Checksum)
	request.URL.RawQuery = values.Encode()

	rawMessage, err := client.doHTTP(&request)
	if err != nil {
		return api.Cluster{}, errors.Wrap(err, "doHTTP")
	}

	if err = json.Unmarshal(rawMessage, &respCluster); err != nil {
		return api.Cluster{}, errors.Wrap(err, "Unmarshal")
	}

	return respCluster, nil
}

func deleteClusterInstance(
	client *clientStruct,
	cluster api.Cluster,
	instance api.Instance,
) (api.Cluster, error) {
	var respCluster api.Cluster
	var request http.Request

	request.Method = "DELETE"
	request.URL = &url.URL{
		Scheme: "http",
		Host:   client.serverAddress,
		Path:   clusterInstancePath(cluster.ClusterKey, instance.Key()),
	}

	values := url.Values{}
	values.Add("checksum", cluster.Checksum.Checksum)
	request.URL.RawQuery = values.Encode()

	rawMessage, err := client.doHTTP(&request)
	if err != nil {
		return api.Cluster{}, errors.Wrap(err, "doHTTP")
	}

	if err = json.Unmarshal(rawMessage, &respCluster); err != nil {
		return api.Cluster{}, errors.Wrap(err, "Unmarshal")
	}

	return respCluster, nil
}

func clusterInstancesPath(clusterKey api.ClusterKey) string {
	return fmt.Sprintf("/v1.0/cluster/%s/instances", url.PathEscape(string(clusterKey)))
}

func clusterInstancePath(clusterKey api.ClusterKey, instanceIdentifier string) string {
	return fmt.Sprintf(
		"/v1.0/cluster/%s/instances/%s",
		url.PathEscape(string(clusterKey)),
		url.PathEscape(instanceIdentifier),
	)
}
