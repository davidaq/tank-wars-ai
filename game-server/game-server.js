const http = require('http');
const express = require('express');
const fs = require('fs');
const path = require('path');

fs.mkdir(path.join(__dirname, 'db'), err => null);

const PORT = process.env.PORT || 8777;

const app = express();
app.use(require('./index-router'));

const server = http.createServer(app);
server.listen(PORT, err => {
  if (err) {
    console.error(err.stack);
  } else {
    console.log(`Listening on port ${PORT}`)
  }
});
