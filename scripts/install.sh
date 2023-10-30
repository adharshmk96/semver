#!/bin/bash

VERSION="1.0.0"

curl -sL https://github.com/adharshmk96/semver/releases/download/v$VERSION/semver_$VERSION_linux_amd64.tar.gz | sudo tar xz -C /usr/local/bin semver
