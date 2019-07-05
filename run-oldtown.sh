#!/bin/bash
set -euxo pipefail

rm -f $(pwd)/oldtown_backend.json
touch $(pwd)/oldtown_backend.json

docker run \
    -p 5555:5555 \
    -e OLDTOWN_LOG_LEVEL=debug \
    -e OLDTOWN_ADDRESS=0.0.0.0:5555 \
    -e OLDTOWN_ORG_KEY=deciphernow \
    -e OLDTOWN_PERSISTER_TYPE=file \
    -e OLDTOWN_PERSISTER_PATH=/control-plane/oldtown_backend.json \
    -v $(pwd)/oldtown_backend.json:/control-plane/oldtown_backend.json \
    deciphernow/oldtown:latest