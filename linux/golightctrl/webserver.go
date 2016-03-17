// (c) Bernhard Tittelbach, 2016
package main

import (
	"encoding/json"
	"net/http"

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
	future chan []byte
}

// handles requests to /cgi-bin/switch.cgi?<switch1>=<state>&<switch2>=<state>&...
// returns json formated state of Ceiling Lights
func webHandleSwitchCGI(w http.ResponseWriter, r *http.Request, retained_lightstate_chan chan JsonFuture) {
	if err := r.ParseForm(); err != nil {
		LogWS_.Print(err)
		return
	}
	ourfuture := make(chan []byte, 1)
	retained_lightstate_chan <- JsonFuture{ourfuture}
	for name, _ := range actionname_map_ {
		v := r.FormValue(name)
		if len(v) == 0 {
			continue
		}
		var err error
		if v == "1" || v == "on" || v == "send" {
			err = SwitchName(name, true)
		} else if v == "0" || v == "off" {
			err = SwitchName(name, false)
		}
		if err != nil {
			LogRF433_.Println(err)
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
	var err error
	for {
		select {
		case <-shutdown_chan:
			return
		case lc := <-lights_changed_chan:
			//prepare and retain json for webHandleSwitchCGI()
			cached_switchcgireply_json, err = json.Marshal(lc)
			if err != nil {
				LogWS_.Print(err)
				cached_switchcgireply_json = nil
			}

		case f := <-retained_lightstate_chan:
			if f.future == nil {
				continue
			}
			//maybe last broadcast was shit or never happened.. get info ourselves
			if cached_switchcgireply_json == nil || len(cached_switchcgireply_json) == 0 {
				cached_switchcgireply_json, err = json.Marshal(ConvertCeilingLightsStateTomap(GetCeilingLightsState(), 1))
				if err != nil {
					LogWS_.Print(err)
					cached_switchcgireply_json = []byte("{}")
				}
			}
			//we can't use non blocking send here, in case webHandleSwitchCGI does not receive yet
			f.future <- cached_switchcgireply_json
		}
	}
}

// glue-code that repackages updates as json
// It is here so we can rewrite the json output format for the webserver if we want
// AND so that conversion to JSON is done only once for every connected websocket
func goJSONMarshalStuffForWebSockClients() {
	shutdown_chan := ps_.SubOnce(PS_SHUTDOWN)
	msgtoall_chan := ps_.Sub(PS_WEBSOCK_ALL)
	lights_changed_chan := ps_.Sub(PS_LIGHTS_CHANGED)
	button_used_chan := ps_.Sub(PS_IRRF433_CHANGED)
	defer ps_.Unsub(msgtoall_chan, PS_WEBSOCK_ALL)
	defer ps_.Unsub(lights_changed_chan, PS_LIGHTS_CHANGED)
	defer ps_.Unsub(button_used_chan, PS_IRRF433_CHANGED)

	for {
		msg := wsMessageOut{}
		select {
		case <-shutdown_chan:
			return
		case lu := <-msgtoall_chan:
			msg.Data = lu
			msg.Ctx = "some_other_message_ctx_example"
		case lc := <-lights_changed_chan:
			msg.Ctx = "ceilinglights"
			msg.Data = lc
		case bu := <-button_used_chan:
			msg.Ctx = "wbp" //web button pressed
			msg.Data = bu
		}
		if len(msg.Ctx) == 0 {
			continue
		}
		if jsonbytes, err := json.Marshal(msg); err == nil {
			ps_.Pub(jsonbytes, PS_WEBSOCK_ALL_JSON)
		} else {
			LogWS_.Println(err)
		}
	}
}

// goroutine responsible for talking TO a websocket client connected to /sock
func goWriteToClientAboutLightStateChanges(ws *websocket.Conn) {
	shutdown_c := ps_.SubOnce(PS_SHUTDOWN)
	udpate_c := ps_.Sub(PS_WEBSOCK_ALL_JSON)
	defer ps_.Unsub(udpate_c, PS_WEBSOCK_ALL_JSON)
	for {
		select {
		case <-shutdown_c:
			LogWS_.Println("goWriteToClientAboutLightStateChanges", ws.RemoteAddr(), "Shutdown")
			return
		case jsonupdate := <-udpate_c:
			LogWS_.Printf("goWriteToClientAboutLightStateChanges %s: %s\n", ws.RemoteAddr(), jsonupdate)
			if err := ws.WriteMessage(websocket.TextMessage, jsonupdate.([]byte)); err != nil {
				LogWS_.Println("goWriteToClientAboutLightStateChanges", ws.RemoteAddr(), "Error", err)
				ps_.Unsub(shutdown_c, "shutdown")
				return
			}
		}
	}
}

// handles requests to /sock WebSocket
// following ctx are handled:
// "switch": {name:..., action:...}
func webHandleWebSocket(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		LogWS_.Println(err)
		return
	}
	LogWS_.Println("Client connected", ws.RemoteAddr())

	//2nd goroutine per client that handles async push info
	//e.g. sends updates about CeilingLight states and maybe about RF Send Actions
	go goWriteToClientAboutLightStateChanges(ws)

	// logged_in := false
	for {
		messageType, msg, err := ws.ReadMessage()
		if err != nil {
			LogWS_.Println("Client disconnected", ws.RemoteAddr(), err)
			return
		}
		if messageType != websocket.TextMessage {
			continue
		}
		LogWS_.Printf("Got message from %s: %s\n", ws.RemoteAddr(), msg)
		var v wsMessage
		if err := json.Unmarshal(msg, &v); err != nil {
			LogWS_.Println("Could not parse message from client ", ws.RemoteAddr(), err, msg)
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
			var err error
			switch action {
			case "1", "on", "send":
				err = SwitchName(name, true)
			case "0", "off":
				err = SwitchName(name, false)
			}
			if err != nil {
				LogRF433_.Println(err)
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
	m.Get("/cgi-bin/switch.cgi", webRedirectToSwitchHTML)
	m.Get("/cgi-bin/rfswitch.cgi", webRedirectToSwitchHTML)
	m.Get("/sock", webHandleWebSocket)
	m.RunOnAddr(EnvironOrDefault("GOLIGHTCTRL_HTTP_INTERFACE", DEFAULT_GOLIGHTCTRL_HTTP_INTERFACE))
}
