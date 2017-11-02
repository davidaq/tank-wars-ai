const { EventEmitter } = require('events');
const UUID = require('uuid');
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
    let createGameHost = false;
    if (!isInit) {
      opt.id = UUID.v4();
      opt.createtime = Date.now();
      opt.games = [];
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
}

module.exports = new GameList();
