// (c) Bernhard Tittelbach, 2016
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/codegangsta/martini"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"github.com/realraum/door_and_sensors/r3events"
)

const (
	ws_ctx_ceilinglights = "ceilinglights"
	ws_ctx_fancylight    = "FancyLight"
	ws_ctx_ledpattern    = "SetPipeLEDsPattern"
	ws_ctx_switch        = "LightCtrlActionOnName"
	ws_ctx_button_used   = "wbp"
)

const (
	recall_sturdyceiling_cgi RetainRecallID = iota
	recall_sturdyceiling_ws  RetainRecallID = iota
	recall_fancyceiling_ws   RetainRecallID = iota
	recall_ledpipe_ws        RetainRecallID = iota
)

const (
	ws_ping_period_      = time.Duration(58) * time.Second
	ws_read_timeout_     = time.Duration(70) * time.Second // must be > than ws_ping_period_
	ws_write_timeout_    = time.Duration(9) * time.Second
	ws_max_message_size_ = int64(512)
)

var wsupgrader = websocket.Upgrader{} // use default options with Origin Check

// handles requests to /cgi-bin/switch.cgi?<switch1>=<state>&<switch2>=<state>&...
// returns json formated state of Ceiling Lights
func webHandleCGISwitch(w http.ResponseWriter, r *http.Request, retained_lightstate_chan chan JsonFuture) {
	defer func() {
		if x := recover(); x != nil {
			LogWS_.Println("webHandleCGISwitch", x)
		}
	}()
	if err := r.ParseForm(); err != nil {
		LogWS_.Print(err)
		return
	}
	ourfuture := make(chan []byte, 2)
	retained_lightstate_chan <- JsonFuture{future: ourfuture}
	for name, _ := range actionname_map_ {
		v := r.FormValue(name)
		if len(v) == 0 {
			continue
		}
		switch_name_chan_ <- r3events.LightCtrlActionOnName{name, v}
	}
	w.Write(<-ourfuture)
	return
}

func SanityCheckWSFancyLight(data *wsMsgFancyLight) error {
	if strings.ContainsAny(data.Name, "/+#!?") {
		return fmt.Errorf("Invalid character in name")
	}
	if (data.Setting.R != nil && *data.Setting.R > 1000) ||
		(data.Setting.G != nil && *data.Setting.G > 1000) ||
		(data.Setting.B != nil && *data.Setting.B > 1000) ||
		(data.Setting.WW != nil && *data.Setting.WW > 1000) ||
		(data.Setting.CW != nil && *data.Setting.CW > 1000) {
		return fmt.Errorf("Luminosity not in valid range [0..1000]")
	}
	if len(data.Setting.Flash.Cc) > 7 || len(data.Setting.Fade.Cc) > 7 {
		return fmt.Errorf("Cc too long")
	}
	return nil
}

// handles requests to /cgi-bin/switch.cgi?<switch1>=<state>&<switch2>=<state>&...
// returns json formated state of Ceiling Lights
func webHandleCGIFancyLight(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if x := recover(); x != nil {
			LogWS_.Println("webHandleCGIFancyLight", x)
		}
	}()
	if err := r.ParseForm(); err != nil {
		LogWS_.Print(err)
		return
	}
	var data wsMsgFancyLight
	data.Name = r.FormValue("name")
	if len(data.Name) == 0 {
		w.Write([]byte("err"))
		return
	}
	if err := json.Unmarshal([]byte(r.FormValue("setting")), &data.Setting); err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	if err := SanityCheckWSFancyLight(&data); err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	MQTT_fancylight_chan_ <- &data
	w.Write([]byte("ok"))
	return
}

// cache Lights Change Update for webHandleCGISwitch()
func goRetainCeilingLightsJSONForLater(retained_lightstate_chan chan JsonFuture) {
	shutdown_chan := ps_.SubOnce(PS_SHUTDOWN)
	lights_changed_chan := ps_.Sub(PS_LIGHTS_CHANGED)
	defer ps_.Unsub(lights_changed_chan, PS_LIGHTS_CHANGED)
	fancylight_update_chan := ps_.Sub(PS_FANCYLIGHT_CHANGED)
	defer ps_.Unsub(fancylight_update_chan, PS_FANCYLIGHT_CHANGED)
	ledpipe_update_chan := ps_.Sub(PS_LEDPIPE_CHANGED)
	defer ps_.Unsub(ledpipe_update_chan, PS_LEDPIPE_CHANGED)

	//TODO FIX ME

	var cached_switchcgireply_json []byte
	var cached_websocketreply_json []byte
	var cached_ledpipe_json []byte
	var cached_fancylight_json map[string][]byte //one json []byte for each clientID.. fancy1, fancy2, etc
	var err error
	updateCache := func(lsm CeilingLightStateMap) {
		cached_switchcgireply_json, err = json.Marshal(lsm)
		if err != nil {
			LogWS_.Print(err)
			cached_switchcgireply_json = []byte("{}")
		}
		cached_websocketreply_json, err = json.Marshal(wsMessageOut{Ctx: ws_ctx_ceilinglights, Data: lsm})
		if err != nil {
			LogWS_.Print(err)
			cached_switchcgireply_json = []byte("{}")
		}
	}
	for {
		select {
		case <-shutdown_chan:
			return
		case lc := <-lights_changed_chan:
			//prepare and retain json for webHandleCGISwitch()
			updateCache(lc.(CeilingLightStateMap))
			//also send update to all Websocket Clients
			ps_.PubNonBlocking(cached_websocketreply_json, PS_WEBSOCK_ALL_JSON)
		case f := <-retained_lightstate_chan:
			if f.future == nil {
				continue
			}
			//maybe last broadcast was shit or never happened.. get info ourselves
			if cached_switchcgireply_json == nil || len(cached_switchcgireply_json) == 0 || cached_websocketreply_json == nil || len(cached_websocketreply_json) == 0 {
				updateCache(ConvertCeilingLightsStateTomap(GetCeilingLightsState(), 1))

			}
			var reply []byte
			if f.wxFormat {
				reply = cached_websocketreply_json
			} else {
				reply = cached_switchcgireply_json
			}
			select {
			case f.future <- reply:
				close(f.future)
			default:
				close(f.future)
			}
		}
	}
}

// glue-code that repackages updates as json
// It is here so we can rewrite the json output format for the webserver if we want
// AND so that conversion to JSON is done only once for every connected websocket
func goJSONMarshalStuffForWebSockClients() {
	shutdown_chan := ps_.SubOnce(PS_SHUTDOWN)
	msgtoall_chan := ps_.Sub(PS_WEBSOCK_ALL)
	//lights_changed_chan := ps_.Sub(PS_LIGHTS_CHANGED)
	button_used_chan := ps_.Sub(PS_IRRF433_CHANGED)
	fancylight_update_chan := ps_.Sub(PS_FANCYLIGHT_CHANGED)
	ledpipe_update_chan := ps_.Sub(PS_LEDPIPE_CHANGED)
	defer ps_.Unsub(msgtoall_chan, PS_WEBSOCK_ALL)
	//defer ps_.Unsub(lights_changed_chan, PS_LIGHTS_CHANGED)
	defer ps_.Unsub(button_used_chan, PS_IRRF433_CHANGED)
	defer ps_.Unsub(fancylight_update_chan, PS_FANCYLIGHT_CHANGED)
	defer ps_.Unsub(ledpipe_update_chan, PS_LEDPIPE_CHANGED)

	//TODO FIX ME

	for {
		msg := wsMessageOut{}
		select {
		case <-shutdown_chan:
			return
		case lu := <-msgtoall_chan:
			msg.Data = lu
			msg.Ctx = "some_other_message_ctx_example"
			//		case lc := <-lights_changed_chan:
			//			msg.Ctx = "ceilinglights"
			//			msg.Data = lc
		case bu := <-button_used_chan:
			msg.Ctx = ws_ctx_button_used //web button pressed
			msg.Data = bu
		}
		if len(msg.Ctx) == 0 {
			continue
		}
		if jsonbytes, err := json.Marshal(msg); err == nil {
			ps_.PubNonBlocking(jsonbytes, PS_WEBSOCK_ALL_JSON)
		} else {
			LogWS_.Println(err)
		}
	}
}

// goroutine responsible for talking TO a websocket client connected to /sock
func goWriteToClientAboutLightStateChanges(ws *websocket.Conn) {
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
				LogWS_.Println("goWriteToClientAboutLightStateChanges", ws.RemoteAddr(), "Error", err)
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

// handles requests to /sock WebSocket
// following ctx are handled:
// "switch": {name:..., action:...}
func webHandleWebSocket(w http.ResponseWriter, r *http.Request, retained_lightstate_chan chan JsonFuture) {
	ws, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		LogWS_.Println(err)
		return
	}
	LogWS_.Println("Client connected", ws.RemoteAddr())

	//send client the inital CeilingLightsState
	ourfuture := make(chan []byte, 2)
	retained_lightstate_chan <- JsonFuture{future: ourfuture, wxFormat: true}
	ws.WriteMessage(websocket.TextMessage, <-ourfuture)

	//2nd goroutine per client that handles async push info
	//e.g. sends updates about CeilingLight states and maybe about RF Send Actions
	// IMPORTANT: After this function runs, WE (THIS FUNCTION) should no longer use ws.WriteMessage(..)
	go goWriteToClientAboutLightStateChanges(ws)

	ws.SetReadLimit(ws_max_message_size_)
	ws.SetReadDeadline(time.Now().Add(ws_read_timeout_))
	// the PongHandler will set the read deadline for next messages if pings arrive
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(ws_read_timeout_)); return nil })

	for {
		var v wsMessage
		err := ws.ReadJSON(&v)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				LogWS_.Printf("webHandleWebSocket Error: %v", err)
			}
			break
		}

		switch v.Ctx {
		case ws_ctx_switch:
			var data r3events.LightCtrlActionOnName
			err = mapstructure.Decode(v.Data, &data)
			if err != nil {
				LogWS_.Printf("%s Data decode error: %s", v.Ctx, err)
				continue
			}
			switch_name_chan_ <- data
		case ws_ctx_fancylight:
			var data wsMsgFancyLight
			err = mapstructure.Decode(v.Data, &data)
			if err != nil {
				LogWS_.Printf("%s Data decode error: %s", v.Ctx, err)
				continue
			}
			if err = SanityCheckWSFancyLight(&data); err != nil {
				LogWS_.Printf("%s Error during SanityCheckWSFancyLight: %s", v.Ctx, err)
				continue
			}
			MQTT_fancylight_chan_ <- &data
		}
	}
}

func webRedirectToSwitchHTML(w http.ResponseWriter, r *http.Request) {
	LogWS_.Printf("%+v", r)
	urlStr := "//" + r.Host + "/switch.html"
	w.Header().Set("Location", urlStr)
	w.WriteHeader(302)
	if r.Method == "GET" {
		note := []byte("<a href=\"" + urlStr + "\">Changed URL</a>.\n")
		w.Write(note)
	}
}

func goRunMartini() {
	m := martini.Classic()
	//m.Use(martini.Static("/var/lib/cloud9/static/"))
	/*	m.Get("/sock", func(w http.ResponseWriter, r *http.Request) {
		goTalkWithClient(w, r, ps)
	})*/
	retained_lightstate_chan := make(chan JsonFuture, 20)
	go goRetainCeilingLightsJSONForLater(retained_lightstate_chan)
	go goJSONMarshalStuffForWebSockClients()

	m.Get("/cgi-bin/mswitch.cgi", func(w http.ResponseWriter, r *http.Request) { webHandleCGISwitch(w, r, retained_lightstate_chan) })
	m.Get("/sock", func(w http.ResponseWriter, r *http.Request) { webHandleWebSocket(w, r, retained_lightstate_chan) })
	m.Get("/cgi-bin/switch.cgi", webRedirectToSwitchHTML)
	m.Get("/cgi-bin/rfswitch.cgi", webRedirectToSwitchHTML)
	m.Get("/cgi-bin/fancylight.cgi", webHandleCGIFancyLight)
	m.RunOnAddr(EnvironOrDefault("GOLIGHTCTRL_HTTP_INTERFACE", DEFAULT_GOLIGHTCTRL_HTTP_INTERFACE))
}
