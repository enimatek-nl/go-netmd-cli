package main

import (
	"fmt"
	"github.com/enimatek-nl/go-netmd-lib"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	version = "0.0.2b"
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
	enc := netmd.DfStereoSP

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
				i, err := ToInt(os.Args[ptr])
				if err != nil {
					log.Fatal(err)
				}
				index = i
			}
			if strings.HasSuffix(a, "d") {
				ptr++
				if len(os.Args) <= ptr {
					log.Fatal("missing encoder (sp, lp2, lp4)")
				}
				switch os.Args[ptr] {
				case "sp":
					enc = netmd.DfStereoSP
				case "lp4":
					enc = netmd.DfLP4
				case "lp2":
					enc = netmd.DfLP2
				}
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
	case "title":
		ptr++
		t := ""
		if len(os.Args) > ptr {
			t = strings.Join(os.Args[ptr:], " ")
		}
		title(md, t, safe)
	case "move":
		ptr++
		if len(os.Args) <= ptr {
			log.Fatal("missing from track #")
		}
		from, err := ToInt(os.Args[ptr])
		if err != nil {
			log.Fatal(err)
		}
		ptr++
		if len(os.Args) <= ptr {
			log.Fatal("missing to track #")
		}
		to, err := ToInt(os.Args[ptr])
		if err != nil {
			log.Fatal(err)
		}
		move(md, from, to, safe)
	case "rename":
		ptr++
		if len(os.Args) <= ptr {
			log.Fatal("missing from track #")
		}
		trk, err := ToInt(os.Args[ptr])
		if err != nil {
			log.Fatal(err)
		}
		ptr++
		t := ""
		if len(os.Args) > ptr {
			t = strings.Join(os.Args[ptr:], " ")
		}
		rename(md, trk, t, safe)
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
		send(md, enc, fn, t)
	case "erase":
		ptr++
		if len(os.Args) <= ptr {
			log.Fatal("missing track #")
		}
		trk, err := ToInt(os.Args[ptr])
		if err != nil {
			log.Fatal(err)
		}
		erase(md, trk, safe)
	case "degroup":
		ptr++
		if len(os.Args) <= ptr {
			log.Fatal("missing track #")
		}
		trk, err := ToInt(os.Args[ptr])
		if err != nil {
			log.Fatal(err)
		}
		degroup(md, trk, safe)
	case "group":
		ptr++
		if len(os.Args) <= ptr {
			log.Fatal("missing to track #")
		}
		trk, err := ToInt(os.Args[ptr])
		if err != nil {
			log.Fatal(err)
		}
		ptr++
		t := ""
		if len(os.Args) > ptr {
			t = strings.Join(os.Args[ptr:], " ")
		}
		group(md, trk, t, safe)
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
	fmt.Println("  list                     List all track data on the disc.")
	fmt.Println("  send [wav] [title]       Send stereo pcm data to the disc.")
	fmt.Println("  title [title]            Rename the disc title.")
	fmt.Println("  rename [number] [title]  Rename the track number.")
	fmt.Println("  move [number] [to]       Move the track number around.")
	fmt.Println("  erase [number]           Erase track number from disc.")
	fmt.Println("  group [to] [title]       Put tracks ungrouped up to track number")
	fmt.Println("                           into a group by the name title.")
	fmt.Println("  degroup [number]         De-group including all tracks of the")
	fmt.Println("                           group that are part of the track number.")
	fmt.Println("Options:")
	fmt.Println("  -v             Verbose logging output.")
	fmt.Println("  -y             Skip confirm questions.")
	fmt.Println("  -d [encoding]  Encoding on disk lp2, lp4 or sp (default: sp)")
	fmt.Println("                 lp-modes known to work only with Sharp NetMD")
	fmt.Println("  -i [index]     Set the NetMD usb device index when multiple")
	fmt.Println("                 devices are connected. [default: 0]")
	fmt.Println("")
}

func send(md *netmd.NetMD, enc netmd.DiscFormat, fn, t string) {
	track, err := md.NewTrack(t, fn)
	track.DiscFormat = enc
	if err != nil {
		log.Fatal(err)
	}
	c := make(chan netmd.Transfer)
	go md.Send(track, c)

	fmt.Println("")
	spinner := []string{"|", "/", "-", "\\"}
	spinIndex := -1
	for {
		res, ok := <-c
		if !ok {
			break
		}
		if res.Error != nil {
			log.Fatal(res.Error)
		}
		switch res.Type {
		case netmd.TtSend:
			i := float32(100) / float32(track.TotalBytes())
			p := int(i * float32(res.Transferred))
			barFill := "=>"
			for j := 0; j < (p / 4); j++ {
				barFill = "=" + barFill
			}
			for j := p / 4; j < 25; j++ {
				barFill = barFill + " "
			}
			fmt.Printf("\rTransferring [%s] %d%% (%d KiB / %d KiB)", barFill, p, res.Transferred/1024, track.TotalBytes()/1024)
		case netmd.TtTrack:
			fmt.Println()
			fmt.Printf("Writing title to new track #%d...\n", res.Track)
		case netmd.TtPoll:
			if spinIndex == -1 {
				fmt.Println()
			}
			spinIndex++
			if spinIndex >= len(spinner) {
				spinIndex = 0
			}
			fmt.Printf("\r%s Waiting for NetMD to finish writing...", spinner[spinIndex])
		}
	}

	fmt.Println("Done.")
}

func list(md *netmd.NetMD) {
	fmt.Println("")
	_, total, available, _ := md.RequestDiscCapacity()
	fmt.Printf("Disc has %s Available of %s\n", ToDateString(available), ToDateString(total))

	discTitle, _ := md.RequestDiscHeader()
	fmt.Printf("RAW Disc Header: %s\n", discTitle)

	root := netmd.NewRoot(discTitle)
	fmt.Println("")
	fmt.Printf("💿 %s\n", root.Title)
	cnt, err := md.RequestTrackCount()
	if err != nil {
		log.Fatal(err)
	}
	var totalDuration uint64
	var group *netmd.Group
	for nr := 0; nr < cnt; nr++ {

		_grp := root.SearchGroup(nr)
		if _grp != group {
			group = _grp
			if group != nil {
				fmt.Printf(" 📁 %s/\n", group.Title)
			}
		}

		if group != nil {
			fmt.Printf(" ")
		}
		fmt.Printf(" 🎵 %d.", nr+1)

		title, _ := md.RequestTrackTitle(nr)
		duration, _ := md.RequestTrackLength(nr)
		fmt.Printf(" %s %s  ", ToDateString(duration), title)

		flag, _ := md.RequestTrackFlag(nr)
		switch flag {
		case netmd.TrackProtected:
			fmt.Print("🔒")
		case netmd.TrackUnprotected:
			fmt.Print("🔓")
		}

		enc, _ := md.RequestTrackEncoding(nr)
		switch enc {
		case netmd.EncLP2:
			fmt.Print("LP2")
		case netmd.EncLP4:
			fmt.Print("LP4")
		case netmd.EncSP:
			fmt.Print("SP")
		}

		fmt.Print("\n")
		totalDuration += duration
	}
	fmt.Println("")
	fmt.Printf(" Total Duration: %s\n", ToDateString(totalDuration))
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

	if !safe || (safe && AskConfirm(fmt.Sprintf("Do you really want to erase track %d - %s?", trk, title))) {
		err := md.EraseTrack(trk - 1)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Track has been erased.")
	} else {
		fmt.Println("Aborted.")
	}
}

func move(md *netmd.NetMD, trk, to int, safe bool) {
	if !safe || (safe && AskConfirm(fmt.Sprintf("Do you really want to move track %d to %d?", trk, to))) {
		err := md.MoveTrack(trk-1, to-1)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Track has been moved.")
	} else {
		fmt.Println("Aborted.")
	}
}

func rename(md *netmd.NetMD, trk int, t string, safe bool) {
	if !safe || (safe && AskConfirm(fmt.Sprintf("Do you really want to rename track %d to '%s'?", trk, t))) {
		err := md.SetTrackTitle(trk-1, t, false)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Track has been renamed.")
	} else {
		fmt.Println("Aborted.")
	}
}

func title(md *netmd.NetMD, t string, safe bool) {
	if !safe || (safe && AskConfirm(fmt.Sprintf("Do you really want to rename the disc to '%s'?", t))) {
		d, err := md.RequestDiscHeader()
		if err == nil {
			r := netmd.NewRoot(d)
			r.Title = t
			err := md.SetDiscHeader(r.ToString())
			if err == nil {
				fmt.Println("Disc has been renamed.")
				return
			}
		}
	}
	fmt.Println("Aborted.")
}

func degroup(md *netmd.NetMD, trk int, safe bool) {
	d, _ := md.RequestDiscHeader()
	r := netmd.NewRoot(d)
	grp := r.SearchGroup(trk - 1)
	if grp != nil {
		if !safe || (safe && AskConfirm(fmt.Sprintf("Do you really want to degroup all tracks in group '%s'?", grp.Title))) {
			n := make([]*netmd.Group, 0)
			for _, g := range r.Groups {
				if g != grp {
					n = append(n, g)
				}
			}
			r.Groups = n
			err := md.SetDiscHeader(r.ToString())
			if err == nil {
				fmt.Println("Group has been removed.")
				return
			}
		}
	}
	fmt.Println("Aborted.")
}

func group(md *netmd.NetMD, to int, t string, safe bool) {
	if !safe || (safe && AskConfirm(fmt.Sprintf("Do you really want to group ungrouped tracks up to track %d?", to))) {
		cnt, err := md.RequestTrackCount()
		if to > cnt {
			fmt.Println("not enough tracks")
			return
		}
		d, err := md.RequestDiscHeader()
		if err == nil {
			r := netmd.NewRoot(d)
			s := 0
			grp := r.SearchGroup(s)
			for grp != nil {
				s++
				grp = r.SearchGroup(s)
			}
			if s > to-1 {
				fmt.Printf("no tracks before %d that are ungrouped\n", to)
				return
			}
			r.AddGroup(t, s+1, to)
			err := md.SetDiscHeader(r.ToString())
			if err == nil {
				fmt.Println("Group has been created.")
				return
			}
		}
	}
	fmt.Println("Aborted.")
}

func ToDateString(s uint64) string {
	hours := s / 3600
	minutes := (s - (3600 * hours)) / 60
	seconds := s - (3600 * hours) - (minutes * 60)
	if hours != 0 {
		return fmt.Sprintf("\u001B[33m%02dh %02dm %02ds\u001B[0m", hours, minutes, seconds)
	} else {
		return fmt.Sprintf("\u001B[33m%02dm %02ds\u001B[0m", minutes, seconds)
	}
}

func ToInt(s string) (int, error) {
	i, err := strconv.ParseInt(s, 10, 8)
	if err != nil {
		return -1, err
	}
	return int(i), nil
}

func AskConfirm(q string) bool {
	var response string

	fmt.Printf("%s (y/n):", q)
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}

	switch strings.ToLower(response) {
	case "y", "yes":
		return true
	case "n", "no":
		return false
	default:
		return AskConfirm(fmt.Sprintf("Please type (y)es or (n)o and then press enter"))
	}
}
