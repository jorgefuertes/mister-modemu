#!/usr/bin/env bash

GHSCRIPTS="https://raw.githubusercontent.com/jorgefuertes/mister-modemu/master/scripts/"
BINDIR="/media/fat/scripts"
STARTSCRIPT="RetroWiki_Modemu_Start.sh"
STOPSCRIPT="RetroWiki_Modemu_Stop.sh"

check() {
	if [[ $1 -eq 0 ]]
	then
		echo "OK"
	else
		echo "FAILED"
		exit 1
	fi
}

echo -n "Start script..."
curl -kLs -o $BINDIR/$STARTSCRIPT $GHSCRIPTS/$STARTSCRIPT
check($?)

echo -n "Stop script..."
curl -kLs -o $BINDIR/$STOPSCRIPT $GHSCRIPTS/$STOPSCRIPT
check($?)

echo "Modemu installed, use Scripts menu to launch or stop it"
