const { EventEmitter } = require('events');

class GameHost extends EventEmitter {
  constructor (id, red, blue, total) {
    super();
  }
  start () {

  }
}

module.exports = GameHost;
