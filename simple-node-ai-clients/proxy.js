const frame = require('./_frame');

frame((state, obj) => {
  state.myTank.forEach(tank => {
    obj[tank.id] = {
      x: state.terain[0].length - 1,
      y: state.terain.length - 1,
      direction: 'up',
      attack: 1,
      travel: 2,
      dodge: 0,
    };
  });
});
