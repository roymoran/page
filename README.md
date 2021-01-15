# page

## Dependencies
- [Golang v1.15.6](https://golang.org/)
## Build, Test, and Run
```bash
# from src directory
# install dependencies
$ go install
# build and output executable to user programs directories
# once built you can execute the command using the executable name
# 'page' from you command line
$ go build -o /usr/local/bin/page
# run tests
$ cd tests
$ go test
# run without building executable
$ go run main.go
```