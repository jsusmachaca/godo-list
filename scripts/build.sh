#!/usr/bin/env bash

ARCH=$(uname -m)

if [[ $ARCH == "x86_64" ]]; then
    GOARCH="amd64"
elif [[ $ARCH == "i686" || $ARCH == "i386" ]]; then
    GOARCH="386"
else
    echo "Arquitectura no soportada: $ARCH"
    exit 1
fi

GOARCH=$GOARCH GOOS=linux go build cmd/web/main.go

mkdir -p src/web/template
mkdir -p src/web/static

cp -r web/template/* src/web/template
cp -r web/static/* src/web/static
mv main src

tar -czf web.gz src 

rm -rf src
