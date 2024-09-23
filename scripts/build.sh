#!/usr/bin/env sh

ARCH=$(uname -m)

case $ARCH in
    x86_64)
        export GOARCH="amd64"
        ;;
    i686|i386)
        export GOARCH="386"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

export GOOS=linux

if ! go build cmd/web/main.go; then
    echo "Build failed"
    exit 1
fi

SRC_DIR="src"
TEMPLATE_DIR="web/template"
STATIC_DIR="web/static"

mkdir -p $SRC_DIR/web/template
mkdir -p $SRC_DIR/web/static

cp -r $TEMPLATE_DIR/* $SRC_DIR/web/template
cp -r $STATIC_DIR/* $SRC_DIR/web/static
mv main $SRC_DIR

tar -czf webapp.tar.gz $SRC_DIR

rm -rf $SRC_DIR
