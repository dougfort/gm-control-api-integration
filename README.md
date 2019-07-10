# oldtown-integration-
Integration tests for the oldtown system

## start oldtown on an empty oldtown_backend.json file
This can be done with ./run-oldtown.sh

## run the integration test
go run .

## if you want you can preserve the data with
curl -X POST localhost:5555/admin/backup
