# Ruckus Mock SSH Server

Mock SSH server written in Go that mocks a subset of Ruckus SSH commands and responses.
Used for testing of [gabe565/pyruckus](https://github.com/gabe565/pyruckus) without a physical device.

## Quick Start

### Install

```shell
go install github.com/gabe565/ruckus-mock-ssh@latest
```

### Run

You can run the mock SSH server with the default configuration with:

```shell
ruckus-mock-ssh
```

To see the available flags, run

```shell
ruckus-mock-ssh --help
```

### Connect via SSH

After running the server, connect to it with the following SSH command:

```shell
ssh localhost -p 2222 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null
```

### Shell Completions

Simple shell completions are set up for command discovery. Once you connect
via SSH, press the Tab key to see available commands.

## Development

### Adding new command responses

Response mocks are hardcoded text files in [`cli/responses`](cli/responses).
To add a new response, add a text file with the command name replaced with underscores.
For example, the response for `show config` is in [`cli/responses/show_config.txt`](cli/responses/show_config.txt).

Command line reference requires a Ruckus account, but is available for free after registration.
[The command line reference for Ruckus Unleashed 200.11 can be found here](https://support.ruckuswireless.com/documents/3946-ruckus-unleashed-200-11-ga-cli-reference-guide).
