#!/bin/zsh
REMOTE_USER=licht
REMOTE_HOST=licht.realraum.at
REMOTE_DIR=/home/licht/bin
REMOTE_CONFDIR=/home/licht/.config

ping -W 1 -c 1 ${REMOTE_HOST} || { OPTIONS=(-o ProxyJump=gw.realraum.at); RSYNCOPTIONS=(-e 'ssh -o ProxyJump=gw.realraum.at')}
export GOOS=linux
export GOARCH=arm
export CGO_ENABLED=0
go build "$@"  && rsync ${RSYNCOPTIONS[@]} -rvp --delay-updates --progress --delete ${PWD:t} gonamectrl ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR} && rsync ${RSYNCOPTIONS[@]} -rvp --delay-updates --progress --delete ${PWD:t} gonamectrl.env ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_CONFDIR} && {echo "Restart Daemon? [Yn]"; read -q && ssh ${OPTIONS[@]} ${REMOTE_USER}@$REMOTE_HOST systemctl --user restart golightctrl.service; return 0}
