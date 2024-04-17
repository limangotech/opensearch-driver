#! /bin/sh

apk add -u --no-cache tzdata make

go mod vendor

make all