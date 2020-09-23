#!/usr/bin/env bash

EXE_NAME="mister-modemu"
VER=$(git describe --tags)
MISTER="root@192.168.1.10"
DEST="/media/fat"

mkdir -p bin
rm -f bin/$EXE_NAME*

i="linux"
j="arm"
o="${EXE_NAME}_${VER}-${i}_${j}"

echo "Building ${i} ${j}"
GOOS=$i GOARCH=$j go build -o bin/$o
if [[ $? -ne 0 ]]
then
    echo "Compilation error!"
    exit 1
fi

scp bin/$o "${MISTER}:${DEST}/mister-modemu.tmp"
ssh $MISTER mv "${DEST}/mister-modemu.tmp" "${DEST}/${o}"
ssh $MISTER "${DEST}/${o} -e dev"
