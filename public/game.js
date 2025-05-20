const gameId = location.pathname.split("/").reverse()[0] ?? "";
const ws = new WebSocket("/ws");

ws.addEventListener("open", function () {
  console.log("connected to ws");
  ws.send(gameId);
});

ws.addEventListener("close", function () {
  location.assign("/");
});

ws.addEventListener("message", handleMessage);

const handlers = {
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
function onPlayerCount([_, count]) {
  console.log("player count changed:", count);
}
