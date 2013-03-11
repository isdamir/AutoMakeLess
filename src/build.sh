#!/bin/sh
APP="AutoMakeLess"
PWD=`pwd`/..
ADD_ON="${PWD}/src/add-on"
export GOPATH=${ADD_ON}:${PWD}

if [ -f ${APP} ]; then
    rm ${APP}
fi

echo "Building ${APP}"
go build .

if [[ -f src ]]; then
    mv ./src ${APP}
	echo "OK"
fi
