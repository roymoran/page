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
PROGRAM_OUTPUT_NAME = page

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
  MACOS_CODESIGN_IDENTITY := $(shell security find-identity -v -p codesigning | grep -o 'Apple Development: Roy Moran (YNA4H679A6)')
  ifeq ($(MACOS_CODESIGN_IDENTITY),)
    MACOS_CODESIGN :=
  else
  # on macos use "security find-identity -v -p codesigning" to list available certificates
    MACOS_CODESIGN := && codesign --sign "Apple Development: Roy Moran (YNA4H679A6)" --timestamp $(PROGRAM_OUTPUT_NAME)
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
	cd src/tests/unit && go test -v
	cd src/tests/integration && go test -v
	cd src/tests/system && go test -v

release: $(BUILD_RELEASE_PATH) page_darwin_amd64.tar.bz2 page_darwin_arm64.tar.bz2 page_linux_amd64.tar.bz2 page_linux_arm64.tar.bz2 page_linux_arm.tar.bz2 page_windows_amd64.zip page_windows_arm64.zip

# macos intel 64-bit and codesign with Apple Developer Certificate
page_darwin_amd64.tar.bz2:
	cd src && env GOOS=darwin GOARCH=amd64 go build -o $(PROGRAM_OUTPUT_NAME) $(MACOS_CODESIGN) && tar -cjvf "page_darwin_amd64.tar.bz2" page && mv page_darwin_amd64.tar.bz2 ../$(BUILD_RELEASE_PATH) && rm $(PROGRAM_OUTPUT_NAME)

# macos arm 64-bit and codesign with Apple Developer Certificate
page_darwin_arm64.tar.bz2:
	cd src && env GOOS=darwin GOARCH=arm64 go build -o $(PROGRAM_OUTPUT_NAME) $(MACOS_CODESIGN) && tar -cjvf "page_darwin_arm64.tar.bz2" $(PROGRAM_OUTPUT_NAME) && mv page_darwin_arm64.tar.bz2 ../$(BUILD_RELEASE_PATH) && rm $(PROGRAM_OUTPUT_NAME)

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

clean:
	$(CLEANUP) $(BUILD_PATH)
