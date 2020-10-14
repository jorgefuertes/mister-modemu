#!/usr/bin/env bash

if [[ "$0" != *scripts/*.sh ]]
then
	echo "Please, execute from project's root directory"
	exit 1
fi

source scripts/build_common.inc.sh

mkdir -p bin
rm -f bin/$EXE_NAME*

i="linux"
j="arm"

echo "Building ${i} ${j}"
GOOS=$i GOARCH=$j go build -ldflags "${FLAGS}" \
	-o "bin/${EXE_NAME}_${VER}-${i}_${j}" \
	cmd/mister-modemu/main.go
if [[ $? -ne 0 ]]
then
    echo "Compilation error!"
    exit 1
fi

rm -f dist/* &> /dev/null
cp "bin/${EXE_NAME}_${VER}-${i}_${j}" "dist/${EXE_NAME}_${VER}-${i}_${j}"
pushd dist &> /dev/null
gzip "${EXE_NAME}_${VER}-${i}_${j}"
echo "RELEASE:"
ls -la *.gz
popd &> /dev/null
