const { EventEmitter } = require('events');
const co = require('co');
const fetch = require('isomorphic-fetch');
const UUID = require('uuid');
const fs = require('fs');
const path = require('path');
const clone = require('clone');

const MapWidth = 50;
const MapHeight = 30;
const Obstacles = 50;
const InitTank = 5;

class GameHost extends EventEmitter {
  constructor (id, red, blue, total) {
    super();
    this.id = id;
    this.red = red;
    this.blue = blue;
    this.total = total;
    this.playRounds();
  }
  playRounds () {
    return co.wrap(function * () {
      for (let i = 0; i < this.total; i++) {
        yield this.playRound(i);
      }
    }).call(this).then(() => null, err => console.error(err.stack));
  }
  playRound (roundNum) {
    return co.wrap(function * () {
      this.roundId = `${this.id}_${roundNum}`;
      this.history = [];
      this.setupTerrain();
      this.spawnTank();
      yield this.callApi('setup');
      let i = 0;
      for (; i < 1000; i++) {
        if (this.blueTank.length === 0 || this.redTank.length === 0) {
          break;
        }
        yield this.callApi('move');
        this.calcState();
      }
      fs.writeFile(path.join(__dirname, 'db', this.roundId + '.json'), JSON.stringify({
        terain: this.terain,
        history: this.history,
      }), err => {
        this.emit('round', {
          blue: this.blueTank.length,
          red: this.redTank.length,
          moves: i,
        });
      });
      yield this.callApi('end');
    }).call(this);
  }
  setupTerrain () {
    this.terain = [];
    for (let y = 0; y < MapHeight; y++) {
      const line = [];
      this.terain.push(line);
      for (let x = 0; x < MapWidth; x++) {
        line.push(0);
      }
    }
    for (let i = 0; i < Obstacles; i++) {
      let x = Math.floor(Math.random() * MapWidth);
      let y = Math.floor(Math.random() * MapHeight);
      this.terain[y][x] = 1;
      this.terain[y][MapWidth - x - 1] = 1;
      switch (Math.floor(Math.random() * 3)) {
        case 0:
          x++;
          if (x < MapWidth) {
            this.terain[y][x] = 1;
            this.terain[y][MapWidth - x - 1] = 1;
          }
          break;
        case 1:
          y++;
          if (y < MapHeight) {
            this.terain[y][x] = 1;
            this.terain[y][MapWidth - x - 1] = 1;
          }
          break;
        default:
          break;
      }
    }
  }
  spawnTank () {
    this.blueTank = [];
    this.redTank = [];
    this.blueBullet = [];
    this.redBullet = [];
    for (let i = 0; i < InitTank; i++) {
      while (true) {
        const x = Math.floor(Math.random() * MapWidth);
        const y = Math.floor(Math.random() * MapHeight);
        if (this.terain[y][x] === 0) {
          this.blueTank.push({ x, y, direction: 'right', id: UUID.v4() });
          this.redTank.push({ x: MapWidth - x - 1, y, direction: 'left', id: UUID.v4() });
          break;
        }
      }
    }
  }
  calcState () {
    const scene = clone(this.terain);
    for (let i = 0; i < this.blueTank.length; i++) {
      const tank = this.blueTank[i];
      scene[tank.y][tank.x] = { t: 'b', i };
    }
    for (let i = 0; i < this.redTank.length; i++) {
      const tank = this.redTank[i];
      scene[tank.y][tank.x] = { t: 'r', i };
    }
    const removedBullet = {};
    this.calcStateMoveBullet(scene, this.blueBullet, removedBullet);
    this.calcStateMoveBullet(scene, this.redBullet, removedBullet);
    this.calcStateMoveTank(scene, this.blueTank, this.blueResp, this.blueBullet);
    this.calcStateMoveTank(scene, this.redTank, this.redResp, this.redBullet);
    this.history.push(clone({
      blueTank: this.blueTank,
      blueBullet: this.blueBullet,
      redTank: this.redTank,
      redBullet: this.redBullet,
    }));
    this.blueBullet = this.blueBullet.filter(bullet => !removedBullet[bullet.id]);
    this.redBullet = this.redBullet.filter(bullet => !removedBullet[bullet.id]);
  }
  calcStateMoveTank (scene, myTank, myResp, myBullet) {
    try {
      for (let i = 0; i < myTank.length; i++) {
        const move = myResp[i];
        if (move) {
          const tank = clone(myTank[i]);
          switch (move) {
            case 'move':
              switch (tank.direction) {
                case 'up':
                  tank.y--;
                  break;
                case 'down':
                  tank.y++;
                  break;
                case 'left':
                  tank.x--;
                  break;
                case 'right':
                  tank.x++;
                  break;
              }
              break;
            case 'left':
              switch (tank.direction) {
                case 'up':
                  tank.direction = 'left';
                  break;
                case 'down':
                  tank.direction = 'right';
                  break;
                case 'left':
                  tank.direction = 'down';
                  break;
                case 'right':
                  tank.direction = 'up';
                  break;
              }
              break;
            case 'right':
              switch (tank.direction) {
                case 'up':
                  tank.direction = 'right';
                  break;
                case 'down':
                  tank.direction = 'left';
                  break;
                case 'left':
                  tank.direction = 'up';
                  break;
                case 'right':
                  tank.direction = 'down';
                  break;
              }
              break;
            case 'fire':
              myBullet.push({
                x: tank.x,
                y: tank.y,
                direction: tank.direction,
                id: UUID.v4(),
              });
              break;
          }
          if (move === 'move') {
            if (tank.x < 0 || tank.x >= MapWidth || tank.y < 0 || tank.y >= MapHeight) {
              continue;
            }
            if (scene[tank.y][tank.x] !== 0) {
              continue;
            }
          }
          Object.assign(myTank[i], tank);
        }
      }
    } catch (err) {
      console.error(err.stack);
    }
  }
  calcStateMoveBullet (scene, myBullet, removedBullet) {
    for (let i = 0; i < myBullet.length; i++) {
      const bullet = myBullet[i];
      for (let j = 0; j < 2; j++) {
        switch (bullet.direction) {
          case 'up':
            bullet.y--;
            break;
          case 'down':
            bullet.y++;
            break;
          case 'left':
            bullet.x--;
            break;
          case 'right':
            bullet.x++;
            break;
        }
        let removeBullet = false;
        if (bullet.x < 0 || bullet.x >= MapWidth || bullet.y < 0 || bullet.y >= MapHeight) {
          removeBullet = true;
        } else {
          const target = scene[bullet.y][bullet.x];
          if (target) {
            if (target === 1) {
              // hit wall
            } else if (target.t === 'r') {
              this.redTank.splice(target.i, 1);
            } else if (target.t === 'b') {
              this.blueTank.splice(target.i, 1);
            }
            removeBullet = true;
          }
        }
        if (removeBullet) {
          removedBullet[bullet.id] = 1;
          break;
        }
      }
    }
  }
  callApi (action) {
    this.blueResp = false;
    this.redResp = false;
    return Promise.all([
      fetch(this.blue, {
        method: 'POST',
        headers: { 'content-type': 'application/json' },
        body: JSON.stringify({
          uuid: this.roundId,
          action: action,
          state: {
            terain: this.terain,
            myTank: this.blueTank,
            myBullet: this.blueBullet,
            opponentTank: this.redTank,
            opponentBullet: this.redBullet,
          },
        }),
      }).then(r => r.json()).then(v => this.blueResp = v, err => this.blueResp = []),
      fetch(this.red, {
        method: 'POST',
        headers: { 'content-type': 'application/json' },
        body: JSON.stringify({
          uuid: this.roundId,
          action: action,
          state: {
            terain: this.terain,
            myTank: this.redTank,
            myBullet: this.redBullet,
            opponentTank: this.blueTank,
            opponentBullet: this.blueBullet,
          },
        }),
      }).then(r => r.json()).then(v => this.redResp = v, err => this.redResp = []),
    ]);
  }
}

module.exports = GameHost;
