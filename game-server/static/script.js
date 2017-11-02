function receiveGameList () {
  document.querySelector('[name="game-red"]').value = location.origin + '/random-player';
  document.querySelector('[name="game-blue"]').value = location.origin + '/random-player';
  const gameList = [];
  const gameMap = {};
  let filteredList = false;
  let filteredName = '';
  const cardTpl = _.template(document.querySelector('#gameCardTemplate').innerHTML);
  const filterGame = _.debounce(() => {
    const perPage = 20;
    const nameFilter = document.querySelector('#gameNameInput').value;
    if (filteredName !== nameFilter) {
      document.querySelector('#gamePageInput').value = 1;
      filteredList = false;
      filteredName = nameFilter;
    }
    if (!filteredList) {
      filteredList = gameList.filter(item => {
        return !item.__del && (nameFilter ? item.title.indexOf(nameFilter) > -1 : true);
      });
      filteredList.reverse();
    }
    $list = document.querySelector('#games');
    $list.innerHTML = '';
    const page = document.querySelector('#gamePageInput').value;
    filteredList.slice((page - 1) * perPage, page * perPage).forEach(item => {
      $el = document.createElement('table');
      $el.innerHTML = cardTpl(item);
      $list.appendChild($el.querySelector('tr'));
    });
  }, 500);
  const sse = new EventSource('/game-list');
  // sse.addEventListener('reset', evt => {
  //   gameList = [];
  //   filteredList = false;
  // });
  sse.addEventListener('game', evt => {
    const game = JSON.parse(evt.data);
    if (gameMap[game.id]) {
      Object.assign(gameMap[game.id], game);
    } else {
      gameMap[game.id] = game;
      gameList.push(game);
    }
    filteredList = false;
    filterGame();
  });
  
  const rmGame = (id) => {
    if (confirm('确认删除？')) {
      fetch(`/game/${id}`, { method: 'delete' });
      gameList.forEach(item => {
        if (item.id === id) {
          item.__del = true;
        }
      });
      filteredList = false;
      filterGame();
    }
  };
  window.filterGame = filterGame;
  window.rmGame = rmGame;
}

function createGame () {
  const data = {};
  ['title', 'total', 'red', 'blue'].forEach(f => {
    data[f] = document.querySelector(`[name="game-${f}"]`).value;
  });
  fetch('/create-game', {
    method: 'post',
    body: JSON.stringify(data),
  });
}

async function setupReplay () {
  const id = location.search.substr(1);
  const $stage = document.querySelector('#stage');
  $stage.innerHTML = 'Loading...';
  const { terain, history } = await fetch(`/db/${id}.json`).then(r => r.json());
  $stage.innerHTML = '';
  let $style = document.createElement('style');
  document.querySelector('head').appendChild($style);
  let cellSize = 0;
  const setDisplaySize = _.throttle(() => {
    cellSize = Math.floor(Math.min((window.innerWidth - 50) / terain[0].length, (window.innerHeight - 200) / terain.length));
    $style.parentElement.removeChild($style);
    $style = document.createElement('style');
    $style.appendChild(document.createTextNode(`
    .cell-size {
      width: ${cellSize}px;
      height: ${cellSize}px;
    }
    #stage {
      width: ${cellSize * terain[0].length}px;
      height: ${cellSize * terain.length}px;
      margin-top: ${Math.floor((window.innerHeight - cellSize * terain.length) / 3)}px;
      margin-left: ${Math.floor((window.innerWidth - cellSize * terain[0].length) / 2)}px;
    }
    .transition {
      transition: all ${Math.max(50, document.querySelector('#interval').value - 0) / 1000}s linear;
    }
    `));
    document.querySelector('head').appendChild($style);
  }, 1000);
  window.onresize = setDisplaySize;
  setInterval(setDisplaySize, 1500);
  setDisplaySize();
  const $terain = document.createElement('div');
  $stage.appendChild($terain);
  $terain.className = 'terain';
  for (let y = 0; y < terain.length; y++) {
    const line = terain[y];
    for (let x = 0; x < line.length; x++) {
      const $cell = document.createElement('div');
      $cell.className = 'cell cell-size cell-' + line[x];
      $terain.appendChild($cell);
    }
    const $br = document.createElement('div');
    $br.className = 'linebreak';
    $terain.appendChild($br);
  }
  const objs = {};
  let hi = 0;
  let paused = false;
  
  window.pause = () => {
    paused = true;
  };
  
  window.resume = () => {
    paused = false;
    hi = document.querySelector('#pos').value - 0;
  };
  
  document.querySelector('#total').innerHTML = history.length;
  while (true) {
    if (paused || hi >= history.length) {
      await new Promise(cb => setTimeout(cb, 1000));
      continue;
    }
    const playInterval = Math.max(50, document.querySelector('#interval').value - 0);
    const state = history[hi];
    const setObj = (type, color) => state => {
      if (!objs[state.id]) {
        const tank = objs[state.id] = {
          $el: document.createElement('div'),
        };
        $stage.appendChild(tank.$el);
      }
      const tank = objs[state.id];
      tank.stamp = hi;
      tank.$el.className = `${type} ${type}-${color} direction-${state.direction} cell-size transition`;
      Object.assign(tank.$el.style, {
        top: (cellSize * state.y) + 'px',
        left: (cellSize * state.x) + 'px',
      });
    };
    state.blueTank.forEach(setObj('tank', 'blue'));
    state.redTank.forEach(setObj('tank', 'red'));
    state.blueBullet.forEach(setObj('bullet', 'blue'));
    state.redBullet.forEach(setObj('bullet', 'red'));
    Object.keys(objs).forEach(id => {
      if (objs[id].stamp !== hi) {
        const obj = objs[id];
        delete objs[id];
        obj.$el.className += ' destroyed';
        setTimeout(() => {
          obj.$el.parentElement.removeChild(obj.$el);
        }, playInterval);
      }
    });
    await new Promise(cb => setTimeout(cb, playInterval));
    hi++;
    document.querySelector('#pos').value = hi;
  }
}