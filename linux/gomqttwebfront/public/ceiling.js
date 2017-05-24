"use strict";

function resizeRoomImg() {
  var room = document.getElementById('room');
  var img = document.getElementById('roommap');
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
  //calc compound RGB from light data
  calcCompoundRGB(fancycolorstate_[fancyid], fancyid);

  //set compound RGB to background-color of button
  var rgbstring = "rgb("+fancycolorstate_[fancyid].compound_r+","+fancycolorstate_[fancyid].compound_g+","+fancycolorstate_[fancyid].compound_b+")";
  var elem = $(".popupselect_trigger[name="+fancyid+"]");
  if (elem) {
    elem.css("background-color",rgbstring);
  }
  //calculcate and set balance/intensity slider
  var cwwwslidedata = calcDayLevelFromColor(data);
  $("input.fancyintensityslider[name="+fancyid+"]").val(Math.floor(cwwwslidedata["intensity"]*1000));
  $("input.fancybalanceslider[name="+fancyid+"]").val(Math.floor((1000-cwwwslidedata["balance"]*1000)/2));

  //on ceiling.js/index.html we only have redshift script on/off buttons
  //redshift uses triggers for each individual light and thus always has an
  //acompaning "s" field in json sent to lamp
  //so it's easy to detect if redshift is on for that light or not.
  var redshift_checkboxes =  $("input.scriptctrl_redshift_checkbox[name='"+fancyid+"']");
  if (redshift_checkboxes.length > 0)
    redshift_checkboxes[0].checked = (data.s && data.s == "redshift");
}

function enableRedShift(event) {
  var participating = Array();
  {
    $(".scriptctrl_redshift_checkbox").each(function(elem){
      if (elem.checked) {
        var lightname = elem.getAttribute("name");
        if (lightname != "ceilingAll"){
          participating.push(lightname);
        }
      }
    });
  }
  if (participating.length > 0) {
    sendMQTT(mqtttopic_activatescript,{"script":"redshift","participating":participating});
  } else {
    //TODO FIXME: propably not what we want.
    //this will switch off all ceiling lights, even if some were not script controlled
    sendMQTT(mqtttopic_activatescript,{script:"off"})
  }
}

var script_running_="off";
function handleExternalActivateScript(data) {
  script_running_ = data.script;
  $("#scriptctrlselect").val(data.script);
  if (data.script == "redshift") {
    // -------- Script redshift ---------
    if (data.participating == undefined || data.participating.length==6) {
      data.participating=["ceiling1","ceiling2","ceiling3","ceiling4","ceiling5","ceiling6","flooddoor"];
    }
    $(".scriptctrl_redshift_checkbox").each(function(elem) {
      var lightname = elem.getAttribute("name");
      elem.checked = (-1 != data.participating.indexOf(lightname));
    });
  } else {
    $(".scriptctrl_redshift_checkbox").each(function(elem) {
      elem.checked = false;
    });
  }
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

$(window).ready(resizeRoomImg);
(function() {
  $(window).on('resize orientationchange', resizeRoomImg);

  webSocketSupport = hasWebSocketSupport();

  //set background color for fancylightpresetbuttons according to ledr=, ledb=, etc.
  $(".fancylightpresetbutton").each(colorFancyLightPresent);

  $(".mqttrawjson").on("click",eventOnRawMqttElement);
  $('.scriptctrl_redshift_checkbox').on("click",enableRedShift);
  popupselect.init({class_option:"popupselect_option"});
  popupselect.addSelectHandlerToAll(eventOnFancyLightPresent);
  $(".fancylightpresetbutton.popupselect_option").each(function(e){popupselect.addSelectHandler(e, eventOnFancyLightPresent)});
  $(".fancylightpresetbutton").not(".popupselect_option").on("click",eventOnFancyLightPresent);
  $(".presetfunctionbutton.popupselect_option").each(function(e){popupselect.addSelectHandler(e, function(event){
    var func = event.target.getAttribute("presetfunc");
    window[func]();
  })});
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
          if (data.Action == "1" || data.Action == "on" || data.Action == "send" || data.Action == 1) {
            buttons[keyid] = true;
          } else {
            buttons[keyid] = false;
          }
          renderButtonStates();
        };
      }(topic, keyid)));
    });

    ws.registerContext(mqtttopic_activatescript,handleExternalActivateScript);
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

  $(".sonoff").on("click",function(event) {
      var id = this.getAttribute('id');
      var topic = mqtttopic_sonoff(id);
      var offset = $(this).offset();
      var relX = (event.pageX - offset.left) / $(this).width();
      var relY = (event.pageY - offset.top) / $(this).height();
      var sendState = relX + relY < 1;
      if (sendState) {
        sendMQTT(topic, "on");
      } else {
        sendMQTT(topic, "off");
      }
  });

  if (webSocketSupport) {
    ws.open(webSocketUrl);
  } else {
    sendMQTT_XHTTP("","");
    setInterval(function() {
      sendMQTT_XHTTP("","");
    }, 30*1000);
  }

})();
