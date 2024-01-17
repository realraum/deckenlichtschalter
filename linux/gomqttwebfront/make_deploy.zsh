#!/bin/zsh
REMOTE_USER=licht
REMOTE_HOST=licht.realraum.at
REMOTE_DIR=/home/licht/gomqttwebfront

ping -W 1 -c 1 ${REMOTE_HOST} || { OPTIONS=(-o ProxyCommand='ssh gw.realraum.at exec nc '$REMOTE_HOST' 22000'); RSYNCOPTIONS=(-e 'ssh -o ProxyCommand="ssh gw.realraum.at exec nc '$REMOTE_HOST' 22000"')}
export GOOS=linux
export GOARCH=arm
export GOARM=6
export CGO_ENABLED=0
#go build "$@"  && rsync ${RSYNCOPTIONS[@]} -rvp --delay-updates --progress --delete ${PWD:t} config.env public ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/  && ssh ${OPTIONS[@]} ${REMOTE_USER}@${REMOTE_HOST} sudo /sbin/setcap 'cap_net_bind_service=+ep' ${REMOTE_DIR}/${PWD:t} && {echo "Restart Daemon? [Yn]"; read -q && ssh ${OPTIONS[@]} ${REMOTE_USER}@$REMOTE_HOST systemctl --user restart gomqttwebfront.service; return 0}
go build "$@" -ldflags "-s"  && rsync ${RSYNCOPTIONS[@]} -rvp --delay-updates --progress --delete ${PWD:t} config.env public ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/ && echo "please run systemctl restart gomqttwebfront.service as root on licht"
