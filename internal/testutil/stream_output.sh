#!/bin/bash

# This script demonstrates real-time output with pauses to show the spinner
# It writes to both stdout and stderr to test both streams

echo "Starting long-running operation with streamed output..."
for i in {1..5}; do
    echo -n "Processing step $i of 5... "
    sleep 1
    echo "done"
    
    # Add a slight delay to allow spinner to be visible
    sleep 0.5
    
    # On step 3, output something to stderr to test error stream
    if [ $i -eq 3 ]; then
        echo "Note: Step $i added diagnostic info" >&2
    fi
done
echo "Operation complete!"