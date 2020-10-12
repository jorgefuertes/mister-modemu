#!/usr/bin/env bash

if [[ "$0" != *scripts/build2mister.sh ]]
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

EXE_NAME="mister-modemu"
VER=$(git describe --tags --contains)
DEST="/media/fat"

WHO=$(whoami)
TIME=$(date +"%d-%m-%Y@%H:%M:%S")

if [[ -f .build ]]
then
	BUILD=$(cat .build)
else
	BUILD=0
fi
BUILD=$(($BUILD + 1))
echo $BUILD > .build

FLAGS="-s -w \
	-X 'github.com/jorgefuertes/mister-modemu/build.version=${VER}' \
	-X 'github.com/jorgefuertes/mister-modemu/build.user=${WHO}' \
	-X 'github.com/jorgefuertes/mister-modemu/build.time=${TIME}' \
	-X 'github.com/jorgefuertes/mister-modemu/build.number=${BUILD}' \
"

mkdir -p bin
rm -f bin/$EXE_NAME*

i="linux"
j="arm"
o="${EXE_NAME}_${VER}-${i}_${j}"

echo "Building ${i} ${j}"
GOOS=$i GOARCH=$j go build -o bin/$o cmd/mister-modemu/main.go
if [[ $? -ne 0 ]]
then
    echo "Compilation error!"
    exit 1
fi

ssh root@$MISTER_IP rm "${DEST}/mister-modemu*" &> /dev/null
scp bin/$o "root@${MISTER_IP}:${DEST}/mister-modemu.tmp"
ssh root@$MISTER_IP mv "${DEST}/mister-modemu.tmp" "${DEST}/${o}"
ssh root@$MISTER_IP "${DEST}/${o} -e dev"
