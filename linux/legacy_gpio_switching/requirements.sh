#!/bin/sh
sudo aptitude install x11-xserver-utils lighttpd zsh python3-pip python-serial iproute git git-man fonts-freefont-ttf ttf-dejavu-core tcpdump python3-setuptools python-numpy midori vim tmux htop xinit xserver-xorg-video-fbdev pciutil xserver-xorg xfonts-100dpi xfonts-75dpi xfonts-scalable python3-pyro4 pyro4
sudo pip3 install pyephem paho-mqtt
sudo aptitude install python3-serpent || { wget http://ftp.at.debian.org/debian/pool/main/s/serpent/python3-serpent_1.8-1_all.deb -O /tmp/python3-serpent_1.8-1_all.deb && dpkg -i /tmp/python3-serpent_1.8-1_all.deb && rm /tmp/python3-serpent_1.8-1_all.deb }
mkdir /var/log/licht
chown www-data:www-data -R /var/log/licht
sed 's/:pi/:licht/' -i /etc/group
adduser licht
adduser licht audio video input users gpio
adduser www-data gpio dialout
adduser realraum
