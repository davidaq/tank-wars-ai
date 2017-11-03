const frame = require('./_frame');

const DIR = {
  up: 0,
  left: 1,
  down: 2,
  right: 3,
};

let init = true;
const temper = {};
const frustrate = {};
const cool = {};

frame((state, moves) => {
  if (init) {
    init = false;
    state.myTank.forEach(my => {
      temper[my.id] = 0;
      frustrate[my.id] = 0;
      cool[my.id] = 0;
    });
  }
  const collide = {};
  state.events.forEach(event => {
    if (event.type.indexOf('collide-') > -1 && frustrate[event.target] <= 0) {
      collide[event.target] = 1;
      temper[event.target]++;
      frustrate[event.target] = temper[event.target] * 2;
    }
  });
  state.myTank.forEach(my => {
    if (collide[my.id]) {
      moves[my.id] = (() => {
        switch (Math.floor(Math.random() * 2)) {
          case 0: return 'left';
          default: return 'right';
        }
      })();
    } else {
      moves[my.id] = (() => {
        switch (Math.floor(Math.random() * 5)) {
          case 1: return 'left';
          case 2: return 'right';
          default: return 'move';
        }
      })();
    }
    if (frustrate[my.id] < -12) {
      temper[my.id]++;
      frustrate[my.id] = temper[my.id] + 3;
    }
    if (frustrate[my.id] > 0) {
      frustrate[my.id]--;
      return;
    }
    let nearest = null;
    let nearestLen = 0;
    state.enemyTank.forEach(enemy => {
      let len = Math.abs(enemy.x - my.x) + Math.abs(enemy.y - my.y);
      if (!nearest || len < nearestLen) {
        nearestLen = len;
        nearest = enemy;
      }
    });
    let moveDir;
    let shoot = false;
    let moveH = nearest.y === my.y || Math.abs(nearest.x - my.x) < Math.abs(nearest.y - my.y);
    if (nearest.x === my.x) {
      moveH = false;
    }
    if (moveH && (state.terain[my.y][my.x - 1] !== 0 || state.terain[my.y][my.x + 1] !== 0)) {
      moveH = false;
    }
    try {
      if (!moveH && (state.terain[my.y - 1][my.x] !== 0 || state.terain[my.y + 1][my.x] !== 0)) {
        moveH = true;
      }
    } catch (err) {
      moveH = true;
    }
    if (moveH) {
      if (nearest.x > my.x) {
        moveDir = 'right';
        shoot = nearest.y === my.y && nearest.x - my.x < 5;
      } else {
        moveDir = 'left';
        shoot = nearest.y === my.y && my.x - nearest.x < 5;
      }
    } else {
      if (nearest.y > my.y) {
        moveDir = 'down';
        shoot = nearest.x === my.x && nearest.y - my.y < 5;
      } else {
        moveDir = 'up';
        shoot = nearest.x === my.x && my.y - nearest.y < 5;
      }
    }
    if (moveDir === my.direction) {
      if (temper[my.id] > 0) {
        cool[my.id]++;
        if (cool[my.id] > 10) {
          cool[my.id] = 0;
          temper[my.id]--;
        }
      } else {
        cool[my.id] = 0;
      }
      if (shoot) {
        moves[my.id] = 'fire';
        frustrate[my.id] -= 4;
      } else {
        moves[my.id] = 'move';
      }
    } else {
      frustrate[my.id] -= 1;
      if ((DIR[moveDir] - DIR[my.direction] + 4) % 4 === 1) {
        moves[my.id] = 'left';
      } else {
        moves[my.id] = 'right';
      }
    }
  });
});
