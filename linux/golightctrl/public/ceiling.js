"use strict";

function addEventListeners(types, listener, useCapture) {
  var typesArray = types.split(' ');
  for (var i = 0; i < typesArray.length; i++) {
    this.addEventListener(typesArray[i], listener, useCapture);
  }
}
Window.prototype.addEventListeners = addEventListeners;
EventTarget.prototype.addEventListeners = addEventListeners;

function resizeRoomImg() {
  var room = document.getElementsByClassName('room')[0];
  var img = room.getElementsByTagName('img')[0];
  var ratio = img.width / img.height;
  var height = (window.innerHeight - room.offsetTop - 10);
  // todo: only set if bigger
  img.style.width = (height * ratio) + 'px';
}

function hasWebSocketSupport() {
  // see: https://github.com/Modernizr/Modernizr/blob/115aa4a583dff0b44409b000165d22f2577ab0a1/feature-detects/websockets.js
  var supports = false;
  try {
    supports = 'WebSocket' in window && window.WebSocket.CLOSING === 2;
  } catch (e) {}
  return supports;
}

var ws = {};
ws.contexts = {};
ws.registerContext = function(ctx, handler) {
  ws.contexts[ctx] = handler;
}

function openWebSocket(webSocketUrl) {
  var webSocket = new WebSocket(webSocketUrl);
  webSocket.onopen = function (event) {
    webSocket.onmessage = function(response){
      //console.log(response);
      var m = $.parseJSON(response.data);
      if (m["ctx"] && m["data"] && typeof(ws.contexts[m.ctx]) == "function") {
        ws.contexts[m.ctx](m.data);
      }
    };
    webSocket.onclose = function(event) {
      console.log('webSocket closed');
      webSocket = null;
    }
    ws.registerContext("ceilinglights",function(data){
      setButtonStates(data);
    })
  };
  return webSocket;
}

function renderButtonStates() {
  var ceilings = document.getElementsByClassName('ceiling');
  for (var i = 0; i < ceilings.length; i++) {
    var id = ceilings[i].getAttribute('id');
    if (!buttons.hasOwnProperty(id)) {
      continue;
    }
    if (buttons[id]) {
      ceilings[i].style.background = 'green';
    } else {
      ceilings[i].style.background = 'red';
    }
  }
}

function setButtonStates(data) {
  for (var key in data) {
    if (data.hasOwnProperty(key) && buttons.hasOwnProperty(key)) {
      buttons[key] = data[key] ? true : false;
    }
  }
}

function switchButton(ceiling, sendState) {
  var url = cgiUrl;
  if (typeof ceiling !== 'undefined' && typeof sendState !== 'undefined') {
    url += '?' + ceiling + '=' + (sendState ? '1' : '0');
  }
  var req = new XMLHttpRequest;
  req.overrideMimeType("application/json");
  req.open("GET", url, true);
  req.onload = function() {
    if (req.status != 200) {
      return;
    }
    var data = JSON.parse(req.responseText);
    setButtonStates(data);
    renderButtonStates();
  };
  req.setRequestHeader("googlechromefix","");
  req.send(null);
}

function switchButtonWebSocket(ceiling, sendState) {
  var message = {
    ctx: 'switch',
    data: {
      name: ceiling,
      action: sendState ? '1' : '0'
    }
  };
  webSocket.send(JSON.stringify(message));

  //todo: remove this when websockets work
  buttons[ceiling] = sendState;
  renderButtonStates();
}

var webSocketUrl = 'ws://licht.realraum.at/sock';
var cgiUrl = '/cgi-bin/mswitch.cgi';
//var cgiUrl = 'fake.json';

var webSocketSupport = null;
var webSocket = null;

var buttons = {
  ceiling1: false,
  ceiling2: false,
  ceiling3: false,
  ceiling4: false,
  ceiling5: false,
  ceiling6: false
};

(function() {
  window.addEventListeners('resize orientationchange', resizeRoomImg);
  resizeRoomImg();

  switchButton();

  webSocketSupport = hasWebSocketSupport();
  if (webSocketSupport) {
    webSocket = openWebSocket(webSocketUrl);
  }

  var ceilings = document.getElementsByClassName('ceiling');
  for (var i = 0; i < ceilings.length; i++) {
    ceilings[i].addEventListener('click', function() {
      var id = this.getAttribute('id');
      if (!buttons.hasOwnProperty(id)) {
        return;
      }
      if (webSocketSupport) {
        switchButtonWebSocket(id, !buttons[id]);
      } else {
        switchButton(id, !buttons[id]);
      }
    });
  }

  if (!webSocketSupport)
  {
    setInterval(function() {
      switchButton();
    }, 1000);
  }
})();
