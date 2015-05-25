#!/usr/bin/python2
# -*- coding: utf-8 -*-

import sys,os

ttydev="/dev/ttyACM0"

rfcodes = {
  "regalleinwand":{"on":"\xa2\xa0\xa8","off":"\xa2\xa0\x28"}, #white remote B 1 
  "bluebar":{"on":"\xa8\xa0\xa8","off":"\xa8\xa0\x28"}, #white remote C 1
  "labortisch":{"on":"\xa2\xa2\xaa","off":"\xa2\xa2\x2a"},
  "couchred":{"on":"\x8a\xa0\x8a","off":"\x8a\xa0\x2a"}, #pollin 00101 a
  "couchwhite":{"on":"\x8a\xa8\x88","off":"\x8a\xa8\x28"}, #pollin 00101 d
  "cxleds":{"on":"\x8a\x88\x8a","off":"\x8a\x88\x2a"}, #pollin 00101 b
  "mashadecke":{"on":"\x8a\x28\x8a","off":"\x8a\x28\x2a"}, #pollin 00101 c
  "boiler":{"on":"\xa0\xa2\xa8","off":"\xa0\xa2\x28"}, #white remote A 2
  "spots":{"on:":"\x00\xaa\x88","off":"\x00\xaa\x28"}, #polling 11110 d
  "lichtpi":{"on:":"\x00\xa2\x8a","off":"\x00\xa2\x2a"},
  "abwasch":{"on:":"\xaa\xa2\xa8","off":"\xaa\xa2\x28"}  #alte jk16 decke vorne
}
#  "jk16decke":{"on:":"\xaa\xa0\xa8","off":"\xaa\xa0\x28"},

multinames = {
  "ambientlights":["bluebar","couchred","couchwhite","regalleinwand","cxleds","abwasch"],
  "all":list(set(rfcodes.keys()) - set(["lichtpi"]))
}

namestoswitch=[]

if len(sys.argv) > 2:
  if sys.argv[1] == "1":
    sys.argv[1] = "on"
  elif sys.argv[1] == "0":
    sys.argv[1] = "off"
  if sys.argv[2] in multinames:
    namestoswitch = multinames[sys.argv[2]]
  elif sys.argv[2] in rfcodes:
    namestoswitch = [sys.argv[2]]
  else:
    sys.exit(1)
  with open(ttydev,"w") as ttyfh:
    for rfname in namestoswitch:
      if sys.argv[1] in rfcodes[rfname]:
        code = rfcodes[rfname][sys.argv[1]]
        ttyfh.write(">"+code)

