const store = { gameId: "", playerId: "" };

const ws = new WebSocket(location.toString() + "/ws");

ws.addEventListener("close", function () {
  location.assign("/");
});

ws.addEventListener("message", handleMessage);

const handlers = {
  "game-id": onGameId,
  "player-id": onPlayerId,
  "player-change": onPlayerChange,
};

const playerList = document.getElementById("players");

/** @param {{data:string}} */
function handleMessage({ data }) {
  const sections = data.split(":");
  if (sections.length !== 2) {
    throw new Error("received invalid message from server");
  }

  const [opcode] = sections;
  (handlers[opcode] ?? unhandler)(sections);
}

/** @param {[string, string]} */
function unhandler([opcode]) {
  console.error("unhandled opcode:", opcode);
}

/** @param {[string, string]} */
function onGameId([_, gameId]) {
  console.log("game id:", gameId);
  store.gameId = gameId;
}

/** @param {[string, string]} */
function onPlayerId([_, playerId]) {
  store.playerId = playerId;
}

/** @param {[string, string]} */
async function onPlayerChange([_, playerId]) {
  const element = document.createElement("div");
  element.innerHTML = await html(`/${store.gameId}/players/${playerId}`);

  playerList.querySelector(`#player-${playerId}`)?.replaceWith(element) ??
    playerList.appendChild(element);
}

/**
 * @param  {number} opcode
 * @param {number[]} payload
 **/
function write(opcode, payload = []) {
  ws.send(new Uint8Array([opcode, ...payload]));
}

/**
 * @param url  {string} url
 * @returns {Promise<string>}
 */
async function html(url) {
  const response = await fetch(url);
  return response.text();
}
