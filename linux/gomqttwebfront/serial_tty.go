// (c) Bernhard Tittelbach, 2013,2015

package main

import (
	"bufio"
	"errors"
	"fmt"
	"syscall"

	"github.com/schleibinger/sio"
)

// ---------- Serial TTY Code -------------

func openTTY(name string, speed uint) (port *sio.Port, err error) {
	switch speed {
	case 1200:
		port, err = sio.Open(name, syscall.B1200)
	case 2400:
		port, err = sio.Open(name, syscall.B2400)
	case 4800:
		port, err = sio.Open(name, syscall.B4800)
	case 9600:
		port, err = sio.Open(name, syscall.B9600)
	case 19200:
		port, err = sio.Open(name, syscall.B19200)
	case 38400:
		port, err = sio.Open(name, syscall.B38400)
	case 57600:
		port, err = sio.Open(name, syscall.B57600)
	case 115200:
		port, err = sio.Open(name, syscall.B115200)
	case 230400:
		port, err = sio.Open(name, syscall.B230400)
	default:
		err = errors.New("Unsupported Baudrate, use 0 to disable setting a baudrate")
	}
	return
}

func serialWriter(in <-chan SerialLine, serial *sio.Port) {
	for totty := range in {
		serial.Write(totty)
	}
	serial.Close()
}

func serialReader(out chan<- SerialLine, serial *sio.Port) {
	linescanner := bufio.NewScanner(serial)
	linescanner.Split(bufio.ScanLines)
	for linescanner.Scan() {
		text := linescanner.Bytes()
		if len(text) == 0 {
			continue
		}
		out <- text
	}
	if err := linescanner.Err(); err != nil {
		panic(err.Error())
	}
	panic(fmt.Sprintf("serial device '%s' has been closed/removed", serial.LocalAddr()))
}

func OpenAndHandleSerial(filename string, serspeed uint) (chan SerialLine, chan SerialLine, error) {
	serial, err := openTTY(filename, serspeed)
	if err != nil {
		return nil, nil, err
	}
	wr := make(chan SerialLine, 1)
	rd := make(chan SerialLine, 20)
	go serialWriter(wr, serial)
	go serialReader(rd, serial)
	return wr, rd, nil
}
