#!/bin/bash

set -e

# Get absolute path to the script directory (https://stackoverflow.com/a/4774063)
SCRIPT_DIR="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
ROOT_DIR="$(realpath $SCRIPT_DIR/..)"

# Build 
cd "$ROOT_DIR" || exit 1  
make build-ie

# Record new video
scenario="interactive"
echo "Recording video..."
if ! docker run --rm -v "$ROOT_DIR:/app" -w "/app/tests" ghcr.io/charmbracelet/vhs "$scenario.tape" --output "$scenario.gif"; then
    echo "ERROR: Failed to record video for $scenario.tape"
    exit 1
fi
echo "Finished recording video"