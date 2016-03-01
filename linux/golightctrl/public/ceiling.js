"use strict";

var baseUrl = '/cgi-bin/mswitch.cgi';

var pc1_state=false;
var pc2_state=false;
var pc3_state=false;
var pc4_state=false;
var pc5_state=false;
var pc6_state=false;

function callBackButtons(req)
{
  if(req.status != 200)
  {
    return;
  }

  var data = JSON.parse(req.responseText);

  if(data['ceiling1']==true)
  {
    pc1_state=true;
    document.getElementById("pc1").style.background="green";
  }
  else
  {
    pc1_state=false;
    document.getElementById("pc1").style.background="red";
  }


  if(data['ceiling2']==true)
  {
    pc2_state=true;
    document.getElementById("pc2").style.background="green";
  }
  else
  {
    pc2_state=false;
    document.getElementById("pc2").style.background="red";
  }

  if(data['ceiling3']==true)
  {
    pc3_state=true;
    document.getElementById("pc3").style.background="green";
  }
  else
  {
    pc3_state=false;
    document.getElementById("pc3").style.background="red";
  }


  if(data['ceiling4']==true)
  {
    pc4_state=true;
    document.getElementById("pc4").style.background="green";
  }
  else
  {
    pc4_state=false;
    document.getElementById("pc4").style.background="red";
  }


  if(data['ceiling5']==true)
  {
    pc5_state=true;
    document.getElementById("pc5").style.background="green";
  }
  else
  {
    pc5_state=false;
    document.getElementById("pc5").style.background="red";
  }

  if(data['ceiling6']==true)
  {
    pc6_state=true;
    document.getElementById("pc6").style.background="green";
  }
  else
  {
    pc6_state=false;
    document.getElementById("pc6").style.background="red";
  }
}

function updateButtons(uri) {
  var req = new XMLHttpRequest;
  req.overrideMimeType("application/json");
  req.open("GET", uri, true);
  req.onload  = function() {callBackButtons(req)};
  req.setRequestHeader("googlechromefix","");
  req.send(null);
}

function sendMultiButton( str ) {
  updateButtons(baseUrl + "?" + str);
}


//function pc(n) switches light n

function pc1()
{
  if(pc1_state===true)
  {
    sendMultiButton("ceiling1=0");
  }
  else
  {
    sendMultiButton("ceiling1=1");
  }
}
function pc2()
{
  if(pc2_state===true)
  {
    sendMultiButton("ceiling2=0");
  }
  else
  {
    sendMultiButton("ceiling2=1");
  }
}
function pc3()
{
  if(pc3_state===true)
  {
    sendMultiButton("ceiling3=0");
  }
  else
  {
    sendMultiButton("ceiling3=1");
  }
}
function pc4()
{
  if(pc4_state===true)
  {
    sendMultiButton("ceiling4=0");
  }
  else
  {
    sendMultiButton("ceiling4=1");
  }
}
function pc5()
{
  if(pc5_state===true)
  {
    sendMultiButton("ceiling5=0");
  }
  else
  {
    sendMultiButton("ceiling5=1");
  }
}

function pc6()
{
  if(pc6_state===true)
  {
    sendMultiButton("ceiling6=0");
  }
  else
  {
    sendMultiButton("ceiling6=1");
  }
}

setInterval("updateButtons(baseUrl);", 30*100 );
updateButtons(baseUrl);
