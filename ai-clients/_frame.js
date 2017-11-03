const fetch = require('isomorphic-fetch');
const co = require('co');

const api = process.argv[2];

module.exports = solver => co(function* () {
  let state = yield fetch(api, { method: 'GET' }).then(r => r.json());
  let i = 0;
  while (!state.ended) {
    process.stdout.write(`  ${i++}\t my: ${state.myTank.length}\t enemy: ${state.enemyTank.length}\t\t\t\r`);
    const moves = {};
    solver(state, moves);
    state = yield fetch(api, { method: 'POST', body: JSON.stringify(moves) }).then(r => r.json());
  }
  process.stdout.write(`  ${i++}\t my: ${state.myTank.length}\t enemy: ${state.enemyTank.length}\t\t\t\n`);
});