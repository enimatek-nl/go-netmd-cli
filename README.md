# go-netmd cli

This is a reference implementation of the [go-netmd-lib](https://github.com/enimatek-nl/go-netmd-lib) as a command-line interface.

# usage
`> netmd-cli help`
```bash
netmd-cli NetMD command line interface.

Author:
  github.com/enimatek-nl
Version:
	0.0.1b
Usage:
  netmd-cli [options] command [arguments...]

Commands:
  list                     List all track data on the disc.
  send [wav] [title]       Send stereo pcm data to the disc.
  title [title]            Rename the disc title.
  rename [number] [title]  Rename the track number.
  move [number] [to]       Move the track number around.
  erase [number]           Erase track number from disc.
Options:
  -v           Verbose logging output.
  -y           Skip confirm questions.
  -i [index]   Set the NetMD usb device index when multiple
               devices are connected. [default: 0]
```

# download

Choose the platform below for which you would like to download the cli and put all files in the same directory.

In case **usblib-1.0** is already installed on the system you will probably not need the extra download.

### windows

[netmd-cli-windows-amd64.exe](https://github.com/enimatek-nl/go-netmd-cli/releases/download/builds/netmd-cli-windows-amd64.exe)

[libusb-1.0.dll](https://github.com/enimatek-nl/go-netmd-cli/releases/download/builds/libusb-1.0.dll)

### macos

### linux