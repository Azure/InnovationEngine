#!/bin/bash

set -eo pipefail

ROOT_DIR=$(git rev-parse --show-toplevel)

# Build 
pushd "$ROOT_DIR"
make build-ie

# Record new video
SCENARIO="interactive"
echo "Recording video..."
if ! docker run --rm -v "$ROOT_DIR:/app" -w "/app/tests" ghcr.io/charmbracelet/vhs "$SCENARIO.tape" --output "$SCENARIO.gif"; then
    echo "ERROR: Failed to record video for $SCENARIO.tape"
    exit 1
fi
echo "Finished recording video"