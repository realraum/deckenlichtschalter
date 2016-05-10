"use strict";

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

function switchButton(light, sendState) {
  var url = cgiUrl;
  if (typeof light !== 'undefined' && typeof sendState !== 'undefined') {
    url += '?' + light + '=' + (sendState ? '1' : '0');
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

function switchButtonWebSocket(light, sendState) {
  var message = {
    name: light,
    action: sendState ? '1' : '0'
  };
  ws.send("switch",message);
}

var webSocketUrl = 'ws://'+window.location.hostname+'/sock';
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
  $(window).on('resize orientationchange', resizeRoomImg);
  resizeRoomImg();

  webSocketSupport = hasWebSocketSupport();
  if (webSocketSupport) {
    openWebSocket(webSocketUrl);
  } else {
    switchButton();
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

  var rfirs = document.getElementsByClassName('rfir');
  for (var i = 0; i < rfirs.length; i++) {
    rfirs[i].addEventListener('click', function(event) {
      var id = this.getAttribute('id');
      var offset = $(this).offset();
      var relX = (event.pageX - offset.left) / $(this).width();
      var relY = (event.pageY - offset.top) / $(this).height();
      var sendState = relX + relY < 1;
      if (webSocketSupport) {
        switchButtonWebSocket(id, sendState);
      } else {
        switchButton(id, sendState);
      }
    });
  }

  if (!webSocketSupport) {
    setInterval(function() {
      switchButton();
    }, 30*1000);
  }
})();
