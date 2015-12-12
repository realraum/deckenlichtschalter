#!/bin/sh
sudo pip3 install pyephem paho-mqtt
sudo aptitude install x11-xserver-utils lighttpd zsh python3-pip python-serial iproute git git-man fonts-freefont-ttf ttf-dejavu-core tcpdump python3-setuptools python-numpy midori vim tmux htop xinit xserver-xorg-video-fbdev pciutil xserver-xorg xfonts-100dpi xfonts-75dpi xfonts-scalable
mkdir /var/log/licht
chown www-data:www-data -R /var/log/licht
sed 'srealraum/:pi/:realraum/' -i /etc/group
adduser realraum
adduser pi audio video input users
adduser www-data gpio dialout
