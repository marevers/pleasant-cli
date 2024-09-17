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
pleasant-cli config serverurl <SERVER URL>
```

The server URL will be saved to the configuration file (default: `$HOME/.pleasant-cli.yaml`).
If you want to use a different config path, add the flag `--config <PATH>` to all commands.

Next, log into the Pleasant Password Server:

```
pleasant-cli login
```

**Note**: this will perform an interactive login. You can also pass your credentials as flags.

```
pleasant-cli login --username <USERNAME> --password <PASSWORD>
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

## License

Pleasant CLI is released under the Apache 2.0 license. See [LICENSE.txt](https://github.com/marevers/pleasant-cli/blob/master/LICENSE.txt)
