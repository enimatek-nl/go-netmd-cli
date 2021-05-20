package main

import (
	"fmt"
	"github.com/enimatek-nl/go-netmd-lib"
	"log"
	"os"
)

const (
	version = "0.0.1b"
)

func main() {

	// TODO: parse...
	cmd := "help"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	md, err := netmd.NewNetMD(0, false)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer md.Close()

	switch cmd {
	case "send":
		fmt.Println("stub")
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
	fmt.Println("	github.com/enimatek-nl")
	fmt.Println("Version:")
	fmt.Printf("	%s\n", version)
	fmt.Println("Usage:")
	fmt.Println("	netmd-cli [options] command [arguments...]")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("	list	List all track data on the MiniDisc")
	fmt.Println("	send")
	fmt.Println("Options:")
	fmt.Println("	-v	verbose output")
	fmt.Println("")
}

func list(md *netmd.NetMD) {
	fmt.Println("")
	_, total, available, _ := md.RequestDiscCapacity()
	fmt.Printf("Disc Capacity is %s Available of %s\n", toDateString(available), toDateString(total))
	discTitle, _ := md.RequestDiscHeader()
	fmt.Printf("RAW Disc Header: %s\n", discTitle)
	fmt.Println("")
	fmt.Println("Tracks:")
	cnt, err := md.RequestTrackCount()
	if err != nil {
		log.Fatal(err)
	}
	for nr := 0; nr < cnt; nr++ {
		title, _ := md.RequestTrackTitle(nr)
		duration, _ := md.RequestTrackLength(nr)
		enc, _ := md.RequestTrackEncoding(nr)
		senc := "SP"
		switch enc {
		case netmd.EncLP2:
			senc = "LP2"
		case netmd.EncLP4:
			senc = "LP4"
		}
		fmt.Printf("  %d. %s [%s] %s\n", nr+1, title, toDateString(duration), senc)
	}
	fmt.Println("")
}

func toDateString(s uint64) string {
	hours := s / 3600
	minutes := (s - (3600 * hours)) / 60
	seconds := s - (3600 * hours) - (minutes * 60)
	if hours != 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	} else {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
}
