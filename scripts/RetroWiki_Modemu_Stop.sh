#!/usr/bin/env bash

echo <<BANNER
 ______         __               ________ __ __     __
|   __ \.-----.|  |_.----.-----.|  |  |  |__|  |--.|__|
|      <|  -__||   _|   _|  _  ||  |  |  |  |    < |  |
|___|__||_____||____|__| |_____||________|__|__|__||__|
_______________________________________________________

      ESP8266 AT Modem Emulator for ZX-Next core
_______________________________________________________

BANNER

EXENAME="mister-modemu"
echo -n "Killing ${EXENAME}..."
if [[ $? -eq 0 ]]
then
      echo "OK"
else
      echo "FAIL"
fi

sleep 1