# gm-control-api-integration-
Integration tests for the gm-control-api system

## start gm-control-api on an empty backend.json file
This can be done with ./run-gm-control-api.sh

## run the integration test
go run .

## if you want you can preserve the data with
curl -X POST localhost:5555/admin/backup
