# canu <!-- omit in toc -->

CLI to switch aws profiles

- [What does canu mean?](#what-does-canu-mean)
- [How to use canu? \[WIP\]](#how-to-use-canu-wip)

## What does canu mean?

The name is a tribute to a friend. It's probably not the best name for an `aws profile switcher` but... well, it means something to me 😊

## How to use canu? [WIP]

The first thing you need to know is: **don't use `go install github.com/belitre/canu`, that would add `canu` as an executable to your `$GOPATH/bin` with name `canu`, and we don't want that!**

* Download the binary for your OS/ARCH from the release page and unpack it in a folder available in your `$PATH`
* Run `_canu install` (more details about flags incoming!)
* Restart your shell, or run the alias returned in the `_canu install` output.
* Use your alias and enjoy!