[![Codacy Badge](https://app.codacy.com/project/badge/Grade/d9b43cca56244010875a13bf8d5a81fa)](https://www.codacy.com?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=SecurityGeekIO/terraform-provider-zpa&amp;utm_campaign=Badge_Grade)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://img.shields.io/badge/license-MIT-blue.svg)
[![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)

Terraform Provider for ☁️Zscaler Private Access☁️
=========================================================================

⚠️  **Attention:** This provider is not affiliated with, nor supported by Zscaler in any way.


- Website: https://www.terraform.io
- Documentation: https://help.zscaler.com/zpa
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)
<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

-	Install [Terraform](https://www.terraform.io/downloads.html) 0.12.x/0.13.x/0.14.x/0.15.x (0.11.x or lower is incompatible)
-	Install [Go](https://golang.org/doc/install) 1.16+ (This will be used to build the provider plugin.)
-	Create a directory, go, follow this [doc](https://github.com/golang/go/wiki/SettingGOPATH) to edit ~/.bash_profile to setup the GOPATH environment variable)

Building The Provider (Terraform v0.12+)
---------------------

Clone repository to: `$GOPATH/src/github.com/SecurityGeekIO/terraform-provider-zpa`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers
$ cd $GOPATH/src/github.com/terraform-providers
$ git clone https://github.com/SecurityGeekIO/terraform-provider-zpa.git
```

To clone on windows
```sh
mkdir %GOPATH%\src\github.com\terraform-providers
cd %GOPATH%\src\github.com\terraform-providers
git clone https://github.com/SecurityGeekIO/terraform-provider-zpa.git
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-zpa
$ make fmt
$ make build
```

To build on Windows
```sh
cd %GOPATH%\src\github.com\terraform-providers\terraform-provider-zpa
go fmt
go install
```

Building The Provider (Terraform v0.13+)
-----------------------

### MacOS / Linux
Run the following command:
```sh
$ make build13
```

### Windows
Run the following commands for cmd:
```sh
cd %GOPATH%\src\github.com\terraform-providers\terraform-provider-zpa
go fmt
go install
xcopy "%GOPATH%\bin\terraform-provider-zpa.exe" "%APPDATA%\terraform.d\plugins\zscaler.com\zpa\zpa\1.0.0\windows_amd64\" /Y
```
Run the following commands if using powershell:
```sh
cd "$env:GOPATH\src\github.com\SecurityGeekIO\terraform-provider-zpa"
go fmt
go install
xcopy "$env:GOPATH\bin\terraform-provider-zpa.exe" "$env:APPDATA\terraform.d\plugins\zscaler.com\zpa\zpa\1.0.0\windows_amd64\" /Y
```

**Note**: For contributions created from forks, the repository should still be cloned under the `$GOPATH/src/github.com/terraform-providers/terraform-provider-zpa` directory to allow the provided `make` commands to properly run, build, and test this project.

Using Zscaler Private Access Provider (Terraform v0.12+)
-----------------------

Activate the provider by adding the following to `~/.terraformrc` on Linux/Unix.
```sh
providers {
  "zpa" = "$GOPATH/bin/terraform-provider-zpa"
}
```
For Windows, the file should be at '%APPDATA%\terraform.rc'. Do not change $GOPATH to %GOPATH%.

In Windows, for terraform 0.11.8 and lower use the above text.

In Windows, for terraform 0.11.9 and higher use the following at '%APPDATA%\terraform.rc'
```sh
providers {
  "zpa" = "$GOPATH/bin/terraform-provider-zpa.exe"
}
```

If the rc file is not present, it should be created

Using Zscaler Private Access Provider (Terraform v0.13+)
-----------------------

For Terraform v0.13+, to use a locally built version of a provider you must add the following snippet to every module
that you want to use the provider in.

```hcl
terraform {
  required_providers {
    zpa = {
      source  = "zscaler.com/zpa/zpa"
      version = "1.0.0"
    }
  }
}
```

Examples
--------

Visit [here](https://github.com/SecurityGeekIO/terraform-provider-zpa/tree/master/website/docs/) for the complete documentation for all resources on github.

Issues
=========
Please feel free to open an issue using [Github Issues](https://github.com/SecurityGeekIO/terraform-provider-zpa/issues) if you run into any problems using this ZPA Terraform provider.


License
=========
MIT License

Copyright (c) 2021 [William Guilherme](https://github.com/willguibr)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
