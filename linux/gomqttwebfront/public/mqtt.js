
var mqtttopic_activatescript = "action/ceilingscripts/activatescript";
var mqtttopic_pipeledpattern = "action/PipeLEDs/pattern";

function mqtttopic_golightctrl(lightname) {
  return "action/GoLightCtrl/"+lightname;
}
function mqtttopic_fancylight(fancyid) {
  return "action/"+fancyid+"/light";
}
function mqtttopic_kajplatsgroup(groupname) {
  return "zigbee2mqtt/w1/"+groupname+"/set";
}
function mqtttopic_sonoff(name) {
  return "action/"+name+"/POWER"
}
function mqtttopic_esphome_r3_action(name) {
  return "action/"+name+"/command"
}
function mqtttopic_esphome_r3_status(name) {
  return "realraum/"+name+"/state"
}
//@arg whg_and_name e.g.: 'w1/OutletBlueLEDBar'
function mqtttopic_zigbee2mqtt_status(whg_and_name)
{
  return "zigbee2mqtt/"+whg_and_name
}
function mqtttopic_zigbee2mqtt_action(whg_and_name)
{
  return mqtttopic_zigbee2mqtt_status(whg_and_name)+"/set"
}
function mqtttopic_wled_action(wled_name)
{
  return "action/wled/"+wled_name+"/api"
}

/**
 * Convert RGB, CW, WW values to CIE 1931 color space (x/y), color_mode, and color_temp
 * @param {number} r - Red value (0-1000)
 * @param {number} g - Green value (0-1000)
 * @param {number} b - Blue value (0-1000)
 * @param {number} cw - Cold white value (0-1000)
 * @param {number} ww - Warm white value (0-1000)
 * @returns {object} - Object with properties: color_mode, color_temp, and color (x/y)
 */

function rgbCwWwToCIE1931(r, g, b, cw, ww) {
  // Normalize values to 0-1 range
  const rNorm = r / 1000;
  const gNorm = g / 1000;
  const bNorm = b / 1000;
  const cwNorm = cw / 1000;
  const wwNorm = ww / 1000;
  var onoffstate = "ON";
  
  // Calculate total white contribution
  const totalWhite = cwNorm + wwNorm;
  
  if (0 == cw+ww+b+g+r) {
    onoffstate = "OFF";
  }

  // If we have significant white values, use color temperature mode
  if (totalWhite > 0.1 && (rNorm + gNorm + bNorm) < 0.3) {
    // Calculate color temperature based on CW/WW ratio
    const tempRatio = ((-1*cwNorm + wwNorm)+1)/2;
    // Map to 2000K-6500K range (Zigbee standard)
    const colorTemp = Math.round(153 + (555 - 153) * tempRatio);

    return {
      brightness: 255,
      state: onoffstate,
      color_mode: "color_temp",
      color_temp: colorTemp,
    };
  }
  // Otherwise use RGB color mode
  else {
    // Convert RGB to XYZ color space
    let x = rNorm * 0.649926 + gNorm * 0.103455 + bNorm * 0.197109;
    let y = rNorm * 0.234327 + gNorm * 0.743075 + bNorm * 0.022598;
    let z = rNorm * 0.000000 + gNorm * 0.053077 + bNorm * 1.035763;
    
    // Convert XYZ to CIE xy
    const sum = x + y + z;
    if (sum === 0) {
      return {
        brightness: 0,
        state: "OFF",        
        color_mode: "xy",
        color: {x: 0, y: 0},
        color_temp: 0
      };
    }
    
    const cieX = x / sum;
    const cieY = y / sum;
    
    // // Scale to 0-65535 range for Zigbee (some devices expect this)
    // const scaledX = Math.round(cieX * 65535);
    // const scaledY = Math.round(cieY * 65535);
    
    return {
      brightness: 255,
      state: onoffstate,
      color_mode: "xy",
      color: {
        x: cieX,
        y: cieY
      },
      color_temp: 0 // Not used in xy mode
    };
  }
}

function cie1931ToRgbCwWw(state) {
  if (!state || state.state === "OFF") {
    return { r: 0, g: 0, b: 0, cw: 0, ww: 0 };
  }

  if (state.color_mode === "color_temp") {
    const t = Math.max(153, Math.min(555, state.color_temp));
    const tempRatio = (t - 153) / (555 - 153); // inverse of colorTemp formula

    // Assumes cwNorm + wwNorm = 1 and r=g=b=0 (this info is lost by the forward fn)
    const wwNorm = tempRatio;
    const cwNorm = 1 - tempRatio;

    return {
      r: 0, g: 0, b: 0,
      cw: Math.round(cwNorm * 1000),
      ww: Math.round(wwNorm * 1000)
    };
  }

  // color_mode === "xy"
  const { x: cieX = 0, y: cieY = 0 } = state.color || {};

  if (cieX === 0 && cieY === 0) {
    return { r: 0, g: 0, b: 0, cw: 0, ww: 0 };
  }

  // Assume luminance Y = 1 (absolute scale unrecoverable)
  const Y = 1;
  const X = (cieX / cieY) * Y;
  const Z = ((1 - cieX - cieY) / cieY) * Y;

  // Inverse of the RGB->XYZ matrix used in the forward function
  let r = 1.612891 * X - 0.202829 * Y - 0.302365 * Z;
  let g = -0.509836 * X + 1.412086 * Y + 0.066062 * Z;
  let b = 0.026082 * X - 0.072345 * Y + 0.962317 * Z;

  // Clip out-of-gamut negatives
  r = Math.max(0, r);
  g = Math.max(0, g);
  b = Math.max(0, b);

  // Normalize brightest channel to 1000 (scale is arbitrary/unrecoverable)
  const maxC = Math.max(r, g, b, 1e-9);
  r = Math.round((r / maxC) * 1000);
  g = Math.round((g / maxC) * 1000);
  b = Math.round((b / maxC) * 1000);

  return { r, g, b, cw: 0, ww: 0 };
}


var mqtt_scriptctrl_scripts_ = ["off","redshift","ceilingsinus","colorfade","randomcolor","wave","sparkle"];
var mqtt_scriptctrl_scripts_uses_loop_ = ["randomcolor","sparkle"];
var mqtt_scriptctrl_scripts_uses_trigger_for_each_light_ = ["redshift"];
var mqtt_scriptctrl_scripts_support_participating_ = ["redshift","randomcolor","wave","colorfade","ceilingsinus","sparkle"];
var mqtt_fancylights_all = ["ceiling1","ceiling2","ceiling3","ceiling4","ceiling5","ceiling6","abwasch","flooddoor", "memberregal"]
var mqtt_fancylights_all_with_ceilingall = ["ceiling1","ceiling2","ceiling3","ceiling4","ceiling5","ceiling6","abwasch","flooddoor","ceilingAll", "memberregal"]
var mqtt_fancylights_ceilingonly = ["ceiling1","ceiling2","ceiling3","ceiling4","ceiling5","ceiling6"]
var mqtt_fancylights_w2realfunk = ["funkbude"]
var mqtt_fancylights_w2r2w2 = []
var mqtt_fancylights_w2tesla = []


const reverseMapping = (obj) => {
    const reversed = {};
    Object.keys(obj).forEach((key) => {
        reversed[obj[key]] = reversed[obj[key]] || [];
        reversed[obj[key]].push(key);
    });
    return reversed;
};

var mqtt_fancylights_kajplats_name = {
  "ceiling4":"lothr_kajplats_g1",
}

var mqtt_kajplats_fancylights_name = reverseMapping(mqtt_fancylights_kajplats_name)

var r3_led_factors_ = {
  "_default_": {
    r_factor:1,
    g_factor:1,
    b_factor:1,
    ww_factor:1,
    cw_factor:1,
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
    g_factor:1, 
    b_factor:1, 
    ww_factor:1,
    cw_factor:1,
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
  "memberregal": {
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
  if (mqtt_fancylights_kajplats_name[name]) {
    var settings = rgbCwWwToCIE1931(R,G,B,CW,WW);
    sendMQTT(mqtttopic_kajplatsgroup(mqtt_fancylights_kajplats_name[name]),settings);
  } else {
    var settings = {r:R,g:G,b:B,cw:CW,ww:WW,fade:{}};
    sendMQTT("action/"+name+"/light",settings);
  }
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

function eventOnEspHomeButton(event) {
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
  sendMQTT(mqtttopic_esphome_r3_action(name),{state:power});
};

function eventOnZigbee2MqttButton(event) {
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
  sendMQTT(mqtttopic_zigbee2mqtt_action(name),{state:power});
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
  ["ceiling1","ceiling2","ceiling3","ceiling4","ceiling5","ceiling6","abwasch","flooddoor","funkbude","ceilingAll", "memberregal"].forEach(function(fancyid) {
    ws.registerContext("action/"+fancyid+"/light",function(fancyid){
      return function(data) {
        fun(fancyid, data);
      }
    }(fancyid));
  });
}

function registerFunctionForKajplatsUpdate(fun) {
  ["lothr_kajplats_g1","lothr_kajplats_g2","lothr_kajplats_g3","lothr_kajplats_g4","lothr_kajplats_g5","lothr_kajplats_g6"].forEach(function(lightid) {
    ws.registerContext("zigbee2mqtt/w1/"+lightid,function(lightid){
      return function(data) {
        fun(lightid, data);
      }
    }(lightid));
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
  sendMQTT(mqtttopic_fancylight("memberregal"), {r:0,g:0,b:0,ww:0,cw:0,fade:{duration:8000}});
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
  sendMQTT(mqtttopic_fancylight("memberregal"), {r:0,g:0,b:0,ww:0,cw:0,fade:{duration:8000}});
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
  sendMQTT(mqtttopic_fancylight("memberregal"), {r:0,g:0,b:0,ww:0,cw:0,fade:{}});
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
    "g":{"amplitude":200,"offset":300,"phase":0},
    "ww":{"amplitude":90,"offset":300,"phase":1},
    "r":{"amplitude":400,"offset":1000,"phase":2},
    "b":{"amplitude":150,"offset":250,"phase":4},
    "cw":{"amplitude":80,"offset":300,"phase":4},
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
  sendMQTT(mqtttopic_fancylight("ceilingAll"), {r:1000,g:500,b:200,ww:1000,cw:1000});
  sendMQTT(mqtttopic_golightctrl("floodtesla"), {Action:"on"});
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
  sendMQTT(mqtttopic_golightctrl("bluebar"), {Action:"off"});
  sendMQTT(mqtttopic_golightctrl("regalleinwand"), {Action:"off"});
  sendMQTT(mqtttopic_golightctrl("couchred"), {Action:"off"});
  sendMQTT(mqtttopic_esphome_r3_action("subtable"),{state:"OFF"});
  sendMQTT(mqtttopic_esphome_r3_action("loeteckenlicht"),{state:"OFF"});
  sendMQTT(mqtttopic_zigbee2mqtt_action("w1/DeckenfluterLoTHRFenster"), {state:"OFF"});
  sendMQTT(mqtttopic_zigbee2mqtt_action("w1/UltraHellSackerl-Membershelf"), {state:"OFF"});
  sendMQTT(mqtttopic_zigbee2mqtt_action("w1/UltraHellSackerl-AudioShelf"), {state:"OFF"});
  sendMQTT(mqtttopic_zigbee2mqtt_action("w1/OutletAuslageW1"), {state:"OFF"});
  sendMQTT(mqtttopic_wled_action("quadrings"),{"on":false});
  sendMQTT(mqtttopic_wled_action("kaltlichtschrank"),{"on":false});
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
