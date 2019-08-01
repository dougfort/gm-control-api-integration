#!/bin/bash
set -euxo pipefail

rm -f $(pwd)/backend.json
touch $(pwd)/backend.json

docker run \
    -p 5555:5555 \
    -e GM_CONTROL_API_LOG_LEVEL=debug \
    -e GM_CONTROL_API_ADDRESS=0.0.0.0:5555 \
    -e GM_CONTROL_API_ORG_KEY=deciphernow \
    -e GM_CONTROL_API_PERSISTER_TYPE=file \
    -e GM_CONTROL_API_PERSISTER_PATH=/control-plane/backend.json \
    -v $(pwd)/backend.json:/control-plane/backend.json \
    deciphernow/gm-control-api:latest