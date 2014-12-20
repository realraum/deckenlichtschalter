#!/bin/zsh

VALID_ONOFF_IDS=(ceiling1 ceiling2 ceiling3 ceiling4 ceiling5 ceiling6)
VALID_RFONOFF_IDS=(regalleinwand labortisch bluebar couchred couchwhite all lichter ambientlights cxleds mashadecke boiler)
VALID_GPIO_IDS=(4 23 18 17 22 21)
local -A IDGPIOMAP
local -A GPIOIDMAP
IDGPIOMAP[ceiling1]=4
IDGPIOMAP[ceiling2]=23
IDGPIOMAP[ceiling3]=18
IDGPIOMAP[ceiling4]=17
IDGPIOMAP[ceiling5]=22
IDGPIOMAP[ceiling6]=21
GPIOPATH=/sys/class/gpio/gpio
SAVESTATE=/var/log/licht/mswitch.state

for k v in ${(kv)IDGPIOMAP}; do
  GPIOIDMAP[$v]=$k
done

local -A GPIOS
local -A RFIDS
for QUERY in `echo $QUERY_STRING | tr '&' ' '`; do
  for VALIDID in $VALID_ONOFF_IDS; do
    if [ "$QUERY" = "$VALIDID=1" ]; then
      GPIOS[$IDGPIOMAP[$VALIDID]]=1
    elif [ "$QUERY" = "$VALIDID=0" ]; then
      GPIOS[$IDGPIOMAP[$VALIDID]]=0
    fi
  done
  for VALIDID in $VALID_RFONOFF_IDS; do
    if [ "$QUERY" = "$VALIDID=1" ]; then
      RFIDS[$VALIDID]=1
    elif [ "$QUERY" = "$VALIDID=0" ]; then
      RFIDS[$VALIDID]=0
    fi
  done
  if [ "$QUERY" = "mobile=1" ]; then
    MOBILE='1'
    NOFLOAT='1'
  elif [ "$QUERY" = "nofloat=1" ]; then
    NOFLOAT='1'
  fi
done


print_gpio_state() {
  GPIO=${IDGPIOMAP[$1]}
  GPIOVALUE=$(cat "${GPIOPATH}${GPIO}/value")
  if [[ $GPIOVALUE == "0" ]]; then
    echo -n "true"
  else
    echo -n "false"
  fi
}

print_gpio_state_10() {
  GPIO=${IDGPIOMAP[$1]}
  GPIOVALUE=$(cat "${GPIOPATH}${GPIO}/value")
  if [[ $GPIOVALUE == "0" ]]; then
    echo -n "1"
  else
    echo -n "0"
  fi
}

gpio_is_on() {
  GPIO=${IDGPIOMAP[$1]}
  GPIOVALUE=$(cat "${GPIOPATH}${GPIO}/value")
  [ "$GPIOVALUE" = "0" ]
}

echo "Content-type: text/html"
echo ""

local -a GPIOSTATES
for CHECKID in $VALID_ONOFF_IDS; do
  VAL=$GPIOS[$IDGPIOMAP[$CHECKID]]
  if [[ $VAL == 1 || $VAL == 0 ]]; then
    [[ $VAL == 1 ]] && VAL=0 || VAL=1
    echo "$VAL" > "${GPIOPATH}${IDGPIOMAP[$CHECKID]}/value"
  fi
  GPIOSTATES+=(\"${CHECKID}\":"$(print_gpio_state $CHECKID)")
  URISTATES+=("${CHECKID}=$(print_gpio_state_10 $CHECKID)")
done
for CHECKID VAL in ${(kv)RFIDS}; do
  logger $VAL $CHECKID
  if [[ $VAL == 1 || $VAL == 0 ]]; then
    /home/realraum/rf433ctl.py $VAL $CHECKID &
  fi
done
JSON_STATE="{${(j:,:)GPIOSTATES}}"
print ${(q)JSON_STATE}
if ((#GPIOS > 0)); then
  print "[$(date +%s),\"$REMOTE_ADDR\",${(q)JSON_STATE}]," >> /var/log/licht/mswitch.log
  echo -n "${(j:&:)URISTATES}">$SAVESTATE
fi

