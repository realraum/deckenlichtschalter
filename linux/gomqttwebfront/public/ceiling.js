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
  var coldwhite_representation  = [1, 0xfa/0xff, 0xc0/0xff];
  var warmwhite_representation = [71/0xff, 171/0xff, 1];

  //save data for next color chooser popup
  fancycolorstate_[fancyid] = data;
  if (data.cw + data.ww == 0)
  {
    fancycolorstate_[fancyid].compound_r = Math.floor(data.r / 4);
    fancycolorstate_[fancyid].compound_g = Math.floor(data.g / 4);
    fancycolorstate_[fancyid].compound_b = Math.floor(data.b / 4);
  } else if (data.r+data.g+data.b == 0)
  {
    fancycolorstate_[fancyid].compound_r = Math.floor((data.cw*coldwhite_representation[0] + data.ww*warmwhite_representation[0]) / 8);
    fancycolorstate_[fancyid].compound_g = Math.floor((data.cw*coldwhite_representation[1] + data.ww*warmwhite_representation[1]) / 8);
    fancycolorstate_[fancyid].compound_b = Math.floor((data.cw*coldwhite_representation[1] + data.ww*warmwhite_representation[1]) / 8);
  } else {
    fancycolorstate_[fancyid].compound_r = Math.floor((data.r/4 + data.cw*coldwhite_representation[0] + data.ww*warmwhite_representation[0]) / 9);
    fancycolorstate_[fancyid].compound_g = Math.floor((data.g/4 + data.cw*coldwhite_representation[1] + data.ww*warmwhite_representation[1]) / 9);
    fancycolorstate_[fancyid].compound_b = Math.floor((data.b/4 + data.cw*coldwhite_representation[1] + data.ww*warmwhite_representation[1]) / 9);
  }
  console.log(fancycolorstate_[fancyid]);
  var rgbstring = "rgb("+fancycolorstate_[fancyid].compound_r+","+fancycolorstate_[fancyid].compound_g+","+fancycolorstate_[fancyid].compound_b+")";
  var elem = $("button.popupselect_trigger[name="+fancyid+"]");
  if (elem) {
    console.log(rgbstring);
    elem.css("background-color",rgbstring);
  }
  if (fancyid=="ceilingAll")
  {
    for (var fid=1; fid<10; fid++)
    {
      fancycolorstate_[fid] = fancycolorstate_["All"];
      elem = $("button.popupselect_trigger[name="+fid+"]");
      if (elem) {
        elem.css("background-color",rgbstring);
      }
    }
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
var topic_namectrl = "action/GoLightCtrl/"

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
  popupselect.init({class_option:"bigpopupselect_option"});
  popupselect.addSelectHandlerToAll(eventOnFancyLightPresent);
  $(".basiclight").on("click",function() {
      var id = this.getAttribute('id');
      var topic = topic_namectrl + id;
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
      var topic = topic_namectrl + keyid;
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
      var topic = topic_namectrl + id;
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
    ws.registerContext(topic_fancy_ceiling1, function(data){console.log(topic_fancy_ceiling1,data);});
    ws.registerContext(topic_fancy_ceiling2, function(data){console.log(topic_fancy_ceiling2,data);});
    ws.registerContext(topic_fancy_ceiling3, function(data){console.log(topic_fancy_ceiling3,data);});
    ws.registerContext(topic_fancy_ceiling4, function(data){console.log(topic_fancy_ceiling4,data);});
    ws.registerContext(topic_fancy_ceiling5, function(data){console.log(topic_fancy_ceiling5,data);});
    ws.registerContext(topic_fancy_ceiling6, function(data){console.log(topic_fancy_ceiling6,data);});
    ws.registerContext(topic_fancy_ceiling7, function(data){console.log(topic_fancy_ceiling7,data);});
    ws.registerContext(topic_fancy_ceiling8, function(data){console.log(topic_fancy_ceiling8,data);});
    ws.registerContext(topic_fancy_ceiling9, function(data){console.log(topic_fancy_ceiling9,data);});  
    ws.open(webSocketUrl);
  } else {
    sendMQTT_XHTTP("","");
    setInterval(function() {
      sendMQTT_XHTTP("","");
    }, 30*1000);
  }

})();
