# Pleasant CLI

Pleasant CLI is a simple CLI that lets you interact with Pleasant Password Server instances through the command line. It uses the [Pleasant Password Server API](https://pleasantpasswords.com/info/pleasant-password-server/m-programmatic-access/restful-api) to communicate with the server.

It should work with any instance that supports API version 5 (`/api/v5/`).

## Installation

Go to [Releases](https://github.com/marevers/pleasant-cli/releases) and download the latest release for your platform.

Extract the archive to a location in your `$PATH`, e.g. on Ubuntu: `/usr/local/bin`.

## Getting started

First, you need to set the server URL for the CLI to connect to.
It must be specified as `<PROTOCOL>://<URL>:<PORT>`.

```
$ pleasant-cli config serverurl <SERVER URL>
```

The server URL will be saved to the configuration file (default: `$HOME/.pleasant-cli.yaml`).
If you want to use a different config path, add the flag `--config <PATH>` to all commands.

Next, log into the Pleasant Password Server:

```
$ pleasant-cli login
```

**Note**: this will perform an interactive login. You can also pass your credentials as flags.

```
$ pleasant-cli login --username <USERNAME> --password <PASSWORD>
```

This will retrieve an access token and save it to a file (default: `$HOME/.pleasant-token.yaml`) for subsequent commands.
If you want to save/use the token in a different location, add the flag `--token <PATH>` to all commands.

## Commands

In order to view all available commands, run the CLI without any arguments.

```
$ pleasant-cli

pleasant-cli is an easy to use CLI that uses the Pleasant Password Server
API to interact with a Pleasant Password Server instance.

To use pleasant-cli, you must first set your server URL by running the following command:
pleasant-cli config serverurl <SERVER URL>

You can then log in by running:
pleasant-cli login

Usage:
  pleasant-cli [command]

Available Commands:
  apply       Applies a configuration to entries or folders
  completion  Generate the autocompletion script for the specified shell
  config      Interact with pleasant-cli configuration
  create      Creates entries or folders
  delete      Archives or deletes entries or folders or user access assignments for them
  get         Gets entries, folders, access levels, server info or password strength
  help        Help about any command
  login       Log in to Pleasant Password Server
  patch       Partially updates entries or folders or adds user access assignments for them
  search      Search for entries and folders matching a query

Flags:
      --config string   config file (default is $HOME/.pleasant-cli.yaml)
  -h, --help            help for pleasant-cli
  -t, --toggle          Help message for toggle
      --token string    token file (default is $HOME/.pleasant-token.yaml)
  -v, --version         version for pleasant-cli

Use "pleasant-cli [command] --help" for more information about a command.
```

To view information on (sub)commands, append `--help` to any command.

```
$ pleasant-cli login --help

Log into Pleasant Password Server with username and password.
Username and password can either be entered interactively or by using flags.

Examples:
pleasant-cli login
pleasant-cli login --username <USERNAME> --password <PASSWORD>

Usage:
  pleasant-cli login [flags]

Flags:
  -h, --help              help for login
  -p, --password string   Password for Pleasant Password Server
  -u, --username string   Username for Pleasant Password Server

Global Flags:
      --config string   config file (default is $HOME/.pleasant-cli.yaml)
      --token string    token file (default is $HOME/.pleasant-token.yaml)
```

## The --path flag

### How does it work?

Pleasant CLI provides a `--path` flag with many of the commands.
This flag allows more human-friendly interaction with the API as it provides a path-based interaction with the server.

For example, an entry can be retrieved by its path in the folder structure like such:

```
pleasant-cli get entry --path Root/Private\ Folders/MyUser/MyEntry
```

Note that certain characters must be escaped (like spaces) and valid paths must start with `Root/`.

### Path autocompletion

In order to make the server a bit more browseable without knowing paths in advance, the flag also supports autocompletion.

There are three main cases of autocompletion behavior:

1. Nothing typed yet
   
   Pleasant CLI autocompletes `Root/` immediately to get the path started.

   Example:

   ```
   $ pleasant-cli get entry --path <TAB> <TAB>
   pleasant-cli get entry --path Root/
   ```
   
2. Path ends with anything but a slash

   Pleasant CLI retrieves all matching entries and/or folders from the parent folder and looks for matching entries.

   Example:

   ```
   $ pleasant-cli get entry --path Root/MyFolder/My <TAB> <TAB>
   Root/MyFolder/MySubfolder/ --- folder
   Root/MyFolder/MyEntry1     --- entry
   Root/MyFolder/MyEntry2     --- entry
   ```

3. Paths ends with a slash

   Pleasant CLI assumes the last past segment is a folder and retrieves all entries and/or folders from it.

   Example:

   ```
   $ pleasant-cli get entry --path Root/MyFolder/ <TAB> <TAB>
   Root/MyFolder/MySubfolder/   --- folder
   Root/MyFolder/AnotherFolder/ --- folder
   Root/MyFolder/MyEntry1       --- entry
   Root/MyFolder/MyEntry2       --- entry
   ```

As cases 2 and 3 require a connection with the server to retrieve the available options, completion can take a while depending on how fast the server reacts.

For `entry` commands, both entries and folders are returned. For `folder` commands, only folders are returned.`

## Clipboard functionality

If you are retrieving a username or password of an entry, it is possible to have Pleasant CLI copy the output directly to your clipboard. This is supported on the following platforms:

- macOS
- Windows
- Linux - needs X11 dev package, install either `libx11-dev`, `xorg-dev` or `libX11-devel` to make it work

It is NOT supported on:

- WSL2

Examples:

```
$ pleasant-cli get entry --path Root/MyFolder/MyEntry --username --clip
$ pleasant-cli get entry --path Root/MyFolder/MyEntry --password --clip
```

## Docker image

A minimal Alpine-based Docker image is provided here:

**Latest version**

```
docker pull ghcr.io/marevers/pleasant-cli:latest
```

**Specific version**

```
docker pull ghcr.io/marevers/pleasant-cli:<version>
```

## Build

### Go

To build with Go, you must have Go installed (version 1.24).

```bash
go build
```

This creates the `pleasant-cli` executable in the current location.

It is also possible to run certain commands directly without first building. For example:

```bash
go run . get entry --id <id>
```

### Nix

To build the project as a *Nix flake* (ensure you have *Nix* installed with flakes enabled):
```bash
nix build
```
This creates a `result` symlink pointing to the built package in your *Nix* store.

To run it directly:
```bash
nix run
```

To enter a development environment with all necessary dependencies:
```bash
nix develop
```

To install it in your current profile:
```bash
nix profile install
```

You may also want to update the locked dependencies in `flake.lock`:
```bash
nix flake update
```

As *Nix* uses hashes to ensure reproducibility between builds, you may need to update the vendor hash in `flake.nix` when Go dependencies change (for example after `go get -u`).

To do so, run `nix build` as usual: this will fail giving the expected new hash that you can use to update the `vendorHash` variable in `flake.nix`. Once updated, the rebuild should work as expected.

For more information about using flakes in *Nix* and *NixOS* environments please have a look at the [documentation](https://nixos.wiki/wiki/flakes).

## Known issues / Limitations

### Search functionality

Having certain character combinations in an entry / folder name can stop the search function or any command using the search API like `get entry --path <path>` or `get folder --path <path>` from working. Currently the only known character combination is a hyphen surrounded by spaces, like in `My - Entry`. The search API call can then run into a timeout. This is a limitation of the Pleasant Password Server API. In order to fix it, remove the hyphen from the name of the entry or folder. This issue may or may not occur, based on the Pleasant Password Server version running.

## Roadmap

Currently there are no specific planned features. Have a feature request? Please open an issue!

## License

Pleasant CLI is released under the Apache 2.0 license. See [LICENSE.txt](https://github.com/marevers/pleasant-cli/blob/master/LICENSE.txt)
