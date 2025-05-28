// Simple Node.js/Express backend for executing whitelisted shell commands
const express = require('express');
const bodyParser = require('body-parser');
const { exec } = require('child_process');
const cors = require('cors');
const path = require('path');
const fs = require('fs');

const app = express();
const PORT = process.env.PORT || 4000;

// Discover the Innovation Engine binary in the expected location
const IE_BIN_PATH = path.resolve(__dirname, '../../../bin/ie');
const ALLOWED_COMMANDS = [];
if (fs.existsSync(IE_BIN_PATH) && fs.statSync(IE_BIN_PATH).isFile()) {
  ALLOWED_COMMANDS.push(IE_BIN_PATH);
  console.log(`Innovation Engine binary found at: ${IE_BIN_PATH}`);
} else {
  console.error(`Innovation Engine binary not found at: ${IE_BIN_PATH}`);
  console.error(
    'Please ensure the binary is built and located in the expected path: ' + IE_BIN_PATH
  );
}

app.use(cors());
app.use(bodyParser.json());

app.post('/api/exec', (req, res) => {
  let { command } = req.body;
  if (!command || typeof command !== 'string') {
    return res.status(400).json({ error: 'Invalid command' });
  }

  // Map 'ie' to the full path if requested
  let cmdName = command.split(' ')[0];
  let rest = command.slice(cmdName.length);
  if (cmdName === 'ie' && fs.existsSync(IE_BIN_PATH)) {
    cmdName = IE_BIN_PATH;
    command = `${IE_BIN_PATH}${rest}`;
  }

  if (!ALLOWED_COMMANDS.includes(cmdName)) {
    return res.status(403).json({ error: 'Command not allowed' });
  }

  exec(command, { timeout: 10000 }, (error, stdout, stderr) => {
    res.json({
      stdout,
      stderr,
      exitCode: error && error.code ? error.code : 0,
    });
  });
});

app.listen(PORT, () => {
  console.log(`Shell exec backend listening on port ${PORT}`);
});
