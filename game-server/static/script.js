function receiveGameList () {
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
        return nameFilter ? item.title.indexOf(nameFilter) > -1 : true;
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
  window.filterGame = filterGame;
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
