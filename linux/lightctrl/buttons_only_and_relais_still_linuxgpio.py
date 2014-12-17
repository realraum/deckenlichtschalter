#!/usr/bin/python2
# -*- coding: utf-8 -*-

import os
import sys
import urllib

mswitchuri = "http://licht.realraum.at/cgi-bin/mswitch.cgi?"
bit_to_query = ["ceiling1=0","ceiling1=1",
				"ceiling2=0","ceiling2=1",
				"ceiling3=0","ceiling3=1",
				"ceiling4=0","ceiling4=1",
				"ceiling5=0","ceiling5=1",
				"ceiling6=0","ceiling6=1",
				"ceiling1=0&ceiling2=0&ceiling3=0&ceiling4=0&ceiling5=0&ceiling6=0",
				"ceiling1=0&ceiling2=0&ceiling3=1&ceiling4=1&ceiling5=0&ceiling6=0",
				"ceiling1=1&ceiling2=1&ceiling3=1&ceiling4=1&ceiling5=1&ceiling6=1"
				]

def touchURL(url):
  try:
    f = urllib.urlopen(url)
    rq_response = f.read()
    logging.debug("touchURL: url: "+url)
    #logging.debug("touchURL: Response "+rq_response)
    f.close()
    return rq_response
  except Exception, e:
    logging.error("touchURL: "+str(e))


def buttonsToUris(btns):
	for i in range(0,len(bit_to_query)):
		if btns & (1 << i) > 0:
			#touchURL(mswitchuri + bit_to_query[i])
			print(mswitchuri + bit_to_query[i])

if len(sys.argv) < 2:
	os.exit(0)

with open(sys.argv[1],"r") as ttydev:
	for line in ttydev:
		if len(line) < 3:
			continue
		buttons_pressed = (ord(line[1]) << 8) | ord(line[2])

		buttonsToUris(buttons_pressed)

