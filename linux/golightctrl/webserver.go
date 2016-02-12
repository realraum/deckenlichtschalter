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

func goHandleSwitchCGI(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		LogWS_.Print(err)
		return
	}
	for name, _ := range actionname_map_ {
		v := r.FormValue(name)
		if len(v) == 0 {
			continue
		}
		if v == "1" || v == "on" || v == "send" {
			SwitchName(name, true)
		} else if v == "0" || v == "off" {
			SwitchName(name, false)
		}
	}
	replydata, err := json.Marshal(ConvertCeilingLightsStateTomap(GetCeilingLightsState(), 1))
	if err != nil {
		LogWS_.Print(err)
		return
	}
	w.Write(replydata)
}

func goHandleWebSocket(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		LogWS_.Println(err)
		return
	}
	LogWS_.Println("Client connected", ws.RemoteAddr())
	//TODO: send Client updates about CeilingLight states and maybe about RF Send Actions
	//go goWriteToWebSocketClient(ws, ps)

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
			switch action {
			case "1", "on", "send":
				SwitchName(name, true)
			case "0", "off":
				SwitchName(name, false)
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
	m.Get("/cgi-bin/mswitch.cgi", goHandleSwitchCGI)
	m.Get("/cgi-bin/switch.cgi", webRedirectToSwitchHTML)
	m.Get("/cgi-bin/rfswitch.cgi", webRedirectToSwitchHTML)
	m.Get("/sock", goHandleWebSocket)
	m.RunOnAddr(EnvironOrDefault("GOLIGHTCTRL_HTTP_INTERFACE", DEFAULT_GOLIGHTCTRL_HTTP_INTERFACE))
}
