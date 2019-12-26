#!/bin/bash

# Install Script for volley for linux/mac/windows (GitBash) enviornment 
# Execute this script with user that sudo is available

# Download volley binary
# https://github.com/TsuyoshiUshio/volley/releases/download/0.0.3/volley-linux-amd64.tgz

VERSION="0.0.6"
BINARY_TYPE=""
ARCHIVE=""
BINARY=""
INSTALL_PATH=/usr/bin

if [ "$(uname)" == "Darwin" ]; then
    # Do something under Mac OS X platform  
    BINARY_TYPE="volley-darwin-amd64" 
    ARCHIVE="${BINARY_TYPE}.tgz"
    BINARY="${BINARY_TYPE}/volley"
elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
    # Do something under GNU/Linux platform
    BINARY_TYPE="volley-linux-amd64" 
    ARCHIVE="${BINARY_TYPE}.tgz"
    BINARY="${BINARY_TYPE}/volley"
elif [ "$(expr substr $(uname -s) 1 10)" == "MINGW32_NT" ]; then
    # Do something under 32 bits Windows NT platform
    BINARY_TYPE="volley-windows-386"
    ARCHIVE="${BINARY_TYPE}.zip"
    BINARY="${BINARY_TYPE}/volley.exe"
elif [ "$(expr substr $(uname -s) 1 10)" == "MINGW64_NT" ]; then
    # Do something under 64 bits Windows NT platform
    BINARY_TYPE="volley-windows-amd64"
    ARCHIVE="${BINARY_TYPE}.zip"
    BINARY="${BINARY_TYPE}/volley.exe"
fi

if type curl > /dev/null 2>&1; then
    echo "curl exist."
else
    # Linux only
    sudo apt-get update -y 
    sudo apt-get install curl -y
fi
echo "Downloading ... : curl -OL  https://github.com/TsuyoshiUshio/volley/releases/download/${VERSION}/${ARCHIVE}"
curl -OL  https://github.com/TsuyoshiUshio/volley/releases/download/${VERSION}/${ARCHIVE}

if [ "$(expr substr $(uname -s) 1 5)" == "MINGW" ]; then
# windows (GitBash) based
    unzip $ARCHIVE
    # TODO put exe file in somewhere already have a path or create directry and add path to the .bashrc
else
# Linux and Mac
    tar xvzf $ARCHIVE 
    if [ -f /usr/bin/volley ]; then
      # remove old version
      sudo rm /usr/bin/volley -f
    fi
    sudo cp $BINARY /usr/bin
    rm -rf $BINARY_TYPE
    rm $ARCHIVE
fi

