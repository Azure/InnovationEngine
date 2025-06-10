#!/bin/bash
# Script to start the app in debug mode

# Kill any existing ts-node processes
pkill -f "ts-node src/server.ts" || true

# Start the server with debugging enabled
echo "Starting server in debug mode on port 4001..."
NODE_ENV=development PORT=4001 node --inspect-brk -r ts-node/register src/server.ts &
SERVER_PID=$!

echo
echo "Server started with PID $SERVER_PID"
echo

# Wait a moment for the server to start
sleep 2

echo
echo "Server started in debug mode! Connect to it using the VSCode debugger."
echo "Debug URL: chrome://inspect or open the VSCode Debug panel."
echo 

# Start the frontend app
echo "Starting frontend app..."
npm run start-no-validation

# When the frontend app exits, kill the server
kill $SERVER_PID
