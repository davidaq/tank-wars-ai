const { EventEmitter } = require('events');
const shortid = require('shortid');
const fs = require('fs');
const path = require('path');
const GameHost = require('./game-host');

class GameList extends EventEmitter {
  constructor () {
    super();
    setInterval(() => {
      this.emit('keepalive');
    }, 15000);
    this.list = [];
    this.map = {};
    this.clientHost = {};
    fs.readFile(path.resolve(__dirname, 'db', 'list.txt'), 'utf-8', (err, data) => {
      if (err) {
        return;
      }
      data.split('\n').forEach(content => {
        if (content) {
          try {
            const opt = JSON.parse(content);
            if (this.map[opt.id]) {
              if (opt.__del) {
                this.list[this.map[opt.id].index] = null;
                delete this.map[opt.id];
              } else if (opt.__game) {
                this.map[opt.id].opt.games.push(opt.__game);
              } else {
                Object.assign(this.map[opt.id].opt, opt);
              }
            } else {
              this.createGame(opt, true);
            }
          } catch (err) {}
        }
      });
      const writer = fs.createWriteStream(path.resolve(__dirname, 'db', 'list.txt.bk'));
      writer.on('finish', () => {
        fs.rename(path.resolve(__dirname, 'db', 'list.txt.bk'), path.resolve(__dirname, 'db', 'list.txt'), err => null);
      });
      this.list.forEach(item => {
        if (item) {
          writer.write(JSON.stringify(item));
          writer.write('\n');
          if (!item.client && item.games.length < item.total) {
            this.createGameHost(item);
          }
        }
      });
      writer.end();
    });
  }

  forEach (cb) {
    this.list.forEach(cb);
  }

  createGame (opt, isInit = false) {
    if (!isInit) {
      opt.id = shortid.generate();
      opt.createtime = Date.now();
      opt.games = [];
    }
    if (opt.client) {
      opt.client = true;
      opt.total = opt.games.length;
      if (!isInit) {
        this.emit('game', opt);
        fs.appendFile(path.resolve(__dirname, 'db', 'list.txt'), JSON.stringify(opt) + '\n', err => null);
      }
    } else {
      opt.client = false;
      if (!isInit) {
        this.emit('game', opt);
        fs.appendFile(path.resolve(__dirname, 'db', 'list.txt'), JSON.stringify(opt) + '\n', err => null);
        setTimeout(() => this.createGameHost(opt), 100);
      }
    }
    this.map[opt.id] = { opt, index: this.list.length, game: null };
    this.list.push(opt);
  }

  interrupt (id) {
    const item = this.map[id];
    if (item && item.game) {
      item.game.interrupt();
    }
  }

  createGameHost (opt) {
    const settings = {};
    if (opt.client) {
      ['id', 'MaxMoves', 'MapWidth', 'MapHeight', 'Obstacles', 'InitTank'].forEach(f => {
        settings[f] = opt[f];
      });
      settings.red = this.clientMoveProvider(opt.id, 'red');
      settings.blue = this.clientMoveProvider(opt.id, 'blue');
      settings.total = opt.total + 1;
      settings.beginRound = opt.total;
    } else {
      ['id', 'total', 'red', 'blue', 'MaxMoves', 'MapWidth', 'MapHeight', 'Obstacles', 'InitTank'].forEach(f => {
        settings[f] = opt[f];
      });
    }
    const game = new GameHost(settings);
    this.map[opt.id].game = game;
    game.on('round', result => {
      this.map[opt.id].game = null;
      opt.games.push(result);
      this.emit('game', opt);
      fs.appendFile(path.resolve(__dirname, 'db', 'list.txt'), JSON.stringify({
        id: opt.id,
        __game: result,
      }) + '\n', err => null);
    });
    return game;
  }

  rmGame (id) {
    const cache = this.map[id];
    if (cache) {
      this.list[cache.index] = null;
      delete this.map[id];
      fs.appendFile(path.resolve(__dirname, 'db', 'list.txt'), JSON.stringify({ id, __del: true }) + '\n', err => null);
      for (let i = 0; i < cache.opt.total; i++) {
        fs.unlink(path.resolve(__dirname, 'db', `${id}_${i}.json.gz`), err => null);
      }
    }
  }

  rename (id, newName) {
    const cache = this.map[id];
    if (cache) {
      fs.appendFile(path.resolve(__dirname, 'db', 'list.txt'), JSON.stringify({ id, title: newName }) + '\n', err => null);
      const opt = cache.opt;
      opt.title = newName;
      this.emit('game', opt);
    }
  }

  hostClients (id, side, cb) {
    if (!this.map[id] || !this.map[id].opt.client) {
      return cb();
    }
    if (!this.clientHost[id]) {
      this.clientHost[id] = {
        ready: {
          red: false,
          blue: false,
        },
        move: {
          red: null,
          blue: null,
        },
      };
      setTimeout(() => this.beginClientHostedGame(id), 5000);
    }
    const host = this.clientHost[id];
    if (host.ready[side]) {
      return cb();
    }
    host.ready[side] = cb;
    if (host.game) {
      host.ready[side] = true;
      cb(host.game.getState(side));
    } else if (host.ready.red && host.ready.blue) {
      this.beginClientHostedGame(id);
    }
  }

  beginClientHostedGame (id) {
    if (!this.map[id]) {
      return;
    }
    const gameOpt = this.map[id].opt;
    const host = this.clientHost[id];
    if (!host || host.game) {
      return;
    }
    const game = this.createGameHost(gameOpt);
    game.on('round', result => {
      delete this.clientHost[id];
      if (host.move.red) {
        host.move.red.cb({ ended: true });
      }
      if (host.move.blue) {
        host.move.blue.cb({ ended: true });
      }
    });
    host.game = game;
    gameOpt.total++;
    this.emit('game', gameOpt);
    if (host.ready.red) {
      host.ready.red(game.getState('red'));
      host.ready.red = true;
    }
    if (host.ready.blue) {
      host.ready.blue(game.getState('red'));
      host.ready.blue = true;
    }
  }

  clientMove (id, side, moves, cb) {
    const host = this.clientHost[id];
    if (!host) {
      return;
    }
    if (host.move[side]) {
      return cb();
    }
    host.move[side] = { moves, cb };
  }

  clientMoveProvider (id, side) {
    const host = this.clientHost[id];
    if (!host) {
      return;
    }
    return (cb) => {
      let timedOut = false;
      const timeout = setTimeout(() => {
        timedOut = true;
      }, 10000);
      const interval = setInterval(() => {
        if (host.move[side]) {
          const { moves, cb: clientCb } = host.move[side];
          host.move[side] = false;
          clearTimeout(timeout);
          clearInterval(interval);
          cb(moves, state => {
            clientCb(state);
          });
        } else if (!host.ready[side] || timedOut) {
          clearTimeout(timeout);
          clearInterval(interval);
          const moves = {};
          host.game.getState(side).myTank.forEach(tank => {
            moves[tank.id] = (() => {
              switch (Math.floor(Math.random() * 7)) {
                case 0: return 'fire';
                case 1: return 'left';
                case 2: return 'right';
                case 3: return 'stay';
                default: return 'move';
              }
            })();
          });
          cb(moves, () => null);
        }
      }, 50);
    };
  }
}

module.exports = new GameList();
