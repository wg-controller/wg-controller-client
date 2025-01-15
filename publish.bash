#!/bin/bash

# Mkdir
mkdir ./prebuilt

# Build for linux/amd64
env GOOS=linux GOARCH=amd64 go build -o ./prebuilt/wg-controller . 