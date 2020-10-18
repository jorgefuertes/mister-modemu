EXE_NAME="mister-modemu"
git describe --tags --contains &> /dev/null
if [[ $? -eq 0 ]]
then
	VER=$(git describe --tags --contains)
else
	VER=$(git describe --tags)
	if [[ $? -ne 0 ]]
	then
		echo "Cannot get version tag, please check git status"
		exit 1
	fi
fi
echo "Version: ${VER}"

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
	-X github.com/jorgefuertes/mister-modemu/internal/build.version=${VER} \
	-X github.com/jorgefuertes/mister-modemu/internal/build.user=${WHO} \
	-X github.com/jorgefuertes/mister-modemu/internal/build.time=${TIME} \
	-X github.com/jorgefuertes/mister-modemu/internal/build.number=${BUILD} \
"
