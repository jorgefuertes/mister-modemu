#!/usr/bin/env bash

if [[ "$0" != *scripts/script2mister.sh ]]
then
	echo "Please, execute from project's root directory"
	exit 1
fi

if [[ -f .mister_ip ]]
then
	MISTER_IP=$(cat .mister_ip)
else
	while [[ 1 ]]
	do
		echo -n "Your Mister's IP? "
		read MISTER_IP
		echo -n "Mister IP ${MISTER_IP}, that's correct? (y/n) "
		read -n1 yn
		echo
		if [[ "$yn" == "y" ]]
		then
			echo $MISTER_IP > .mister_ip
			break
		fi
	done
fi

scp scripts/RetroWiki_Modemu*.sh root@$MISTER_IP:/media/fat/scripts/.
if [[ $? -eq 0 ]]
then
	echo "OK"
else
	echo "FAILED"
	exit 1
fi
