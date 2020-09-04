
var mqtttopic_activatescript = "action/ceilingscripts/activatescript";
var mqtttopic_pipeledpattern = "action/PipeLEDs/pattern";
function mqtttopic_golightctrl(lightname) {
  return "action/GoLightCtrl/"+lightname;
}
function mqtttopic_fancylight(fancyid) {
  return "action/"+fancyid+"/light";
}
function mqtttopic_sonoff(name) {
  return "action/"+name+"/POWER"
}

var mqtt_scriptctrl_scripts_ = ["off","redshift","ceilingsinus","colorfade","randomcolor","wave","sparkle"];
var mqtt_scriptctrl_scripts_uses_loop_ = ["randomcolor","sparkle"];
var mqtt_scriptctrl_scripts_uses_trigger_for_each_light_ = ["redshift"];
var mqtt_scriptctrl_scripts_support_participating_ = ["redshift","randomcolor","wave","colorfade","ceilingsinus","sparkle"];
var mqtt_fancylights_all = ["ceiling1","ceiling2","ceiling3","ceiling4","ceiling5","ceiling6","abwasch","flooddoor"]
var mqtt_fancylights_all_with_ceilingall = ["ceiling1","ceiling2","ceiling3","ceiling4","ceiling5","ceiling6","abwasch","flooddoor","ceilingAll"]
var mqtt_fancylights_ceilingonly = ["ceiling1","ceiling2","ceiling3","ceiling4","ceiling5","ceiling6"]
var mqtt_fancylights_w2realfunk = ["funkbude"]
var mqtt_fancylights_w2r2w2 = []
var mqtt_fancylights_w2tesla = []


var r3_led_factors_ = {
  "_default_": {
    r_factor:1,
    g_factor:5, //green 5 times as bright as red
    b_factor:10, //blue 2 times as bright as green
    ww_factor:22, //yes warmwhite is about 22 times as bright as red
    cw_factor:18,
  },
  "flooddoor": {
    r_factor:4,
    g_factor:4,
    b_factor:4,
    ww_factor:12,
    cw_factor:12,
  },
  "ceiling1": {
    r_factor:1,
    g_factor:5, //green 5 times as bright as red
    b_factor:10, //blue 2 times as bright as green
    ww_factor:22, //yes warmwhite is about 22 times as bright as red
    cw_factor:18,
  },
  "ceiling2": {
    r_factor:1,
    g_factor:5, //green 5 times as bright as red
    b_factor:10, //blue 2 times as bright as green
    ww_factor:22, //yes warmwhite is about 22 times as bright as red
    cw_factor:18,
  },
  "ceiling3": {
    r_factor:1,
    g_factor:5, //green 5 times as bright as red
    b_factor:10, //blue 2 times as bright as green
    ww_factor:22, //yes warmwhite is about 22 times as bright as red
    cw_factor:18,
  },
  "ceiling4": {
    r_factor:1,
    g_factor:5, //green 5 times as bright as red
    b_factor:10, //blue 2 times as bright as green
    ww_factor:22, //yes warmwhite is about 22 times as bright as red
    cw_factor:18,
  },
  "ceiling5": {
    r_factor:1,
    g_factor:5, //green 5 times as bright as red
    b_factor:10, //blue 2 times as bright as green
    ww_factor:22, //yes warmwhite is about 22 times as bright as red
    cw_factor:18,
  },
  "ceiling6": {
    r_factor:1,
    g_factor:5, //green 5 times as bright as red
    b_factor:10, //blue 2 times as bright as green
    ww_factor:22, //yes warmwhite is about 22 times as bright as red
    cw_factor:18,
  },
  "abwasch": {
    r_factor:4,
    g_factor:4,
    b_factor:4,
    ww_factor:12,
    cw_factor:12,
  },
  "funkbude": {
    r_factor:4,
    g_factor:4,
    b_factor:4,
    ww_factor:12,
    cw_factor:12,
  },
};

function getr3ledfactors(name) {
  console.log(name);
  if (r3_led_factors_[name])
    return r3_led_factors_[name];
  else
    return r3_led_factors_["_default_"];
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

function eventOnRawMqttElement(event) {
  var topic = event.target.getAttribute("topic");
  var payloadobj = JSON.parse(event.target.getAttribute("payload"));
  if (payloadobj) {
    sendMQTT(topic, payloadobj);
  }
}

function eventOnFancyLightPresent(event) {
  var name = event.target.getAttribute("name");
  if (!name) { return;  }
  var R = parseInt(event.target.getAttribute("ledr")) || 0;
  var G = parseInt(event.target.getAttribute("ledg")) || 0;
  var B = parseInt(event.target.getAttribute("ledb")) || 0;
  var CW = parseInt(event.target.getAttribute("ledcw")) || 0;
  var WW = parseInt(event.target.getAttribute("ledww")) || 0;
  var settings = {r:R,g:G,b:B,cw:CW,ww:WW,fade:{}};
  sendMQTT("action/"+name+"/light",settings);
};

function eventOnSonOffButton(event) {
  var name = event.target.getAttribute("name");
  if (!name) { return;  }
  var power = event.target.getAttribute("power");
  if (!power && event.target.getAttribute("type")=="checkbox")
  {
    if (event.target.checked)
      power="ON";
    else
      power="OFF";
  }
  if (power != "ON" && power != "OFF" && power != "TOGGLE") {return;}
  sendMQTT(mqtttopic_sonoff(name),power);
};

function colorFancyLightPresent(elem) {
  var R = parseInt(elem.getAttribute("ledr")) || 0;
  var G = parseInt(elem.getAttribute("ledg")) || 0;
  var B = parseInt(elem.getAttribute("ledb")) || 0;
  var CW = parseInt(elem.getAttribute("ledcw")) || 0;
  var WW = parseInt(elem.getAttribute("ledww")) || 0;
  var settings = {r:R,g:G,b:B,cw:CW,ww:WW,fade:{}};
  var name = elem.getAttribute("name");
  calcCompoundRGB(settings, name);
  elem.style.backgroundColor="rgb("+settings.compound_r+","+settings.compound_g+","+settings.compound_b+")";
}

//takes function with signature (fancyid, data)
//and calls it if fancy light updates externally
function registerFunctionForFancyLightUpdate(fun) {
  ["ceiling1","ceiling2","ceiling3","ceiling4","ceiling5","ceiling6","abwasch","flooddoor","funkbude","ceilingAll"].forEach(function(fancyid) {
    ws.registerContext("action/"+fancyid+"/light",function(fancyid){
      return function(data) {
        fun(fancyid, data);
      }
    }(fancyid));
  });
}

function calcDayLevelFromColor(data)
{
  var day_factor;
  if (data.cw+data.ww == 0) {
    day_factor = 0.0;
  } else {
    day_factor = data.cw * 1.0 / (data.cw+data.ww) - data.r / 1000.0;
    day_factor = Math.min(1.0,Math.max(-1.0,  day_factor ));
  }

  var numvalues = 1;
  var value = Math.min(1.0,(data.ww + (data.r/3.0) + data.cw) / 1000.0);
  /*
  if (day_factor < 0.0) {
    value += data.r / day_factor / -1000.0;
    numvalues += 1;
  }
  */
  if (day_factor > 0.0) {
    value += data.cw / day_factor / 1000.0;
    numvalues += 1;
  }
  //calc average
  value /= numvalues;
  return {"balance":day_factor, "intensity":value};
}

function calcColorFromDayLevel(day_factor, value)
{
  var day_factor = Math.min(1.0,Math.max(-1.0,day_factor));
  var r = 1000 * value * Math.max(0.0, -1.0 * day_factor);
  var b = 0;
  var cw = 1000 * value * Math.max(0.0, day_factor);
  var ww = Math.max(0,1000 * value - cw - (r/3));
  return {"r":Math.trunc(r), "b":Math.trunc(b), "cw":Math.trunc(cw), "ww":Math.trunc(ww)};
}

function calcCompoundRGB(data, name)
{
  var warmwhite_representation = [255.0/255, 250.0/255, 192/255];
  var coldwhite_representation = [220.0/255, 220.0/255, 255/255];

  //fill data with zero if missing
  data.r  = data.r  || 0;
  data.g  = data.g  || 0;
  data.b  = data.b  || 0;
  data.ww = data.ww || 0;
  data.cw = data.cw || 0;

  var magn_orig = Math.sqrt(data.r*data.r+data.g*data.g+data.b*data.b+data.ww*data.ww+data.cw*data.cw);

  if (magn_orig == 0) {
    data.compound_r = 0;
    data.compound_g = 0;
    data.compound_b = 0;
    return;
  }

  var ledfactors = getr3ledfactors(name);

  var r = data.r*ledfactors["r_factor"];
  var g = data.g*ledfactors["g_factor"];
  var b = data.b*ledfactors["b_factor"];
  var cw = data.cw*ledfactors["cw_factor"];
  var ww = data.ww*ledfactors["ww_factor"];

  //vector magnitude
  var magn_new = Math.sqrt(r*r+g*g+b*b+cw*cw+ww*ww);
  var scale = magn_orig/magn_new;

  r *= scale;
  g *= scale;
  b *= scale;
  ww *= scale;
  cw *= scale;

  r += ww*warmwhite_representation[0] + cw*coldwhite_representation[0];
  g += ww*warmwhite_representation[1] + cw*coldwhite_representation[1];
  b += ww*warmwhite_representation[2] + cw*coldwhite_representation[2];

  var maximum = Math.max(1000,r,g,b);

  //now fit to box
  r = r * 255.0 / maximum;
  g = g * 255.0 / maximum;
  b = b * 255.0 / maximum;

  data.compound_r = Math.min(255,Math.floor(r));
  data.compound_g = Math.min(255,Math.floor(g));
  data.compound_b = Math.min(255,Math.floor(b));
}

function calcCeilingValuesFrom(data,r,g,b,name)
{
  var magn_orig = Math.sqrt(r*r+g*g+b*b);

  var ledfactors = getr3ledfactors(name);

  r = r/ledfactors["r_factor"];
  g = g/ledfactors["g_factor"];
  b = b/ledfactors["b_factor"];
  var magn_new = Math.sqrt(r*r+g*g+b*b);

  //scale color vector to original magnitude
  var scale = magn_orig/magn_new;
  r *= scale;
  g *= scale;
  b *= scale;

  //fit into 255 by 255 by 255 box
  var fitting = Math.max(255,r,g,b);

  data.r = Math.trunc(r * 1000 / fitting)
  data.g = Math.trunc(g * 1000 / fitting)
  data.b = Math.trunc(b * 1000 / fitting)
}

function ShowWaitingForConnection() {
  $("div.waitingoverlay").css("display","initial");
}

function ShowConnectionEstablished() {
  $("div.waitingoverlay").css("display","none");
}

function ceilingPreset_BeamerTalkMode()
{
  sendMQTT(mqtttopic_activatescript, {script:"off"});
  sendMQTT(mqtttopic_golightctrl("basiclight1"), {Action:"off"});
  sendMQTT(mqtttopic_golightctrl("basiclight2"), {Action:"off"});
  sendMQTT(mqtttopic_golightctrl("basiclight3"), {Action:"on"});
  sendMQTT(mqtttopic_golightctrl("basiclight4"), {Action:"on"});
  sendMQTT(mqtttopic_golightctrl("basiclight5"), {Action:"off"});
  sendMQTT(mqtttopic_golightctrl("basiclight6"), {Action:"off"});
  sendMQTT(mqtttopic_fancylight("ceiling1"), {r:0,g:0,b:0,ww:0,cw:0,fade:{duration:8000}});
  sendMQTT(mqtttopic_fancylight("ceiling2"), {r:0,g:0,b:0,ww:0,cw:900,fade:{duration:8000}});
  sendMQTT(mqtttopic_fancylight("ceiling3"), {r:0,g:0,b:0,ww:0,cw:1000,fade:{duration:8000}});
  sendMQTT(mqtttopic_fancylight("ceiling4"), {r:0,g:0,b:0,ww:0,cw:1000,fade:{duration:8000}});
  sendMQTT(mqtttopic_fancylight("ceiling5"), {r:0,g:0,b:0,ww:0,cw:900,fade:{duration:8000}});
  sendMQTT(mqtttopic_fancylight("ceiling6"), {r:0,g:0,b:0,ww:0,cw:0,fade:{duration:8000}});
  sendMQTT(mqtttopic_fancylight("abwasch"), {r:0,g:660,b:0,ww:500,cw:500,fade:{duration:8000}});
  sendMQTT(mqtttopic_fancylight("flooddoor"), {r:0,g:0,b:0,ww:800,cw:800,fade:{duration:8000}});
  sendMQTT(mqtttopic_golightctrl("floodtesla"), {Action:"off"});
  sendMQTT(mqtttopic_golightctrl("subtable"), {Action:"on"});
}

function ceilingPreset_BeamerTalkPauseMode()
{
  sendMQTT(mqtttopic_activatescript, {script:"off"});
  sendMQTT(mqtttopic_golightctrl("basiclight1"), {Action:"off"});
  sendMQTT(mqtttopic_golightctrl("basiclight2"), {Action:"off"});
  sendMQTT(mqtttopic_golightctrl("basiclight3"), {Action:"on"});
  sendMQTT(mqtttopic_golightctrl("basiclight4"), {Action:"on"});
  sendMQTT(mqtttopic_golightctrl("basiclight5"), {Action:"off"});
  sendMQTT(mqtttopic_golightctrl("basiclight6"), {Action:"off"});
  sendMQTT(mqtttopic_fancylight("ceiling1"), {r:0,g:0,b:0,ww:0,cw:500,fade:{duration:8000}});
  sendMQTT(mqtttopic_fancylight("ceiling2"), {r:0,g:0,b:0,ww:1000,cw:1000,fade:{duration:8000}});
  sendMQTT(mqtttopic_fancylight("ceiling3"), {r:800,g:0,b:0,ww:1000,cw:1000,fade:{duration:8000}});
  sendMQTT(mqtttopic_fancylight("ceiling4"), {r:800,g:0,b:0,ww:1000,cw:1000,fade:{duration:8000}});
  sendMQTT(mqtttopic_fancylight("ceiling5"), {r:0,g:0,b:0,ww:1000,cw:1000,fade:{duration:8000}});
  sendMQTT(mqtttopic_fancylight("ceiling6"), {r:0,g:0,b:0,ww:0,cw:500,fade:{duration:8000}});
  sendMQTT(mqtttopic_fancylight("abwasch"), {r:0,g:1000,b:0,ww:1000,cw:800,fade:{duration:8000}});
  sendMQTT(mqtttopic_fancylight("flooddoor"), {r:800,g:0,b:0,ww:1000,cw:1000,fade:{duration:8000}});
  sendMQTT(mqtttopic_golightctrl("floodtesla"), {Action:"on"});
  sendMQTT(mqtttopic_golightctrl("subtable"), {Action:"on"});
}

function ceilingPreset_BeamerMovieMode()
{
  sendMQTT(mqtttopic_activatescript, {script:"off"});
  sendMQTT(mqtttopic_golightctrl("basiclightAll"), {Action:"off"});
  sendMQTT(mqtttopic_fancylight("ceiling1"), {r:0,g:0,b:0,ww:0,cw:0,fade:{}});
  sendMQTT(mqtttopic_fancylight("ceiling2"), {r:0,g:0,b:0,ww:0,cw:0,fade:{}});
  sendMQTT(mqtttopic_fancylight("ceiling3"), {r:50,g:0,b:0,ww:100,cw:0,fade:{}});
  sendMQTT(mqtttopic_fancylight("ceiling4"), {r:50,g:0,b:0,ww:100,cw:0,fade:{}});
  sendMQTT(mqtttopic_fancylight("ceiling5"), {r:0,g:0,b:0,ww:0,cw:0,fade:{}});
  sendMQTT(mqtttopic_fancylight("ceiling6"), {r:0,g:0,b:0,ww:0,cw:0,fade:{}});
  sendMQTT(mqtttopic_fancylight("flooddoor"), {r:0,g:0,b:0,ww:0,cw:0,fade:{}});
  sendMQTT(mqtttopic_fancylight("abwasch"), {r:0,g:0,b:0,ww:0,cw:0,fade:{}});
  sendMQTT(mqtttopic_golightctrl("floodtesla"), {Action:"off"});
  sendMQTT(mqtttopic_golightctrl("subtable"), {Action:"off"});
}

function ceilingPreset_RedShiftMost()
{
  sendMQTT(mqtttopic_golightctrl("basiclightAll"), {Action:"off"});
  sendMQTT(mqtttopic_fancylight("ceilingAll"), {r:0,g:0,b:0,ww:0,cw:0});
  sendMQTT(mqtttopic_activatescript, {script:"redshift",participating:["ceiling1","ceiling2","ceiling3","ceiling6"],"value":0.99});
}

function ceilingPreset_AlienSky()
{
  sendMQTT(mqtttopic_golightctrl("basiclightAll"), {Action:"off"});
  sendMQTT(mqtttopic_fancylight("ceilingAll"), {r:0,g:0,b:0,ww:0,cw:0});
  sendMQTT(mqtttopic_activatescript, {"script":"ceilingsinus",
    "g":{"amplitude":200,"offset":200,"phase":0},
    "ww":{"amplitude":90,"offset":100,"phase":1},
    "r":{"amplitude":400,"offset":900,"phase":2},
    "b":{"amplitude":150,"offset":150,"phase":4},
    "cw":{"amplitude":80,"offset":100,"phase":4},
    "fadeduration":3000}
    );
  sendMQTT(mqtttopic_golightctrl("subtable"), {Action:"on"});
}

function ceilingPreset_DimRandomColor()
{
  sendMQTT(mqtttopic_golightctrl("basiclightAll"), {Action:"off"});
  sendMQTT(mqtttopic_fancylight("ceilingAll"), {r:0,g:0,b:0,ww:0,cw:0});
  sendMQTT(mqtttopic_activatescript, {"script":"randomcolor","value":0.3});
}

function ceilingPreset_SuperFullEverything()
{
  sendMQTT(mqtttopic_activatescript, {script:"off"});
  sendMQTT(mqtttopic_golightctrl("basiclightAll"), {Action:"on"});
  sendMQTT(mqtttopic_fancylight("ceilingAll"), {r:1000,g:1000,b:1000,ww:1000,cw:1000});
  //sendMQTT(mqtttopic_golightctrl("floodtesla"), {Action:"on"});
  sendMQTT(mqtttopic_golightctrl("subtable"), {Action:"on"});
}

function ceilingPreset_AlmostEverything()
{
  sendMQTT(mqtttopic_activatescript, {script:"off"});
  sendMQTT(mqtttopic_golightctrl("basiclightAll"), {Action:"on"});
  sendMQTT(mqtttopic_fancylight("ceilingAll"), {r:1000,g:500,b:200,ww:1000,cw:1000});
  sendMQTT(mqtttopic_golightctrl("subtable"), {Action:"on"});
}

function ceilingPreset_MostBasic()
{
  sendMQTT(mqtttopic_activatescript, {script:"off"});
  sendMQTT(mqtttopic_golightctrl("basiclight1"), {Action:"on"});
  sendMQTT(mqtttopic_golightctrl("basiclight2"), {Action:"on"});
  sendMQTT(mqtttopic_golightctrl("basiclight3"), {Action:"on"});
  sendMQTT(mqtttopic_golightctrl("basiclight4"), {Action:"off"});
  sendMQTT(mqtttopic_golightctrl("basiclight5"), {Action:"off"});
  sendMQTT(mqtttopic_golightctrl("basiclight6"), {Action:"on"});
  sendMQTT(mqtttopic_fancylight("ceilingAll"), {r:0,g:0,b:0,ww:0,cw:0});
}

function ceilingPreset_MixedForWork()
{
  sendMQTT(mqtttopic_activatescript, {script:"off"});
  sendMQTT(mqtttopic_golightctrl("basiclight1"), {Action:"on"});
  sendMQTT(mqtttopic_golightctrl("basiclight2"), {Action:"on"});
  sendMQTT(mqtttopic_golightctrl("basiclight3"), {Action:"off"});
  sendMQTT(mqtttopic_golightctrl("basiclight4"), {Action:"off"});
  sendMQTT(mqtttopic_golightctrl("basiclight5"), {Action:"off"});
  sendMQTT(mqtttopic_golightctrl("basiclight6"), {Action:"on"});
  sendMQTT(mqtttopic_fancylight("ceiling1"), {r:0,g:0,b:0,ww:0,cw:0});
  sendMQTT(mqtttopic_fancylight("ceiling2"), {r:0,g:0,b:0,ww:0,cw:0});
  sendMQTT(mqtttopic_fancylight("ceiling3"), {r:0,g:0,b:0,ww:1000,cw:0});
  sendMQTT(mqtttopic_fancylight("ceiling4"), {r:0,g:0,b:0,ww:600,cw:0});
  sendMQTT(mqtttopic_fancylight("ceiling5"), {r:0,g:0,b:0,ww:600,cw:0});
  sendMQTT(mqtttopic_fancylight("ceiling6"), {r:0,g:0,b:0,ww:0,cw:0});
  sendMQTT(mqtttopic_fancylight("flooddoor"), {r:0,g:0,b:0,ww:600,cw:300});
}

function ceilingPreset_AllOff()
{
  sendMQTT(mqtttopic_activatescript, {script:"off"});
  sendMQTT(mqtttopic_golightctrl("basiclightAll"), {Action:"off"});
  sendMQTT(mqtttopic_fancylight("ceilingAll"), {r:0,g:0,b:0,ww:0,cw:0});
  sendMQTT(mqtttopic_golightctrl("floodtesla"), {Action:"off"});
  sendMQTT(mqtttopic_golightctrl("bluebar"), {Action:"off"});
  sendMQTT(mqtttopic_golightctrl("regalleinwand"), {Action:"off"});
  sendMQTT(mqtttopic_golightctrl("couchred"), {Action:"off"});
  sendMQTT(mqtttopic_golightctrl("subtable"), {Action:"off"});
}

function ceilingPreset_ColorWave()
{
  sendMQTT(mqtttopic_golightctrl("basiclightAll"), {Action:"off"});
  sendMQTT(mqtttopic_activatescript, {"script":"wave","colourlist":[
      {r:1000,g:0,b:0,ww:0,cw:0},
      {r:800,g:0,b:100,ww:0,cw:0},
      {r:0,g:0,b:300,ww:0,cw:0},
      {r:0,g:500,b:100,ww:0,cw:0},
      {r:0,g:800,b:0,ww:0,cw:0},
      {r:800,g:200,b:0,ww:0,cw:0},
    ], "fadeduration":5000}
    );
}

function ceilingPreset_BlueWave()
{
  sendMQTT(mqtttopic_golightctrl("basiclightAll"), {Action:"off"});
  sendMQTT(mqtttopic_fancylight("ceilingAll"), {r:0,g:0,b:0,ww:0,cw:0});
  sendMQTT(mqtttopic_activatescript, {"script":"wave","colourlist":[
      {"r":200,"g":0,"b":1000,"cw":0,"ww":0},
      {"r":0,"g":0,"b":0,"cw":50,"ww":50},
      {"r":0,"g":0,"b":0,"cw":50,"ww":50},
      {"r":0,"g":0,"b":0,"cw":50,"ww":50},
    ], "fadeduration":2000,
      "reversed":1}
    );
}

function ceilingPreset_SkyWithClouds()
{
  sendMQTT(mqtttopic_golightctrl("basiclightAll"), {Action:"off"});
  sendMQTT(mqtttopic_fancylight("ceilingAll"), {r:0,g:0,b:0,ww:0,cw:0});
  sendMQTT(mqtttopic_activatescript, {"script":"ceilingsinus","value":1.0}
    );
}
