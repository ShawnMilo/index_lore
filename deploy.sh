#!/usr/bin/env bash

if [[ "$BINARYSERVER" == "" ]]; then
    echo "Not uploading; no server info provided."
else
    # Look how easy cross-compilation is in Go -- just environment variables!
    GOARCH=amd64 GOOS=linux go build -o /tmp/index_lore_Linux64
    GOARCH=amd64 GOOS=darwin go build -o /tmp/index_lore_OSX64
    GOARCH=amd64 GOOS=windows go build -o /tmp/index_lore_Win64

    scp /tmp/index_lore_Linux64 $BINARYSERVER
    scp /tmp/index_lore_OSX64 $BINARYSERVER
    scp /tmp/index_lore_Win64 $BINARYSERVER
fi

git push origin master
