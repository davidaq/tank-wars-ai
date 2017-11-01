const { Router } = require('express');
const serveStatic = require('serve-static');
const path = require('path');
const concat = require('concat-stream');
const gameList = require('./game-list');

const router = module.exports = Router();

router.get('/', (req, res, next) => {
  req.url = '/static/index.html';
  next();
});

router.use('/static', serveStatic(path.resolve(__dirname, 'static')));

router.get('/game-list', (req, res) => {
  res.writeHead(200, {
    'content-type': 'text/event-stream',
  });
  res.write(`event: reset\ndata: reset\n\n`);
  gameList.forEach((item) => {
    res.write(`event: game\ndata: ${JSON.stringify(item)}\n\n`);
  });
  const gameHandler = item => {
    if (req.socket.destroyed) {
      gameList.removeListener('game', gameHandler);
      gameList.removeListener('keepalive', keepaliveHandler);
      return;
    }
    res.write(`event: game\ndata: ${JSON.stringify(item)}\n\n`);
  };
  const keepaliveHandler = () => {
    if (req.socket.destroyed) {
      gameList.removeListener('game', gameHandler);
      gameList.removeListener('keepalive', keepaliveHandler);
      return;
    }
    res.write(`event: keepalive\ndata: ${Date.now()}\n\n`);
  };
  gameList.on('game', gameHandler);
  gameList.on('keepalive', keepaliveHandler);
});

router.post('/create-game', (req, res) => {
  req.pipe(concat(buffer => {
    gameList.createGame(JSON.parse(buffer));
  }));
  res.end();
});
