#!/bin/bash
# build binary distributions for linux/amd64 and darwin/amd64
# usage: dist.sh <arch> <target os> [git hash]
set -e

# This is how we want to name the binary output
NAME=tamtam
PACKAGE="github.com/eelcocramer/tamtam"

declare -a goos
if [ -z "$1" ]; then
    goos=$(go env GOOS)
else
    goos=$1
fi

declare -a arch
if [ -z "$2" ]; then
    arch=$(go env GOARCH)
else
    arch=$2
fi

# First argument holds
if [ -z "$3" ]; then
    GIT_HASH=$(git rev-parse HEAD)
else
    GIT_HASH=$3
fi

# These are the values we want to pass for Version and BuildTime
VERSION=$(awk '{print $1}' < VERSION)
BUILD_TIME=$(date +%FT%T%z)

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
echo "working dir $DIR"
mkdir -p $DIR/dist

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS="-X $PACKAGE/cmd.Version=$VERSION -X $PACKAGE/cmd.BuildTime=$BUILD_TIME -X $PACKAGE/cmd.GitHash=$GIT_HASH"

declare -a goarm
if [ $arch = 'arm' ] ; then
    if [ $goos = 'darwin' ] ; then
        goarm=7
    else
        if [ $goos = 'linux' ] ; then
            goarm=5
        else
            goarm=6
        fi
    fi
fi

echo "... building v$VERSION for $goos/$arch$goarm"
TARGET="$NAME-v$VERSION-$goos-$arch"
if [ $goos = 'windows' ] ; then
    BINARY=$NAME.exe
else
    BINARY=$NAME
fi

BUILD=$(mktemp -d -t $NAME.XXXX)
GOOS=$goos GOARCH=$arch GOARM=$goarm CGO_ENABLED=0 go build -ldflags "$LDFLAGS" -o $BUILD/$TARGET/$BINARY || exit 1
mv $BUILD/$TARGET/$BINARY $DIR
