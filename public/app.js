prependChild = function(container, element) {
    if (container.firstChild)
        return container.insertBefore(element, container.firstChild);
    else
        return container.appendChild(element);
};

document.getElementById("usr").onkeyup = function(e) {
  if (e.keyCode == 13) {
    join();
  }
};

function join() {
  var username = document.getElementById("usr").value;
  var proto = location.protocol == "http:" ? "ws" : "wss";
  var host = window.location.host
  var url = proto + "://" +  host + "/ws?username=" + encodeURIComponent(username);
  var ws = new WebSocket(url);
  window.onbeforeunload = function() { ws.close() };

  ws.onerror = function(e) {
    alert("Connection error")
  }

  ws.onopen = function() {
    var msgs = document.getElementById("messages");
    document.getElementById("join").style.display = "none";
    document.getElementById("chat").style.display = "block";

    ws.onmessage = function(e) {
      var msg = JSON.parse(e.data);

      var node = document.createElement("p");

      var k = document.createElement("span");
      k.setAttribute("class", msg["Prefix"] ? "channel" : "user");
      k.appendChild(document.createTextNode(msg["Prefix"]));
      node.appendChild(k);

      node.appendChild(document.createTextNode(' '));

      var v = document.createElement("span");
      v.setAttribute("class", "message");
      v.appendChild(document.createTextNode(msg["Params"]));
      node.appendChild(v);

      prependChild(msgs, node);
    };

    var roomInsert = document.getElementById("room");
    document.getElementById("send").onkeyup = function(e) {
      if (e.keyCode == 13) {
        ws.send(
          JSON.stringify({
              Command: "privmsg",
              Params: [roomInsert.value, e.target.value],
          })
        );
        e.target.value = "";
      }
    }
  }
}
