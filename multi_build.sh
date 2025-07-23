#!/bin/bash

rm ./build/*

GOOS=darwin GOARCH=arm64 go build ./src -o build/godot-beautifier-mac-apple_sillicon
GOOS=darwin GOARCH=amd64 go build ./src -o build/godot-beautifier-mac-intel
GOOS=windows GOARCH=amd64 go build ./src -o build/godot-beautifier-win.exe
GOOS=linux GOARCH=amd64 go build ./src -o build/godot-beautifier-linux
