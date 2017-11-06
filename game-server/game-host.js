const { EventEmitter } = require('events');
const co = require('co');
const fetch = require('isomorphic-fetch');
const shortid = require('shortid');
const fs = require('fs');
const path = require('path');
const zlib = require('zlib');
const clone = require('clone');
const Random = require('./prng');

class GameHost extends EventEmitter {
  constructor (settings) {
    super();
    this.id = '';
    this.red = '';
    this.blue = '';
    this.total = 1;
    this.beginRound = 0;
    this.MaxMoves = 1000;
    this.MapWidth = 50;
    this.MapHeight = 20;
    this.Obstacles = 35;
    this.InitTank = 5;
    Object.assign(this, settings);
    this.FriendlyFire = !!this.FriendlyFire;
    this.StaticMap = !!this.StaticMap;
    this.InitTank = Math.max(1, this.InitTank - 0);
    this.TankHP = Math.max(1, this.TankHP - 0);
    this.MaxMoves = Math.max(1, this.MaxMoves - 0);
    this.TankSpeed = Math.max(1, this.TankSpeed - 0);
    this.BulletSpeed = Math.max(this.TankSpeed + 1, this.BulletSpeed - 0);
    const tankW = Math.ceil(Math.sqrt(this.InitTank));
    const tankH = Math.ceil(this.InitTank / tankW);
    this.MapWidth = Math.max(tankW * 2, this.MapWidth - 0);
    this.MapHeight = Math.max(tankH * 2, this.MapHeight - 0);
    this.MapWidth += 1 - (this.MapHeight & 1);
    this.MapHeight += 1 - (this.MapHeight & 1);
    let remain = this.MapWidth * this.MapHeight - tankW * tankH * 2 - 3;
    this.Obstacles = Math.min(this.Obstacles - 0, remain);
    remain -= this.Obstacles - 3;
    this.Forests = Math.min(this.Forests - 0, remain);
    this.playRounds();
  }
  playRounds () {
    return co.wrap(function * () {
      for (let i = this.beginRound; i < this.total; i++) {
        yield this.playRound(i);
      }
    }).call(this).then(() => null, err => {
      console.error(err.stack);
      this.emit('round', {
        blue: this.blueTank.length,
        red: this.redTank.length,
        moves: i,
      });
    });
  }
  playRound (roundNum) {
    return co.wrap(function * () {
      this.roundId = `${this.id}_${roundNum}`;
      this.history = [];
      this.stepsMoved = 0;
      this.blueEvents = [];
      this.redEvents = [];
      this.random = this.StaticMap ? new Random(this.id) : new Random(this.roundId);
      this.setupTerrain();
      this.spawnTank();
      yield this.callApi('setup');
      let i = 0;
      for (; i < this.MaxMoves; i++) {
        this.stepsMoved = i;
        if (this.blueTank.length === 0 || this.redTank.length === 0) {
          break;
        }
        yield this.callApi('move');
        this.blueEvents = [];
        this.redEvents = [];
        this.calcState();
        if (i % 10 === 0) {
          console.log('GAME', this.roundId, i, {
            blueTank: this.blueTank.length,
            redTank: this.redTank.length,
          });
        }
        this.emit('state');
      }
      const fwriter = fs.createWriteStream(path.join(__dirname, 'db', this.roundId + '.json.gz'));
      const writer = zlib.createGzip();
      writer.pipe(fwriter);
      writer.end(JSON.stringify({
        terain: this.terain,
        history: this.history,
        rounds: i,
        bulletSpeed: this.BulletSpeed,
        tankSpeed: this.TankSpeed,
      }));
      writer.on('error', err => null);
      fwriter.on('error', err => null);
      yield new Promise(cb => {
        fwriter.on('finish', () => {
          this.emit('round', {
            blue: this.blueTank.length,
            red: this.redTank.length,
            moves: i,
          });
          cb();
        });
      });
      yield this.callApi('end');
    }).call(this);
  }
  interrupt () {
    this.MaxMoves = 0;
  }
  setupTerrain () {
    this.terain = [];
    for (let y = 0; y < this.MapHeight; y++) {
      const line = [];
      this.terain.push(line);
      for (let x = 0; x < this.MapWidth; x++) {
        line.push(0);
      }
    }
    let x = Math.floor(this.random.nextFloat() * this.MapWidth / 2);
    let y = Math.floor(this.random.nextFloat() * this.MapHeight);
    const tankW = Math.ceil(Math.sqrt(this.InitTank));
    const tankH = Math.ceil(this.InitTank / tankW);
    for (let i = 0; i < this.Obstacles;) {
      switch (Math.floor(this.random.nextFloat() * 4)) {
        case 0:
          x++;
          break;
        case 1:
          y++;
          break;
        case 2:
          x--;
          break;
        case 3:
          y--;
          break;
        default:
          x = Math.floor(this.random.nextFloat() * this.MapWidth / 2);
          y = Math.floor(this.random.nextFloat() * this.MapHeight);
          break;
      }
      if (x >= 0 && x < this.MapWidth && y >= 0 && y < this.MapHeight && this.terain[y][x] === 0 && x >= tankW && y >= tankH) {
        this.terain[y][x] = 1;
        this.terain[this.MapHeight - y - 1][this.MapWidth - x - 1] = 1;
        i += 2;
      } else {
        x = Math.floor(this.random.nextFloat() * this.MapWidth / 2);
        y = Math.floor(this.random.nextFloat() * this.MapHeight);
      }
    }
    for (let i = 0; i < this.Forests;) {
      while (true) {
        const x = Math.floor(this.random.nextFloat() * this.MapWidth / 2);
        const y = Math.floor(this.random.nextFloat() * this.MapHeight);
        if (x >= 0 && x < this.MapWidth && y >= 0 && y < this.MapHeight && this.terain[y][x] === 0 && x >= tankW && y >= tankH) {
          this.terain[y][x] = 2;
          this.terain[this.MapHeight - y - 1][this.MapWidth - x - 1] = 2;
          i += 2;
          break;
        }
      }
    }
  }
  spawnTank () {
    this.blueTank = [];
    this.redTank = [];
    this.blueBullet = [];
    this.redBullet = [];
    const tankW = Math.ceil(Math.sqrt(this.InitTank));
    for (let i = 0; i < this.InitTank; i++) {
      const x = Math.floor(i % tankW);
      const y = Math.floor(i / tankW);
      this.blueTank.push({ hp: this.TankHP, color: 'blue', x, y, direction: 'down', id: shortid.generate() });
      this.redTank.push({ hp: this.TankHP, color: 'red', x: this.MapWidth - x - 1, y: this.MapHeight - y - 1, direction: 'up', id: shortid.generate() });
    }
  }
  calcState () {
    const scene = clone(this.terain);
    for (let i = 0; i < this.blueTank.length; i++) {
      const tank = this.blueTank[i];
      scene[tank.y][tank.x] = { tank: 'blue', i };
    }
    for (let i = 0; i < this.redTank.length; i++) {
      const tank = this.redTank[i];
      scene[tank.y][tank.x] = { tank: 'red', i };
    }
    this.history.push({
      blueMove: this.blueResp,
      redMove: this.redResp,
    });
    for (let i = 0; i < this.BulletSpeed; i++) {
      this.calcStateMoveBullet(scene, this.blueBullet);
      this.calcStateMoveBullet(scene, this.redBullet);
      this.history.push(clone({
        blueTank: this.blueTank,
        blueBullet: this.blueBullet,
        redTank: this.redTank,
        redBullet: this.redBullet,
      }));
    }
    for (let i = 0; i < this.TankSpeed; i++) {
      this.calcStateMoveTank(scene, this.blueTank, this.blueResp, this.blueBullet);
      this.calcStateMoveTank(scene, this.redTank, this.redResp, this.redBullet);
      this.history.push(clone({
        blueTank: this.blueTank,
        blueBullet: this.blueBullet,
        redTank: this.redTank,
        redBullet: this.redBullet,
      }));
    }
    this.blueTank = this.blueTank.filter(v => !!v);
    this.redTank = this.redTank.filter(v => !!v);
  }
  calcStateMoveTank (scene, myTank, myResp, myBullet) {
    if (!myResp) {
      myResp = {};
    }
    try {
      for (let i = 0; i < myTank.length; i++) {
        if (!myTank[i]) {
          continue;
        }
        const tank = clone(myTank[i]);
        const move = myResp[tank.id];
        if (move) {
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
            case 'fire-up':
            case 'fire-left':
            case 'fire-down':
            case 'fire-right':
              const direction = move.replace(/^fire-?/, '') || tank.direction;
              myBullet.push({
                x: tank.x,
                y: tank.y,
                direction,
                id: shortid.generate(),
                from: tank.id,
                color: tank.color,
              });
              break;
          }
          if (move === 'move') {
            if (tank.x < 0 || tank.x >= this.MapWidth || tank.y < 0 || tank.y >= this.MapHeight) {
              this[tank.color + 'Events'].push({
                type: 'collide-wall',
                target: tank.id,
              });
              continue;
            }
            if (scene[tank.y][tank.x] !== 0 && scene[tank.y][tank.x] !== 2) {
              this[tank.color + 'Events'].push({
                type: typeof scene[tank.y][tank.x] == 'number' ? 'collide-obstacle' : 'collide-tank',
                target: tank.id,
              });
              continue;
            }
            const oTank = myTank[i];
            scene[tank.y][tank.x] = scene[oTank.y][oTank.x];
            scene[oTank.y][oTank.x] = this.terain[oTank.y][oTank.x];
          }
          Object.assign(myTank[i], tank);
        }
      }
    } catch (err) {
      console.error(err.stack);
    }
  }
  calcStateMoveBullet (scene, myBullet) {
    for (let i = 0; i < myBullet.length; i++) {
      const bullet = myBullet[i];
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
      if (bullet.x < 0 || bullet.x >= this.MapWidth || bullet.y < 0 || bullet.y >= this.MapHeight) {
        removeBullet = true;
      } else {
        const target = scene[bullet.y][bullet.x];
        if (target === 1) {
          removeBullet = true;
        } else if (target.tank) {
          removeBullet = true;
          const isFriendlyFire = bullet.color === target.tank;
          if (!this.FriendlyFire && isFriendlyFire) {
            removeBullet = false;
          } else {
            let blueEventType;
            let redEventType;
            if (isFriendlyFire) {
              if (bullet.color === 'red') {
                redEventType = 'me-hit-me';
                blueEventType = 'enemy-hit-enemy';
              } else {
                blueEventType = 'me-hit-me';
                redEventType = 'enemy-hit-enemy';
              }
            } else {
              if (bullet.color === 'red') {
                redEventType = 'me-hit-enemy';
                blueEventType = 'enemy-hit-me';
              } else {
                blueEventType = 'me-hit-enemy';
                redEventType = 'enemy-hit-me';
              }
            }
            scene[bullet.y][bullet.x] = this.terain[bullet.y][bullet.x];
            const hitTank = this[target.tank + 'Tank'][target.i];
            this.redEvents.push({
              type: redEventType,
              from: bullet.from,
              target: hitTank.id,
            });
            this.blueEvents.push({
              type: blueEventType,
              from: bullet.from,
              target: hitTank.id,
            });
            hitTank.hp--;
            if (hitTank.hp == 0) {
              this[target.tank + 'Tank'][target.i] = null;
            }
          }
        }
      }
      if (removeBullet) {
        myBullet.splice(i, 1);
        i--;
      }
    }
  }
  getState (side) {
    const ended = this.blueTank.length === 0 || this.redTank.length === 0 || this.stepsMoved + 1 >= this.MaxMoves;
    if (side === 'blue') {
      return {
        terain: this.terain,
        myTank: this.blueTank,
        myBullet: this.blueBullet,
        enemyTank: this.redTank,
        enemyBullet: this.redBullet,
        events: this.blueEvents,
        ended,
      };
    } else {
      return {
        terain: this.terain,
        myTank: this.redTank,
        myBullet: this.redBullet,
        enemyTank: this.blueTank,
        enemyBullet: this.blueBullet,
        events: this.redEvents,
        ended,
      };
    }
  }
  callApi (action, side) {
    if (!side) {
      return new Promise((resolve, reject) => {
        const interval = setInterval(() => {
          if (this.MaxMoves === 0) {
            clearInterval(interval);
            resolve({});
            resolve = null;
            reject = null;
          }
        }, 1000);
        Promise.all([
          this.callApi(action, 'red'),
          this.callApi(action, 'blue'),
        ]).then(v => {
          clearInterval(interval);
          resolve && resolve(v);
        }, err => {
          clearInterval(interval);
          reject && reject(err);
        });
      });
    }
    return co.wrap(function * () {
      this[side + 'Resp'] = false;
      if (typeof this[side] === 'function') {
        if (action === 'move') {
          this[side + 'Resp'] = yield new Promise(cb => {
            this[side]((moves, waitCalc) => {
              this.once('state', () => {
                waitCalc(this.getState(side));
              });
              cb(moves);
            });
          });
        }
      } else {
        this[side + 'Resp'] = yield fetch(this[side], {
          method: 'POST',
          headers: { 'content-type': 'application/json' },
          body: JSON.stringify({
            uuid: this.roundId,
            action: action,
            state: this.getState(side),
          }),
        }).then(r => r.json()).catch(err => []);
      }
    }).call(this);
  }
}

module.exports = GameHost;
