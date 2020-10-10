#!/usr/bin/env bash

EXE_NAME="mister-modemu"
VER=$(git describe --tags --contains)

WHO=$(whoami)
TIME=$(date +"%d-%m-%Y@%H:%M:%S")
OS_LIST=(darwin linux)
ARCH_LIST=(amd64 arm arm64)

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

for i in "${OS_LIST[@]}"
do
    for j in "${ARCH_LIST[@]}"
    do
		if [[ "$i" == "darwin" && "$j" != "amd64" ]]
		then
			echo "Refusing to build ${i}/${j}"
		else
	        echo "Building ${i} ${j}"
	        GOOS=$i GOARCH=$j go build -ldflags="${FLAGS}" -o "bin/${EXE_NAME}_${VER}-${i}_${j}"
	        if [[ $? -ne 0 ]]
	        then
	            echo "Compilation error!"
	            exit 1
	        fi
		fi
    done
done

mkdir -p dist
rm -f dist/*.tar.gz
tar -czvf "dist/${EXE_NAME}_${VER}.tar.gz" -C bin .
