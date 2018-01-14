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

function populatedivsonfoff(elem, names) {
  Object.keys(names).forEach(function(lightname) {
    $(elem).append('\
      <div class="switchbox">\
          <span class="alignbuttonsleft">\
            <div class="onoffswitch">\
                <input type="checkbox" class="onoffswitch-checkbox sonoff_checkbox" name="'+lightname+'" id="'+lightname+'ctonoff">\
                  <label class="onoffswitch-label" for="'+lightname+'ctonoff">\
                      <span class="onoffswitch-inner"></span>\
                      <span class="onoffswitch-switch"></span>\
                  </label>\
              </div>\
          </span>\
          <div class="switchnameright">'+names[lightname]+'</div>\
      </div>\
      <br>');
  });
}

function populatedivfancyswitchboxes(elem, names) {
  Object.keys(names).forEach(function(lightname) {
    var extrahtml="";
    if (names[lightname].hasbasic) {
      var basiclightname = lightname.replace("ceiling","basiclight");
      extrahtml += '\
      <span class="alignbuttonsright" style="margin:1ex; margin-right:2ex;">\
        BasicLight On: <div class="onoffswitch">\
            <input type="checkbox" class="onoffswitch-checkbox basiclight_checkbox" name="'+basiclightname+'" id="'+basiclightname+'basiconoff">\
            <label class="onoffswitch-label" for="'+basiclightname+'basiconoff">\
               <span class="onoffswitch-inner"></span>\
               <span class="onoffswitch-switch"></span>\
            </label>\
        </div>\
      </span>';
    }
    if (names[lightname].allbasic) {
      extrahtml += '<span class="alignbuttonsright" style="margin:1ex; margin-right:2ex;">\
        All BasicLights:\
          <button class="onbutton" lightname="basiclightAll" action="on">On</button>\
          <button class="offbutton" lightname="basiclightAll" action="off">Off</button>\
      </span>';
    }
    if (names[lightname].uv) {
      extrahtml += '<span class="alignbuttonsleft">\
      <input type="range" min="0" max="1000" step="1" value="0" class="fancyuvslider" name="'+lightname+'"> UV Intensity\
      </span>';
    }

    $(elem).append('<div class="switchbox">\
            <div style="width:100%; font-weight: bold; color:white; background-color: black;">'+names[lightname].desc+'</div>\
            <span class="alignbuttonsleft">\
            <button class="fancylightcolourtempselectorbutton leftalignroundedbutton" name="'+lightname+'">PickColor</button>\
            <button class="popupselect_trigger" optionsid="fancycolorquickoptions1" optionscopyattr="name" name="'+lightname+'" style="background-color:black;"></button>\
              <div style="display:inline-block; text-align:right; padding-top:0em; padding-right:0.5em;">Script<br/>Controlled:</div>\
              <div class="onoffswitch">\
                  <input type="checkbox" class="onoffswitch-checkbox scriptctrl_checkbox" name="'+lightname+'" id="'+lightname+'ctonoff">\
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
          '+extrahtml+'\
          </div>\
          <br/>');
  });
}

function sendYmhButton( btn ) {
  document.getElementById('indicator').style.backgroundColor="red";
  document.getElementById('commandlabel').innerHTML=btn;
  sendMQTT(mqtttopic_golightctrl(btn),{Action:"send"});
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
//   case 48: sendYmhButton( 'ymhtest' ); break;     // 0
//   case 49: sendYmhButton( 'ymhcd' ); break;       // 1
//   case 50: sendYmhButton( 'ymhtuner' ); break;    // 2
//   case 51: sendYmhButton( 'ymhtape' ); break;     // 3
//   case 52: sendYmhButton( 'ymhwdtv' ); break;     // 4
//   case 53: sendYmhButton( 'ymhsattv' ); break;    // 5
//   case 54: sendYmhButton( 'ymhvcr' ); break;      // 6
//   case 55: sendYmhButton( 'ymh7' ); break;        // 7
//   case 56: sendYmhButton( 'ymhaux' ); break;      // 8
//   case 57: sendYmhButton( 'ymhextdec' ); break;   // 9
  }
}

document.onkeydown = remoteKeyboard;

var fancycolorstate_={};
function handleExternalFancySetting(fancyid, data)
{
  //save data for next color chooser popup
  //note that not all fields my be set so we use keep previous data
  if (fancycolorstate_[fancyid])
    $.extend(fancycolorstate_[fancyid], data);
  else
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
  var cwwwslidedata = calcDayLevelFromColor(fancycolorstate_[fancyid]);
  $("input.fancyintensityslider[name="+fancyid+"]").val(Math.floor(cwwwslidedata["intensity"]*1000));
  $("input.fancybalanceslider[name="+fancyid+"]").val(Math.floor((1000-cwwwslidedata["balance"]*1000)/2));

  //activate Script Buttons if detected that light is controlled by script
  //i.e. the json includes a trigger sequence number "sq"
  // only works for lights being used as triggers
  $("input.scriptctrl_checkbox[name='"+fancyid+"']")[0].checked = (data.s && data.s != "off");

  //set UV slider if data present
  if (data.uv && data.uv >= 0 && data.uv <= 1000)
    $("input.fancyuvslider[name="+fancyid+"]").val(data.uv);
}

function updateColdWarmWhiteBalanceIntensity(event)
{
  var fancyid = event.target.getAttribute("name");
  var intensity = parseInt($("input.fancyintensityslider[name="+fancyid+"]")[0].value,10) / 1000.0;
  var balance = (1000 - parseInt($("input.fancybalanceslider[name="+fancyid+"]")[0].value,10)*2) / 1000.0;
  sendMQTT(mqtttopic_fancylight(fancyid),calcColorFromDayLevel(balance, intensity));
}

function updateUVIntensity(event)
{
  var fancyid = event.target.getAttribute("name");
  var uv = parseInt($("input.fancyuvslider[name="+fancyid+"]")[0].value,10);
  sendMQTT(mqtttopic_fancylight(fancyid),{uv:uv});
}



function handleChangeScriptCtrl(event) {
  var participating = Array();
  if (event && event.target.getAttribute("name") == "ceilingAll") {
    if (event.target.checked)
      participating=mqtt_fancylights_all.slice();
    else
      participating = Array();
  } else {
    $(".scriptctrl_checkbox").each(function(elem){
      if (elem.checked) {
        var lightname = elem.getAttribute("name");
        if (lightname != "ceilingAll"){
          participating.push(lightname);
        }
      }
    });
  }
  //don't switch off if scriptselect is "off" and we want to pre-toggle checkboxes
  if (event && $(event.target).hasClass('scriptctrl_checkbox') && $("#scriptctrlselect").val() == "off")
    return;

  // if checkboxes selected and script selected, enable script
  if (participating.length > 0) {
    sendMQTT(mqtttopic_activatescript,{"script":$("#scriptctrlselect").val(),"participating":participating,"value":parseInt($("#scriptctrlfancyintensityslider").val(),10)/1000.0});
    return;
  }

  //don't switch off if no checkbox is checked yet but script got selected
  if (event && event.target == document.getElementById("scriptctrlselect") && $("#scriptctrlselect").val() != "off")
    return;

  //otherwise switch off
  sendMQTT(mqtttopic_activatescript,{script:"off"});
}

var script_running_="off";
function handleExternalActivateScript(data) {
  script_running_ = data.script;  $("#scriptctrlselect").val(data.script);
  if (data.value) {
    $("#scriptctrlfancyintensityslider").val(Math.floor(data.value*1000));
  }
  if (data.participating == undefined || data.participating.length>=mqtt_fancylights_all.length) {
    data.participating=mqtt_fancylights_all.slice(); //copy array
  }

  $(".scriptctrl_checkbox").each(function(elem) {
    var target = elem.getAttribute("name");
    elem.checked = (-1 != data.participating.indexOf(target));
  });

  if (data.script == "redshift") {
    // -------- Script redshift ---------
    //add stuff here for this script
  } else if (data.script == "randomcolor") {
    // -------- Script randomcolor ---------
    //add stuff here for this script
  } else if (data.script == "colorfade") {
    // -------- Script colorfade ---------
    //add stuff here for this script
  } else if (data.script == "ceilingsinus") {
    // -------- Script ceilingsinus ---------
    //add stuff here for this script
  } else if (data.script == "wave") {
    // -------- Script wave ---------
    //add stuff here for this script
  } else {
    // -------- Script OFF ---------
    $(".scriptctrl_checkbox").each(function(elem) {
      elem.checked = false;
    });
  }
}

var fancycolorpicker_apply_name="ceiling1";
function popupFancyColorPicker(event) {
  var x = event.target.offsetLeft;
  var y = event.target.offsetTop;
  $("#fancycolorpicker").css("left",x+"px").css("top",y+"px").css("visibility","visible");
  fancycolorpicker_apply_name = event.target.getAttribute("name");
  colorandtempicker_set_ledfactorname(fancycolorpicker_apply_name);
}


function setLedPipePattern(data) {
  sendMQTT(mqtttopic_pipeledpattern, data);
}

var webSocketUrl = 'ws://'+window.location.hostname+'/sock';
var cgiUrl = '/cgi-bin/mswitch.cgi';
//var cgiUrl = 'fake.json';

var webSocketSupport = null;

populatedivrfswitchboxes(document.getElementById("divrfswitchboxes"), {
  "regalleinwand":"LEDs Regal Leinwand",
  "floodtesla":"TESLA Deckenfluter",
  "bluebar":"Blaue LEDs Bar",
  "couchwhite":"Couch White LEDs",
  "cxleds":"CX Gang LEDs",
  "mashadecke":"MaSha Werkstatt Decke",
  "floodoutside":"Außenfluter",
  "logo":"r3 Logo Auslage",
  "allrf":"Alle Funksteckdosen",
  "ambientlights":"Ambient Lichter",
  "boilerolga":"Warmwasser OLGA"
});

populatedivsonfoff(document.getElementById("divrfswitchboxes"), {
  "couchred":"Couch Ducks & Red-LEDs",
  "tesla":"TESLA Labortisch",
  "lothrboiler":"Warmwasser LoTHR",
  "druckerstrom":"Laserjet8000",
});

populatedivfancyswitchboxes(document.getElementById("divfancylightswitchboxes"), {
"ceilingAll":{desc:"Alle Deckenlichter", hasbasic:false, allbasic:true, uv:false},
"ceiling1":{desc:"Decke Leinwand", hasbasic:true, allbasic:false, uv:false},
"ceiling2":{desc:"Decke Durchgang", hasbasic:true, allbasic:false, uv:false},
"ceiling3":{desc:"Decke Küche", hasbasic:true, allbasic:false, uv:false},
"ceiling4":{desc:"Decke Lasercutter", hasbasic:true, allbasic:false, uv:false},
"ceiling5":{desc:"Decke Eingang", hasbasic:true, allbasic:false, uv:false},
"ceiling6":{desc:"Decke Tesla", hasbasic:true, allbasic:false, uv:false},
"flooddoor":{desc:"Deckenfluter Tür", hasbasic:false, allbasic:false, uv:false},
"abwasch":{desc:"Abwasch/Bar", hasbasic:false, allbasic:false, uv:true},
});

(function() {
  webSocketSupport = hasWebSocketSupport();

  var topics_to_subscribe = {};

  var onbtns = [].slice.call(document.getElementsByClassName('onbutton'));
  var offbtns = [].slice.call(document.getElementsByClassName('offbutton'));
  var onoffbtn = Array.prototype.concat(onbtns,offbtns);
  for (var i = 0; i < onoffbtn.length; i++) {
    var lightname = onoffbtn[i].getAttribute("lightname");
    var action = onoffbtn[i].getAttribute("action");
    if (lightname) {
      var topic = mqtttopic_golightctrl(lightname);
      onoffbtn[i].addEventListener('click', function(topic, action) {
        return function() {  sendMQTT(topic,{Action:action});  };
      }(topic, action));
      topics_to_subscribe[topic] = lightname;
    }
  }
  $(".basiclight_checkbox").on("click",function(event){
    var elemname = event.target.getAttribute("name");
    var topic = mqtttopic_golightctrl(elemname);
    var action = "off";
    if (event.target.checked) {
      action = "on";
    }
    sendMQTT(topic,{Action:action});
  });
  if (webSocketSupport) {
    // register MQTT Update Handler: Basiclights
    $(".basiclight_checkbox").each(function(oelem){
      var topic = mqtttopic_golightctrl(oelem.getAttribute("name"));
      ws.registerContext(topic, function(elem) {
        return function(data) {
          if (data.Action == "1" || data.Action == "on" || data.Action == "send" || data.Action == 1)
            elem.checked = true;
          else
            elem.checked = false;

          //now check if all are checked and thus also check the "all"-checkbox
          var allelem=undefined;
          var checked=true;
          $(".basiclight_checkbox").each(function(elem){
            if (elem.getAttribute("name")=="basiclightAll")
              allelem=elem;
            else
              checked=checked && elem.checked;
          });
          allelem.checked=checked;
        }
      }(oelem));
    });
    // register MQTT Update Handler: RF433 Poweroutlets
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
    $(".sonoff_checkbox").each(function(oelem){
      var topic = mqtttopic_sonoff(oelem.getAttribute("name"));
      ws.registerContext(topic, function(elem) {
        return function(data) {
          console.log(data);
          if (data == "on")
            elem.checked = true;
          else if (data == "toggle")
            elem.checked = !elem.checked;
          else
            elem.checked = false;
        }
      }(oelem));
    });

    // register MQTT Update Handler: Fancy Lights
    registerFunctionForFancyLightUpdate(handleExternalFancySetting);
    // register MQTT Update Handler: ScriptCtrl
    ws.registerContext(mqtttopic_activatescript, handleExternalActivateScript);
  }

  //set background color for fancylightpresetbuttons according to ledr=, ledb=, etc.
  $(".fancylightpresetbutton").each(colorFancyLightPresent);

  popupselect.init();
  popupselect.addSelectHandlerToAll(eventOnFancyLightPresent);
  $(".fancylightpresetbutton").not(".popupselect_option").on("click",eventOnFancyLightPresent);
  $('.scriptctrl_checkbox').on("click",handleChangeScriptCtrl);
  $('#scriptctrlfancyintensityslider').on("change",handleChangeScriptCtrl);
  $("#scriptctrlselect").on("change", handleChangeScriptCtrl);
  $("input.fancyintensityslider").on("change",updateColdWarmWhiteBalanceIntensity)
  $("input.fancybalanceslider").on("change",updateColdWarmWhiteBalanceIntensity)
  $("input.fancyuvslider").on("change",updateUVIntensity)
  $("input.sonoff_checkbox").on("click",eventOnSonOffButton);
  $(".fancylightcolourtempselectorbutton").on("click",popupFancyColorPicker);
  $(document).on("click",function(event){
    if (!document.getElementById("fancycolorpicker").contains(event.target) &&
      $(".fancylightcolourtempselectorbutton").has(event.target).length==0)
    {
      $("#fancycolorpicker").css("visibility","hidden");
    }
  });
  $("#fancycolorpicker_close_button").on("click",function(event){$("#fancycolorpicker").css("visibility","hidden")});
  $("#fancycolorpicker_apply_button").on("click",function(event){
      var R = parseInt($('#R input').val()) || 0;
      var G = parseInt($('#G input').val()) || 0;
      var B = parseInt($('#B input').val()) || 0;
      var CW = parseInt($('#CW input').val()) || 0;
      var WW = parseInt($('#WW input').val()) || 0;
      var settings = {r:R,g:G,b:B,cw:CW,ww:WW,fade:{}};
      sendMQTT(mqtttopic_fancylight(fancycolorpicker_apply_name),settings);
  });
  //draw color picker canvases
  init_colour_temp_picker();
  //define setLedPipePattern(objdata)

  $(".ledpipepresetbutton").on('click', function() {
    var pipepattern = this.getAttribute("pipepattern");
    if (!pipepattern) { return; }
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

  if (webSocketSupport) {
    ws.onopen = ShowConnectionEstablished;
    ws.ondisconnect = ShowWaitingForConnection;
    ws.open(webSocketUrl);
  } else {
    updateButtons("/cgi-bin/mswitch.cgi");
    setInterval("updateButtons(\"/cgi-bin/mswitch.cgi\");", 30*1000);
  }

})();
