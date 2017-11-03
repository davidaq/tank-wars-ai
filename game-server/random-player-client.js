const fetch = require('isomorphic-fetch');
const co = require('co');

const api = process.argv[2];

co(function* () {
  let state = yield fetch(api, { method: 'GET' }).then(r => r.json(), err => ({ myTank: [] }));
  let i = 0;
  while (!state.ended) {
    console.log(i++, state.events.map(v => v.type));
    const moves = {};
    state.myTank.forEach(function(tank) {
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
    state = yield fetch(api, { method: 'POST', body: JSON.stringify(moves) }).then(r => r.json());
  }
});
