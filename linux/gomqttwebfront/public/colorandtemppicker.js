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

function rgb2hsv (r,g,b) {
	var computedH = 0;
	var computedS = 0;
	var computedV = 0;

	//remove spaces from input RGB values, convert to int
	var r = parseInt( (''+r).replace(/\s/g,''),10 );
	var g = parseInt( (''+g).replace(/\s/g,''),10 );
	var b = parseInt( (''+b).replace(/\s/g,''),10 );

	if ( r==null || g==null || b==null ||
	 isNaN(r) || isNaN(g)|| isNaN(b) ) {
	alert ('Please enter numeric RGB values!');
	return;
	}
	if (r<0 || g<0 || b<0 || r>255 || g>255 || b>255) {
	alert ('RGB values must be in the range 0 to 255.');
	return;
	}
	r=r/255; g=g/255; b=b/255;
	var minRGB = Math.min(r,Math.min(g,b));
	var maxRGB = Math.max(r,Math.max(g,b));

	// Black-gray-white
	if (minRGB==maxRGB) {
	computedV = minRGB;
	return [0,0,computedV];
	}

	// Colors other than black-gray-white:
	var d = (r==minRGB) ? g-b : ((b==minRGB) ? r-g : b-r);
	var h = (r==minRGB) ? 3 : ((b==minRGB) ? 1 : 5);
	computedH = 60*(h - d/(maxRGB - minRGB));
	computedS = (maxRGB - minRGB)/maxRGB;
	computedV = maxRGB;
	return [computedH,computedS,computedV];
}

function hsv2rgb(h,s,v) {
	var range=0;
	var f=0.0;
	//let's use the same, slighlty adjusted ranges, as for the canvas drawing
	if (h > 0.835) {range=5; f=(h-0.835)/(1-0.835)}
	else if (h > 0.665) {range=4;f=(h-0.665)/(0.835-0.665)}
	else if (h > 0.470) {range=3;f=(h-0.470)/(0.665-0.470)}
	else if (h > 0.333) {range=2;f=(h-0.333)/(0.470-0.333)}
	else if (h > 0.166) {range=1;f=(h-0.166)/(0.333-0.166)}
	else {range = 0; f=h/0.166}

	var v = v * 255.0;
	var p = v * (1 - s);
	var q = v * (1 - f * s);
	var t = v * (1 - (1 - f) * s);
	// console.log(range, h, f);
	switch (range)
	{
	    case 0:
	        return [v, t, p];
	    case 1:
	        return [q, v, p];
	    case 2:
	        return [p, v, t];
	    case 3:
	        return [p, q, v];
	    case 4:
	        return [t, p, v];
	}
	return [v, p, q];
}

var colortemppickerstate = {
	cw:0, ww:0, r:0, g:0, b:0,
	compound_r:0,
	compound_g:0,
	compound_b:0,
};

function drawcolourtemppicker(elemid) {
    var canvas = document.getElementById(elemid);
    var canvas2d = canvas.getContext('2d');
    // console.log(canvas);

    quadGradient(canvas, canvas2d, {
      topLeft: [1,1,1,1],
      topRight: [1, 0xfa/0xff, 0xc0/0xff, 1],
      bottomLeft: [71/0xff, 171/0xff, 1, 1],
      bottomRight: [0, 0, 0, 1]
    });

	var pickcolour = function(event) {
	  if (event.type == "mousemove" && event.buttons != 1) {return;}
	  // getting user coordinates
	  //var untransY = event.pageY - this.offsetTop;
	  var transX = (typeof event.offsetX == "number") ? event.offsetX : event.layerX || 0;
	  var transY = (typeof event.offsetY == "number") ? event.offsetY : event.layerY || 0;
		  // getting image data and RGB values
	  // var CWwhole = Math.trunc((rotatedwidth - x)*1000/rotatedwidth);
	  // var WWwhole = Math.trunc(x*1000/rotatedwidth);
	  var squarewidth = canvas.offsetWidth;
	  // var diamondwidth = Math.sqrt(squarewidth*squarewidth*2);
	  var CW = Math.trunc(Math.max(0,Math.min(1000,(squarewidth - transX)   *1000/squarewidth)));
	  var WW = Math.trunc(Math.max(0,Math.min(1000,(squarewidth - transY)   *1000/squarewidth))); //0...1000
	  //var brightness = (diamondwidth - untransY) *1000/diamondwidth;
	  // colortemppickerstate.intensity = 1000 - Math.trunc(Math.sqrt(transX*transX+transY*transY)*1000/diamondwidth);
	  // colortemppickerstate.balance =  Math.trunc(WW*1000.0/(WW+CW));

	  // making the color the value of the input
	  // var img_data = canvas2d.getImageData(transX, transY, 1, 1);
	  // $('#cwwwoverrgbcolor').css("background-color","rgb("+img_data.data[0]+","+img_data.data[1]+","+img_data.data[2]+")").css("opacity",brightness/1000.0);
	  // $('#cwwwcolor').css("background-color","rgb("+img_data.data[0]+","+img_data.data[1]+","+img_data.data[2]+")");
	  colortemppickerstate.ww = WW;
	  colortemppickerstate.cw = CW;
	  updateColorPreview(colortemppicker_ledfactor_name_);
	};
	$(canvas).on("click",pickcolour);
	$(canvas).on("mousemove",pickcolour);
}

function updateColorPreview(name) {
	calcCompoundRGB(colortemppickerstate, name);
	//var bi = calcDayLevelFromColor(colortemppickerstate);
	$("#cwwwrgbcolor").css("background-color","rgb("+colortemppickerstate.compound_r+","+colortemppickerstate.compound_g+","+colortemppickerstate.compound_b+")");
	$('#R input').val(colortemppickerstate.r);
	$('#G input').val(colortemppickerstate.g);
	$('#B input').val(colortemppickerstate.b);
	$('#R div.colorlevel').css("width",colortemppickerstate.r/10+"%");
	$('#G div.colorlevel').css("width",colortemppickerstate.g/10+"%");
	$('#B div.colorlevel').css("width",colortemppickerstate.b/10+"%");
	$('#CW input').val(colortemppickerstate.cw);
	$('#WW input').val(colortemppickerstate.ww);
	$('#CW div.colorlevel').css("width",colortemppickerstate.cw/10+"%");
	$('#WW div.colorlevel').css("width",colortemppickerstate.ww/10+"%");
	var intensity = (colortemppickerstate.ww+colortemppickerstate.cw)/2;
	var balance =  Math.trunc(colortemppickerstate.ww*1000.0/(colortemppickerstate.ww+colortemppickerstate.cw));
	$('#WhiteBrightness input').val(intensity);
	$('#WhiteTemp input').val(balance);
}

var lower_black_percent = 0.95;
var colortemppicker_ledfactor_name_ = "";

function colorandtempicker_set_ledfactorname(name) {
	colortemppicker_ledfactor_name_ = name;
}

function drawcolourpicker(elemid) {
	var canvas = document.getElementById(elemid);
	var canvas2d = canvas.getContext('2d');

	rainbowHSLpicker(canvas, canvas2d);
	var pickcolour = function(event){
	  if (event.type == "mousemove" && event.buttons != 1) {return;}
	  // getting user coordinates
	  var x = (typeof event.offsetX == "number") ? event.offsetX : event.layerX || 0;
	  var y = (typeof event.offsetY == "number") ? event.offsetY : event.layerY || 0;
	  // getting image data and RGB values
	  //var img_data = canvas2d.getImageData(0,0, this.width,this.height);
	  //var pxoffset = (y*img_data.width+x)*4;
	  //let's  cheat a bit... make upper bit H(S)L picker and lower part how we think HS(V) picker should work
	  if (y <= this.height/2 || y > this.height*lower_black_percent) { //get HSL from image
		  var img_data = canvas2d.getImageData(x, y, 1, 1);
		  var pxoffset=0;
		  var R = img_data.data[pxoffset+0];
		  var G = img_data.data[pxoffset+1];
		  var B = img_data.data[pxoffset+2];
	  } else { //calc HSV from coordiantes
	  	heightminusblack = Math.trunc(this.height*lower_black_percent);
	  	var h = x / this.width;
	  	var s = 1.0;
	  	var v = (heightminusblack - y) / heightminusblack * 2.15;
	  	var rgb = hsv2rgb(h,s,v);
	  	var R = Math.trunc(rgb[0]);
	  	var G = Math.trunc(rgb[1]);
	  	var B = Math.trunc(rgb[2]);
	  }
	  calcCeilingValuesFrom(colortemppickerstate, R, G, B, colortemppicker_ledfactor_name_);
	  updateColorPreview(colortemppicker_ledfactor_name_);
	};
	$(canvas).on("click",pickcolour);
	$(canvas).on("mousemove",pickcolour);
}

function init_colour_temp_picker() {
	drawcolourtemppicker("pickcolourtemp");
	drawcolourpicker("pickcolour");

	changecolourlevel = function(event){
		if (this.value < 0 || this.value > 1000) {return;};
		var variable = this.parentNode.id.toLowerCase();
		colortemppickerstate[variable] = this.value;
		updateColorPreview(colortemppicker_ledfactor_name_);
	}
	$('#R input').on("change",changecolourlevel);
	$('#G input').on("change",changecolourlevel);
	$('#B input').on("change",changecolourlevel);
	$('#WW input').on("change",changecolourlevel);
	$('#CW input').on("change",changecolourlevel);
	changetextfromlevel = function(event){
	  if (event.type == "mousemove" && event.buttons != 1) {return;}
	  var x = (typeof event.offsetX == "number") ? event.offsetX : event.layerX || 0;
	  // var y = event.pageY - this.offsetTop;
	  var promille = Math.trunc(Math.max(0,Math.min(1000,x*1000/this.offsetWidth)));
	  var variable = this.parentNode.id.toLowerCase();
	  colortemppickerstate[variable] = promille;
	  updateColorPreview(colortemppicker_ledfactor_name_);
	};
	$('#R div.colorlevelcontainer').on("mousemove",changetextfromlevel);
	$('#G div.colorlevelcontainer').on("mousemove",changetextfromlevel);
	$('#B div.colorlevelcontainer').on("mousemove",changetextfromlevel);
	$('#WW div.colorlevelcontainer').on("mousemove",changetextfromlevel);
	$('#CW div.colorlevelcontainer').on("mousemove",changetextfromlevel);
	$('#R div.colorlevelcontainer').on("click",changetextfromlevel);
	$('#G div.colorlevelcontainer').on("click",changetextfromlevel);
	$('#B div.colorlevelcontainer').on("click",changetextfromlevel);
	$('#WW div.colorlevelcontainer').on("click",changetextfromlevel);
	$('#CW div.colorlevelcontainer').on("click",changetextfromlevel);
	updateColorPreview(colortemppicker_ledfactor_name_);
}

function rainbowHSLpicker(canvas,ctx) {
    var w = canvas.width;
    var h = canvas.height;
    var gradient;

  	gradient = ctx.createLinearGradient(0,0,w,0);
  	gradient.addColorStop(0,"rgba(255, 0, 0, 1)");
  	gradient.addColorStop(0.166,"rgba(255, 255, 0, 1)");
  	gradient.addColorStop(0.333,"rgba(0, 255, 0, 1)");
  	gradient.addColorStop(0.47,"rgba(0, 255, 255, 1)");
  	gradient.addColorStop(0.665,"rgba(0, 0, 255, 1)");
  	gradient.addColorStop(0.835,"rgba(255, 0, 255, 1)");
  	gradient.addColorStop(1.00,"rgba(255, 0, 0, 1)");
  	ctx.fillStyle = gradient;
    ctx.fillRect(0,0,w,h);
  	gradient = ctx.createLinearGradient(0,0,0,h);
  	gradient.addColorStop(0,"rgba(255, 255, 255, 1)");
  	gradient.addColorStop(0.05,"rgba(255, 255, 255, 1)");
  	gradient.addColorStop(0.35,"rgba(255, 255, 255, 0)");
  	gradient.addColorStop(0.65,"rgba(255,255,255, 0)");
  	gradient.addColorStop(lower_black_percent,"rgba(0, 0, 0, 1)");
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