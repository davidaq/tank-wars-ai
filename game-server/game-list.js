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
    fs.readFile(path.resolve(__dirname, 'db', 'list.txt'), 'utf-8', (err, data) => {
      if (err) {
        return;
      }
      data.split('\n').forEach(content => {
        try {
          this.createGame(JSON.parse(content), true);
        } catch (err) {}
      });
    });
  }

  forEach (cb) {
    this.list.forEach(cb);
  }

  createGame (opt, isInit = false) {
    opt.id = UUID.v4();
    opt.createtime = Date.now();
    opt.redWin = 0;
    opt.blueWin = 0;
    opt.tie = 0;
    this.list.push(opt);
    if (!isInit) {
      this.emit('game', opt);
      fs.appendFile(path.resolve(__dirname, 'db', 'list.txt'), JSON.stringify(opt) + '\n');
    }
    const game = new GameHost(opt.id, opt.red, opt.blue, opt.total);
    game.on('round', winner => {
      if (winner == 'red') {
        opt.redWin++;
      } else if (winner == 'blue') {
        opt.blueWin++;
      } else {
        opt.tie++;
      }
      this.emit('game', opt);
    });
  }
}

module.exports = new GameList();
