// (c) Bernhard Tittelbach, 2016
package main

import (
	"encoding/json"
	"net/http"

	"github.com/codegangsta/martini"
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
	for name, _ := range rfcode_map {
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

//TODO: upgrade to WebSocket

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
	m.RunOnAddr(EnvironOrDefault("GOLIGHTCTRL_HTTP_INTERFACE", DEFAULT_GOLIGHTCTRL_HTTP_INTERFACE))
}
