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
  var ceilings = document.getElementsByClassName('basiclight');
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

var fancycolorstate_={};
function handleExternalFancySetting(fancyid, data)
{
  //save data for next color chooser popup
  fancycolorstate_[fancyid] = data;
  calcCompoundRGB(fancycolorstate_[fancyid]);

  var rgbstring = "rgb("+fancycolorstate_[fancyid].compound_r+","+fancycolorstate_[fancyid].compound_g+","+fancycolorstate_[fancyid].compound_b+")";
  var elem = $(".popupselect_trigger[name="+fancyid+"]");
  if (elem) {
    elem.css("background-color",rgbstring);
  }
  var cwwwslidedata = calcDayLevelFromColor(data);
  $("input.fancyintensityslider[name="+fancyid+"]").val(Math.floor(cwwwslidedata["intensity"]*1000));
  $("input.fancybalanceslider[name="+fancyid+"]").val(Math.floor((1000-cwwwslidedata["balance"]*1000)/2));
}

var webSocketUrl = 'ws://'+window.location.hostname+'/sock';
var cgiUrl = '/cgi-bin/fallback.cgi';

var webSocketSupport = null;

var buttons = {
  basiclight1: false,
  basiclight2: false,
  basiclight3: false,
  basiclight4: false,
  basiclight5: false,
  basiclight6: false
};

(function() {
  $(window).on('resize orientationchange', resizeRoomImg);
  resizeRoomImg();

  webSocketSupport = hasWebSocketSupport();


  $(".mqttrawjson").on("click",eventOnRawMqttElement);
  $(".fancylightpresetbutton").on("click",eventOnFancyLightPresent);
  popupselect.init({class_option:"popupselect_option"});
  popupselect.addSelectHandlerToAll(eventOnFancyLightPresent);
  $(".basiclight").on("click",function() {
      var id = this.getAttribute('id');
      var topic = mqtttopic_golightctrl(id);
      if (!buttons.hasOwnProperty(id)) {
        return;
      }
      if (buttons[id]) {
        sendMQTT(topic, {Action:"off"});
      } else {
        sendMQTT(topic, {Action:"on"});
      }
    });
  if (webSocketSupport) {
    $(".basiclight").each(function(elem) {
      var keyid = $(elem).attr('id');
      var topic = mqtttopic_golightctrl(keyid);
      ws.registerContext(topic, (function(topic,keyid) {
        return function(data) {
          console.log(topic, data);
          if (data.Action == "1" || data.Action == "on" || data.Action == "send") {
            buttons[keyid] = true;
          } else {
            buttons[keyid] = false;
          }
          renderButtonStates();
        };
      }(topic, keyid)));
    });

    registerFunctionForFancyLightUpdate(handleExternalFancySetting);
  }

  var rfirs = document.getElementsByClassName('rfir');
  for (var i = 0; i < rfirs.length; i++) {
    rfirs[i].addEventListener('click', function(event) {
      var id = this.getAttribute('id');
      var topic = mqtttopic_golightctrl(id);
      var offset = $(this).offset();
      var relX = (event.pageX - offset.left) / $(this).width();
      var relY = (event.pageY - offset.top) / $(this).height();
      var sendState = relX + relY < 1;
      if (sendState) {
        sendMQTT(topic, {Action:"on"});
      } else {
        sendMQTT(topic, {Action:"off"});
      }
    });
  }

  if (webSocketSupport) {
    ws.open(webSocketUrl);
  } else {
    sendMQTT_XHTTP("","");
    setInterval(function() {
      sendMQTT_XHTTP("","");
    }, 30*1000);
  }

})();
