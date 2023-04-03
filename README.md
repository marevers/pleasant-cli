# Pleasant CLI

Pleasant CLI is a simple CLI that lets you interact with Pleasant Password Server instances through the command line. It uses the [Pleasant Password Server API](https://pleasantpasswords.com/info/pleasant-password-server/m-programmatic-access/restful-api) to communicate with the server.

It should work with any instance that supports API version 5 (`/api/v5/`).

## Installation

TO DO

## Getting started

First, you need to set the server URL for the CLI to connect to.

```bash
pleasant-cli config serverurl <SERVER URL>
```

The server URL will be saved to the configuration file (default: $HOME/.pleasant-cli.yaml). If you want to use a different config path, you can add the flag `--config <PATH>` to any command.

Next, log into the Pleasant Password Server:

```bash
pleasant-cli login
```

**Note**: this will perform an interactive login. You can also pass your credentials as flags.

```
pleasant-cli login --username <USERNAME> --password <PASSWORD>
```

This will retrieve an access token and save it to the configuration file for subsequent commands.

## Commands

In order to view all available commands, run the CLI without any arguments.

```bash
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
  completion  Generate the autocompletion script for the specified shell
  config      Interact with pleasant-cli configuration
  create      Creates entries or folders
  get         Gets entries, folders or access levels
  help        Help about any command
  login       Log in to Pleasant Password Server
  search      Search for entries and folders matching a query

Flags:
      --config string   config file (default is $HOME/.pleasant-cli.yaml)
  -h, --help            help for pleasant-cli
  -t, --toggle          Help message for toggle

Use "pleasant-cli [command] --help" for more information about a command.
```

## License

Pleasant CLI is released under the Apache 2.0 license. See [LICENSE.txt](https://github.com/marevers/pleasant-cli/blob/master/LICENSE.txt)
