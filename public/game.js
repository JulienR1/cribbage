const gameId = location.pathname.split("/").reverse()[0] ?? "";
const ws = new WebSocket("/ws");

ws.addEventListener("open", function () {
  console.log("connected to ws");
  ws.send(gameId);
});

ws.addEventListener("close", function () {
  location.assign("/");
});
