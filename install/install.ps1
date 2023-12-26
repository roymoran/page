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
$BINARY_URL = "https://github.com/$REPO_USER/$REPO_NAME/releases/download/$VERSION/page_${OS}_${ARCH}.zip"


# Define the local binary path and zip path
$LOCAL_BINARY_DIR = "C:\Program Files\page"
$LOCAL_BINARY_TARBALL = "$LOCAL_BINARY_DIR\page.zip"

# Create directory if not exists
New-Item -ItemType Directory -Force -Path $LOCAL_BINARY_DIR

# Check if directory was created
if (-not (Test-Path $LOCAL_BINARY_DIR)) {
    Write-Host "Failed to create required directory for program. Please run PowerShell as an administrator."
    $success = $false
}

# Check if the directory is already in PATH, if not add it
if ($success) {
    $Path = [System.Environment]::GetEnvironmentVariable("Path", [System.EnvironmentVariableTarget]::Machine)
    if (-not ($Path.Split(';') -contains $LOCAL_BINARY_DIR)) {
        $NewPath = $Path + ';' + $LOCAL_BINARY_DIR
        [System.Environment]::SetEnvironmentVariable("Path", $NewPath, [System.EnvironmentVariableTarget]::Machine)
        $env:Path = [System.Environment]::GetEnvironmentVariable("Path", [System.EnvironmentVariableTarget]::Machine)
    }
}

# Download and extract the binary
if ($success) {
    try {
        Invoke-WebRequest -Uri $BINARY_URL -OutFile $LOCAL_BINARY_TARBALL
    } catch {
        Write-Host "Failed to download the binary. Please check your internet connection and try again."
        $success = $false
    }
}

# Extract the binary
if ($success) {
    try {
        tar -xf $LOCAL_BINARY_TARBALL -C $LOCAL_BINARY_DIR
    } catch {
        Write-Host "Failed to extract the binary. Please ensure 'tar' is installed and try again."
        $success = $false
    }
}

# Clean up
if ($success) {
    Remove-Item $LOCAL_BINARY_TARBALL
    Write-Host "Installation complete. 'page' command is now available."
}
