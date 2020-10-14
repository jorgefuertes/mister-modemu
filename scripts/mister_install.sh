#!/usr/bin/env bash

GHSCRIPTS="https://raw.githubusercontent.com/jorgefuertes/mister-modemu/master/scripts/"

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
curl -kLs $GHSCRIPTS/RetroWiki_Modemu_Start.sh
check($?)

echo -n "Stop script..."
curl -kLs $GHSCRIPTS/RetroWiki_Modemu_Stop.sh
check($?)

echo "Modemu installed, use Scripts menu to launch or stop it"
