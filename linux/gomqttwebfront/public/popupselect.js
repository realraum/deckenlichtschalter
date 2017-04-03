//// (c) Bernhard Tittelbach, 2017
//// License: MIT

var popupselect = {
	target:undefined,
	onselect:{},
	options: {
		class_triggerpopup:"popupselect_trigger",
		class_popupoverlay:"popupselect_overlay",
		class_option:"popupselect_option",
		attribute_specifiying_elementid_with_options:"optionsid",
	},

	popupselectOpen:function(event) {
		var x = event.pageX;
		var y = event.pageY;
		popupselect.target = event.target;
		var optionselementid = event.target.getAttribute("optionsid");
		var oelem = $(document.getElementById(optionselementid));
		oelem.css("left",x - oelem.width()/2).css("top",y - oelem.height()/2).css("visibility","visible").animate(1000);
	},

	popupselectSelect:function(event) {
		//console.log(event.target);
		$("."+popupselect.options.class_popupoverlay).css("visibility","hidden").animate(500);
		var selectedbtn = $(event.target);
		if (selectedbtn.hasClass(popupselect.options.class_option))
		{
			$(popupselect.target).html(selectedbtn.html());
			$(popupselect.target).attr("style",selectedbtn.attr("style"));
			var fun = popupselect.onselect[event.target];
			if (fun) {fun(event);}
		}
	},

	init: function(options) {
		if (options)
		{
			$.extend(this.options,options)
		}
		$("."+this.options.class_triggerpopup).mousedown(this.popupselectOpen);
		//$("button.hoverselector").click(popupHoverSelector);
		$(document).mouseup(this.popupselectSelect);
	},

	addSelectHandler:function(elem, func) {
		if ($(elem).hasClass(this.options.class_option) && typeof(func) == "function")
		{
			this.onselect[elem] = func;
		}
	},

	addSelectHandlerToAll:function(func) {
		if (typeof(func) == "function")
		{
			$("."+this.options.class_option).each(function(idx, elem) {popupselect.onselect[elem] = func;});
		}

	},
}
