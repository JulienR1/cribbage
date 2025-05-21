const params = new URLSearchParams();
const playerId = localStorage.getItem("player-id");
if (playerId) {
  params.append("player-id", playerId);
}
const ws = new WebSocket(location.toString() + "/ws?" + params.toString());

ws.addEventListener("close", function () {
  location.assign("/");
});

ws.addEventListener("message", handleMessage);

const handlers = {
  "game-id": onGameId,
  "player-id": onPlayerId,
  "player-count": onPlayerCount,
};

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
}

/** @param {[string, string]} */
function onPlayerId([_, playerId]) {
  localStorage.setItem("player-id", playerId);
}

/** @param {[string, string]} */
function onPlayerCount([_, count]) {
  console.log("player count changed:", count);
}
