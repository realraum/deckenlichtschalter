// (c) Bernhard Tittelbach, 2015

package main

import "os"
import "log"

type NullWriter struct{}

func (n *NullWriter) Write(p []byte) (int, error) { return len(p), nil }

var (
	LogMain_  *log.Logger
	LogWS_    *log.Logger
	LogGPIO_  *log.Logger
	LogRF433_ *log.Logger
	LogMQTT_  *log.Logger
)

func init() {
	LogMain_ = log.New(&NullWriter{}, "", 0)
	LogWS_ = log.New(&NullWriter{}, "", 0)
	LogGPIO_ = log.New(&NullWriter{}, "", 0)
	LogRF433_ = log.New(&NullWriter{}, "", 0)
	LogMQTT_ = log.New(&NullWriter{}, "", 0)
}

func LogEnable(logtypes ...string) {
	for _, logtype := range logtypes {
		switch logtype {
		case "GPIO":
			LogGPIO_ = log.New(os.Stderr, "GPIO ", log.LstdFlags)
		case "MAIN":
			LogMain_ = log.New(os.Stderr, logtype+" ", log.LstdFlags)
		case "WS":
			LogWS_ = log.New(os.Stderr, logtype+" ", log.LstdFlags)
		case "RF433":
			LogRF433_ = log.New(os.Stderr, logtype+" ", log.LstdFlags)
		case "MQTT":
			LogMQTT_ = log.New(os.Stderr, "MQTT"+" ", log.LstdFlags)
		case "ALL":
			LogGPIO_ = log.New(os.Stderr, "GPIO ", log.LstdFlags)
			LogMain_ = log.New(os.Stderr, "MAIN"+" ", log.LstdFlags)
			LogWS_ = log.New(os.Stderr, "WS"+" ", log.LstdFlags)
			LogRF433_ = log.New(os.Stderr, "RF433"+" ", log.LstdFlags)
			LogMQTT_ = log.New(os.Stderr, "MQTT"+" ", log.LstdFlags)
		}
	}
}
