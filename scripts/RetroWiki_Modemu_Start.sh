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

BINDIR="/media/fat/retrowiki-bin"
EXENAME="mister-modemu"
GHAPI="https://api.github.com/repos/jorgefuertes/mister-modemu"

pidof $EXENAME &> /dev/null
if [[ $? -eq 0 ]]
then
      echo "mister-modemu is already running!"
      echo "Please stop it before"
      exit 1
fi

echo "Making retrowiki-bin dir"
mkdir -p $BINDIR
echo -n "Local release..."
if [[ -f $BINDIR/$EXENAME ]]
then
      LOCAL=$($BINDIR/$EXENAME -v)
else
      LOCAL="none"
fi
echo $LOCAL
echo -n "Latest release..."
LATEST=$(curl -ksL $GHAPI/releases|jq -r '.[0].tag_name')
echo $LATEST

if [[ $LOCAL != $LATEST ]]
then
      echo "Downloading new release ${LATEST}"
      LATEST_URL=$(curl -ksL $GHAPI/releases|jq -r '.[0].assets[0].browser_download_url')
      curl -kL --progress -o $BINDIR/modemu-latest.gz $LATEST_URL
      if [[ $? -eq 0 ]]
      then
            echo "Decrunching ${LATEST}..."
            pushd $BINDIR &> /dev/null
            gunzip modemu-latest.gz
            mv modemu-latest mister-modemu
            popd &> /dev/null
      else
            echo "Error downloading ${LATEST}"
      fi
else
      echo "Up to date"
fi

echo "Launching modemu"
$BINDIR/mister-modemu &
