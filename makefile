ifeq ($(OS),Windows_NT)
  ifeq ($(shell uname -s),) # not in a bash-like shell
	CLEANUP = del /F /Q
	MKDIR = mkdir
  else # in a bash-like shell, like msys
	CLEANUP = rm -f -r
	MKDIR = mkdir -p
  endif
	TARGET_EXTENSION=.exe
else
	CLEANUP = rm -f -r
	MKDIR = mkdir -p
	TARGET_EXTENSION=out
endif

#######################################
# build output paths
#######################################
# Build output root path
BUILD_PATH = build/
BUILD_RELEASE_PATH = build/bins/

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
	cd src && go build -o ../$(BUILD_PATH)page

test:
	cd src/tests && go test

release: $(BUILD_RELEASE_PATH) page_darwin_amd64.tar.bz2 page_linux_amd64.tar.bz2 page_linux_arm64.tar.bz2 page_linux_arm.tar.bz2 page_windows_amd64.tar.bz2

# macos intel 64-bit
page_darwin_amd64.tar.bz2:
	cd src && env GOOS=darwin GOARCH=amd64 go build -o page && tar -cjvf "page_darwin_amd64.tar.bz2" page && mv page_darwin_amd64.tar.bz2 ../$(BUILD_RELEASE_PATH) && rm page

# linux intel 64-bit
page_linux_amd64.tar.bz2:
	cd src && env GOOS=linux GOARCH=amd64 go build -o page && tar -cjvf "page_linux_amd64.tar.bz2" page && mv page_linux_amd64.tar.bz2 ../$(BUILD_RELEASE_PATH) && rm page

# linux arm 64-bit
page_linux_arm64.tar.bz2:
	cd src && env GOOS=linux GOARCH=arm64 go build -o page && tar -cjvf "page_linux_arm64.tar.bz2" page && mv page_linux_arm64.tar.bz2 ../$(BUILD_RELEASE_PATH) && rm page

# linux arm 32-bit
page_linux_arm.tar.bz2:
	cd src && env GOOS=linux GOARCH=arm go build -o page && tar -cjvf "page_linux_arm.tar.bz2" page && mv page_linux_arm.tar.bz2 ../$(BUILD_RELEASE_PATH) && rm page

# windows intel 64-bit
page_windows_amd64.tar.bz2:
	cd src && env GOOS=windows GOARCH=amd64 go build -o page.exe && tar -cjvf "page_windows_amd64.tar.bz2" page.exe && mv page_windows_amd64.tar.bz2 ../$(BUILD_RELEASE_PATH) && rm page.exe

$(BUILD_PATH):
	$(MKDIR) $@

$(BUILD_RELEASE_PATH):
	$(MKDIR) $@

clean:
	$(CLEANUP) $(BUILD_PATH)
