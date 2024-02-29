ifeq ($(OS),Windows_NT)
  ifeq ($(shell uname -s),) # not in a bash-like shell
	CLEANUP = del /F /Q
	MKDIR = mkdir
  else # in a bash-like shell, like msys
	CLEANUP = rm -f -r
	MKDIR = mkdir -p
  endif
else
	CLEANUP = rm -f -r
	MKDIR = mkdir -p
endif

#######################################
# build output paths
#######################################
# Build output root path
BUILD_PATH = build/
BUILD_RELEASE_PATH = build/bins/
BUILD_RELEASE_PKGS_PATH = build/bins/pkg/
PROGRAM_OUTPUT_NAME = page
MACOS_ARM64_PACKAGE_PATH = build/bins/pkg/arm64/
MACOS_AMD64_PACKAGE_PATH = build/bins/pkg/amd64/
DEVELOPER_ID_INSTALLER := "Developer ID Installer: Roy Moran (KQ859WKQZ6)"
NOTARIZE_KEYCHAIN_PROFILE := AppNotaryService
PKG_IDENTIFIER := com.pagecli.page
VERSION := 1.0.0

#######################################
# Check if codesign programs are available
#######################################
MACOS_CODESIGN_AVAILABLE := $(shell command -v codesign 2> /dev/null)
# provided by osslsigncode package https://github.com/mtrojnar/osslsigncode
WINDOWS_CODESIGN_AVAILABLE := $(shell command -v osslsigncode 2> /dev/null)
LINUX_CODESIGN_AVAILABLE := $(shell command -v gpg 2> /dev/null)

#######################################
# code signing
#######################################
ifeq ($(MACOS_CODESIGN_AVAILABLE),)
  MACOS_CODESIGN :=
else
  MACOS_CODESIGN_IDENTITY := $(shell security find-identity -v -p codesigning | grep -o 'Developer ID Application: Roy Moran (KQ859WKQZ6)')
  ifeq ($(MACOS_CODESIGN_IDENTITY),)
    MACOS_CODESIGN :=
  else
  # on macos use "security find-identity -v -p codesigning" to list available certificates
    MACOS_CODESIGN := && codesign --sign "Developer ID Application: Roy Moran (KQ859WKQZ6)" --timestamp --options runtime --verbose $(PROGRAM_OUTPUT_NAME)
	MACOS_PKGBUILD_AMD64 := && pkgbuild --root ../$(MACOS_AMD64_PACKAGE_PATH) --identifier $(PKG_IDENTIFIER) --version $(VERSION) --install-location /usr/local/bin ../$(MACOS_AMD64_PACKAGE_PATH)page_darwin_amd64_unsigned.pkg
	MACOS_PKGBUILD_ARM64 := && pkgbuild --root ../$(MACOS_ARM64_PACKAGE_PATH) --identifier $(PKG_IDENTIFIER) --version $(VERSION) --install-location /usr/local/bin ../$(MACOS_ARM64_PACKAGE_PATH)page_darwin_arm64_unsigned.pkg
	MACOS_PKGSIGN_AMD64 := && productsign --sign "Developer ID Installer: Roy Moran (KQ859WKQZ6)" ../$(MACOS_AMD64_PACKAGE_PATH)page_darwin_amd64_unsigned.pkg ../$(MACOS_AMD64_PACKAGE_PATH)page_darwin_amd64.pkg
	MACOS_PKGSIGN_ARM64 := && productsign --sign "Developer ID Installer: Roy Moran (KQ859WKQZ6)" ../$(MACOS_ARM64_PACKAGE_PATH)page_darwin_arm64_unsigned.pkg ../$(MACOS_ARM64_PACKAGE_PATH)page_darwin_arm64.pkg
	MACOS_PKGNOTARIZE_AMD64 := && xcrun notarytool submit ../$(MACOS_AMD64_PACKAGE_PATH)page_darwin_amd64.pkg --keychain-profile $(NOTARIZE_KEYCHAIN_PROFILE) --wait && xcrun stapler staple ../$(MACOS_AMD64_PACKAGE_PATH)page_darwin_amd64.pkg && spctl --assess --verbose --type install ../$(MACOS_AMD64_PACKAGE_PATH)page_darwin_amd64.pkg
	MACOS_PKGNOTARIZE_ARM64 := && xcrun notarytool submit ../$(MACOS_ARM64_PACKAGE_PATH)page_darwin_arm64.pkg --keychain-profile $(NOTARIZE_KEYCHAIN_PROFILE) --wait && xcrun stapler staple ../$(MACOS_ARM64_PACKAGE_PATH)page_darwin_arm64.pkg && spctl --assess --verbose --type install ../$(MACOS_ARM64_PACKAGE_PATH)page_darwin_arm64.pkg
  endif
endif

ifeq ($(WINDOWS_CODESIGN_AVAILABLE),)
  WINDOWS_CODESIGN :=
else
ifeq ($(CODESIGN_PFX_PASS),)
  $(error CODESIGN_PFX_PASS is not set)
endif
# on macos use "security find-identity -v -p codesigning" to list available certificates
  WINDOWS_CODESIGN := && osslsigncode sign -pkcs12 "/Volumes/home/devkeys/codesigningcerts/windows/selfsigned/certificate.pfx" -pass "$(CODESIGN_PFX_PASS)" -n "$(PROGRAM_OUTPUT_NAME)" -i https://pagecli.com/ -in $(PROGRAM_OUTPUT_NAME).exe -out $(PROGRAM_OUTPUT_NAME)-signed.exe && mv $(PROGRAM_OUTPUT_NAME)-signed.exe $(PROGRAM_OUTPUT_NAME).exe
endif

ifeq ($(LINUX_CODESIGN_AVAILABLE),)
  LINUX_CODESIGN :=
else
  LINUX_CODESIGN_KEY_PRESENT := $(shell gpg --list-keys | grep -o 'roy.moran@icloud.com')
  ifeq ($(LINUX_CODESIGN_KEY_PRESENT),)
    LINUX_CODESIGN :=
  else
  # gpg-agent must be running and passphrase caching for signing key configured 
    LINUX_CODESIGN := && gpg --detach-sign --armor $(PROGRAM_OUTPUT_NAME) && rm $(PROGRAM_OUTPUT_NAME).asc
  endif
endif

.PHONY: info release fmt tests clean

default: page

info:
	@echo
	@echo following targets are available:
	@echo		make                      - build page cli and output to build directory
	@echo		make tests                - run application tests using go test
	@echo		make release              - build and package page cli for all supported architectures
	@echo		make clean                - remove build output directory forcibly and recursively
	@echo

page: $(BUILD_PATH)
	cd src && go build -o ../$(BUILD_PATH)$(PROGRAM_OUTPUT_NAME)

test:
	cd src/tests/unit && PAGE_CLI_TEST=true go test -v
	cd src/tests/integration && PAGE_CLI_TEST=true go test -v
	cd src/tests/system && PAGE_CLI_TEST=true go test -v

release: $(BUILD_RELEASE_PATH) $(MACOS_ARM64_PACKAGE_PATH) $(MACOS_AMD64_PACKAGE_PATH) page_darwin_amd64.tar.bz2 page_darwin_arm64.tar.bz2 page_linux_amd64.tar.bz2 page_linux_arm64.tar.bz2 page_linux_arm.tar.bz2 page_windows_amd64.zip page_windows_arm64.zip rmpkg

# macos intel 64-bit and codesign with Apple Developer Certificate
page_darwin_amd64.tar.bz2:
	cd src && env GOOS=darwin GOARCH=amd64 go build -o $(PROGRAM_OUTPUT_NAME) $(MACOS_CODESIGN) && mv $(PROGRAM_OUTPUT_NAME) ../$(MACOS_AMD64_PACKAGE_PATH) $(MACOS_PKGBUILD_AMD64) $(MACOS_PKGSIGN_AMD64) $(MACOS_PKGNOTARIZE_AMD64) && mv ../$(MACOS_AMD64_PACKAGE_PATH)page_darwin_amd64.pkg ../$(BUILD_RELEASE_PATH)

# macos arm 64-bit and codesign with Apple Developer Certificate
page_darwin_arm64.tar.bz2:
	cd src && env GOOS=darwin GOARCH=arm64 go build -o $(PROGRAM_OUTPUT_NAME) $(MACOS_CODESIGN) && mv $(PROGRAM_OUTPUT_NAME) ../$(MACOS_ARM64_PACKAGE_PATH) $(MACOS_PKGBUILD_ARM64) $(MACOS_PKGSIGN_ARM64) $(MACOS_PKGNOTARIZE_ARM64) && mv ../$(MACOS_ARM64_PACKAGE_PATH)page_darwin_arm64.pkg ../$(BUILD_RELEASE_PATH)

# linux intel 64-bit
page_linux_amd64.tar.bz2:
	cd src && env GOOS=linux GOARCH=amd64 go build -o $(PROGRAM_OUTPUT_NAME) $(LINUX_CODESIGN) && tar -cjvf "page_linux_amd64.tar.bz2" $(PROGRAM_OUTPUT_NAME) && mv page_linux_amd64.tar.bz2 ../$(BUILD_RELEASE_PATH) && rm $(PROGRAM_OUTPUT_NAME)

# linux arm 64-bit
page_linux_arm64.tar.bz2:
	cd src && env GOOS=linux GOARCH=arm64 go build -o $(PROGRAM_OUTPUT_NAME) $(LINUX_CODESIGN) && tar -cjvf "page_linux_arm64.tar.bz2" $(PROGRAM_OUTPUT_NAME) && mv page_linux_arm64.tar.bz2 ../$(BUILD_RELEASE_PATH) && rm $(PROGRAM_OUTPUT_NAME)

# linux arm 32-bit
page_linux_arm.tar.bz2:
	cd src && env GOOS=linux GOARCH=arm go build -o $(PROGRAM_OUTPUT_NAME) $(LINUX_CODESIGN) && tar -cjvf "page_linux_arm.tar.bz2" $(PROGRAM_OUTPUT_NAME) && mv page_linux_arm.tar.bz2 ../$(BUILD_RELEASE_PATH) && rm $(PROGRAM_OUTPUT_NAME)

# windows intel 64-bit and codesign with osslsigncode
page_windows_amd64.zip:
	cd src && env GOOS=windows GOARCH=amd64 go build -o $(PROGRAM_OUTPUT_NAME).exe $(WINDOWS_CODESIGN) && tar -a -c -f "page_windows_amd64.zip" $(PROGRAM_OUTPUT_NAME).exe && mv page_windows_amd64.zip ../$(BUILD_RELEASE_PATH) && rm $(PROGRAM_OUTPUT_NAME).exe

# windows arm 64-bit and codesign with osslsigncode
page_windows_arm64.zip:
	cd src && env GOOS=windows GOARCH=arm64 go build -o $(PROGRAM_OUTPUT_NAME).exe $(WINDOWS_CODESIGN) && tar -a -c -f "page_windows_arm64.zip" $(PROGRAM_OUTPUT_NAME).exe && mv page_windows_arm64.zip ../$(BUILD_RELEASE_PATH) && rm $(PROGRAM_OUTPUT_NAME).exe

$(BUILD_PATH):
	$(MKDIR) $@

$(BUILD_RELEASE_PATH):
	$(MKDIR) $@

$(MACOS_AMD64_PACKAGE_PATH):
	$(MKDIR) $@

$(MACOS_ARM64_PACKAGE_PATH):
	$(MKDIR) $@

clean:
	$(CLEANUP) $(BUILD_PATH)

rmpkg:
	$(CLEANUP) $(BUILD_RELEASE_PKGS_PATH)
