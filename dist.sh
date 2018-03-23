#!/bin/bash

# build binary distributions for linux/amd64 and darwin/amd64
set -e

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
echo "working dir $DIR"
mkdir -p $DIR/dist
glide install

# This is how we want to name the binary output
NAME=tamtam

# These are the values we want to pass for Version and BuildTime
VERSION=`cat VERSION | awk '{print $1}'`
BUILD_TIME=`date +%FT%T%z`

if [ -z "$2" ]; then
    GIT_HASH=`git rev-parse HEAD`
else
    GIT_HASH=$2
fi

PACKAGE="github.com/eelcocramer/tamtam"

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS="-X $PACKAGE/cmd.Version=$VERSION -X $PACKAGE/cmd.BuildTime=$BUILD_TIME -X $PACKAGE/cmd.GitHash=$GIT_HASH"
echo $LDFLAGS

os=$(go env GOOS)
arch=$(go env GOARCH)
goversion=$(go version | awk '{print $3}')

declare -a targets
if [ -z "$1" ]; then
    targets=(windows linux darwin)
else
    targets+=($1)
fi

for os in "${targets[@]}"; do
    echo "... building v$VERSION for $os/$arch"
    TARGET="$NAME-v$VERSION-$os"
    if [ $os = 'windows' ] ; then
        BINARY=$NAME.exe
    else
        BINARY=$NAME
    fi

    BUILD=$(mktemp -d -t $NAME.XXXX)
    GOOS=$os GOARCH=$arch CGO_ENABLED=0 go build -ldflags "$LDFLAGS" -o $BUILD/$TARGET/$BINARY || exit 1

    # Only tar if no specific target was specified
    if [ -z "$1" ]; then
        pushd $BUILD
        tar czvf $TARGET.tar.gz $TARGET
        mv $TARGET.tar.gz $DIR/dist
        popd
    else
        mv $BUILD/$TARGET/$BINARY .
    fi
done
