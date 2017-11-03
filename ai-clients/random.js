const frame = require('./_frame');

frame((state, moves) => {
  state.myTank.forEach(tank => {
    moves[tank.id] = (() => {
      switch (Math.floor(Math.random() * 7)) {
        case 0: return 'fire';
        case 1: return 'left';
        case 2: return 'right';
        case 3: return 'stay';
        default: return 'move';
      }
    })();
  });
});
