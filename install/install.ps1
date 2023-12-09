# PowerShell script for installing pagecli on Windows

# GitHub repository details
$REPO_USER = "roymoran"
$REPO_NAME = "page"

# Fetch the latest release tag
$LATEST_RELEASE_API = "https://api.github.com/repos/$REPO_USER/$REPO_NAME/releases/latest"
$VERSION = Invoke-RestMethod -Uri $LATEST_RELEASE_API | Select-Object -ExpandProperty tag_name

# Construct the download URL
$OS = "windows"
$ARCH = if ([System.Environment]::Is64BitProcess) { "amd64" } else { "arm64" } 
$BINARY_URL = "https://github.com/$REPO_USER/$REPO_NAME/releases/download/$VERSION/page_${OS}_${ARCH}.tar.bz2"


# Define the local binary path
$LOCAL_BINARY_PATH = "C:\Program Files\page\page.exe"

# Create directory if not exists
New-Item -ItemType Directory -Force -Path "C:\Program Files\page"

# Download and extract the binary
Invoke-WebRequest -Uri $BINARY_URL -OutFile "$LOCAL_BINARY_PATH.tar.bz2"
tar -xf "$LOCAL_BINARY_PATH.tar.bz2" -C "C:\Program Files\page"

# Clean up
Remove-Item "$LOCAL_BINARY_PATH.tar.bz2"

Write-Host "Installation complete. 'page' command is now available."
