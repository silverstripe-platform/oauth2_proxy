#!/bin/bash
# build binary distributions for linux/amd64 and darwin/amd64
set -e

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
echo "working dir $DIR"
mkdir -p $DIR/dist
dep ensure || exit 1

os=$(go env GOOS)
arch=$(go env GOARCH)
version=$(cat $DIR/version.go | grep "const VERSION" | awk '{print $NF}' | sed 's/"//g')
goversion=$(go version | awk '{print $3}')
sha1sum=()

# echo "... running tests"
# ./test.sh

for os in windows linux darwin; do
    echo "... building v$version for $os/$arch"
    EXT=
    if [ $os = windows ]; then
        EXT=".exe"
    fi
    BUILD=$(mktemp -d ${TMPDIR:-/tmp}/oauth2_proxy.XXXXXX)
    TARGET="oauth2_proxy-$version.$os-$arch.$goversion"
    FILENAME="oauth2-proxy$EXT"
    GOOS=$os GOARCH=$arch CGO_ENABLED=0 \
    GO111MODULE=auto \
        go build -ldflags="-s -w" -o $BUILD/$TARGET/$FILENAME || exit 1
    pushd $BUILD/$TARGET
    mv $FILENAME $TARGET
    mv $TARGET $FILENAME
    cd .. && tar czvf $TARGET.tar.gz $TARGET
    sha1sum+=("$(shasum -a 1 $TARGET.tar.gz || exit 1)")
    mv $TARGET.tar.gz $DIR/dist
    popd
done

checksum_file="sha1sum.txt"
cd $DIR/dist
if [ -f $checksum_file ]; then
    rm $checksum_file
fi
touch $checksum_file
for checksum in "${sha1sum[@]}"; do
    echo "$checksum" >> $checksum_file
done
