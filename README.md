[![Test](https://github.com/zscaler/terraform-provider-zpa/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/zscaler/terraform-provider-zpa/actions/workflows/test.yml)
[![Release](https://github.com/zscaler/terraform-provider-zpa/actions/workflows/release.yml/badge.svg?branch=master)](https://github.com/zscaler/terraform-provider-zpa/actions/workflows/release.yml)
[![Zscaler Community](https://img.shields.io/badge/zscaler-community-blue)](https://community.zscaler.com/)
[![Slack](https://img.shields.io/badge/Join%20Our%20Community-Slack-blue)](https://forms.gle/3iMJvVmJDvmUy36q9)

<a href="https://terraform.io">
    <img src="https://raw.githubusercontent.com/hashicorp/terraform-website/master/public/img/logo-text.svg" alt="Terraform logo" title="Terraform" height="50" />
</a>

<a href="https://www.zscaler.com/">
    <img src="https://www.zscaler.com/themes/custom/zscaler/logo.svg" alt="Zscaler logo" title="Zscaler" height="50" />
</a>

Terraform Provider for ☁️Zscaler Private Access☁️
=========================================================================

- Website: [https://www.terraform.io](https://registry.terraform.io/providers/zscaler/zpa/latest)
- Documentation: https://help.zscaler.com/zpa
- Zscaler Community: [Zscaler Community](https://community.zscaler.com/)

Support Disclaimer
-------
!> **Disclaimer:** This Terraform provider is community supported. Although this provider is supported by Zscaler employees, it is **NOT** supported by Zscaler support. Please open all enhancement requests and issues on [Github Issues](https://github.com/zscaler/terraform-provider-zpa/issues) for support.

Requirements
------------

-	Install [Terraform](https://www.terraform.io/downloads.html) 0.12.x/0.13.x/0.14.x/0.15.x (0.11.x or lower is incompatible)
-	Install [Go](https://golang.org/doc/install) 1.16+ (This will be used to build the provider plugin.)
-	Create a directory, go, follow this [doc](https://github.com/golang/go/wiki/SettingGOPATH) to edit ~/.bash_profile to setup the GOPATH environment variable)

Building The Provider (Terraform v0.12+)
---------------------

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-zpa`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers
$ cd $GOPATH/src/github.com/terraform-providers
$ git clone https://github.com/terraform-providers/terraform-provider-zpa.git
```

To clone on windows
```sh
mkdir %GOPATH%\src\github.com\terraform-providers
cd %GOPATH%\src\github.com\terraform-providers
git clone https://github.com/zscaler/terraform-provider-zpa.git
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
xcopy "%GOPATH%\bin\terraform-provider-zpa.exe" "%APPDATA%\terraform.d\plugins\zscaler.com\zpa\zpa\2.0.5\windows_amd64\" /Y
```
Run the following commands if using powershell:
```sh
cd "$env:GOPATH\src\github.com\terraform-providers\terraform-provider-zpa"
go fmt
go install
xcopy "$env:GOPATH\bin\terraform-provider-zpa.exe" "$env:APPDATA\terraform.d\plugins\zscaler.com\zpa\zpa\2.0.5\windows_amd64\" /Y
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
      version = "2.0.5"
    }
  }
}
```

Examples
--------

Visit [here](https://github.com/zscaler/terraform-provider-zpa/tree/master/docs) for the complete documentation for all resources on github.

Examples [here] (https://github.com/zscaler/terraform-provider-zpa/tree/master/examples) for the complete list of examples on github.

Issues
=========

Please feel free to open an issue using [Github Issues](https://github.com/zscaler/terraform-provider-zpa/issues) if you run into any problems using this ZPA Terraform provider.

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.16+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-zpa
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

License
=========

MIT License

=======

Copyright (c) 2022 [Zscaler](https://github.com/zscaler)

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
