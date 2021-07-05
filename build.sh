set -x

sh ./prebuild.sh

go build -o bin/ipblock_server rest_server/main.go

cd bin
./ipblock_server -c=config.yml