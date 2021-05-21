package main

import (
	"fmt"
	"github.com/enimatek-nl/go-netmd-cli/cli"
	"github.com/enimatek-nl/go-netmd-lib"
	"log"
	"os"
	"strings"
)

const (
	version = "0.0.1b"
)

func main() {

	if len(os.Args) <= 1 {
		log.Fatal("not enough parameters, try 'help'")
	}

	cmd := "help"
	ptr := 1
	debug := false
	safe := true
	index := 0

	// check options
	for ; ptr < len(os.Args); ptr++ {
		a := os.Args[ptr]
		if strings.HasPrefix(a, "-") {
			if strings.HasSuffix(a, "v") {
				debug = true
			}
			if strings.HasSuffix(a, "y") {
				safe = false
			}
			if strings.HasSuffix(a, "i") {
				ptr++
				if len(os.Args) <= ptr {
					log.Fatal("missing index #")
				}
				i, err := cli.ToInt(os.Args[ptr])
				if err != nil {
					log.Fatal(err)
				}
				index = i
			}
		} else {
			break
		}
	}

	if len(os.Args) > ptr {
		cmd = os.Args[ptr]
	}

	md, err := netmd.NewNetMD(index, debug)
	if cmd != "help" { // only if 'help' is requested skip device errors
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		defer md.Close()
	}

	switch cmd {
	case "send":
		ptr++
		if len(os.Args) <= ptr {
			log.Fatal("missing filename")
		}
		fn := os.Args[ptr]
		t := fn // default to filename if no title is given
		ptr++
		if len(os.Args) > ptr {
			t = strings.Join(os.Args[ptr:], " ")
		}
		send(md, fn, t)
	case "erase":
		ptr++
		if len(os.Args) <= ptr {
			log.Fatal("missing track #")
		}
		trk, err := cli.ToInt(os.Args[ptr])
		if err != nil {
			log.Fatal(err)
		}
		erase(md, trk, safe)
	case "list":
		list(md)
	default: // help
		help()
	}

}

func help() {
	fmt.Println("")
	fmt.Println("netmd-cli NetMD command line interface.")
	fmt.Println("")
	fmt.Println("Author:")
	fmt.Println("  github.com/enimatek-nl")
	fmt.Println("Version:")
	fmt.Printf("	%s\n", version)
	fmt.Println("Usage:")
	fmt.Println("  netmd-cli [options] command [arguments...]")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  list                List all track data on the disc.")
	fmt.Println("  erase [number]      Erase track number from disc.")
	fmt.Println("  send [wav] [title]  Send stereo pcm data to the disc.")
	fmt.Println("Options:")
	fmt.Println("  -v           Verbose logging output.")
	fmt.Println("  -y           Skip confirm questions.")
	fmt.Println("  -i [index]   Set the NetMD usb device index when multiple")
	fmt.Println("               devices are connected. [default: 0]")
	fmt.Println("")
}

func send(md *netmd.NetMD, fn, t string) {
	track, err := md.NewTrack(t, fn, netmd.WfPCM, netmd.DfStereoSP)
	if err != nil {
		log.Fatal(err)
	}
	err = md.Send(track)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Track has been send.")
}

func list(md *netmd.NetMD) {
	fmt.Println("")
	_, total, available, _ := md.RequestDiscCapacity()
	fmt.Printf("Disc Capacity is %s Available of %s\n", cli.ToDateString(available), cli.ToDateString(total))
	discTitle, _ := md.RequestDiscHeader()
	fmt.Printf("RAW Disc Header: %s\n", discTitle)
	fmt.Println("")
	fmt.Println("Tracks:")
	cnt, err := md.RequestTrackCount()
	if err != nil {
		log.Fatal(err)
	}
	var totalDuration uint64
	for nr := 0; nr < cnt; nr++ {
		title, _ := md.RequestTrackTitle(nr)
		duration, _ := md.RequestTrackLength(nr)
		totalDuration += duration
		enc, _ := md.RequestTrackEncoding(nr)
		senc := "SP"
		switch enc {
		case netmd.EncLP2:
			senc = "LP2"
		case netmd.EncLP4:
			senc = "LP4"
		}
		fmt.Printf("  %d. %s [%s] %s\n", nr+1, title, cli.ToDateString(duration), senc)
	}
	fmt.Println("")
	fmt.Printf("  Total Duration: %s\n", cli.ToDateString(totalDuration))
	fmt.Println("")
}

func erase(md *netmd.NetMD, trk int, safe bool) {
	cnt, err := md.RequestTrackCount()
	if err != nil {
		log.Fatal(err)
	}
	if trk > cnt {
		fmt.Printf("Track %d does not exist because disc only has %d tracks.\n", trk, cnt)
		return
	}
	title, _ := md.RequestTrackTitle(trk - 1)
	fmt.Printf("Do you really want to erase track %d - %s?\n", trk, title)
	if !safe || (safe && cli.AskConfirm()) {
		err := md.EraseTrack(trk - 1)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Track has been erased.")
	} else {
		fmt.Println("Aborted.")
	}
}
