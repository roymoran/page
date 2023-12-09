## [v0.1.0-alpha.12] - 2023-12-09

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
