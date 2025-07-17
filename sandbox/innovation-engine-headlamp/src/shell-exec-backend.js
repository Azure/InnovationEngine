// server.js
const express = require("express");
const http = require("http");
const WebSocket = require("ws");
const pty = require("node-pty");

const app = express();
const server = http.createServer(app);
const wssTerm = new WebSocket.Server({ noServer: true });
const wssResize = new WebSocket.Server({ noServer: true });

app.use(express.static("public")); // serve React build

server.on("upgrade", (req, socket, head) => {
  if (req.url === "/ws/term") {
    wssTerm.handleUpgrade(req, socket, head, ws => {
      handleTerminalSocket(ws);
    });
  } else if (req.url === "/ws/resize") {
    wssResize.handleUpgrade(req, socket, head, ws => {
      ws.on("message", msg => {
        const { cols, rows } = JSON.parse(msg);
        if (ptyProcess) {
          ptyProcess.resize(cols, rows);
        }
      });
    });
  }
});

let ptyProcess = null;

function handleTerminalSocket(ws) {
  ptyProcess = pty.spawn("bash", [], {
    name: "xterm-color",
    cols: 80,
    rows: 24,
    cwd: process.env.HOME,
    env: process.env,
  });

  ptyProcess.on("data", data => ws.send(data));
  ws.on("message", msg => ptyProcess.write(msg));
  ws.on("close", () => ptyProcess.kill());
}

server.listen(4001, () => {
  console.log("Server running at http://localhost:4001");
});
