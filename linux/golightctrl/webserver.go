// (c) Bernhard Tittelbach, 2016
package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/codegangsta/martini"
	"github.com/gorilla/websocket"
)

type wsMessage struct {
	Ctx  string                 `json:"ctx"`
	Data map[string]interface{} `json:"data"`
}

type wsMessageOut struct {
	Ctx  string      `json:"ctx"`
	Data interface{} `json:"data"`
}

type jsonButtonUsed struct {
	Name string `json:"name"`
}

type JsonFuture struct {
	future   chan []byte
	wxFormat bool
}

const (
	ws_ping_period_      = time.Duration(58) * time.Second
	ws_read_timeout_     = time.Duration(70) * time.Second // must be > than ws_ping_period_
	ws_write_timeout_    = time.Duration(9) * time.Second
	ws_max_message_size_ = int64(512)
)

var wsupgrader = websocket.Upgrader{} // use default options with Origin Check

// handles requests to /cgi-bin/switch.cgi?<switch1>=<state>&<switch2>=<state>&...
// returns json formated state of Ceiling Lights
func webHandleSwitchCGI(w http.ResponseWriter, r *http.Request, retained_lightstate_chan chan JsonFuture) {
	defer func() {
		if x := recover(); x != nil {
			LogWS_.Println("webHandleSwitchCGI", x)
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
		if v == "1" || v == "on" || v == "send" {
			SwitchNameAsync(name, true)
		} else if v == "0" || v == "off" {
			SwitchNameAsync(name, false)
		}
	}
	w.Write(<-ourfuture)
	return
}

// cache Lights Change Update for webHandleSwitchCGI()
func goRetainCeilingLightsJSONForLater(retained_lightstate_chan chan JsonFuture) {
	shutdown_chan := ps_.SubOnce(PS_SHUTDOWN)
	lights_changed_chan := ps_.Sub(PS_LIGHTS_CHANGED)
	defer ps_.Unsub(lights_changed_chan, PS_LIGHTS_CHANGED)

	var cached_switchcgireply_json []byte
	var cached_websocketreply_json []byte
	var err error
	updateCache := func(lsm CeilingLightStateMap) {
		cached_switchcgireply_json, err = json.Marshal(lsm)
		if err != nil {
			LogWS_.Print(err)
			cached_switchcgireply_json = []byte("{}")
		}
		cached_websocketreply_json, err = json.Marshal(wsMessageOut{Ctx: "ceilinglights", Data: lsm})
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
			//prepare and retain json for webHandleSwitchCGI()
			updateCache(lc.(CeilingLightStateMap))
			//also send update to all Websocket Clients
			ps_.Pub2(false, cached_websocketreply_json, PS_WEBSOCK_ALL_JSON)
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
	defer ps_.Unsub(msgtoall_chan, PS_WEBSOCK_ALL)
	//defer ps_.Unsub(lights_changed_chan, PS_LIGHTS_CHANGED)
	defer ps_.Unsub(button_used_chan, PS_IRRF433_CHANGED)

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
			msg.Ctx = "wbp" //web button pressed
			msg.Data = bu
		}
		if len(msg.Ctx) == 0 {
			continue
		}
		if jsonbytes, err := json.Marshal(msg); err == nil {
			ps_.Pub2(false, jsonbytes, PS_WEBSOCK_ALL_JSON)
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
		case "switch":
			switchname, inmap := v.Data["name"]
			if inmap == false {
				LogWS_.Print("open/close ctx without name")
				continue
			}
			switchaction, inmap := v.Data["action"]
			if inmap == false {
				LogWS_.Print("open/close ctx without action")
				continue
			}
			name, typeok := switchname.(string)
			if typeok == false {
				LogWS_.Print("name not a string")
				continue
			}
			action, typeok := switchaction.(string)
			if typeok == false {
				LogWS_.Print("action not a string")
				continue
			}
			switch action {
			case "1", "on", "send":
				SwitchNameAsync(name, true)
			case "0", "off":
				SwitchNameAsync(name, false)
			}
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

	m.Get("/cgi-bin/mswitch.cgi", func(w http.ResponseWriter, r *http.Request) { webHandleSwitchCGI(w, r, retained_lightstate_chan) })
	m.Get("/sock", func(w http.ResponseWriter, r *http.Request) { webHandleWebSocket(w, r, retained_lightstate_chan) })
	m.Get("/cgi-bin/switch.cgi", webRedirectToSwitchHTML)
	m.Get("/cgi-bin/rfswitch.cgi", webRedirectToSwitchHTML)
	m.RunOnAddr(EnvironOrDefault("GOLIGHTCTRL_HTTP_INTERFACE", DEFAULT_GOLIGHTCTRL_HTTP_INTERFACE))
}
