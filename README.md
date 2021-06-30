# go-netmd cli

This is a reference implementation of the [go-netmd-lib](https://github.com/enimatek-nl/go-netmd-lib) as a command-line interface.

# Downloads

Choose the platform below for which you would like to download netmd-cli.

Each download is continuously build by [github actions](https://github.com/enimatek-nl/go-netmd-cli/actions), when ARM based workflows are available I will add these as well.

|Platform |Arch |Link |
--- | --- | ---
|Windows |amd64 |[netmd-cli-windows-amd64.exe](https://github.com/enimatek-nl/go-netmd-cli/releases/download/builds/netmd-cli-windows-amd64.exe) |
|Linux |amd64 |[netmd-cli-linux-amd64](https://github.com/enimatek-nl/go-netmd-cli/releases/download/builds/netmd-cli-linux-amd64)
|macOs _*_ |intel |[netmd-cli-macos-intel](https://github.com/enimatek-nl/go-netmd-cli/releases/download/builds/netmd-cli-macos-intel)
* For macOS you will need [brew](https://brew.sh) and run `brew install libusb`

# Usage
`> netmd-cli help`
```shell
Usage:
  netmd-cli [options] command [arguments...]

Commands:
  list                     List all track data on the disc.
  send [wav] [title]       Send stereo pcm data to the disc.
  title [title]            Rename the disc title.
  rename [number] [title]  Rename the track number.
  move [number] [to]       Move the track number around.
  erase [number]           Erase track number from disc.
  group [to] [title]       Put tracks ungrouped up to track number
                           into a group by the name title.
  degroup [number]         De-group including all tracks of the
                           group that are part of the track number.
Options:
  -v             Verbose logging output.
  -y             Skip confirm questions.
  -d [encoding]  Encoding on disk lp2, lp4 or sp (default: sp)
                 lp-modes known to work only with Sharp NetMD
  -i [index]     Set the NetMD usb device index when multiple
                 devices are connected. [default: 0]
```
# Record PCM Through NetMD
You can send raw PCM data in a WAV-container to the NetMD by default encoded in ATRAC (SP) on the device using the chip on the unit.

If you own a one of the units tested positive below you can use the `-d lp2` or `-d lp4` flag to save the PCM as ATRAC3 on the disc.

|Brand |Type |LP2/LP4 |
--- | --- | ---
|Sharp |IM-DR410 |✅ |
|Sharp |IM-DR420 |✅ |
|Sony |MDS-JB980 |✅ |
|Sony |MDS-JE780 |✅ |
|Sony |MZ-N710 |❌ |
|Sony |MZ-N510 |❌ |
|Sony |MZ-RH910 |❌ |
|Sony |MZ-NH600 |❌ |

The list is not yet complete, so if your device is not listed and if you hear only silence after the transfer was completed you device does not support ATRAC3 through NetMD (USB). You could always use the solution mentioned below to encode the PCM to ATRAC3 first before transfer.

In each case you will need to prepare your source (mp3, aac, flac etc.) the WAV yourself like so:
```shell
ffmpeg -i mytrack.flac -f wav -ar 44100 -ac 2 mytrack.wav
netmd-cli send mytrack.wav
```
It's recommended to use [ffmpeg](https://ffmpeg.org).

# Use LP2 Without NetMD Encoding Support
You will need to create the LP2 files yourself on the host machine.
For this you can use [atracdenc](https://github.com/dcherednik/atracdenc) created by Daniil Cherednik and put the LP2 encoded track into a WAV-container.

### example
First convert your mp3, flac etc. into a stereo wav file:
```shell
ffmpeg -i mytrack.flac -f wav -ar 44100 -ac 2 out.wav
```
Now encode the wav into an ATRAC3 file with a bitrate of 128 for LP2:
```shell
atracdenc -e atrac3 -i out.wav -o out.aea --bitrate 128
```
Put the AEA file into a WAV container like so:
```shell
ffmpeg -i out.aea -f wav -c:a copy mytrack_lp2.wav
```
And send the wav to the NetMD:
```shell
netmd-cli send mytrack_lp2.wav
```
