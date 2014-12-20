#!/bin/sh

for QUERY in `echo $QUERY_STRING | tr '&' ' '`; do
  for VALUE in `echo $QUERY | tr '=' ' '`; do
    if [ "$VALUE" = "id" ]; then
      ID='?'
    elif [ "$ID" = "?" ]; then
      ID=$VALUE
    elif [ "$VALUE" = "power" ]; then
      POWER='?'
    elif [ "$POWER" = "?" ]; then
      POWER=$VALUE
    elif [ "$VALUE" = "mobile" ]; then
      MOBILE='1'
      NOFLOAT='1'
    elif [ "$VALUE" = "nofloat" ]; then
      NOFLOAT='1'
    fi
    i=$i+1
  done
done

VALID_ONOFF_IDS="ceiling1 ceiling2 ceiling3 ceiling4 ceiling5 ceiling6"
VALID_RFONOFF_IDS="regalleinwand labortisch bluebar couchred couchwhite all ambientlights cxleds mashadecke boiler"
VALID_SEND_IDS=""

DESC_ceiling1="Decke E-Labor (SSW)"
DESC_ceiling2="Decke Leinwand (S)"
DESC_ceiling3="Decke Eingang (W)"
DESC_ceiling4="Decke Durchgang (O)"
DESC_ceiling5="Decke Auslage (N)"
DESC_ceiling6="Decke K&uuml;che (NNO)"

DESC_regalleinwand="LEDs Regal Leinwand"
DESC_bluebar="Blaue LEDs Bar"
DESC_labortisch="Labortisch"
DESC_couchred="LEDs Couch Red"
DESC_couchwhite="LEDS Couch White"
DESC_cxleds="CX Leds"
DESC_mashadecke="MaSha Decke"
DESC_ambientlights="Ambient Lichter"
DESC_boiler="Warmwasser K&uuml;che"
DESC_all="Alle Funksteckdosen"
DESC_ymhpoweron="Receiver On (off+tgl)"
DESC_ymhpoweroff="Receiver Off"
DESC_ymhpower="Receiver On/Off"
DESC_ymhvolup="VolumeUp"
DESC_ymhvoldown="VolumeDown"
DESC_ymhcd="Input CD"
DESC_ymhwdtv="Input S/PDIF Wuerfel"
DESC_ymhtuner="Input Tuner"
DESC_ymhvolmute="Mute"
DESC_ymhmenu="Menu"
DESC_ymhplus="+"
DESC_ymhminus="-"
DESC_ymhtest="Test"
DESC_ymhtimelevel="Time/Levels"
DESC_ymheffect="DSP Effect Toggle"
DESC_ymhprgup="DSP Up"
DESC_ymhprgdown="DSP Down"
DESC_ymhtunplus="Tuner +"
DESC_ymhtunminus="Tuner -"
DESC_ymhtunabcde="Tuner ABCDE"
DESC_ymhtape="Tape"
DESC_ymhvcr="VCR"
DESC_ymhextdec="ExtDec Toggle"
DESC_seep="Sleep Modus"
DESC_panicled="HAL9000 says hi"
DESC_blueled="Blue Led"
DESC_moviemode="Movie Mode"

echo "Content-type: text/html"
echo ""
echo "<html>"
echo "<head>"
echo "<title>Realraum RF and Relay Power</title>"
echo '<script type="text/javascript">'

echo 'function callbackUpdateButtons(req) {
  if (req.status != 200) {
    return;
  }
  var data = JSON.parse(req.responseText);
  for (var keyid in data) {
    on_btn = document.getElementById("onbtn_"+keyid);
    off_btn = document.getElementById("offbtn_"+keyid);
    if (on_btn && off_btn)
    {
      on_btn.className = "onbutton";
      off_btn.className = "offbutton";
      if (data[keyid])
      { on_btn.className += " enableborder"; }
      else
      { off_btn.className += " enableborder"; }
    }
  }
}'

echo 'function updateButtons(uri) {
  var req = new XMLHttpRequest;
  req.overrideMimeType("application/json");
  req.open("GET", uri, true);
  req.onload  = function() {callbackUpdateButtons(req)};
  req.setRequestHeader("googlechromefix","");
  req.send(null);
}'

echo 'function sendMultiButton( str ) {
 url = "/cgi-bin/mswitch.cgi?"+str;
  updateButtons(url);
}'

echo 'setInterval("updateButtons(\"/cgi-bin/mswitch.cgi\");", 30*1000);'
echo 'updateButtons("/cgi-bin/mswitch.cgi");'

echo '</script>'
echo '<style>'
echo 'div.switchbox {'
echo '    float:left;'
echo '    margin:2px;'
#echo '    max-width:236px;'
echo '    max-width:300px;'
echo '    font-size:10pt;'
echo '    border:1px solid black;'
#echo '    height: 32px;'
echo '    padding:0;'
echo '}'

echo 'div.switchnameleft {'
echo '    width:12em; display:inline-block; vertical-align:middle; margin-left:3px;'
echo '}'

echo 'span.alignbuttonsright {'
echo '    top:0px; float:right; display:inline-block; text-align:right; padding:0;'
echo '}'

echo 'div.switchnameright {'
echo '    width:12em; display:inline-block; vertical-align:middle; float:right; display:inline-block; margin-left:1ex; margin-right:3px; margin-top:3px; margin-bottom:3px;'
echo '}'

echo 'span.alignbuttonsleft {'
echo '    float:left; text-align:left; padding:0;'
echo '}'

echo '.onbutton {'
echo '    font-size:11pt;'
echo '    width: 40px;'
echo '    height: 32px;'
echo '    background-color: lime;'
echo '    margin: 0px;'
echo '}'

echo '.offbutton {'
echo '    font-size:11pt;'
echo '    width: 40px;'
echo '    height: 32px;'
echo '    background-color: red;'
echo '    margin: 0px;'
echo '}'

echo '.sendbutton {'
echo '    font-size:11pt;'
echo '    width: 40px;'
echo '    height: 32px;'
#echo '    background-color: grey;'
echo '    margin: 0px;'
echo '}'

echo '.enableborder {
    font-weight: bold;
    font-variant: small-caps;
    border-style: inset;'
echo '}'
echo '</style>'
echo "</head>"
echo "<body>"
#echo "<h1>Realraum rf433ctl</h1>"
#echo "<div style=\"float:left; border:1px solid black;\">"
echo "<div style=\"float:left;\">"
echo "<div style=\"float:left; border:1px solid black; margin-right:2ex; margin-bottom:2ex;\">"
for DISPID in $VALID_ONOFF_IDS; do
  NAME="$(eval echo -n \$DESC_$DISPID)"
  [ -z "$NAME" ] && NAME=$DISPID

  echo "<div class=\"switchbox\">"
  echo "<span class=\"alignbuttonsleft\">"
  if gpio_is_on $DISPID; then
  echo " <button id=\"onbtn_$DISPID\" class=\"onbutton enableborder\" onClick='sendMultiButton(\"$DISPID=1\");'>On</button>"
  echo " <button id=\"offbtn_$DISPID\" class=\"offbutton\" onClick='sendMultiButton(\"$DISPID=0\");'>Off</button>"
  else
  echo " <button id=\"onbtn_$DISPID\" class=\"onbutton\" onClick='sendMultiButton(\"$DISPID=1\");'>On</button>"
  echo " <button id=\"offbtn_$DISPID\" class=\"offbutton enableborder\" onClick='sendMultiButton(\"$DISPID=0\");'>Off</button>"
  fi
  echo "</span>"
  echo -n "<div class=\"switchnameright\">$NAME</div>"
#  echo -n "<div class=\"switchnameright\">$NAME ("
#  print_gpio_state $DISPID
#  echo ")</div>"
  echo "</div>"

  if [ "$NOFLOAT" = "1" ]; then
    echo "<br/>"
  fi
done

#Alle
echo "<div class=\"switchbox\">"
echo "<span class=\"alignbuttonsleft\">"
echo -n " <button class=\"onbutton\" onClick='sendMultiButton(\""
for DISPID in $VALID_ONOFF_IDS; do echo -n "$DISPID=1&"; done
echo "\");'>On</button>"
echo -n " <button class=\"offbutton\" onClick='sendMultiButton(\""
for DISPID in $VALID_ONOFF_IDS; do echo -n "$DISPID=0&"; done
echo "\");'>Off</button>"
echo "</span>"
echo -n "<div class=\"switchnameright\">Alle</div>"
echo "</div>"

if [ "$NOFLOAT" = "1" ]; then
  echo "<br/>"
fi

#Pattern1
echo "<div class=\"switchbox\">"
echo -n " <button class=\"sendbutton\" onClick='sendMultiButton(\"ceiling1=0&ceiling2=0&ceiling3=1&ceiling4=1&ceiling5=0&ceiling6=0\");'>[&nbsp;|&nbsp;]</button>"
echo -n " <button class=\"sendbutton\" onClick='sendMultiButton(\"ceiling1=1&ceiling2=1&ceiling3=0&ceiling4=0&ceiling5=1&ceiling6=1\");'>[|&nbsp;|]</button>"
echo -n " <button class=\"sendbutton\" onClick='sendMultiButton(\"ceiling1=1&ceiling2=0&ceiling3=0&ceiling4=0&ceiling5=1&ceiling6=1\");'>[|&nbsp;.]</button>"
echo -n " <button class=\"sendbutton\" onClick='sendMultiButton(\"ceiling1=0&ceiling2=1&ceiling3=0&ceiling4=1&ceiling5=0&ceiling6=1\");'>[***]</button>"
echo -n " <button class=\"sendbutton\" onClick='sendMultiButton(\"ceiling1=1&ceiling2=0&ceiling3=0&ceiling4=0&ceiling5=0&ceiling6=1\");'>[*&nbsp;.]</button>"
echo -n " <button class=\"sendbutton\" onClick='sendMultiButton(\"ceiling1=0&ceiling2=0&ceiling3=0&ceiling4=0&ceiling5=1&ceiling6=1\");'>[|&nbsp;&nbsp;]</button>"
echo "</div>"

if [ "$NOFLOAT" = "1" ]; then
  echo "<br/>"
fi
echo "</div>"

if [ "$MOBILE" != "1" -a -n "$VALID_SEND_IDS" ]; then

echo "<div style=\"float:left; border:1px solid black; margin-right:2ex; margin-bottom:2ex;\">"

ITEMCOUNT=0

echo "</div>"
fi
echo "</div>"


echo "<div style=\"float:left;\">"
echo "<div style=\"float:left; border:1px solid black; margin-right:2ex; margin-bottom:2ex;\">"
for DISPID in $VALID_RFONOFF_IDS; do
  NAME="$(eval echo -n \$DESC_$DISPID)"
  [ -z "$NAME" ] && NAME=$DISPID

  echo "<div class=\"switchbox\">"
  echo "<span class=\"alignbuttonsleft\">"
  echo " <button class=\"onbutton\" onClick='sendMultiButton(\"$DISPID=1\");'>On</button>"
  echo " <button class=\"offbutton\" onClick='sendMultiButton(\"$DISPID=0\");'>Off</button>"
  echo "</span>"
  echo "<div class=\"switchnameright\">$NAME</div>"
  echo "</div>"
  
  if [ "$NOFLOAT" = "1" ]; then
    echo "<br/>"
  fi

done
echo "</div>"
echo "</div>"

echo "</body>"
echo "</html>"
