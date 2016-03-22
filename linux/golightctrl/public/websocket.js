// (c) Bernhard Tittelbach
"use strict";

function hasWebSocketSupport() {
  // see: https://github.com/Modernizr/Modernizr/blob/115aa4a583dff0b44409b000165d22f2577ab0a1/feature-detects/websockets.js
  var supports = false;
  try {
    supports = 'WebSocket' in window && window.WebSocket.CLOSING === 2;
  } catch (e) {}
  return supports;
}

window.addEventListener('beforeunload', function() {
	ws.close();
});


var ws = {};
ws.onopen=false;
ws.ondisconnect=false;
ws.retrytimer=false;
ws.contexts = {};
ws.retryinterval_ms = 1500;

ws.waitandreconnect = function(uri) {
	if (ws.retrytimer == false)
	{
		ws.retrytimer = setTimeout(function () {
			console.log('attempting reconnect!');
			ws.ws = null;
			ws.open(uri);
			ws.retrytimer=false;
		}, ws.retryinterval_ms);
	}
}

ws.stopreconnecting = function() {
	if (ws.retrytimer) {
		clearTimeout(ws.retrytimer);
		ws.retrytimer=false;
	}
}

ws.open = function(uri) {
	ws.stopreconnecting();
	ws.ws=new WebSocket(uri);
	ws.ws.onmessage = function(response){
		var m = JSON.parse(response.data);
		if (m["ctx"] && m["data"] && typeof(ws.contexts[m.ctx]) == "function") {
			ws.contexts[m.ctx](m.data);
		}
	}
 	ws.ws.onopen = function(){
		ws.ws.onclose = function(){
			console.log('Connection to server lost. reconnecting...');
			if (typeof(ws["ondisconnect"]) == "function"){
				ws.ondisconnect();
			}
			ws.waitandreconnect(uri);
		}
		if (typeof(ws["onopen"]) == "function"){
			ws.onopen();
		}
 	}
 	ws.ws.onerror = function() {
 		ws.waitandreconnect(uri);
 	}
}

ws.close = function() {
	ws.stopreconnecting();
	ws.ws.onclose=undefined;
	ws.ws.close()
	ws.ws=null;
}

ws.registerContext = function(ctx, handler) {
	ws.contexts[ctx] = handler;
}


ws.send = function(ctx, data) {
	if (ws.ws) {
		var m = {ctx: ctx, data:data};
		ws.ws.send(JSON.stringify(m));
	}
}

ws.isopen = function() {
	return ws.ws && ws.ws.readyState == 1;
}