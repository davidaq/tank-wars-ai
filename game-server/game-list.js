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
    } else {
      opt.client = false;
      let createGameHost = false;
      if (!isInit) {
        this.emit('game', opt);
        fs.appendFile(path.resolve(__dirname, 'db', 'list.txt'), JSON.stringify(opt) + '\n', err => null);
        createGameHost = true;
      } else {
        createGameHost = opt.games.length < opt.total;
      }
      if (createGameHost) {
        const game = new GameHost(opt.id, opt.red, opt.blue, opt.total);
        game.on('round', result => {
          opt.games.push(result);
          this.emit('game', opt);
          fs.appendFile(path.resolve(__dirname, 'db', 'list.txt'), JSON.stringify({
            id: opt.id,
            __game: result,
          }) + '\n', err => null);
        });
      }
    }
    this.map[opt.id] = { opt, index: this.list.length };
    this.list.push(opt);
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

  hostClients (id, side, cb) {
    if (!this.map[id]) {
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
    if (host.game) {
      host.ready[side] = true;
      cb(host.game.getState(side));
    } else if (host.ready.red && host.ready.blue) {
      host.ready[side] = cb;
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
    const game = new GameHost(id, this.clientMoveProvider(id, 'red'), this.clientMoveProvider(id, 'blue'), gameOpt.total + 1, gameOpt.total);
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
    host.move[side] = { moves,  cb };
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
      }, 3000);
      const interval = setInterval(() => {
        if (host.move[side]) {
          clearTimeout(timeout);
          clearInterval(interval);
          cb(host.move[side].moves, state => {
            host.move[side].cb(state);
          });
        } else if (!host.ready[side] || timedOut) {
          clearTimeout(timeout);
          clearInterval(interval);
          const moves = {};
          host.game.getState(side).myTanks.forEach(tank => {
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
