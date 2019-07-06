package main

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	api "github.com/deciphernow/gm-control-api"
)

type Model struct {
	Zone     api.Zone
	Cluster1 api.Cluster
	Domain   api.Domain
}

func (model *Model) loadZone(logger zerolog.Logger, client *clientStruct) error {
	logger.Debug().Msg("verifying that zone does not exist before test")
	zones, err := queryZoneByName(client)
	if err != nil {
		return errors.Wrap(err, "queryZoneByName")
	}
	if len(zones) != 0 {
		return errors.Errorf("zone found before test: %+v", zones)
	}
	logger.Debug().Msg("creating zone")
	model.Zone, err = createZone(client)
	if err != nil {
		return errors.Wrap(err, "createZone")
	}
	logger.Debug().Msg("verifying that zone exists")
	zones, err = queryZoneByName(client)
	if err != nil {
		return errors.Wrap(err, "queryZoneByName")
	}
	if len(zones) != 1 {
		return errors.Errorf("wrong number of zones found: %+v", zones)
	}

	return nil
}

func (model *Model) loadCluster(logger zerolog.Logger, client *clientStruct) error {
	logger.Debug().Msg("verifying that cluster does not exist before test")
	clusters, err := queryClusterByName(client)
	if err != nil {
		return errors.Wrap(err, "queryClustersByName")
	}
	if len(clusters) != 0 {
		return errors.Errorf("cluster found before test: %+v", clusters)
	}
	logger.Debug().Msg("creating cluster")
	model.Cluster1, err = createCluster(client, model.Zone)
	if err != nil {
		return errors.Wrap(err, "createCluster")
	}
	logger.Debug().Msg("verifying that cluster exists")
	clusters, err = queryClusterByName(client)
	if err != nil {
		return errors.Wrap(err, "queryClusterByName")
	}
	if len(clusters) != 1 {
		return errors.Errorf("wrong number of clusters found: %+v", clusters)
	}

	return nil
}

func (model *Model) loadDomain(logger zerolog.Logger, client *clientStruct) error {
	logger.Debug().Msg("verifying that domain does not exist before test")
	domains, err := queryDomainByName(client)
	if err != nil {
		return errors.Wrap(err, "queryClustersByName")
	}
	if len(domains) != 0 {
		return errors.Errorf("domain found before test: %+v", domains)
	}
	logger.Debug().Msg("creating domain")
	model.Domain, err = createDomain(client, model.Zone)
	if err != nil {
		return errors.Wrap(err, "createDomain")
	}
	logger.Debug().Msg("verifying that domain exists")
	domains, err = queryDomainByName(client)
	if err != nil {
		return errors.Wrap(err, "queryDomainByName")
	}
	if len(domains) != 1 {
		return errors.Errorf("wrong number of domains found: %+v", domains)
	}

	return nil
}

func (model *Model) getZone(logger zerolog.Logger, client *clientStruct) error {
	logger.Debug().Msg("getting the zone object")
	zone2, err := getZoneByKey(client, model.Zone.ZoneKey)
	if err != nil {
		return errors.Wrap(err, "getZoneByKey")
	}
	if !zone2.Equals(model.Zone) {
		return errors.Errorf(
			"zone object mismatch: zone: %+v; zone2: %+v",
			model.Zone,
			zone2,
		)
	}

	return nil
}

func (model *Model) modifyCluster(logger zerolog.Logger, client *clientStruct) error {
	logger.Debug().Msg("editing the cluster object")
	maxConnections := 42
	model.Cluster1.CircuitBreakers = &api.CircuitBreakers{MaxConnections: &maxConnections}
	cluster2, err := editCluster(client, model.Cluster1)
	if err != nil {
		return errors.Wrap(err, "editCluster")
	}
	if *cluster2.CircuitBreakers.MaxConnections != *model.Cluster1.CircuitBreakers.MaxConnections {
		return errors.Errorf(
			"CircuitBreakers mismatch: cluster: %+v; cluster2: %+v",
			model.Cluster1,
			cluster2,
		)
	}
	model.Cluster1 = cluster2

	logger.Debug().Msg("adding a cluster instance")
	instance := api.Instance{Host: "localhost", Port: 42}
	cluster3, err := putClusterInstance(client, model.Cluster1, instance)
	if err != nil {
		return errors.Wrap(err, "putClusterInstance")
	}
	if len(cluster3.Instances) != 1 || !cluster3.Instances[0].Equals(instance) {
		return errors.Errorf(
			"cluster instances mismatch: cluster: %+v; cluster3: %+v",
			model.Cluster1,
			cluster3,
		)
	}
	model.Cluster1 = cluster3

	logger.Debug().Msg("deleting a cluster instance")
	cluster4, err := deleteClusterInstance(client, model.Cluster1, instance)
	if err != nil {
		return errors.Wrap(err, "deleteClusterInstance")
	}
	if len(cluster4.Instances) != 0 {
		return errors.Errorf(
			"cluster instance not deleted: cluster: %+v, cluster4: %+v",
			model.Cluster1,
			cluster4,
		)
	}
	model.Cluster1 = cluster4

	logger.Debug().Msg("getting the cluster object")
	cluster5, err := getClusterByKey(client, model.Cluster1.ClusterKey)
	if err != nil {
		return errors.Wrap(err, "getClusterByKey")
	}
	if !cluster5.Equals(model.Cluster1) {
		return errors.Errorf(
			"cluster object mismatch:  cluster: %+v; cluster5: %+v",
			model.Cluster1,
			cluster5,
		)
	}

	return nil
}

func (model *Model) modifyDomain(logger zerolog.Logger, client *clientStruct) error {
	logger.Debug().Msg("editing the domain object")
	const testPort = 333
	model.Domain.Port = testPort
	domain2, err := editDomain(client, model.Domain)
	if err != nil {
		return errors.Wrap(err, "editDomain")
	}
	if domain2.Port != model.Domain.Port {
		return errors.Errorf(
			"Port mismatch: domain: %+v; domain2: %+v",
			model.Domain,
			domain2,
		)
	}
	model.Domain = domain2

	logger.Debug().Msg("getting the domain object")
	domain3, err := getDomainByKey(client, model.Domain.DomainKey)
	if err != nil {
		return errors.Wrap(err, "getDomainByKey")
	}
	if !domain3.Equals(model.Domain) {
		return errors.Errorf(
			"domain object mismatch:  domain: %+v; domain3: %+v",
			model.Domain,
			domain3,
		)
	}

	return nil
}

func (model *Model) deleteCluster(logger zerolog.Logger, client *clientStruct) error {
	logger.Debug().Msg("deleting cluster")
	err := deleteCluster(client, model.Cluster1)
	if err != nil {
		return errors.Wrap(err, "deleteCluster")
	}
	logger.Debug().Msg("verifying that cluster does not exist after test")
	clusters, err := queryClusterByName(client)
	if err != nil {
		return errors.Wrap(err, "queryClusterByName")
	}
	if len(clusters) != 0 {
		return errors.Errorf("cluster found after delete: %+v", fmt.Sprintf("%+v", clusters))
	}

	model.Cluster1 = api.Cluster{}
	return nil
}

func (model *Model) deleteDomain(logger zerolog.Logger, client *clientStruct) error {
	logger.Debug().Msg("deleting domain")
	err := deleteDomain(client, model.Domain)
	if err != nil {
		return errors.Wrap(err, "deleteDomain")
	}
	logger.Debug().Msg("verifying that domain does not exist after test")
	domains, err := queryDomainByName(client)
	if err != nil {
		return errors.Wrap(err, "queryDomainByName")
	}
	if len(domains) != 0 {
		return errors.Errorf("domain found after delete: %+v", fmt.Sprintf("%+v", domains))
	}

	model.Domain = api.Domain{}
	return nil
}

func (model *Model) deleteZone(logger zerolog.Logger, client *clientStruct) error {
	logger.Debug().Msg("deleting zone")
	err := deleteZone(client, model.Zone)
	if err != nil {
		return errors.Wrap(err, "deleteZone")
	}
	logger.Debug().Msg("verifying that zone does not exist after test")
	zones, err := queryZoneByName(client)
	if err != nil {
		return errors.Wrap(err, "queryZoneByName")
	}
	if len(zones) != 0 {
		return errors.Errorf("zone found after delete: %+v", zones)
	}

	return nil
}
