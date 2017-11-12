function receiveGameList () {
  document.querySelector('[name="game-red"]').value = location.origin + '/random-player';
  document.querySelector('[name="game-blue"]').value = location.origin + '/random-player';
  document.querySelector('[name="game-title"]').value = (() => {
    const candidates = [
      '宇宙的答案',
      '世界的尽头',
      '无主之地',
      '孤岛惊魂',
      '世界之门',
      '死胡同',
    ];
    return candidates[Math.floor(Math.random() * candidates.length)];
  })();
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
        return !item.__del && (nameFilter ? item.title.indexOf(nameFilter) > -1 || item.id.indexOf(nameFilter) > -1 : true);
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
  const sse = new EventSource('/game/-events');
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
      fetch(`/game/${id}`, { method: 'DELETE' });
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

function interruptGame (id) {
  fetch(`/game/${id}/interrupt`, { method: 'GET' });
}

function rename(id, $el) {
  const newTitle = prompt('重命名', $el.innerHTML);
  if (newTitle) {
    fetch(`/game/${id}/name`, { method: 'POST', body: newTitle });
  }
}

function toggleAcceptClient(accept) {
  document.querySelectorAll('.accept-client-disable input').forEach(ele => {
    ele.disabled = accept;
  });
}

function toggleCustom(isCustom) {
  document.querySelector('#mapeditor').style.display = isCustom ? 'block' : 'none';
  if (isCustom) {
    validateMap();
  }
}

const validateMap = _.throttle(() => {
  $el = document.querySelector('[name="game-CustomMapEdit"]');
  $err = document.querySelector('#maperror');
  $err.innerHTML = '';
  let width = 0;
  let obstacles = 0;
  let forests = 0;
  const val = $el.value.split('\n').map(line => {
    line = line.trim();
    if (line) {
      line = line.split(/[^0-9]*/).map(v => v | 0);
      line.forEach(v => {
        if (v == 1) {
          obstacles++;
        } else if (v == 2) {
          forests++;
        }
      });
      if (width == 0) {
        width = line.length;
      } else if (width != line.length) {
        width = -1;
      }
      return line;
    }
  }).filter(v => !!v)
  if (width === -1) {
    $err.innerHTML = '地图每一行长度必须一样';
    return;
  }
  const height = val.length;
  document.querySelector('[name="game-MapWidth"]').value = width;
  document.querySelector('[name="game-MapHeight"]').value = height;
  document.querySelector('[name="game-Obstacles"]').value = obstacles;
  document.querySelector('[name="game-Forests"]').value = forests;
  document.querySelector('[name="game-CustomMapValue"]').value = JSON.stringify(val);
}, 1000);

function createGame () {
  const data = {};
  [
    'title', 'total', 'red', 'blue',
    'MapWidth', 'MapHeight', 'InitTank', 'TankHP', 'TankSpeed', 'BulletSpeed', 'FlagTime',
    'Forests', 'Obstacles', 'MaxMoves', 'CustomMapValue',
    'TankScore', 'FlagScore',
  ].forEach(f => {
    data[f] = document.querySelector(`[name="game-${f}"]`).value;
  });
  ['client', 'StaticMap', 'FriendlyFire', 'CustomMap'].forEach(f => {
    data[f] = document.querySelector(`[name="game-${f}"]`).checked;
  });
  fetch('/game', {
    method: 'post',
    body: JSON.stringify(data),
  });
}

async function setupReplay () {
  const id = location.search.substr(1);
  const $stage = document.querySelector('#stage');
  $stage.innerHTML = 'Loading...';
  window.replay = await fetch(`/db/${id}.json`).then(r => r.json());
  const { terain, history, bulletSpeed, tankSpeed } = replay;
  const framesPerRound = bulletSpeed + tankSpeed + 2;
  $stage.innerHTML = '';
  let $style = document.createElement('style');
  document.querySelector('head').appendChild($style);
  const cellSize = 50;
  const setDisplaySize = _.throttle(() => {
    const prefCellSize = Math.floor(Math.min((window.innerWidth - 50) / terain[0].length, (window.innerHeight - 200) / terain.length));
    $style.parentElement.removeChild($style);
    $style = document.createElement('style');
    const playInterval = Math.min(2000, Math.max(50, document.querySelector('#interval').value - 0));
    $style.appendChild(document.createTextNode(`
    .cell-size {
      width: ${cellSize}px;
      height: ${cellSize}px;
    }
    #stagewrap {
      width: ${prefCellSize * terain[0].length}px;
      height: ${prefCellSize * terain.length}px;
      margin-top: ${Math.floor((window.innerHeight - prefCellSize * terain.length) / 3)}px;
      margin-left: ${Math.floor((window.innerWidth - prefCellSize * terain[0].length) / 2)}px;
    }
    #stage {
      width: ${cellSize * terain[0].length}px;
      height: ${cellSize * terain.length}px;
      transform: scale(${prefCellSize / cellSize});
      transform-origin: top left;
    }
    .transition {
      transition: all ${playInterval / 1000}s linear;
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
      $cell.onmouseover = () => {
        document.title = `(${x},${y})`;
      };
    }
    const $br = document.createElement('div');
    $br.className = 'linebreak';
    $terain.appendChild($br);
  }
  const $cell = document.createElement('div');
  $cell.className = 'cell cell-size flag';
  $terain.appendChild($cell);
  const objs = {};
  let hi = 0;
  let paused = false;
  
  window.pause = (isPause = true) => {
    paused = isPause;
    hi = document.querySelector('#pos').value * framesPerRound;
    document.querySelector('#paused').checked = isPause;
  };
  
  document.querySelector('#total').innerHTML = history.length / framesPerRound;
  document.querySelector('#pos').max = history.length / framesPerRound;
  while (true) {
    if (paused || hi >= history.length) {
      await new Promise(cb => setTimeout(cb, 1000));
      continue;
    }
    const playInterval = Math.min(2000, Math.max(50, document.querySelector('#interval').value - 0));
    const ext = history[hi];
    if (ext.flagWait === 0) {
      document.querySelector('.flag').style.display = 'block';
    } else {
      document.querySelector('.flag').style.display = 'none';
    }
    document.querySelector('#flags').innerHTML = `蓝：${ext.blueFlag} 红：${ext.redFlag}`;
    const stamp = hi++;
    const setObj = (type, color) => state => {
      if (!state) {
        return;
      }
      if (!objs[state.id]) {
        const obj = objs[state.id] = {
          $el: document.createElement('div'),
          $mark: document.createElement('div'),
        };
        $stage.appendChild(obj.$el);
        obj.$el.appendChild(obj.$mark);
        obj.$el.title = state.id;
        obj.$mark.className = 'mark';
      }
      const obj = objs[state.id];
      obj.stamp = stamp;
      let direction = state.direction;
      if (direction === 'down' && obj.$el.className.indexOf('direction-left') > -1) {
        direction = 'pre-down-left';
        setTimeout(() => {
          obj.$el.className = `${type} ${type}-${color} direction-${state.direction} cell-size transition`;
        }, 10);
      } else if (direction === 'left' && obj.$el.className.indexOf('direction-down') > -1) {
        direction = 'pre-left-down';
        setTimeout(() => {
          obj.$el.className = `${type} ${type}-${color} direction-${state.direction} cell-size transition`;
        }, 10);
      }
      obj.$el.className = `${type} ${type}-${color} direction-${direction} cell-size transition`;
      Object.assign(obj.$el.style, {
        top: (cellSize * state.y) + 'px',
        left: (cellSize * state.x) + 'px',
      });
      if (state.hp) {
        obj.$mark.innerHTML = state.hp;
      }
    };
    for (let j = 0; j < bulletSpeed; j++) {
      const state = history[hi++];
      state.blueBullet.forEach(setObj('bullet', 'blue'));
      state.redBullet.forEach(setObj('bullet', 'red'));
    }
    for (let j = 0; j < tankSpeed; j++) {
      const state = history[hi++];
      state.blueTank.forEach(setObj('tank', 'blue'));
      state.redTank.forEach(setObj('tank', 'red'));
    }
    const state = history[hi++];
    state.blueBullet.forEach(setObj('bullet', 'blue'));
    state.redBullet.forEach(setObj('bullet', 'red'));
    Object.keys(objs).forEach(id => {
      if (objs[id].stamp !== stamp) {
        const obj = objs[id];
        delete objs[id];
        obj.$el.className += ' destroyed';
        setTimeout(() => {
          obj.$el.parentElement.removeChild(obj.$el);
        }, 600);
      }
    });
    await new Promise(cb => setTimeout(cb, playInterval * (document.querySelector('#lag').checked ? 2 : 1)));
    document.querySelector('#pos').value = hi / framesPerRound;
  }
}