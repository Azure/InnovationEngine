#!/bin/bash
# Script to start both the server and the frontend app together

# Kill any existing ts-node processes
pkill -f "ts-node src/server.ts" || true

# Start the server with PORT=4001
echo "Starting server on port 4001..."
PORT=4001 NODE_ENV=development npx ts-node src/server.ts &
SERVER_PID=$!

# Wait a moment for the server to start
sleep 2

# Check if the server is running
if ! curl -s http://localhost:4001/api/health >/dev/null; then
  echo "Server failed to start. Exiting."
  kill $SERVER_PID
  exit 1
fi
echo "Server started successfully!"

# Start the frontend app
echo "Starting frontend app..."
npm run start-no-validation

# When the frontend app exits, kill the server
kill $SERVER_PID
