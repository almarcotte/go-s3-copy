#!/usr/bin/env bash

OS="darwin linux"
ARCH="amd64"
REPO="gnumast/go-s3-copy"

version=$1

if [[ ! -n "$version" ]]; then
    echo "Usage: release.sh <tag> <token>"
    exit 1
fi

echo "This will create tag v$version. Are you sure? (Y/N)"
read confirm

if [ "${confirm}" != "Y" ] && [ "${confirm}" != "y" ]; then
    exit
fi

mkdir -p dist

for os in ${OS}; do
    for arch in ${ARCH}; do
        GOOS=${os} GOARCH=${arch} go build -o dist/go-s3-copy main.go
        tar -cvzf dist/go-s3-copy_${os}_${arch}.tar.gz dist/go-s3-copy
        rm dist/go-s3-copy
    done
done

git tag -a v${version} -m "Release of version ${version}"
git push --tags
