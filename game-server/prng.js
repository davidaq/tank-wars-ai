// pseudo random number generator

class Random {
  constructor (seed = `${Date.now()}`) {
    if (typeof seed === 'string') {
      const v = seed;
      seed = 0;
      for (let i = 0; i < v.length; i++) {
        seed = (seed << 8) | v.charCodeAt(i);
      }
    }
    this._seed = seed % 2147483647;
    if (this._seed <= 0) this._seed += 2147483646;
  }

  next () {
    return this._seed = this._seed * 16807 % 2147483647;
  }

  nextFloat () {
    return (this.next() - 1) / 2147483646;
  }
}

module.exports = Random;
