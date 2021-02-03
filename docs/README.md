# Docs
## Development
### Dependencies
- [Golang v1.15.6](https://golang.org/)
### Build, Test, and Run
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
## Page Definition File
The goal of the definition file is to describe only the essential information required to create and deploy a page with the cli. The simplest form of the definition file is what is generated when running `page init`. We should use smart defaults for any key that isn't included or isn't expanded on. Any modifications to the definition file should aim to be compatible with current versions of the definition file. If there is an incompatible change made to the definition file, the value of `version` must reflect that by incrementing it like so - 0, 1, 2, etcetera.

### Current Version
0

### Examples
```yaml
# version - Page config template version
version: "v0"
# unexpanded host uses default host info/config
host: "aws"
# unexpanded registrar uses default registrar info/config
registrar: "namecheap"
# unexpanded domain uses default registrar info/config
domain: "example.com"
# template - uniform resource locator where
# page template is located and accessible
# via git clone
template: "https://github.com/roymoran/index"
```

The example definition file below uses expanded keys to provide extra flexibility for use cases/scenarios that might come up.
```yaml
# version - Page config template version
version: "0"
# template - uniform resource locator where
# page template is located and accessible
# via git clone
template:
  url: "https://github.com/roymoran/page"

# expanded domain 
# name - name of domain for this page
# name - domain name with tld -> example.com (domain may or may not exist on account)
# registrar -  domain name registrar
# username/password - credentials for provided domain name registrar
domain:
  name: 'example.com'
  registrar: 'namecheap'
  username: username
  password: password
  token: token

# expanded host
# host provider where page is
# deployed
host:
  name: 'page'
  username: 'username'
  password: 'password'

```

## Architectural Decision Records
Significant software design decisions that are open to dicussion or have already been decided on. Each is formatted as question with **embolded answer** if it has been decided on and additional bullets for elaboration of decision.

Definition file format in **Yaml** or JSON?
- Allowance of comments so that its easy to include instructions for each key/value. Comments also allow us to link to external site for further instruction.
- Readablity for both technical and non-technical users, the less syntax the better.

Best approach for setting up infrastructure for initial configuration? Terraform, Pulumi, or directly using cloud SDKs? This question applies to self-hosted options (deploying on platforms that require standing up your own infrastructure) as opposed to hosting on platforms like GitHub Pages where they offer hosting as a service.
- TBD