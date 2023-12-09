# Developing

For fast debugging just build the application and execute it.

```bash
# compile and output application to build directory
$ make
# navigate to executable
$ cd build
# execute
$ ./page
```

[Delve go debugger](https://github.com/go-delve/delve) is used for step through debugging. Use VSCode debugging panel for getting started quickly. Otherwise see go-delve/delve repo for instructions.

To debug terraform problems or see deprecation warnings use terraform directly from the command line.

```bash
# see basic output
$ terraform apply
# verbose logging
$ TF_LOG=debug terraform apply
```

# Testing

Test coverage is limited but growing. Tests will focus on unit testing, integration testing, and system testing. Here is what those mean in this project.

- Unit testing: This means testing individual functions or group of functions inside the page codebase. These tests should run quickly (in matter of seconds) and shouldn't interact with any external resources like the filesystem or network. A couple examples that might fall under this category are commands like `page`.
- Integration testing: This means testing commands that interact with external resources like the file system. These are somewhat longer running than unit test but each test should complete in under a minute so these won't include commands that interact with external network resources that may be long-running like a request to create infrastructure on AWS.
- System testing: This means testing the commands that include interactions with external resources like a specific host or registrar. These tests may be longest-running and may complete in several minutes. This means testing for DNS changes or new infrastucture creation when executing commands like `page up`

```bash
$ make test
```

# Releasing

Releases are managed through git tags. All binaries are compiled using the [makefile](../makefile) and is handled by the build system.

```bash
# create tag
$  git tag -a v0.1.0-alpha.12 -m "v0.1.0-alpha.12"
# push tag
$ git push origin --tags
```
