const { Router } = require('express');
const serveStatic = require('serve-static');
const path = require('path');
const fs = require('fs');
const concat = require('concat-stream');
const gameList = require('./game-list');

const router = module.exports = Router();

router.get('/', (req, res, next) => {
  req.url = '/static/index.html';
  next();
});

router.use('/static', serveStatic(path.resolve(__dirname, 'static')));
router.use('/db', (req, res) => {
  res.writeHead(200, {
    'content-encoding': 'gzip',
  });
  const read = fs.createReadStream(path.resolve(__dirname, 'db', req.url.substr(1) + '.gz'));
  read.pipe(res);
  read.on('error', () => {
    res.end();
  });
});

router.get('/game-list', (req, res) => {
  res.writeHead(200, {
    'content-type': 'text/event-stream',
  });
  res.write(`event: reset\ndata: reset\n\n`);
  gameList.forEach((item) => {
    if (item) {
      res.write(`event: game\ndata: ${JSON.stringify(item)}\n\n`);
    }
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

router.delete('/game/:id', (req, res) => {
  gameList.rmGame(req.params.id);
  res.end();
});

const pat = ['fire', 'fire', 'move', 'move'];
let c = 0;
router.post('/random-player', (req, res) => {
  req.pipe(concat(buffer => {
    const msg = JSON.parse(buffer);
    if (msg.action === 'move') {
      const resp = msg.state.myTank.map(() => {
        switch (Math.floor(Math.random() * 7)) {
          case 0: return 'fire';
          case 1: return 'left';
          case 2: return 'right';
          case 3: return 'stay';
          default: return 'move';
        }
      });
      c++;
      res.end(JSON.stringify(resp));
    } else {
      res.end('[]');
    }
  }));
});
