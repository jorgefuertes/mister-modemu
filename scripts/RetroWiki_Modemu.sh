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

BINDIR="/media/fat/retrowiki-bin"

echo "Making retrowiki-bin dir"
mkdir -p $BINDIR
echo "Checking for latest release"
LATEST_URI=$(curl -s https://github.com/jorgefuertes/mister-modemu/releases/latest)
echo "Latest is: ${LATEST}"
