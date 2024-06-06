# canu <!-- omit in toc -->

CLI to switch aws profiles

- [What does canu mean?](#what-does-canu-mean)
- [How to use canu? \[WIP\]](#how-to-use-canu-wip)

## What does canu mean?

The name is a tribute to a friend. It's probably not the best name for an `aws profile switcher` but... well, it means something to me 😊

I do this for fun and on my free time, so if you have any issues or suggestions, please add an issue to the repository, and if I have time I'll try to do it, but keep in mind I can be months without checking the repository again! So feel free to fork the repository if you want and do all the changes you want/need!

## How to use canu? [WIP]

The first thing you need to know is: **don't use `go install github.com/belitre/canu`, that would add `canu` as an executable to your `$GOPATH/bin` with name `canu`, and we don't want that!**

* Download the binary for your OS/ARCH from the release page and unpack it in a folder available in your `$PATH` (I'll use for the example `$GOPATH/bin`):
  ```
  curl -sL https://github.com/belitre/canu/releases/download/1.0.0/canu-1.0.0-darwin-arm64.tar.gz | tar -zxf - -C $GOPATH/bin
 
  ```
* Run `_canu install` (more details about flags incoming!)
* Restart your shell, or run the alias returned in the `_canu install` output.
* Run `canu` and enjoy!