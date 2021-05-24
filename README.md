# go-netmd cli

This is a reference implementation of the [go-netmd-lib](https://github.com/enimatek-nl/go-netmd-lib) as a command-line interface.

# Downloads

Choose the platform below for which you would like to download netmd-cli.

Each download is continuously build by [github actions](https://github.com/enimatek-nl/go-netmd-cli/actions), when ARM based workflows are available I will add these as well.

### Windows

Download Link: [ [netmd-cli-windows-amd64.exe](https://github.com/enimatek-nl/go-netmd-cli/releases/download/builds/netmd-cli-windows-amd64.exe) ]

### Linux

Download Link: [ [netmd-cli-linux-amd64](https://github.com/enimatek-nl/go-netmd-cli/releases/download/builds/netmd-cli-linux-amd64) ]

### macOS

You will need [brew](https://brew.sh) and run `brew install libusb libusb-compat`

Download Link: [ [netmd-cli-macos-intel](https://github.com/enimatek-nl/go-netmd-cli/releases/download/builds/netmd-cli-macos-intel) ]

# Usage
`> netmd-cli help`
```shell
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
# Convert PCM to ATRAC SP
You can send raw PCM data in a WAV-container to the NetMD. This will be encoded to ATRAC on the device using the chip on the unit (eg. using Type-R etc.).
You will need to prepare your source (mp3, aac, flac etc.) the WAV yourself like so:
```shell
ffmpeg -i mytrack.flac -ac 2 -ar 44100 -f s16be mytrack.wav
```

# Send ATRAC3 (LP2)
It's possible to send LP2 tracks to the NetMD. But you will need to create them yourself on the host machine. For this you can use [atracdenc](https://github.com/dcherednik/atracdenc) created by Daniil Cherednik.
You will need to put the ATRAC3 encoded track into a WAV-container. For all these steps it's recommended to use [ffmpeg](https://ffmpeg.org).

### example
First convert your mp3, flac etc. into a stereo wav file:
```shell
ffmpeg -i mytrack.flac -f wav -ar 44100 -ac 2 out.wav
```
Now encode the wav into an ATRAC3 file with a bitrate of 128 for LP2:
```shell
atracdenc -e atrac3 -i out.wav -o out.aea --bitrate 128
```
Last step is putting the AEA file into a WAV container like so:
```shell
ffmpeg -i out.aea -f wav -c:a copy mytrack_lp2.wav
```
