# bootstrap

bootstrap is intended to automatically download the correct binary from a GitHub repository.

## Quick start
1. Download the latest archive from [here](https://github.com/gearboxworks/bootstrap/releases/latest)
2. [Untar](https://www.howtogeek.com/362203/what-is-a-tar.gz-file-and-how-do-i-open-it/) the download file.
3. Place `bootstrap` executable binary in a directory located in your PATH.<br>_(e.g. `/usr/local/bin` on macOS and Linux)._
4. Execute `bootstrap install` - This will create placeholder symlinks for the default available commands.


## Usage: bootstrap
```
Usage:
	bootstrap [command] <args>

Where [command] is one of:
	help	- Help about any command
	selfupdate	- bootstrap - Update version of executable.
	version	- bootstrap - Self-manage executable.

Use bootstrap help [command] for more information about a command.
```


## Usage: bootstrap version
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

Download latest buildtools binary from repo.

```
./bootstrap --bin buildtools version update
```

