const frame = require('./_frame');

frame((state, obj) => {
  state.myTank.forEach(tank => {
    let dist = 100;
    let target = null;
    state.enemyTank.forEach(enemy => {
      let d = Math.abs(enemy.x - tank.x) + Math.abs(enemy.y - tank.y);
      if (d < dist) {
        dist = d;
        target = enemy;
      }
    })
    obj[tank.id] = {
      x: 0,//target.x,
      y: 0,//target.y,
      action: 'travel-with-dodge',
    };
  });
  console.log(obj);
});
