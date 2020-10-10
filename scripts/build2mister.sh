#!/usr/bin/env bash

EXE_NAME="mister-modemu"
VER=$(git describe --tags --contains)
MISTER="root@192.168.1.10"
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

scp bin/$o "${MISTER}:${DEST}/mister-modemu.tmp"
ssh $MISTER mv "${DEST}/mister-modemu.tmp" "${DEST}/${o}"
ssh $MISTER "${DEST}/${o} -e dev"
