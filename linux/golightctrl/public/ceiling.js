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
  var width =  img.getAttribute('width');
  var height =  img.getAttribute('height');
  width = width === null || width === '' ? img.width : width;
  height = height === null || height === '' ? img.height : height;
  var ratio = width / height;
  var height = (window.innerHeight - room.offsetTop - 10);
  // todo: only set if bigger
  img.style.width = (height * ratio) + 'px';
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
  renderButtonStates();
}

function switchButton(ceiling, sendState) {
  console.log("switchButton");
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
  };
  req.setRequestHeader("googlechromefix","");
  req.send(null);
}

function switchButtonWebSocket(ceiling, sendState) {
  var message = {
    name: ceiling,
    action: sendState ? '1' : '0'
  };
  console.log("switchButtonWebSocket");
  ws.send("switch",message);
}

var webSocketUrl = 'ws://licht.realraum.at/sock';
var cgiUrl = '/cgi-bin/mswitch.cgi';
//var cgiUrl = 'fake.json';

var webSocketSupport = null;

var buttons = {
  ceiling1: false,
  ceiling2: false,
  ceiling3: false,
  ceiling4: false,
  ceiling5: false,
  ceiling6: false
};

function openWebSocket(webSocketUrl) {
  ws.registerContext("ceilinglights",setButtonStates);
  ws.open(webSocketUrl);
}

(function() {
  window.addEventListeners('resize orientationchange', resizeRoomImg);
  resizeRoomImg();

  switchButton();

  webSocketSupport = hasWebSocketSupport();
  if (webSocketSupport) {
    openWebSocket(webSocketUrl);
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

  if (!webSocketSupport) {
    setInterval(function() {
      switchButton();
    }, 1500);
  }
})();
