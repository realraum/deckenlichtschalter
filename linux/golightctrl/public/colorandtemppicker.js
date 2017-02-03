// This blog post was a big help:
// http://www.nixtu.info/2014/12/html5-canvas-gradients-rectangle.html

// http://www.javascripter.net/faq/rgbtohex.htm
function rgbToHex(R,G,B) {return toHex(R)+toHex(G)+toHex(B)}

function toHex(n) {
  n = parseInt(n,10);
  if (isNaN(n)) return "00";
  n = Math.max(0,Math.min(n,255));
  return "0123456789ABCDEF".charAt((n-n%16)/16)  + "0123456789ABCDEF".charAt(n%16);
}

function drawcolourtemppicker(elemid) {
    var canvas = document.getElementById(elemid);
    var canvas2d = canvas.getContext('2d');
    console.log(canvas);

    quadGradient(canvas, canvas2d, {
      topLeft: [1,1,1,1],
      topRight: [1, 0xfa/0xff, 0xc0/0xff, 1],
      bottomLeft: [71/0xff, 171/0xff, 1, 1],
      bottomRight: [0, 0, 0, 1]
    });

	var pickcolour = function(event){
	  if (event.type == "mousemove" && event.buttons != 1) {return;}
	  // getting user coordinates
	  var untransX = event.pageX - this.offsetLeft;
	  var untransY = event.pageY - this.offsetTop;
	  var transX = (typeof event.offsetX == "number") ? event.offsetX : event.layerX || 0;
	  var transY = (typeof event.offsetY == "number") ? event.offsetY : event.layerY || 0;

	  // getting image data and RGB values
	  // var CWwhole = Math.trunc((rotatedwidth - x)*1000/rotatedwidth);
	  // var WWwhole = Math.trunc(x*1000/rotatedwidth);
	  var squarewidth = $(canvas).width();
	  var diamondwidth = Math.sqrt(squarewidth*squarewidth*2);
	  var CW = (squarewidth - transX)   *1000/squarewidth;
	  var WW = (squarewidth - transY)   *1000/squarewidth; //0...1000
	  var brightness = (diamondwidth - untransY) *1000/diamondwidth;
	  var ctempmix = WW*1000/CW/2;	
	  // making the color the value of the input
	  $('#CW input').val(Math.trunc(CW));
	  $('#WW input').val(Math.trunc(WW));
	  $('#CW div.colorlevel').css("width",CW/10+"%");
	  $('#WW div.colorlevel').css("width",WW/10+"%");
	  $('#WhiteBrightness input').val(brightness);
	  $('#WhiteTemp input').val(ctempmix);
	  var img_data = canvas2d.getImageData(transX, transY, 1, 1);
	  $('#cwwwcolor').css("background-color","rgb("+img_data.data[0]+","+img_data.data[1]+","+img_data.data[2]+")").css("opacity",brightness/1000.0);
	};
	$(canvas).click(pickcolour);
	$(canvas).mousemove(pickcolour);

}


function drawcolourpicker(elemid) {
	var canvas = document.getElementById(elemid);
	var canvas2d = canvas.getContext('2d');

	rainbowWhiteBlackGradient(canvas, canvas2d);
	var pickcolour = function(event){
	  if (event.type == "mousemove" && event.buttons != 1) {return;}
	  // getting user coordinates
	  var x = event.pageX - this.offsetLeft;
	  var y = event.pageY - this.offsetTop;
	  // getting image data and RGB values
	  //var img_data = canvas2d.getImageData(0,0, this.width,this.height);
	  //var pxoffset = (y*img_data.width+x)*4;
	  var img_data = canvas2d.getImageData(x, y, 1, 1);
	  var pxoffset=0;
	  var R = img_data.data[pxoffset+0];
	  var G = img_data.data[pxoffset+1];
	  var B = img_data.data[pxoffset+2];
	  // making the color the value of the input
	  $('#R input').val(Math.trunc(R*1000/255));
	  $('#G input').val(Math.trunc(G*1000/255));
	  $('#B input').val(Math.trunc(B*1000/255));
	  $('#R div.colorlevel').css("width",R*100/255+"%");
	  $('#G div.colorlevel').css("width",G*100/255+"%");
	  $('#B div.colorlevel').css("width",B*100/255+"%");
	  $('#rgbcolor').css("background-color","rgb("+R+","+G+","+B+")");
	};
	$(canvas).click(pickcolour);
	$(canvas).mousemove(pickcolour);
}

function init_colour_temp_picker() {
	drawcolourtemppicker("pickcolourtemp");
	drawcolourpicker("pickcolour");

	changecolourlevel = function(event){if (this.value < 0 || this.value > 1000) {return;}; $(this).siblings().find('.colorlevel').css("width",this.value/10+"%");}
	$('#R input').change(changecolourlevel);
	$('#G input').change(changecolourlevel);
	$('#B input').change(changecolourlevel);
	$('#WW input').change(changecolourlevel);
	$('#CW input').change(changecolourlevel);
	changetextfromlevel = function(event){
	  if (event.type == "mousemove" && event.buttons != 1) {return;}
	  var x = event.pageX - this.offsetLeft;
	  // var y = event.pageY - this.offsetTop;
	  var promille = Math.trunc(x*1000/this.offsetWidth);
	  $(this).siblings("input").val(promille);
	  $(this).find('.colorlevel').css("width",x);
	};
	$('#R div.colorlevelcontainer').mousemove(changetextfromlevel);
	$('#G div.colorlevelcontainer').mousemove(changetextfromlevel);
	$('#B div.colorlevelcontainer').mousemove(changetextfromlevel);
	$('#WW div.colorlevelcontainer').mousemove(changetextfromlevel);
	$('#CW div.colorlevelcontainer').mousemove(changetextfromlevel);
	$('#R div.colorlevelcontainer').click(changetextfromlevel);
	$('#G div.colorlevelcontainer').click(changetextfromlevel);
	$('#B div.colorlevelcontainer').click(changetextfromlevel);
	$('#WW div.colorlevelcontainer').click(changetextfromlevel);
	$('#CW div.colorlevelcontainer').click(changetextfromlevel);
}

function rainbowWhiteBlackGradient(canvas,ctx) { 
    var w = canvas.width;
    var h = canvas.height;
    var gradient;
  
  	gradient = ctx.createLinearGradient(0,0,w,0);
  	gradient.addColorStop(0,"rgba(255, 0, 0, 1)");
  	gradient.addColorStop(0.15,"rgba(255, 255, 0, 1)");
  	gradient.addColorStop(0.30,"rgba(0, 255, 0, 1)");
  	gradient.addColorStop(0.50,"rgba(0, 255, 255, 1)");
  	gradient.addColorStop(0.65,"rgba(0, 0, 255, 1)");
  	gradient.addColorStop(0.80,"rgba(255, 0, 255, 1)");
  	gradient.addColorStop(1.00,"rgba(255, 0, 0, 1)");
  	ctx.fillStyle = gradient;
    ctx.fillRect(0,0,w,h);
  	gradient = ctx.createLinearGradient(0,0,0,h);
  	gradient.addColorStop(0,"rgba(255, 255, 255, 1)");
  	gradient.addColorStop(0.05,"rgba(255, 255, 255, 1)");
  	gradient.addColorStop(0.40,"rgba(255, 255, 255, 0)");
  	gradient.addColorStop(0.60,"rgba(255,255,255, 0)");
  	gradient.addColorStop(0.95,"rgba(0, 0, 0, 1)");
  	gradient.addColorStop(1.0,"rgba(0, 0, 0, 1)");
  	ctx.fillStyle = gradient;
    ctx.fillRect(0,0,w,h);  	
}

function quadGradient(canvas, ctx, corners) { 
    var w = canvas.width;
    var h = canvas.height;
    var gradient, startColor, endColor, fac;
  
    for(var i = 0; i < h; i++) {
        gradient = ctx.createLinearGradient(0, i, w, i);
        fac = i / (h - 1);

        startColor = arrayToRGBA(
          lerp(corners.topLeft, corners.bottomLeft, fac)
        );
        endColor = arrayToRGBA(
          lerp(corners.topRight, corners.bottomRight, fac)
        );
      
        gradient.addColorStop(0, startColor);
        gradient.addColorStop(1, endColor);
      
        ctx.fillStyle = gradient;
        ctx.fillRect(0, i, w, i);
    }
}

function arrayToRGBA(arr) {
    var ret = arr.map(function(v) {
        // map to [0, 255] and clamp
        return Math.max(Math.min(Math.round(v * 255), 255), 0);
    });

    // alpha should retain its value
    ret[3] = arr[3];
  
    return 'rgba(' + ret.join(',') + ')';
}

function lerp(a, b, fac) {
    return a.map(function(v, i) {
        return v * (1 - fac) + b[i] * fac;
    });
}