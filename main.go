package main

import (
	"github.com/enimatek-nl/go-netmd-lib"
	"log"
)

func main() {
	md, err := netmd.NewNetMD(0, false)
	if err != nil {
		log.Fatal(err)
	}
	defer md.Close()
}
