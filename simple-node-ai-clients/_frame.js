const fetch = require('isomorphic-fetch');
const co = require('co');

const api = process.argv[2] || 'http://localhost:8776';

module.exports = solver => co(function* () {
  let state;
  while (true) {
    try {
      state = yield fetch(api, { method: 'GET' }).then(r => r.json());
      break;
    } catch (err) {
      console.log('Not ready...');
      yield new Promise(cb => setTimeout(cb), 1000);
    }
  }
  let i = 0;
  while (!state.ended) {
    process.stdout.write(`  ${i++}\t my: ${state.myTank.length}\t enemy: ${state.enemyTank.length}\t\t\t\r`);
    const moves = {};
    solver(state, moves);
    state = yield fetch(api, { method: 'POST', body: JSON.stringify(moves) }).then(r => r.json());
  }
  process.stdout.write(`  ${i++}\t my: ${state.myTank.length}\t enemy: ${state.enemyTank.length}\t\t\t\n`);
}).catch(err => console.error(err.stack));