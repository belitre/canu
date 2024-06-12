# canu <!-- omit in toc -->

CLI to switch aws profiles

I do this for fun on my free time, so if you have any issues or suggestions, please add an issue to the repository, and if I have time I'll try to do it, but keep in mind I can be months without checking the repository again! So feel free to fork the repository if you want and do all the changes you want/need!

- [What does canu mean?](#what-does-canu-mean)
- [How to use canu?](#how-to-use-canu)
  - [Installation](#installation)
    - [Custom installation](#custom-installation)
      - [Examples](#examples)
  - [Uninstall](#uninstall)

## What does canu mean?

The name is a tribute to a friend. It's probably not the best name for an `aws profile switcher` but... well, it means something to me 😊

## How to use canu?

The first thing you need to know is: **don't use `go install github.com/belitre/canu`, that would add `canu` as an executable to your `$GOPATH/bin` with name `canu`, and we don't want that!**

Also, `canu` will only read profiles from the aws profiles config file that match the pattern `[profile profile-name]` with `sso_start_url`, `sso_account_id` and `sso_role_name` defined. That means it will not show the `default` profile, or any other profile that is not using SSO. For example, if your config file looks like this:
```
[default]
region = us-east-1

[profile test]
region = us-east-2

[profile admin]
sso_start_url = https://d-93456h53g77.awsapps.com/start
sso_region = us-east-1
sso_account_id = 987654321
sso_role_name = MyRoleName
region = us-east-1
```
`canu` will only show the profile `admin` in the menu.

### Installation

* Download the binary for your OS/ARCH from the release page and unpack it in a folder available in your `$PATH` (I'll use for the example `$GOPATH/bin`):
  ```
  curl -sL https://github.com/belitre/canu/releases/download/1.1.0/canu-darwin-arm64.tar.gz | tar -zxf - -C $GOPATH/bin
  ```
  I recommend to download the `_canu` binary in a folder where your user has permissions to write. Of course, as a user, you can download the binary to a folder that requires root permissions, but then you will need to run the `install` command with root permissions as well, and unless you are using `sudo`, which should keep the correct value for your user in the `$HOME` env var, you may need to set the path where you want `canu` to write the alias, to be sure it writes it to the user shell config file, not the root one. More information about this in the next section.
* Run `_canu install`. Running this command, `canu` will install with the default settings. That means it will use as default config path for your aws config: `$HOME/.aws/config`, and it will use `$HOME/.bash_profile` to add the alias. If you have different settings, please check the next section to learn how to customise your `canu` installation.
* Restart your shell, or run the alias returned in the `_canu install` output.
* Run `canu` and enjoy!

What happens during the installation is:
* `canu` will create a shell script in the same folder the binary is downloaded with the name of the alias (by default this name is `canu`, but you can set a different name, check the next section to learn how).
* `canu` will add an alias to your `$HOME/.bash_profile` that will source the generated shell script, this is how `canu` sets the selected profile in the `AWS_PROFILE` environment variable.

#### Custom installation

Once you downloaded the binary, you can check the available flags running `_canu install -h`:

```
installs canu aliases

Usage:
  _canu install [flags]

Flags:
      --alias-name string     name of the alias (default "canu")
  -c, --config-path string    Path with the AWS config file with the defined profiles (default "/Users/miguel/.aws/config")
  -e, --exclude stringArray   Profiles containing any of the excludes will be ignored. Exclude takes preference over include (ignores upper/lower case)
  -h, --help                  help for install
  -i, --include stringArray   Profiles containing any of the includes will be available to select (ignores upper/lower case)
      --shell-config string   path to the shell config script (default "/Users/miguel/.bash_profile")
      --skip-alias            if provided it won't add an alias to the shell config script
  -s, --sort                  If provided it will sort the profiles, by default the profiles will use the same order you have in your aws config file
```

There are two different kind of flags here:

* Flags related with the installation:
  ```
  --alias-name string     name of the alias (default "canu")
  --shell-config string   path to the shell config script (default "/Users/miguel/.bash_profile")
  --skip-alias            if provided it won't add an alias to the shell config script
  ```
  * `--alias-name`: This is the name of the alias `canu` will create, also, it will be the name of the shell script generated to set the correct value for the `AWS_PROFILE` environment variable.
  * `--shell-config`: The path to your shell profile config file. By default `canu` will use `$HOME/.bash_profile`. If you use `zsh` you can run the install command with `--shell-config $HOME/.zshrc`
  * `--skip-alias`: If this flag is provided `canu` won't add the alias to your shell profile config file.
* The other flags will change the behaviour when running the generated alias:
  ```
  -c, --config-path string    Path with the AWS config file with the defined profiles (default "/Users/miguel/.aws/config")
  -e, --exclude stringArray   Profiles containing any of the excludes will be ignored. Exclude takes preference over include (ignores upper/lower case)
  -i, --include stringArray   Profiles containing any of the includes will be available to select (ignores upper/lower case)
  -s, --sort                  If provided it will sort the profiles, by default the profiles will use the same order you have in your aws config file
  ```
  * `-c, --config-path`: The path where the aws config for your profiles is located, by default it will use `$HOME/.aws/config`
  * `-e, --exclude`: Flag to exclude profiles, multiple values can be provided. `canu` will exclude the profiles containing any of the provided values with this flag. This flag takes preference over `--include`
  * `-i, --include`: Flag to include profiles, multiple values can be provided. `canu` will include only the profiles containing any of provided values with this flag.
  * `-s, --sort`: If this flag is provided `canu` will sort the profiles alphabetically before showing the menu. By default the order of the profiles will be the same one you have in your aws profiles config file.

##### Examples

* Install `canu` with alias `profile-switcher`, using as shell `zsh`:
  ```
  _canu install --shell-config /Users/my-user/.zshrc --alias-name profile-switcher
  ```
* Install `canu` with default values for installation, excluding profiles with the string `test` and `dev`, sorting the list and using a different path for the aws profiles config file while running:
  ```
  _canu install -e test -e dev -s -c /Users/my-user/.aws/custom-config
  ```
* Install `canu` with all custom settings:
  ```
  _canu install --alias-name profile-switcher --shell-config /Users/my-user/.zshrc -c /Users/my-user/.aws/custom-config -e test -e dev -i account1 -i account2 -s
  ```
  This will install `canu` with the alias `profile-switcher` using the `zsh` shell config file to create the alias, excluding profiles containing `test` and `dev`, and including profiles containing `account1` and `account2`, and sorting the profiles before showing the menu.
* Install multiple aliases for different profiles:
  ```
  _canu install --alias-name canu-dev -i dev
  _canu install --alias-name canu-prod -i prod
  ```
  This will create two aliases, one called `canu-dev` that will include only profiles containing `dev`, and another one called `canu-prod` that will include only profiles containing `prod`.

### Uninstall

The uninstall process will delete the script created by `_canu install` and remove the alias added if it's found in the file.

The available flags for uninstall are:
```
./bin/_canu uninstall -h
uninstalls canu aliases

Usage:
  _canu uninstall [flags]

Flags:
      --alias-name string     name of the alias (default "canu")
  -h, --help                  help for uninstall
      --shell-config string   path to the shell config script (default "/Users/miguel/.bash_profile")
```

If you used the default installation, you can just run:
```
_canu uninstall
```

In case you used a different alias or a different shell config file to add the alias you can run:
```
_canu uninstall --alias-name profile-switcher --shell-config /Users/my-user/.zshrc
```