const frame = require('./_frame');

let fired;

frame((state, obj) => {
  console.log(state.radar);
  if (!fired) {
    fired = {};
    state.myTank.forEach(tank => {
      fired[tank.id] = 0;
    });
  }
  state.myTank.forEach(tank => {
    let fire = state.radar.fire[tank.id];
    fire = [fire.up, fire.left, fire.down, fire.right].filter(v => v);
    if (fire.length > 0) {
      obj[tank.id] = fire[0].action;
      return;
    }
    let dist = 100;
    let target = null;
    state.enemyTank.forEach(enemy => {
      let d = Math.abs(enemy.x - tank.x) + Math.abs(enemy.y - tank.y);
      if (d < dist) {
        dist = d;
        target = enemy;
      }
    });
    if (!target) {
      return;
    }
    // if (fired[tank.id] == 0) {
    //   if (target.x == tank.x) {
    //     if (Math.abs(target.y - tank.y) < 10) {
    //       fired[tank.id] = 4;
    //       obj[tank.id] = {
    //         action: target.y < tank.y ? 'fire-up' : 'fire-down',
    //       };
    //       return;
    //     }
    //   } else if (target.y == tank.y) {
    //     if (Math.abs(target.x - tank.x) < 10) {
    //       fired[tank.id] = 4;
    //       obj[tank.id] = {
    //         action: target.x < tank.x ? 'fire-left' : 'fire-right',
    //       };
    //       return;
    //     }
    //   }
    // } else {
    //   fired[tank.id]--;
    // }
    obj[tank.id] = {
      // x: Math.floor(state.terain[0].length / 2),
      // y: Math.floor(state.terain.length / 2),
      x: target.x,
      y: target.y,
      action: 'travel-with-dodge',
    };
  });
});
