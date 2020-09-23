#!/usr/bin/env bash

EXE_NAME="mister-modemu"
VER=$(git describe --tags)

OS_LIST=(darwin linux)
ARCH_LIST=(amd64 arm arm64)

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
	        GOOS=$i GOARCH=$j go build -o "bin/${EXE_NAME}_${VER}-${i}_${j}"
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
