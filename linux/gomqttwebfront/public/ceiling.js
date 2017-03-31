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

function sendMQTT_XHTTP(ctx, data) {
  var req = new XMLHttpRequest;
  req.open("POST", cgiUrl, true);
  req.onload = function() {
    if (req.status != 200) {
      return;
    }
    var data = JSON.parse(req.responseText);
    setButtonStates(data);
  };
  var param = "Ctx=" + encodeURIComponent(ctx);
  params = params + "&Data="+encodeURIComponent(data);
  params = params.replace(/%20/g, '+');
  req.overrideMimeType("application/json");
  req.setRequestHeader("googlechromefix","");
  req.setRequestHeader("Content-length", params.length);
  req.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
  req.setRequestHeader("Connection", "close");
  req.send(params);
}

function sendMQTT(ctx, data) {
  if (webSocketSupport) {
    ws.send(ctx,data);
  } else {
    sendMQTT_XHTTP(ctx, data);
  }
}

var webSocketUrl = 'ws://'+window.location.hostname+'/sock';
var cgiUrl = '/cgi-bin/fallback.cgi';

var webSocketSupport = null;

var topic_fancy_ceiling1 = "action/ceiling1/light"
var topic_fancy_ceiling2 = "action/ceiling2/light"
var topic_fancy_ceiling3 = "action/ceiling3/light"
var topic_fancy_ceiling4 = "action/ceiling4/light"
var topic_fancy_ceiling5 = "action/ceiling5/light"
var topic_fancy_ceiling6 = "action/ceiling6/light"
var topic_fancy_ceiling7 = "action/ceiling7/light"
var topic_fancy_ceiling8 = "action/ceiling8/light"
var topic_fancy_ceiling9 = "action/ceiling9/light"
var topic_namectrl = "action/GoNameCtrl/name"

var buttons = {
  ceiling1: false,
  ceiling2: false,
  ceiling3: false,
  ceiling4: false,
  ceiling5: false,
  ceiling6: false
};

function renderCeilingButtonUpdate(data) {
  for (var keyid in data) {
    btn = document.getElementById(keyid);
    if ($(btn).hasClass("rfir"))
    {
      var origclass = snd_btn.className;
      snd_btn.className += " activatedbtn";

      setTimeout(function(){
        snd_btn.className = origclass;
      },900);
    }
  }
}

function openWebSocket(webSocketUrl) {
  ws.registerContext("ceilinglights",setButtonStates);
  ws.registerContext("FancyLight",setFancyLightStates);
  ws.registerContext("wbp",renderCeilingButtonUpdate);
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
