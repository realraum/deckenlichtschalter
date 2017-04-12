// function renderCeilingButtonUpdate(data) {
//   for (var keyid in data) {
//     on_btn = document.getElementById("onbtn_"+keyid);
//     off_btn = document.getElementById("offbtn_"+keyid);
//     if (on_btn && off_btn)
//     {
//       on_btn.className = "onbutton";
//       off_btn.className = "offbutton";
//       if (data[keyid])
//       { on_btn.className += " enableborder"; }
//       else
//       { off_btn.className += " enableborder"; }
//     }
//   }
// }

function renderRFIRButtonUpdate(snd_btn) {
  if (snd_btn)
  {
    var origclass = snd_btn.className;
    snd_btn.className += " activatedbtn";

    setTimeout(function(){
      $(".activatedbtn").removeClass("activatedbtn");
    },700);
  }
}

function populatedivrfswitchboxes(elem, names) {
  Object.keys(names).forEach(function(lightname) {
    $(elem).append('\
      <div class="switchbox">\
          <span class="alignbuttonsleft">\
          <button class="onbutton" lightname="'+lightname+'" action="on">On</button>\
          <button class="offbutton" lightname="'+lightname+'" action="off">Off</button>\
          </span>\
          <div class="switchnameright">'+names[lightname]+'</div>\
      </div>\
      <br>');
  });
}

function populatedivfancyswitchboxes(elem, names) {
  Object.keys(names).forEach(function(lightname) {
    var targetid = lightname.substr(7,1);
    $(elem).append('        <div class="switchbox">\
            <div style="width:100%; font-weight: bold; color:white; background-color: black;">'+names[lightname]+'</div>\
            <span class="alignbuttonsleft">\
            <button class="popupselect_trigger" optionsid="fancycolorquickoptions1" optionscopyattr="name" name="'+lightname+'" style="background-color:black;"></button>\
            <button class="fancylightcolourtempselectorbutton leftalignroundedbutton" name="'+lightname+'">PickColor</button>\
              RedShift:\
              <div class="onoffswitch">\
                  <input type="checkbox" class="onoffswitch-checkbox scriptctrl_redshift_checkbox" target="'+targetid+'" id="'+lightname+'ctonoff">\
                  <label class="onoffswitch-label" for="'+lightname+'ctonoff">\
                      <span class="onoffswitch-inner"></span>\
                      <span class="onoffswitch-switch"></span>\
                  </label>\
              </div>\
            </span>\
            <span class="alignbuttonsleft">\
              <input type="range" min="0" max="1000" step="1" value="0" class="fancyintensityslider" name="'+lightname+'"> Light Intensity\
            </span>\
            <span class="alignbuttonsleft">\
              <input type="range" min="0" max="1000" step="1" value="500" class="fancybalanceslider" name="'+lightname+'"> Color Temp.\
            </span>\
          </div>\
          <br/>');
  });
}

function sendYmhButton( btn ) {
  document.getElementById('indicator').style.backgroundColor="red";
  document.getElementById('commandlabel').innerHTML=btn;
  sendMQTT("action/GoLightCtrl/"+btn,{Action:"send"});
  document.getElementById('indicator').style.backgroundColor="white";
  document.getElementById('commandlabel').innerHTML='&nbsp;';
}

function remoteKeyboard( e ) {
  e = e || window.event;
  switch( e.keyCode )
  {
    //case 81: sendYmhButton( ( e.altKey ? 'REBOOT' : 'w' ) ); break;       // ALT-Q = reboot, Q = Power
    case 38: sendYmhButton( 'ymhprgup' ); break; //up
    case 37: sendYmhButton( 'ymhminus' ); break; //left
    //case 13: sendYmhButton( '' ); break;  //enter 
    case 39: sendYmhButton( 'ymhplus' ); break; //right
    case 40: sendYmhButton( 'ymhprgdown' ); break; //down
    //case  8: sendYmhButton( 'T_back' ); break; //backspace
    //case 27: sendYmhButton( 't_stop' ); break;    // ESC
    case 77: sendYmhButton( 'ymhmenu' ); break;     // M
    case 84: sendYmhButton( 'ymhtimelevel' ); break;        // T
    //case 32: sendYmhButton( '' ); break;  // Space
    case 83: sendYmhButton( 'ymhsleep' ); break;            // S
    case 36: sendYmhButton( 'ymhmute' ); break; //pos1
    case 33: sendYmhButton( 'ymhvolup' ); break;    // P = 80, PAGEUP = 33
    case 43: sendYmhButton( 'ymhvolup' ); break;    // +
    case 34: sendYmhButton( 'ymhvoldown' ); break;  // N = 78, PAGEDOWN = 34
    case 45: sendYmhButton( 'ymhvoldown' ); break;  // -
    case 46: sendYmhButton( 'ymhtunplus' ); break;  // .
    case 44: sendYmhButton( 'ymhtunminus' ); break; // ,
    case 59: sendYmhButton( 'ymhtunabcde' ); break; // ;
    case 69: sendYmhButton( 'ymheffect' ); break;   // E
    case 80: sendYmhButton( 'ymhpower' ); break;    // P
    case 48: sendYmhButton( 'ymhtest' ); break;     // 0
    case 49: sendYmhButton( 'ymhcd' ); break;       // 1
    case 50: sendYmhButton( 'ymhtuner' ); break;    // 2
    case 51: sendYmhButton( 'ymhtape' ); break;     // 3
    case 52: sendYmhButton( 'ymhwdtv' ); break;     // 4
    case 53: sendYmhButton( 'ymhsattv' ); break;    // 5
    case 54: sendYmhButton( 'ymhvcr' ); break;      // 6
    case 55: sendYmhButton( 'ymh7' ); break;        // 7
    case 56: sendYmhButton( 'ymhaux' ); break;      // 8
    case 57: sendYmhButton( 'ymhextdec' ); break;   // 9
  }
}

document.onkeydown = remoteKeyboard;

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

function calcColorFromDayLevel(day_factor, value)
{
    var day_factor = Math.min(1.0,Math.max(-1.0,day_factor))
    var r = 1000 * value * Math.max(0.0, -1.0 * day_factor)
    var b = 0
    var cw = 1000 * value * Math.max(0.0, day_factor)
    var ww = Math.max(0,1000 * value - cw - (r/3))
    return {"r":Math.trunc(r), "b":Math.trunc(b), "cw":Math.trunc(cw), "ww":Math.trunc(ww)}
}

function updateColdWarmWhiteBalanceIntensity(event)
{
  var fancyid = event.target.getAttribute("name");
  var intensity = parseInt($("input.fancyintensityslider[name="+fancyid+"]")[0].value,10) / 1000.0;
  var balance = (1000 - parseInt($("input.fancybalanceslider[name="+fancyid+"]")[0].value,10)*2) / 1000.0;
  sendMQTT("action/"+fancyid+"/light",calcColorFromDayLevel(balance, intensity));
}

function enableRedShift() {
  var participating = Array();
  $(".scriptctrl_redshift_checkbox").each(function(elem){
    if ($(elem).checked) {
      participating.append(parseInt(elem.getAttribute("target")))
    }
  });
  if (participating.length > 0) {
    sendMQTT("action/ceilingscripts/activatescript",{"script":"redshift","participating":participating})
  } else {
    //TODO FIXME: propably not what we want.
    //this will switch off all ceiling lights, even if some were not script controlled
    sendMQTT("action/ceilingscripts/activatescript",{script:"off"})
  }
}

var fancycolorpicker_apply_name="ceiling1";
function popupFancyColorPicker(event) {
  var x = event.pageX;
  var y = event.pageY;
  $("#fancycolorpicker").css("left",x).css("top",y).css("visibility","visible");
  fancycolorpicker_apply_name = event.target.getAttribute("name");
}


function setLedPipePattern(data) {
  sendMQTT("action/PipeLEDs/pattern",data);
}

var webSocketUrl = 'ws://'+window.location.hostname+'/sock';
var cgiUrl = '/cgi-bin/mswitch.cgi';
//var cgiUrl = 'fake.json';

var webSocketSupport = null;

populatedivrfswitchboxes(document.getElementById("divrfswitchboxes"), {
  "regalleinwand":"LEDs Regal Leinwand",
  "labortisch":"TESLA Labortisch",
  "floodtesla":"TESLA Deckenfluter",
  "bluebar":"Blaue LEDs Bar",
  "abwasch":"Licht Waschbecken",
  "couchred":"LEDs Couch Red",
  "couchwhite":"LEDS Couch White",
  "cxleds":"CX Gang LEDs",
  "mashadecke":"MaSha Werkstatt Decke",
  "allrf":"Alle Funksteckdosen",
  "ambientlights":"Ambient Lichter",
  "boiler":"Warmwasser Küche",
  "boilerolga":"Warmwasser OLGA"
});

populatedivrfswitchboxes(document.getElementById("divbasiclightswitchboxes"), {
"ceiling1":"Decke Leinwand",
"ceiling2":"Decke Durchgang",
"ceiling3":"Decke Küche",
"ceiling4":"Decke Lasercutter",
"ceiling5":"Decke Eingang",
"ceiling6":"Decke Tesla",
"ceilingAll":"All BasicLights",
});

populatedivfancyswitchboxes(document.getElementById("divfancylightswitchboxes"), {
"ceiling1":"Decke Leinwand",
"ceiling2":"Decke Durchgang",
"ceiling3":"Decke Küche",
"ceiling4":"Decke Lasercutter",
"ceiling5":"Decke Eingang",
"ceiling6":"Decke Tesla",
"ceilingAll":"All FancyLights",
});

(function() {
  webSocketSupport = hasWebSocketSupport();
  if (webSocketSupport) {
    ws.open(webSocketUrl);
  } else {
    updateButtons("/cgi-bin/mswitch.cgi");
    setInterval("updateButtons(\"/cgi-bin/mswitch.cgi\");", 30*1000);
  }

  var topics_to_subscribe = {};

  var onbtns = [].slice.call(document.getElementsByClassName('onbutton'));
  var offbtns = [].slice.call(document.getElementsByClassName('offbutton'));
  var onoffbtn = Array.prototype.concat(onbtns,offbtns);
  for (var i = 0; i < onoffbtn.length; i++) {
    var lightname = onoffbtn[i].getAttribute("lightname");
    var action = onoffbtn[i].getAttribute("action");
    if (lightname) {
      var topic = "action/GoLightCtrl/"+lightname;
      onoffbtn[i].addEventListener('click', function(topic, action) {
        return function() {  sendMQTT(topic,{Action:action});  };
      }(topic, action));
      topics_to_subscribe[topic] = lightname;
    }
  }
  if (webSocketSupport) {
    Object.keys(topics_to_subscribe).forEach(function(topic) {
      var lightname = topics_to_subscribe[topic];
      ws.registerContext(topic, function(topic) {
        var btns = $("button[lightname="+lightname+"]");
        return function(data) {
          btns.each(function(elem) {
            if (elem.getAttribute("action") == data.Action) 
            {
              renderRFIRButtonUpdate(elem);
            }
          });
        };
      }(topic));
    });

    registerFunctionForFancyLightUpdate(handleExternalFancySetting);
  }

  popupselect.init();
  $(".fancylightpresetbutton").on("click",eventOnFancyLightPresent);
  popupselect.addSelectHandlerToAll(eventOnFancyLightPresent);
  $('.scriptctrl_redshift_checkbox').on("click",enableRedShift);
  $("input.fancyintensityslider").on("change",updateColdWarmWhiteBalanceIntensity)
  $("input.fancybalanceslider").on("change",updateColdWarmWhiteBalanceIntensity)
  $(".fancylightcolourtempselectorbutton").on("click",popupFancyColorPicker);
  $("#fancycolorpicker_close_button").on("click",function(event){$("#fancycolorpicker").css("visibility","hidden")});
  $("#fancycolorpicker_apply_button").on("click",function(event){
      var R = parseInt($('#R input').val()) || 0;
      var G = parseInt($('#G input').val()) || 0;
      var B = parseInt($('#B input').val()) || 0;
      var CW = parseInt($('#CW input').val()) || 0;
      var WW = parseInt($('#WW input').val()) || 0;
      var settings = {r:R,g:G,b:B,cw:CW,ww:WW,fade:{}};
      sendMQTT("action/"+fancycolorpicker_apply_name+"/light",settings);
  });
  //draw color picker canvases
  init_colour_temp_picker();
  //define setLedPipePattern(objdata)

  $(".ledpipepresetbutton").on('click', function() {
    var pipepattern = this.getAttribute("pipepattern");
    if (!pipepattern) { return;  }
    var hue = parseInt(this.getAttribute("pipehue")) || undefined;
    var brightness = parseInt(this.getAttribute("pipebrightness")) || undefined;
    var speed = parseInt(this.getAttribute("pipespeed")) || undefined;
    var arg = parseInt(this.getAttribute("pipearg")) || undefined;
    var arg1 = parseInt(this.getAttribute("pipearg1")) || undefined;
    var data = {pattern:pipepattern, hue:hue, brightness:brightness, speed:speed, arg:arg, arg1:arg1};
    setLedPipePattern(data);
    if (brightness) {
      document.getElementById("pipebrightness").value=brightness;
    }
    if (hue) {
      document.getElementById("pipehue").value=hue;
    }
    if (speed) {
      document.getElementById("pipespeed").value=speed;
    }      
  });

})();
