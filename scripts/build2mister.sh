#!/usr/bin/env bash

check() {
	if [[ $1 -ne 0 ]]
	then
		echo "> ERROR !!!"
		pwd
		exit 1
	fi
}

if [[ "$0" != *scripts/*.sh ]]
then
	echo "> Please, execute from project's root directory !!!"
	exit 1
fi

source scripts/build_common.inc.sh

if [[ -f .mister_ip ]]
then
	MISTER_IP=$(cat .mister_ip)
else
	while [[ 1 ]]
	do
		echo -n "> Your Mister's IP? "
		read MISTER_IP
		echo -n "> Mister IP ${MISTER_IP}, that's correct? (y/n) "
		read -n1 yn
		echo
		if [[ "$yn" == "y" ]]
		then
			echo $MISTER_IP > .mister_ip
			break
		fi
	done
fi

DEST="/media/fat/.tmp"

mkdir -p bin
rm -f bin/$EXE_NAME*

i="linux"
j="arm"

echo "> Building ${i} ${j}"
LOCAL_EXE_NAME="bin/${EXE_NAME}_${VER}-${i}_${j}"
GOOS=$i GOARCH=$j go build -ldflags "${FLAGS}" \
	-o $LOCAL_EXE_NAME \
	cmd/mister-modemu/main.go
check $?

rm -f dist/* &> /dev/null
cp $LOCAL_EXE_NAME "dist/${EXE_NAME}_${VER}-${i}_${j}"
pushd dist &> /dev/null
gzip "${EXE_NAME}_${VER}-${i}_${j}"
echo -n "> RELEASE: "
ls -la *.gz
popd &> /dev/null

echo "> Making tmp dir"
ssh root@$MISTER_IP mkdir -p $DEST
check $?
echo "> Copying to mister"
scp -q $LOCAL_EXE_NAME root@$MISTER_IP:$DEST/$EXE_NAME.tmp
check $?
echo "> Moving .tmp to ${EXE_NAME}"
ssh root@$MISTER_IP mv "${DEST}/${EXE_NAME}.tmp" "${DEST}/${EXE_NAME}"
check $?
echo "> Executing in dev mode"
ssh root@$MISTER_IP "${DEST}/${EXE_NAME} -e dev"
