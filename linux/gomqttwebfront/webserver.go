// (c) Bernhard Tittelbach, 2016
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/btittelbach/pubsub"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/websocket"
	"github.com/hexadecy/nocache"
	"github.com/realraum/door_and_sensors/r3events"
)

var (
	topics_other = []string{r3events.ACT_PIPELEDS_PATTERN, r3events.ACT_YAMAHA_SEND,
		"action/GoLightCtrl/all",
		"action/GoLightCtrl/allrf",
		"action/GoLightCtrl/ambientlights",
		"action/GoLightCtrl/ymhpoweroff",
		"action/GoLightCtrl/ymhpower",
		"action/GoLightCtrl/ymhpoweron",
		"action/GoLightCtrl/ymhcd",
		"action/GoLightCtrl/ymhtuner",
		"action/GoLightCtrl/ymhtape",
		"action/GoLightCtrl/ymhwdtv",
		"action/GoLightCtrl/ymhsattv",
		"action/GoLightCtrl/ymhvcr",
		"action/GoLightCtrl/ymh7",
		"action/GoLightCtrl/ymhaux",
		"action/GoLightCtrl/ymhextdec",
		"action/GoLightCtrl/ymhtest",
		"action/GoLightCtrl/ymhtunabcde",
		"action/GoLightCtrl/ymheffect",
		"action/GoLightCtrl/ymhtunplus",
		"action/GoLightCtrl/ymhtunminus",
		"action/GoLightCtrl/ymhvolup",
		"action/GoLightCtrl/ymhvoldown",
		"action/GoLightCtrl/ymhvolmute",
		"action/GoLightCtrl/ymhmenu",
		"action/GoLightCtrl/ymhplus",
		"action/GoLightCtrl/ymhminus",
		"action/GoLightCtrl/ymhtimelevel",
		"action/GoLightCtrl/ymhprgdown",
		"action/GoLightCtrl/ymhprgup",
		"action/GoLightCtrl/ymhsleep",
		"action/GoLightCtrl/ymhp5",
		"action/GoLightCtrl/bluebar",
		"action/GoLightCtrl/couchwhite",
		"action/GoLightCtrl/couchred",
		"action/GoLightCtrl/abwasch",
		"action/GoLightCtrl/cxleds",
		"action/GoLightCtrl/spots",
		"action/GoLightCtrl/regalleinwand",
		"action/GoLightCtrl/labortisch",
		"action/GoLightCtrl/floodtesla",
		"action/GoLightCtrl/laserball",
		"action/GoLightCtrl/logo",
		"action/GoLightCtrl/boilerolga",
		"action/GoLightCtrl/fancyvortrag",
		"action/ceilingscripts/activatescript"}
	topic_fancy_ceiling_all    = "action/ceilingAll/light"
	topic_basic_ceiling_all    = "action/GoLightCtrl/basiclightAll"
	topic_oldbasic_ceiling_all = "action/GoLightCtrl/ceilingAll"
	topics_fancy_ceiling       = []string{
		"action/ceiling1/light",
		"action/ceiling2/light",
		"action/ceiling3/light",
		"action/ceiling4/light",
		"action/ceiling5/light",
		"action/ceiling6/light",
		"action/abwasch/light",
		"action/flooddoor/light",
		"action/funkbude/light",
		"action/ducttape-ledstrip/light",
	}
	topics_basic_ceiling = []string{
		"action/GoLightCtrl/basiclight1",
		"action/GoLightCtrl/basiclight2",
		"action/GoLightCtrl/basiclight3",
		"action/GoLightCtrl/basiclight4",
		"action/GoLightCtrl/basiclight5",
		"action/GoLightCtrl/basiclight6",
	}
	topics_oldbasic_ceiling = []string{
		"action/GoLightCtrl/ceiling1",
		"action/GoLightCtrl/ceiling2",
		"action/GoLightCtrl/ceiling3",
		"action/GoLightCtrl/ceiling4",
		"action/GoLightCtrl/ceiling5",
		"action/GoLightCtrl/ceiling6",
	}
	topics_sonoff_action = []string{
		"action/mashadecke/POWER",
		"action/couchred/POWER",
		"action/olgaboiler/POWER",
		"action/lothrboiler/POWER",
		"action/hallwaylight/POWER",
		"action/r2w2whiteboard/POWER",
		"action/twang/POWER",
	}
	topics_esphome_state = []string{
		"realraum/olgadecke/state",
		"realraum/subtable/state",
		"realraum/mashacompressor/state",
	}
	topics_esphome_command = []string{
		"action/olgadecke/command",
		"action/subtable/command",
		"action/mashacompressor/command",
	}
	topics_zigbee2mqtt_state = []string{
		"zigbee2mqtt/w1/OutletBlueLEDBar",
		"zigbee2mqtt/w1/OutletAuslageW1",
	}
	topics_zigbee2mqtt_action = []string{
		"zigbee2mqtt/w1/OutletBlueLEDBar/set",
		"zigbee2mqtt/w1/OutletAuslageW1/set",
	}

	ws_allowed_ctx_all = append(
		append(
			append(
				append(
					append(
						append(
							append(
								append(
									append(
										append(
											append(topics_other, topics_fancy_ceiling...),
											topics_basic_ceiling...),
										topics_oldbasic_ceiling...),
									topic_basic_ceiling_all),
								topic_fancy_ceiling_all),
							topic_oldbasic_ceiling_all),
						topics_sonoff_action...),
					topics_esphome_command...),
				topics_esphome_state...),
			topics_zigbee2mqtt_state...),
		topics_zigbee2mqtt_action...)
	ws_allowed_ctx_sendtoclientonconnect = append(append(append(append(append(topics_other, topics_fancy_ceiling...), topics_basic_ceiling...), topics_sonoff_action...), topics_esphome_state...), topics_zigbee2mqtt_state...)
)

const (
	ws_ping_period_      = time.Duration(58) * time.Second
	ws_read_timeout_     = time.Duration(70) * time.Second // must be > than ws_ping_period_
	ws_write_timeout_    = time.Duration(9) * time.Second
	ws_max_message_size_ = int64(512)
)

var wsupgrader = websocket.Upgrader{} // use default options with Origin Check

//Atomizing because we take a CeilingAll msg and split it and send on its parts Ceiling1 .. Ceiling9
func goAtomizeCeilingAll(ps_ *pubsub.PubSub, atomized_wsout_chan chan<- wsMessage) {
	shutdown_chan := ps_.SubOnce(PS_SHUTDOWN)
	msgtoall_chan := ps_.Sub(PS_WEBSOCK_ALL)
	defer ps_.Unsub(msgtoall_chan, PS_WEBSOCK_ALL)
	sendnonblockingToAtomizedWSOutChan := func(wsmsg wsMessage) {
		select {
		case atomized_wsout_chan <- wsmsg: //just pointer. should be ok to use webmsg.Data multiple times since we never change single bytes
		default:
			LogWS_.Printf("ERROR: goAtomizeCeilingAll can't write to full atomized_wsout_chan")
		}
	}
	for {
		select {
		case <-shutdown_chan:
			return
		case webmsg_i := <-msgtoall_chan:
			if webmsg, castok := webmsg_i.(wsMessage); castok {
			SWITCHCTX:
				switch webmsg.Ctx {
				case topic_fancy_ceiling_all:
					for _, tp := range topics_fancy_ceiling {
						sendnonblockingToAtomizedWSOutChan(wsMessage{Ctx: tp, Data: webmsg.Data}) //just pointer. should be ok to use webmsg.Data multiple times since we never change single bytes
					}
				case topic_basic_ceiling_all:
					for _, tp := range topics_basic_ceiling {
						sendnonblockingToAtomizedWSOutChan(wsMessage{Ctx: tp, Data: webmsg.Data}) //just pointer. should be ok to use webmsg.Data multiple times since we never change single bytes
					}
				case topic_oldbasic_ceiling_all:
					for _, tp := range topics_basic_ceiling { //convert oldbasic to new basic
						sendnonblockingToAtomizedWSOutChan(wsMessage{Ctx: tp, Data: webmsg.Data}) //just pointer. should be ok to use webmsg.Data multiple times since we never change single bytes
					}
				default:
					for idx, topicmatch := range topics_oldbasic_ceiling {
						if webmsg.Ctx == topicmatch {
							sendnonblockingToAtomizedWSOutChan(wsMessage{topics_basic_ceiling[idx], webmsg.Data})
							break SWITCHCTX
						}
					}
					sendnonblockingToAtomizedWSOutChan(webmsg)
				}
			}
		}
	}

}

// glue-code that repackages updates as json
// It is here so we can rewrite the json output format for the webserver if we want
// AND so that conversion to JSON is done only once for every connected websocket
func goJSONMarshalStuffForWebSockClientsAndRetain(getretained_chan chan JsonFuture) {
	shutdown_chan := ps_.SubOnce(PS_SHUTDOWN)
	atomized_wsout_chan := make(chan wsMessage, 400)
	retained_json_map := make(map[string][]byte, len(ws_allowed_ctx_sendtoclientonconnect))

	go goAtomizeCeilingAll(ps_, atomized_wsout_chan) //subscribes to PS_WEBSOCK_ALL and gives us possibly replaced wsMessage structs

	for {
		select {
		case <-shutdown_chan:
			return

		case webmsg := <-atomized_wsout_chan:
			LogWS_.Println("goJSONMarshalStuffForWebSockClientsAndRetain", webmsg)
			if webjson, err := json.Marshal(webmsg); err == nil {
				ps_.Pub(webjson, PS_WEBSOCK_ALL_JSON)
				retained_json_map[webmsg.Ctx] = webjson
			} else {
				LogWS_.Println(err)
			}

		case f := <-getretained_chan:
			if f.future == nil {
				continue
			}
			var reply = make(OurFutures, len(f.what))
			idx := 0
			for _, rettopic := range f.what {
				jsonbytes, inmap := retained_json_map[rettopic]
				if inmap {
					reply[idx] = jsonbytes
					idx++
				} else if f.omitempty == false {
					reply[idx] = []byte{'{', '}'}
					idx++
				}
			}
			select {
			case f.future <- reply[:idx]:
				close(f.future)
			default:
				close(f.future)
			}
		}
	}
}

// goroutine responsible for talking TO a websocket client connected to /sock
func goWriteToClient(ws *websocket.Conn) {
	shutdown_c := ps_.SubOnce(PS_SHUTDOWN)
	udpate_c := ps_.Sub(PS_WEBSOCK_ALL_JSON)
	ticker := time.NewTicker(ws_ping_period_)
	defer ps_.Unsub(udpate_c, PS_WEBSOCK_ALL_JSON)
	for {
		select {
		case <-shutdown_c:
			ws.SetWriteDeadline(time.Now().Add(ws_write_timeout_))
			ws.WriteMessage(websocket.CloseMessage, []byte{})
			return
		case jsonupdate, isopen := <-udpate_c:
			if !isopen {
				ws.SetWriteDeadline(time.Now().Add(ws_write_timeout_))
				ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			ws.SetWriteDeadline(time.Now().Add(ws_write_timeout_))
			if err := ws.WriteMessage(websocket.TextMessage, jsonupdate.([]byte)); err != nil {
				LogWS_.Println("goWriteToClient", ws.RemoteAddr(), "Error", err)
				ps_.Unsub(shutdown_c, "shutdown")
				return
			}
		case <-ticker.C:
			ws.SetWriteDeadline(time.Now().Add(ws_write_timeout_))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func SanityCheckWSFancyLight(data *wsMsgFancyLight) error {
	if strings.ContainsAny(data.Name, "/+#!?") {
		return fmt.Errorf("Invalid character in name")
	}
	if data.Setting != nil {
		if (data.Setting.R != nil && *data.Setting.R > 1000) ||
			(data.Setting.G != nil && *data.Setting.G > 1000) ||
			(data.Setting.B != nil && *data.Setting.B > 1000) ||
			(data.Setting.WW != nil && *data.Setting.WW > 1000) ||
			(data.Setting.CW != nil && *data.Setting.CW > 1000) {
			return fmt.Errorf("Luminosity not in valid range [0..1000]")
		}
		if (data.Setting.Flash != nil && len(data.Setting.Flash.Cc) > 7) || (data.Setting.Fade != nil && len(data.Setting.Fade.Cc) > 7) {
			return fmt.Errorf("Cc too long")
		}
	}
	if data.AdvSetting != nil {
		if data.AdvSetting.WIntensity != nil && (*data.AdvSetting.WIntensity > 1000 || *data.AdvSetting.WIntensity < 0) {
			return fmt.Errorf("WIntensity not in valid range [0..1000]")
		}
		if data.AdvSetting.WBalance != nil && (*data.AdvSetting.WBalance > 500 || *data.AdvSetting.WBalance < -500) {
			return fmt.Errorf("WBalance not in valid range [-500..500]")
		}
		if data.AdvSetting.HSV != nil {
			if data.AdvSetting.HSV.H > 1.0 || data.AdvSetting.HSV.H < 0.0 || data.AdvSetting.HSV.S > 1.0 || data.AdvSetting.HSV.S < 0.0 || data.AdvSetting.HSV.V > 1.0 || data.AdvSetting.HSV.V < 0.0 {
				return fmt.Errorf("HSV must be in range [0..1.0]")
			}
		}
	}
	return nil
}

func SanityCheckPipeLedPattern(data *r3events.SetPipeLEDsPattern) error {
	if len(data.Pattern) > 42 {
		return fmt.Errorf("pattern name too long")
	}
	if data.Hue != nil && (*data.Hue < -1 || *data.Hue > 0xff) {
		return fmt.Errorf("Hue value -0x01..0xFF")
	}
	if data.EffectHue != nil && (*data.EffectHue < -1 || *data.EffectHue > 0xff) {
		return fmt.Errorf("EffectHue value -0x01..0xFF")
	}
	if data.Speed != nil && (*data.Speed < 0 || *data.Speed > 0xff) {
		return fmt.Errorf("Speed value 0x00..0xFF")
	}
	if data.Brightness != nil && (*data.Brightness < 0 || *data.Brightness > 100) {
		return fmt.Errorf("Brightness value 0..100")
	}
	if data.EffectBrightness != nil && (*data.EffectBrightness < 0 || *data.EffectBrightness > 100) {
		return fmt.Errorf("EffectBrightness value 0..100")
	}
	if (data.Arg != nil && *data.Arg>>32 != 0) || data.Arg1 != nil && *data.Arg1>>32 != 0 {
		return fmt.Errorf("Args are at most 32bit values")
	}
	return nil
}

// handles requests to /cgi-bin/switch.cgi and accepts GET/POST Fields "Ctx" and "Data"
// returns json formated state of Everything
func webHandleCGICtxData(w http.ResponseWriter, r *http.Request, retained_json_chan chan JsonFuture) {
	defer func() {
		if x := recover(); x != nil {
			LogWS_.Println("webHandleCGISwitch", x)
		}
	}()
	if err := r.ParseForm(); err != nil {
		LogWS_.Print(err)
		return
	}

	ctx_a, ctx_inmap := r.Form["Ctx"]
	data_a, data_inmap := r.Form["Data"]

	if ctx_inmap && data_inmap && ctx_a != nil && data_a != nil && len(ctx_a) == 1 && len(data_a) == 1 {
		ctx := ctx_a[0]
		if stringInSlice(ctx, ws_allowed_ctx_all) {
			//TODO: sanity check json payload that goes from web to MQTT
			//TODO: then sanity check specific structs
			MQTT_sendmsg_chan_ <- MQTTOutboundMsg{topic: ctx, msg: data_a[0]}
		}
	}

	ourfuture := make(chan OurFutures, 2)
	retained_json_chan <- JsonFuture{future: ourfuture, omitempty: true, what: ws_allowed_ctx_sendtoclientonconnect}
	futures := <-ourfuture
	w.Write([]byte{'['})
	w.Write(bytes.Join(futures, []byte{','}))
	w.Write([]byte{']'})
	return
}

// handles requests to /sock WebSocket
// following ctx are handled:
// "switch": {name:..., action:...}
func webHandleWebSocket(w http.ResponseWriter, r *http.Request, retained_json_chan chan JsonFuture) {
	ws, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		LogWS_.Println(err)
		return
	}
	LogWS_.Println("Client connected", ws.RemoteAddr())

	//send client the inital known states
	ourfuture := make(chan OurFutures, 2)
	retained_json_chan <- JsonFuture{future: ourfuture, omitempty: true, what: ws_allowed_ctx_sendtoclientonconnect}
	for _, f := range <-ourfuture {
		ws.WriteMessage(websocket.TextMessage, f)
	}
	//2nd goroutine per client that handles async push info
	//e.g. sends updates about CeilingLight states and maybe about RF Send Actions
	// IMPORTANT: After this function runs, WE (THIS FUNCTION) should no longer use ws.WriteMessage(..)
	go goWriteToClient(ws)

	ws.SetReadLimit(ws_max_message_size_)
	ws.SetReadDeadline(time.Now().Add(ws_read_timeout_))
	// the PongHandler will set the read deadline for next messages if pings arrive
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(ws_read_timeout_)); return nil })

	for {
		var v wsMessage
		err := ws.ReadJSON(&v)
		if err != nil {
			if _, iswserr := err.(*websocket.CloseError); iswserr {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
					LogWS_.Printf("webHandleWebSocket Error: %v", err)
				}
				break
			} else if err, ok := err.(net.Error); ok && err.Timeout() {
				LogWS_.Printf("goChatWithClientAboutCardList Timeout: %v", err)
				break
			} else {
				LogWS_.Printf("webHandleWebSocket nonfatal Error: %v", err)
			}
		}
		LogWS_.Printf("webHandleWebSocket Gotmsg: %+v", v)
		if stringInSlice(v.Ctx, ws_allowed_ctx_all) {
			//TODO: sanity check json payload that goes from web to MQTT
			//TODO: then sanity check specific structs
			// if err = SanityCheckWSFancyLight(&data); err != nil {
			// 	LogWS_.Printf("%s Error during SanityCheckWSFancyLight: %s", v.Ctx, err)
			// 	continue
			// }
			// if err := SanityCheckPipeLedPattern(&data); err != nil {
			// 	LogWS_.Printf("%s Error during SanityCheckPipeLedPattern: %s", v.Ctx, err)
			// 	continue
			// }
			MQTT_sendmsg_chan_ <- MQTTOutboundMsg{topic: v.Ctx, msg: v.Data}
		}
	}
	LogWS_.Println("webHandleWebSocket terminating", ws.RemoteAddr())
}

func webRedirectToFallbackHTML(w http.ResponseWriter, r *http.Request) {
	LogWS_.Printf("%+v", r)
	urlStr := "//" + r.Host + "/switch.html"
	w.Header().Set("Location", urlStr)
	w.WriteHeader(302)
	if r.Method == "GET" {
		note := []byte("<a href=\"" + urlStr + "\">Changed URL</a>.\n")
		w.Write(note)
	}
}

func goRunWebserver() {
	static := nocache.NoCacheStatic(negroni.NewStatic(http.Dir("public")))
	negroni_recovery_on_panic := negroni.NewRecovery()
	negroni_recovery_on_panic.PrintStack = false
	negroni_recovery_on_panic.PanicHandlerFunc = func(x *negroni.PanicInformation) { panic(x) }
	logger := negroni.NewLogger()
	n := negroni.New(negroni_recovery_on_panic, logger, static)
	// n := negroni.Classic() // Includes some default middlewares

	retained_json_chan := make(chan JsonFuture, 20)
	go goJSONMarshalStuffForWebSockClientsAndRetain(retained_json_chan)

	mux := http.NewServeMux()
	mux.HandleFunc("/sock", func(w http.ResponseWriter, r *http.Request) { webHandleWebSocket(w, r, retained_json_chan) })
	mux.HandleFunc("/cgi-bin/fallback.cgi", func(w http.ResponseWriter, r *http.Request) { webHandleCGICtxData(w, r, retained_json_chan) })
	mux.HandleFunc("/cgi-bin/rfswitch.cgi", webRedirectToFallbackHTML)
	mux.HandleFunc("/cgi-bin/mswitch.cgi", webRedirectToFallbackHTML)
	mux.HandleFunc("/cgi-bin/fancylight.cgi", webRedirectToFallbackHTML)
	mux.HandleFunc("/cgi-bin/ledpipe.cgi", webRedirectToFallbackHTML)
	n.UseHandler(mux)
	if err := http.ListenAndServe(EnvironOrDefault("GOMQTTWEBFRONT_HTTP_INTERFACE", DEFAULT_GOMQTTWEBFRONT_HTTP_INTERFACE), n); err != nil {
		panic(err)
	}
}
