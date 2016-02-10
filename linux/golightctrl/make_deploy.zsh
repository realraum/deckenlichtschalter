#!/bin/zsh
REMOTE_HOST=realraum@licht.realraum.at
REMOTE_DIR=/home/realraum/golightctrl

#ping -W 1 -c 1 $REMOTE_HOST || OPTIONS=(-e "ssh -o ProxyCommand='ssh gw.realraum.at exec nc %h %p'")
export GOOS=linux
export GOARCH=arm
export CGO_ENABLED=0
go build "$@"  && rsync ${OPTIONS[@]} -rv --delay-updates --progress ${PWD:t} config.env public --delete ${REMOTE_HOST}:${REMOTE_DIR}/  && ssh ${REMOTE_HOST/realraum@/root@} sudo /sbin/setcap 'cap_net_bind_service=+ep' ${REMOTE_DIR}/${PWD:t}&& ssh $REMOTE_HOST systemctl --user restart golightctrl.service
