# Docs

## Page Definition File
The goal of the definition file is to describe only the essential information required to create and deploy a page with the cli. The simplest form of the definition file is what is generated when running `page init`. We should use smart defaults for any key that isn't included or isn't expanded on. Any modifications to the definition file should aim to be compatible with current versions of the definition file. If there is an incompatible change made to the definition file, the value of `version` must reflect that by incrementing it like so - 0, 1, 2, etcetera.

### Current Version
0

### Examples
```yaml
# version - Page config template version
version: "0"
# template - uniform resource locator where page 
# template is located and accessible
template: "https://github.com/roymoran/page"
# unexpanded domain uses default registrar info/config
domain: 'example.com'
# unexpanded host uses default host info/config
host: 'page'
```

The example definition file below uses expanded keys to provide extra flexibility for use cases/scenarios that might come up.
```yaml
# version - Page config template version
version: "0"
# template - uniform resource locator where page 
# template is located and accessible
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