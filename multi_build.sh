#!/bin/bash

rm -r ./build

GOOS=darwin GOARCH=arm64 go build -o build/godot-beautifier-mac-arm64
GOOS=darwin GOARCH=amd64 go build -o build/godot-beautifier-mac-amd64
GOOS=windows GOARCH=amd64 go build -o build/godot-beautifier-win.exe
GOOS=linux GOARCH=amd64 go build -o build/godot-beautifier-linux
