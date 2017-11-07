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

router.get('/game/-events', (req, res) => {
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

router.post('/game', (req, res) => {
  req.pipe(concat(buffer => {
    gameList.createGame(JSON.parse(buffer));
  }));
  res.end();
});

router.get('/game/:id/interrupt', (req, res) => {
  gameList.interrupt(req.params.id);
  res.end();
});

router.delete('/game/:id', (req, res) => {
  gameList.rmGame(req.params.id);
  res.end();
});

router.post('/game/:id/name', (req, res) => {
  req.pipe(concat(buffer => {
    gameList.rename(req.params.id, buffer.toString());
  }));
  res.end();
});

// try to take side in a new match
router.get('/game/:id/match/:side', (req, res) => {
  if (['red', 'blue'].indexOf(req.params.side) === -1) {
    res.writeHead(400);
    res.end('side must be red or blue');
    return;
  }
  gameList.hostClients(req.params.id, req.params.side, state => {
    if (!state) {
      res.writeHead(404);
      res.end('game not found or match already started');
    } else {
      res.end(JSON.stringify(state));
    }
  });
});

// emit command and receive results
router.post('/game/:id/match/:side', (req, res) => {
  req.pipe(concat(buffer => {
    const moves = JSON.parse(buffer);
    gameList.clientMove(req.params.id, req.params.side, moves, state => {
      if (!state) {
        res.writeHead(404);
        res.end('game not found or move already set');
      } else {
        res.end(JSON.stringify(state));
      }
    });
  }));
});

router.post('/random-player', (req, res) => {
  req.pipe(concat(buffer => {
    const msg = JSON.parse(buffer);
    if (msg.action === 'move') {
      const resp = {};
      msg.state.myTank.forEach(tank => {
        resp[tank.id] = (() => {
          return 'fire-right';
          switch (Math.floor(Math.random() * 10)) {
            // case 1: return 'left';
            // case 2: return 'right';
            // case 3: return 'stay';
            // case 4: return 'fire-up';
            // case 5: return 'fire-left';
            // case 6: return 'fire-down';
            // case 7: return 'fire-right';
            // default: return 'move';
          }
        })();
      });
      res.end(JSON.stringify(resp));
    } else {
      res.end('');
    }
  }));
});
