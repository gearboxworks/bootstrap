# bootstrap

bootstrap is intended to automatically download the correct binary from a GitHub repository.

```
Usage:
	bootstrap [command] <args>

Where [command] is one of:
	help	- Help about any command
	selfupdate	- bootstrap - Update version of executable.
	version	- bootstrap - Self-manage executable.

Use bootstrap help [command] for more information about a command.
```


## bootstrap version

```
Usage:
	bootstrap version [command] <args>

Where [command] is one of:
	example		- bootstrap - Self-manage help examples.
	install		- bootstrap - Install placeholder for all supported apps.
	selfupdate	- bootstrap - Update version of this executable.
	version		- bootstrap - Self-manage executables.
```

## Examples

Create symlink placeholders for all supported binaries. Once you execute, it'll install.

```
bootstrap install
```

Rename or symlink the bootstrap binary and automatically fetch and replace the symlink with the binary from repo.

```
ln -s bootstrap launch
./launch version update
```


Download latest buildtools binary from repo.

```
./bootstrap --bin buildtools version update
```

