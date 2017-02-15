#!/usr/bin/python2
# -*- coding: utf-8 -*-

import os
import sys
import urllib
import serial

mswitchuri = "http://licht.realraum.at/cgi-bin/mswitch.cgi?"
bit_to_query = ["ceiling6=1","ceiling6=0",
				"ceiling4=1","ceiling4=0",
				"ceiling2=1&regalleinwand=1","ceiling2=0&regalleinwand=0",
				"ceiling5=1","ceiling5=0",
				"ceiling3=1","ceiling3=0",
				"ceiling1=1","ceiling1=0",
				"ceiling1=0&ceiling2=0&regalleinwand=0&ceiling3=0&ceiling4=0&ceiling5=0&ceiling6=0&cxleds=0",
				"ceiling1=0&ceiling2=0&regalleinwand=0&ceiling3=1&ceiling4=1&ceiling5=0&ceiling6=0&cxleds=1",
				"ceiling1=1&ceiling2=1&regalleinwand=1&ceiling3=1&ceiling4=1&ceiling5=1&ceiling6=1&cxleds=1"
				]

def touchURL(url):
  try:
    f = urllib.urlopen(url)
    rq_response = f.read()
    #logging.debug("touchURL: url: "+url)
    #logging.debug("touchURL: Response "+rq_response)
    f.close()
    return rq_response
  except Exception, e:
    logging.error("touchURL: "+str(e))


def buttonsToUris(btns):
	for i in range(0,len(bit_to_query)):
		if btns & (1 << i) > 0:
			touchURL(mswitchuri + bit_to_query[i])
			#print(mswitchuri + bit_to_query[i])

if len(sys.argv) < 2:
	sys.exit(0)

ttydev = serial.Serial(sys.argv[1], baudrate=115200)
ttydev.write("0");
while True:
	line = ttydev.read(4)
	#print  map(hex,map(ord,line))
	buttons_pressed = (ord(line[1]) << 8) | ord(line[2])
	buttonsToUris(buttons_pressed)

