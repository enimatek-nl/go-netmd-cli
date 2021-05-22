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

Choose the platform below for which you would like to download netmd-cli.

Each download is continuously build by [github actions](https://github.com/enimatek-nl/go-netmd-cli/actions), when ARM based workflows are available I will add these as well.

### windows

Download Link: [ [netmd-cli-windows-amd64.exe](https://github.com/enimatek-nl/go-netmd-cli/releases/download/builds/netmd-cli-windows-amd64.exe) ]

### linux

Download Link: [ [netmd-cli-linux-amd64](https://github.com/enimatek-nl/go-netmd-cli/releases/download/builds/netmd-cli-linux-amd64) ]

### macos

You will need [brew](https://brew.sh) and run `brew install libusb libusb-compat`

Download Link: [ [netmd-cli-macos-intel](https://github.com/enimatek-nl/go-netmd-cli/releases/download/builds/netmd-cli-macos-intel) ]


