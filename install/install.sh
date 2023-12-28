#!/bin/sh

# GitHub repository details
REPO_USER="roymoran"
REPO_NAME="page"

# GitHub API URL to get the latest release tag
LATEST_RELEASE_API="https://api.github.com/repos/$REPO_USER/$REPO_NAME/releases/latest"

# Fetch the latest release tag
VERSION=$(curl -s $LATEST_RELEASE_API | awk -F '"' '/tag_name/{print $4}')

# Base URL for downloading binaries
BASE_URL="https://github.com/$REPO_USER/$REPO_NAME/releases/download/$VERSION"

# Detect OS and Architecture
OS=$(uname -s)
ARCH=$(uname -m)

# Map OS and Architecture to binary naming convention
case "$OS" in
    Darwin)
        OS="darwin"
        ;;
    Linux)
        OS="linux"
        ;;
    CYGWIN*|MINGW32*|MSYS*|MINGW*)
        OS="windows"
        ;;
esac

case "$ARCH" in
    aarch64|arm64)
        ARCH="arm64"
        ;;
    amd64|x86_64)
        ARCH="amd64"
        ;;
    arm*)
        ARCH="arm"
        ;;
esac

# Construct the download URL
BINARY_URL="$BASE_URL/page_${OS}_${ARCH}.tar.bz2"

# Define the local binary path
LOCAL_BINARY_PATH="/usr/local/bin/page"

# Download and extract the binary
curl -L "$BINARY_URL" -o "page.tar.bz2"
if [ $? -ne 0 ]; then
    exit 1
fi

# Extract the binary and check if the extraction was successful
tar -xjf "page.tar.bz2"
if [ $? -ne 0 ]; then
    rm "page.tar.bz2"
    exit 1
fi
rm "page.tar.bz2"

# Make the binary executable
chmod +x page

# Move the binary to the desired location using sudo for necessary permissions
sudo mv page "$LOCAL_BINARY_PATH"

echo "Installation complete. 'page' command is now available."
