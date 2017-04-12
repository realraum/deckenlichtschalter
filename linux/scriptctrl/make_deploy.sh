#!/bin/zsh
cd "${0:h}"
rsync -vr --delay-updates --delete --include "**/*.py" --exclude "__pycache__" . licht@licht.realraum.at:scriptctrl/
