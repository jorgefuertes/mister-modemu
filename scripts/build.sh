#!/usr/bin/env bash

if [[ "$0" != *scripts/*.sh ]]
then
	echo "Please, execute from project's root directory"
	exit 1
fi

source scripts/build_common.inc.sh

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
	        GOOS=$i GOARCH=$j go build -ldflags "${FLAGS}" \
				-o "bin/${EXE_NAME}_${VER}-${i}_${j}" \
				cmd/mister-modemu/main.go
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
