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
    this.TankScore = Math.max(1, this.TankScore - 0);
    this.FlagScore = Math.max(1, this.FlagScore - 0);
    this.FlagTime = Math.max(1, this.FlagTime - 0);
    this.TankHP = Math.max(1, this.TankHP - 0);
    this.MaxMoves = Math.max(1, this.MaxMoves - 0);
    this.TankSpeed = Math.max(1, this.TankSpeed - 0);
    this.BulletSpeed = Math.max(this.TankSpeed + 1, this.BulletSpeed - 0);
    const tankW = Math.ceil(Math.sqrt(this.InitTank));
    const tankH = Math.ceil(this.InitTank / tankW);
    this.MapWidth = Math.max(tankW * 2, this.MapWidth - 0);
    this.MapHeight = Math.max(tankH * 2, this.MapHeight - 0);
    if (this.MapWidth != settings.MapWidth || this.MapHeight != settings.MapHeight) {
      this.CustomMap = false;
    }
    let remain = (this.MapWidth - 3) * (this.MapHeight - 3) - (tankW + 1) * (tankH + 1) * 3 - 5;
    if (remain < 0) {
      this.CustomMap = false;
    }
    if (!this.CustomMap) {
      this.MapWidth += 1 - (this.MapHeight & 1);
      this.MapHeight += 1 - (this.MapHeight & 1);
    }
    this.Obstacles = Math.min(this.Obstacles - 0, remain);
    remain -= this.Obstacles - 5;
    this.Forests = Math.min(this.Forests - 0, remain);
    this.playRounds();
  }
  playRounds () {
    return co.wrap(function * () {
      for (let i = this.beginRound; i < this.total; i++) {
        try {
          yield this.playRound(i);
        } catch (err) {
          console.error(err.stack);
        }
        this.emit('round', {
          blue: this.blueTank.length > 0 ? this.blueTank.length * this.TankScore + this.blueFlag * this.FlagScore : 0,
          red: this.redTank.length > 0 ? this.redTank.length * this.TankScore + this.redFlag * this.FlagScore : 0,
          moves: this.stepsMoved,
        });
      }
    }).call(this);
  }
  playRound (roundNum) {
    return co.wrap(function * () {
      this.roundId = `${this.id}_${roundNum}`;
      this.history = [];
      this.stepsMoved = 0;
      this.blueEvents = [];
      this.redEvents = [];
      this.flagWait = 0;
      this.redFlag = 0;
      this.blueFlag = 0;
      this.random = this.StaticMap ? new Random(this.id) : new Random(this.roundId);
      console.info('setup terain');
      this.setupTerrain();
      console.info('setup tank');
      this.spawnTank();
      console.info('calc first time state');
      yield* this.calcState();
      console.info('call player setup');
      yield this.callApi('setup');
      let i = 0;
      for (; i < this.MaxMoves; i++) {
        this.stepsMoved = i;
        if (this.blueTank.length === 0 || this.redTank.length === 0) {
          break;
        }
        console.info('call player move');
        yield this.callApi('move');
        this.blueEvents = [];
        this.redEvents = [];
        console.info('calc state');
        yield* this.calcState();
        if (i % 10 === 0) {
          console.info('GAME', this.roundId, i, {
            blueTank: this.blueTank.length,
            redTank: this.redTank.length,
          });
        }
        this.emit('state');
      }
      console.info('round end');
      const fwriter = fs.createWriteStream(path.join(__dirname, 'db', this.roundId + '.json.gz'));
      const writer = zlib.createGzip();
      writer.pipe(fwriter);
      writer.end(JSON.stringify({
        terain: this.terain,
        history: this.history,
        rounds: i,
        bulletSpeed: this.BulletSpeed + 2,
        tankSpeed: this.TankSpeed,
      }));
      writer.on('error', err => null);
      fwriter.on('error', err => null);
      yield new Promise(cb => fwriter.on('finish', () => cb()));
      yield this.callApi('end');
    }).call(this);
  }
  interrupt () {
    this.MaxMoves = 0;
  }
  setupTerrain () {
    this.terain = [];
    this.flagX = Math.floor(this.MapWidth / 2);
    this.flagY = Math.floor(this.MapHeight / 2);
    if (this.CustomMap) {
      this.terain = JSON.parse(this.CustomMapValue);
    } else {
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
        if (x >= 0 && x < this.MapWidth && y >= 0 && y < this.MapHeight && this.terain[y][x] === 0 && x >= tankW && y >= tankH && x !== this.flagX && y != this.flagY) {
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
          if (x >= 0 && x < this.MapWidth && y >= 0 && y < this.MapHeight && this.terain[y][x] === 0 && x >= tankW && y >= tankH && x !== this.flagX && y != this.flagY) {
            this.terain[y][x] = 2;
            this.terain[this.MapHeight - y - 1][this.MapWidth - x - 1] = 2;
            i += 2;
            break;
          }
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
    let M = 0;
    const maxM = Math.min(this.MapWidth, this.MapHeight);
    while (M < maxM && this.blueTank.length < this.InitTank) {
      for (let i = 0; i < M; i++) {
        this.assignTank(i, M);
        this.assignTank(M, i);
      }
      this.assignTank(M, M);
      M++;
    }
  }
  assignTank (x, y) {
    if (this.terain[this.MapHeight - y - 1][x] !== 0 || this.terain[y][this.MapWidth - x - 1] !== 0 || this.blueTank.length >= this.InitTank) {
      return false;
    }
    this.blueTank.push({ bullet: '', hp: this.TankHP, color: 'blue', x, y: this.MapHeight - y - 1, direction: 'down', id: shortid.generate() });
    this.redTank.push({ bullet: '', hp: this.TankHP, color: 'red', x: this.MapWidth - x - 1, y, direction: 'up', id: shortid.generate() });
    return true;
  }
  *calcState () {
    if (this.flagWait > 0) {
      this.flagWait--;
    }
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
      flagWait: this.flagWait,
      blueFlag: this.blueFlag,
      redFlag: this.redFlag,
    });
    for (let i = 0; i < this.BulletSpeed; i++) {
      this.calcStateMoveBullet(scene, this.blueBullet);
      this.calcStateMoveBullet(scene, this.redBullet);
      this.history.push(clone({
        blueBullet: this.blueBullet,
        redBullet: this.redBullet,
      }));
    }
    this.calcStateMoveTank(scene, 'blue', [], [], false, true);
    this.calcStateMoveTank(scene, 'red', [], [], false, true);
    this.history.push(clone({
      blueBullet: this.blueBullet,
      redBullet: this.redBullet,
    }));
    this.calcStateMoveBullet(scene, this.blueBullet, true);
    this.calcStateMoveBullet(scene, this.redBullet, true);
    this.history.push(clone({
      blueBullet: this.blueBullet,
      redBullet: this.redBullet,
    }));
    const bullets = {};
    this.blueBullet.forEach((bullet, i) => {
      const k = `${bullet.x},${bullet.y}`;
      if (!bullets[k]) {
        bullets[k] = [];
      }
      bullets[k].push({ set: this.blueBullet, i });
    });
    this.redBullet.forEach((bullet, i) => {
      const k = `${bullet.x},${bullet.y}`;
      if (!bullets[k]) {
        bullets[k] = [];
      }
      bullets[k].push({ set: this.redBullet, i });
    });

    for (let i = 0; i < this.TankSpeed; i++) {
      let advances = [];
      const forbid = {};
      this.calcStateMoveTank(scene, 'blue', advances, forbid, i > 0, false);
      this.calcStateMoveTank(scene, 'red', advances, forbid, i > 0, false);
      let taken = {};
      const take = (index, val) => {
        if (!taken[index]) {
          taken[index] = [];
        }
        taken[index].push(val);
      };
      advances = advances.filter(item => {
        const { oTank, tank } = item;
        const posIndex = `${tank.x},${tank.y}`;
        const oPosIndex = `${oTank.x},${oTank.y}`
        const forbidState = forbid[posIndex];
        if (forbidState && (!forbidState.direction || forbidState.direction === tank.direction)) {
          this[tank.color + 'Events'].push({
            type: 'collide-tank',
            target: tank.id,
          });
          take(oPosIndex);
          return false;
        } else if (tank.x < 0 || tank.x >= this.MapWidth || tank.y < 0 || tank.y >= this.MapHeight) {
          this[tank.color + 'Events'].push({
            type: 'collide-wall',
            target: tank.id,
          });
          take(oPosIndex);
          return false;
        } else if (this.terain[tank.y][tank.x] === 1) {
          this[tank.color + 'Events'].push({
            type: 'collide-obstacle',
            target: tank.id,
          });
          take(oPosIndex);
          return false;
        } else {
          take(posIndex, item);
          return true;
        }
      });
      let chaining = true;
      while (chaining) {
        chaining = false;
        for (const posIndex of Object.keys(taken)) {
          const takev = taken[posIndex];
          if (takev.length > 1) {
            taken[posIndex] = [];
            chaining = true;
            for (const item of takev) {
              if (item) {
                this[item.tank.color + 'Events'].push({
                  type: 'collide-tank',
                  target: item.tank.id,
                });
                item.invalid = true;
                take(`${item.oTank.x},${item.oTank.y}`);
              } else {
                if (taken[posIndex].length && !taken[posIndex][0]) {
                  throw new Error('Tank movement system error');
                }
                take(posIndex);
              }
            }
          }
        }
      }
      advances = advances.filter(v => !v.invalid);
      advances.forEach(item => {
        const { oTank, tank } = item;
        Object.assign(oTank, tank);
      });
      this.blueTank.forEach((tank, i) => {
        if (tank) {
          const bulletList = bullets[`${tank.x},${tank.y}`];
          if (bulletList) {
            for (const bullet of bulletList) {
              this.hitTank(scene, bullet.set[bullet.i], this.blueTank, i);
              bullet.set[bullet.i] = null;
            }
          }
        }
      });
      this.redTank.forEach((tank, i) => {
        if (tank) {
          const bulletList = bullets[`${tank.x},${tank.y}`];
          if (bulletList) {
            for (const bullet of bulletList) {
              this.hitTank(scene, bullet.set[bullet.i], this.redTank, i);
              bullet.set[bullet.i] = null;
            }
          }
        }
      });
      this.history.push(clone({
        blueTank: this.blueTank,
        redTank: this.redTank,
      }));
      if (this.flagWait === 0) {
        for (const tank of this.blueTank) {
          if (tank && tank.x === this.flagX && tank.y === this.flagY) {
            this.flagWait = this.FlagTime;
            this.blueFlag++;
            this.blueEvents.push({
              type: 'my-flag',
              target: tank.id,
            });
            this.redEvents.push({
              type: 'enemy-flag',
              target: tank.id,
            });
          }
        }
        for (const tank of this.redTank) {
          if (tank && tank.x === this.flagX && tank.y === this.flagY) {
            this.flagWait = this.FlagTime;
            this.redFlag++;
            this.redEvents.push({
              type: 'my-flag',
              target: tank.id,
            });
            this.blueEvents.push({
              type: 'enemy-flag',
              target: tank.id,
            });
          }
        }
      }
    }
    this.blueBullet = this.blueBullet.filter(v => !!v);
    this.redBullet = this.redTank.redBullet(v => !!v);
    this.blueTank = this.blueTank.filter(v => v.hp > 0);
    this.redTank = this.redTank.filter(v => v.hp > 0);
    this.history.push(clone({
      blueBullet: this.blueBullet,
      redBullet: this.redBullet,
      blueTank: this.blueTank,
      redTank: this.redTank,
    }));
  }
  calcStateMoveTank (scene, color, advances, forbid, moveOnly, fireOnly) {
    const myResp = this[color + 'Resp'] || {};
    const myTank = this[color + 'Tank'];
    try {
      for (let i = 0; i < myTank.length; i++) {
        let tank = myTank[i];
        if (!tank) {
          continue;
        }
        const move = myResp[tank.id] || 'stay';
        let skip = false;
        switch (move) {
          case 'fire':
          case 'fire-up':
          case 'fire-left':
          case 'fire-down':
          case 'fire-right':
            skip = !fireOnly;
            break;
          default:
            break;
        }
        if (skip) {
          continue;
        }
        if (move !== 'move') {
          forbid[`${tank.x},${tank.y}`] = { tank };
          if (moveOnly) {
            continue;
          }
        } else {
          forbid[`${tank.x},${tank.y}`] = {
            tank,
            direction: (() => {
              switch (tank.direction) {
                case 'up': return 'down';
                case 'down': return 'up';
                case 'left': return 'right';
                case 'right': return 'left';
                default: return 'all';
              }
            })(),
          };
        }
        switch (move) {
          case 'move': {
            const oTank = tank;
            tank = clone(tank);
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
            advances.push({ oTank, tank });
            break;
          }
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
          case 'back':
            switch (tank.direction) {
              case 'up':
                tank.direction = 'down';
                break;
              case 'down':
                tank.direction = 'up';
                break;
              case 'left':
                tank.direction = 'right';
                break;
              case 'right':
                tank.direction = 'left';
                break;
            }
            break;
          case 'fire':
          case 'fire-up':
          case 'fire-left':
          case 'fire-down':
          case 'fire-right':
            if (!tank.bullet) {
              const direction = move.replace(/^fire-?/, '') || tank.direction;
              tank.bullet = shortid.generate();
              this[tank.color + 'Bullet'].push({
                x: tank.x,
                y: tank.y,
                direction,
                id: tank.bullet,
                from: tank.id,
                color: tank.color,
                round: this.stepsMoved,
              });
            }
            break;
        }
      }
    } catch (err) {
      console.error(err.stack);
    }
  }
  calcStateMoveBullet (scene, myBullet, newOnly) {
    for (let i = 0; i < myBullet.length; i++) {
      const bullet = myBullet[i];
      if (newOnly && bullet.round !== this.stepsMoved) {
        continue;
      }
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
          this.hitTank(scene, bullet, this[target.tank + 'Tank'], target.i);
        }
      }
      if (removeBullet) {
        const fromTank = this[bullet.color + 'Tank'].filter(v => v && v.id === bullet.from)[0];
        if (fromTank) {
          fromTank.bullet = '';
        }
        myBullet.splice(i, 1);
        i--;
      }
    }
  }
  hitTank (scene, bullet, tankSet, tankI) {
    const hitTank = tankSet[tankI];
    const isFriendlyFire = bullet.color === hitTank.color;
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
    const fromTank = this[bullet.color + 'Tank'].filter(v => v && v.id == bullet.from)[0];
    if (fromTank) {
      fromTank.bullet = '';
    }
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
  }
  getState (side) {
    const ended = this.blueTank.length === 0 || this.redTank.length === 0 || this.stepsMoved + 1 >= this.MaxMoves;
    const params = {
      tankScore: this.TankScore,
      flagScore: this.FlagScore,
      flagTime: this.FlagTime,
      tankSpeed: this.TankSpeed,
      bulletSpeed: this.BulletSpeed,
      flagX: this.flagX,
      flagY: this.flagY,
    };
    if (side === 'blue') {
      return {
        terain: this.terain,
        myTank: this.blueTank,
        myBullet: this.blueBullet,
        enemyTank: this.redTank.filter(obj => this.terain[obj.y][obj.x] !== 2),
        enemyBullet: this.redBullet.filter(obj => this.terain[obj.y][obj.x] !== 2),
        events: this.blueEvents,
        flagWait: this.flagWait,
        myFlag: this.blueFlag,
        enemyFlag: this.redFlag,
        params,
        ended,
      };
    } else {
      return {
        terain: this.terain,
        myTank: this.redTank,
        myBullet: this.redBullet,
        enemyTank: this.blueTank.filter(obj => this.terain[obj.y][obj.x] !== 2),
        enemyBullet: this.blueBullet.filter(obj => this.terain[obj.y][obj.x] !== 2),
        events: this.redEvents,
        flagWait: this.flagWait,
        myFlag: this.redFlag,
        enemyFlag: this.blueFlag,
        params,
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
