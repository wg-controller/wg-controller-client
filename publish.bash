#!/bin/bash

# Get the current git tag
IMAGE_TAG=$(git describe --tags --abbrev=0)

# If no tag is found, exit with an error
if [ -z "$IMAGE_TAG" ]; then
    echo "Error: No git tag found. Exiting."
    exit 1
fi

echo "Using image tag: $IMAGE_TAG"

# Create prebuilt directory
mkdir -p ./prebuilt

# Build for linux/amd64
env GOOS=linux GOARCH=amd64 go build -ldflags "-X 'main.IMAGE_TAG=$IMAGE_TAG'" -o ./prebuilt/wg-controller-linux . 

# Build for linux/arm64
env GOOS=linux GOARCH=amd64 go build -ldflags "-X 'main.IMAGE_TAG=$IMAGE_TAG'" -o ./prebuilt/wg-controller-linuxarm64 . 

# Check if the release already exists and delete it
EXISTING_RELEASE=$(gh release view "$IMAGE_TAG" --json tagName -q ".tagName")
if [ "$EXISTING_RELEASE" == "$IMAGE_TAG" ]; then
    echo "Release with tag '$IMAGE_TAG' exists, deleting it..."
    gh release delete "$IMAGE_TAG" -y
fi

# Create release
gh release create $IMAGE_TAG ./prebuilt/wg-controller-linux ./prebuilt/wg-controller-linuxarm64 --title $IMAGE_TAG --notes "Release $IMAGE_TAG"

# Check if release "latest" already exists and delete it
EXISTING_RELEASE=$(gh release view "latest" --json tagName -q ".tagName")
if [ "$EXISTING_RELEASE" == "latest" ]; then
    echo "Release with tag 'latest' exists, deleting it..."
    gh release delete "latest" -y
fi

# Create release "latest"
gh release create latest ./prebuilt/wg-controller-linux ./prebuilt/wg-controller-linuxarm64 --title "Latest Release" --notes "Release $IMAGE_TAG"