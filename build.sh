#!/bin/sh

set -e
set -x

go get -v -u github.com/mitchellh/gox
gox
