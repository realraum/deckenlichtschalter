// (c) Bernhard Tittelbach

var ws = {};
ws.onopen=false;
ws.contexts = {};

ws.open = function(uri) {
	ws.ws=new WebSocket(uri);
 	ws.ws.onopen = function(){
		ws.ws.onmessage = function(response){
			//console.log(response);
			var m = $.parseJSON(response.data);
			if (m["ctx"] && m["data"] && typeof(ws.contexts[m.ctx]) == "function") {
				ws.contexts[m.ctx](m.data);
			}
		}
		ws.ws.onclose = function(){
			alert("Connection to server lost");
			window.location.reload()
		}
		if (typeof(ws["onopen"]) == "function"){
			ws.onopen();
		}
 	}
}

ws.close = function() {
	ws.ws.onclose=undefined;
	ws.ws.close()
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