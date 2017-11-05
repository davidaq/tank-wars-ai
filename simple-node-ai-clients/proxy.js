const frame = require('./_frame');

frame((state, obj) => {
  state.myTank.forEach(tank => {
    obj[tank.id] = {};
  });
});
