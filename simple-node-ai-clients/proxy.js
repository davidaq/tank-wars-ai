const frame = require('./_frame');

frame((state, obj) => {
  state.myTank.forEach((tank, i) => {
    obj[tank.id] = {
      x: state.enemyTank[0].x,
      y: state.enemyTank[0].y,
      attack: 10000,
      travel: 10000,
      dodge: 0,
    };
  });
});
