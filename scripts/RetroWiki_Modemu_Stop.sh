#!/usr/bin/env bash

echo "
 ______         __               ________ __ __     __
|   __ \.-----.|  |_.----.-----.|  |  |  |__|  |--.|__|
|      <|  -__||   _|   _|  _  ||  |  |  |  |    < |  |
|___|__||_____||____|__| |_____||________|__|__|__||__|
_______________________________________________________

      ESP8266 AT Modem Emulator for ZX-Next core
_______________________________________________________

"

EXENAME="mister-modemu"

pidof $EXENAME > /dev/null
if [[ $? -ne 0 ]]
then
      echo "$EXENAME its not running"
      exit 1
fi

echo -n "Killing ${EXENAME}..."
killall $EXENAME
if [[ $? -eq 0 ]]
then
      echo "OK"
else
      echo "FAIL"
fi
