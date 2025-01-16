#!/bin/bash

# Mkdir
mkdir -p ./prebuilt

# Get the current git tag
IMAGE_TAG=$(git describe --tags --abbrev=0)

# If no tag is found, exit with an error
if [ -z "$IMAGE_TAG" ]; then
    echo "Error: No git tag found. Exiting."
    exit 1
fi

echo "Using image tag: $IMAGE_TAG"

# Build for linux/amd64
env GOOS=linux GOARCH=amd64 go build -ldflags "-X 'main.IMAGE_TAG=$IMAGE_TAG'" -o ./prebuilt/wg-controller . 