#!/usr/bin/env bash

if [[ "$0" != *scripts/*.sh ]]
then
	echo "Please, execute from project's root directory"
	exit 1
fi

source scripts/build_common.inc.sh

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

DEST="/media/fat/retrowiki-bin"

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
pushd dist
gzip "${EXE_NAME}_${VER}-${i}_${j}"
echo "RELEASE:"
ls -la *.gz
popd

echo "Copying to mister at ${MISTER_IP}"
ssh root@$MISTER_IP mkdir -p $DEST
ssh root@$MISTER_IP rm "${DEST}/${EXE_NAME}" &> /dev/null
scp "bin/${EXE_NAME}_${VER}-${i}_${j}" "root@${MISTER_IP}:${DEST}/${EXE_NAME}.tmp"
ssh root@$MISTER_IP mv "${DEST}/${EXE_NAME}.tmp" "${DEST}/${EXE_NAME}"
ssh root@$MISTER_IP "${DEST}/${EXE_NAME} -e dev"
