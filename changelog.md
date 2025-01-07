## [v0.3.3] - Small Changes

Use your terminal to install this release for Linux or Windows:

```bash
# for macos/linux
$ curl -sSL https://raw.githubusercontent.com/roymoran/page/main/install/install.sh | sh
# for windows powershell
$ Invoke-Expression (New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/roymoran/page/main/install/install.ps1')
```

On macOS download and run the .pkg installer.

Changes:

* Improve analytics reporting
* Add logging service
* Print out commands in consistent order

## [v0.3.2] - Forgive Me for I Have Sinned...

Use your terminal to install this release for Linux or Windows:

```bash
# for macos/linux
$ curl -sSL https://raw.githubusercontent.com/roymoran/page/main/install/install.sh | sh
# for windows powershell
$ Invoke-Expression (New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/roymoran/page/main/install/install.ps1')
```

On macOS download and run the .pkg installer.

Changes:

* Avoid CNAME record at apex for Namecheap, this causes failure when attempting to validate domain ownership via DNS challenge.

## [v0.3.1] - Better. Faster. Stronger...

Use your terminal to install this release for Linux or Windows:

```bash
# for macos/linux
$ curl -sSL https://raw.githubusercontent.com/roymoran/page/main/install/install.sh | sh
# for windows powershell
$ Invoke-Expression (New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/roymoran/page/main/install/install.ps1')
```

On macOS download and run the .pkg installer.

Changes:

* Revert init.go
* Validate Host and Registrar valid

## [v0.3.0] - Another day, another improvement...

Changes:

* Introduce PAGE_CLI_TEST flag for better testing
* Change `template` key in definition file to `files`
* Remove payment link output

## [v0.1.2] - Another day, another bug...

Install this release:

```bash
# for macos/linux
$ curl -sSL https://raw.githubusercontent.com/roymoran/page/main/install/install.sh | sh
# for windows powershell
$ Invoke-Expression (New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/roymoran/page/main/install/install.ps1')
```

Changes:

* Fix errors in ProccessDefinitionFile and add templateFieldErrors helper function to surface template file error messages

## [v0.1.1-alpha.14] - Another bug bites the dust...

Install this release:

```bash
# for macos/linux
$ curl -sSL https://raw.githubusercontent.com/roymoran/page/main/install/install.sh | sh
# for windows powershell
$ Invoke-Expression (New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/roymoran/page/main/install/install.ps1')
```

Changes:

* Fix error when copying over files from a template path

## [v0.1.0-alpha.14] - No need to be alarmed, just small fixes...

Install this release:

```bash
# for macos/linux
$ curl -sSL https://raw.githubusercontent.com/roymoran/page/main/install/install.sh | sh
# for windows powershell
$ Invoke-Expression (New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/roymoran/page/main/install/install.ps1')
```

Changes:

* Fix error when deploying more than 1 webpage under aws host.
* Move webpage templates to [https://gitlab.com/page-templates](https://gitlab.com/page-templates)

## [v0.1.0-alpha.13] - Now with AI powers...

Install this release:

```bash
# for macos/linux
$ curl -sSL https://raw.githubusercontent.com/roymoran/page/main/install/install.sh | sh
# for windows powershell
$ Invoke-Expression (New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/roymoran/page/main/install/install.ps1')
```

Changes:

* run `page build` to generate a webpage for you. Use prompts such as "Create a webpage with a headline 'Hello, World!' centered both vertically and horizontally." Your page is generated in the current directory and can be used as the template in `page.yml` - just provide the path to the generated webpage. See this [video demo](https://youtu.be/kgzQIeom6g8)

Example `page.yml` with path to generate webpage.

```yaml
.
.
.
template: "/Users/roymoran/Documents/tempsite-xyhm"
```

## [v0.1.0-alpha.12] - Alpha 12 has landed...

Install this release:

```bash
# for macos/linux
$ curl -sSL https://raw.githubusercontent.com/roymoran/page/main/install/install.sh | sh
# for windows powershell
$ Invoke-Expression (New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/roymoran/page/main/install/install.ps1')
```

Changes since last release.

* Now run `page up` for certificate renewals
* Updated the `page up` command to include more descriptive and user-friendly output sequences for each deployment step, including visual animations and completion indicators for Host, Certificate, Domain, and Website Files.
* Enhanced `page.yml` to accept local file paths for the `Template` property, for greater flexibility in template sourcing. Also include improved parsing of the domain.
* Now install page cli with single command `curl -sSL https://raw.githubusercontent.com/roymoran/page/main/install/install.sh | sh` for macOS/Linux `Invoke-Expression (New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/roymoran/page/main/install/install.ps1')` on Windows Powershell
* Introducing [short videos](https://youtube.com/playlist?list=PLSqMEKs-lT4qVtG7-jSJj9_ZvsUMZTxAJ) to guide in configuring tool and deploy your first website.
* Introducing a [templates directory](./templates/).
* Add support for additional 64-bit ARM systems including Apple Silicon and Windows
