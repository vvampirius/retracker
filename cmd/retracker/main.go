package main

import (
	Core "github.com/vvampirius/retracker/core"
	"github.com/vvampirius/retracker/core/common"
	"flag"
	"fmt"
	"syscall"
	"os"
)

const VERSION  = 0.1

func PrintRepo(){
	fmt.Fprintln(os.Stderr, "\n# https://github.com/vvampirius/retracker")
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		PrintRepo()
	}
	listen := flag.String("l", ":80", "Listen address:port")
	age := flag.Float64("a", 180, "Keep 'n' minutes peer in memory")
	debug := flag.Bool("d", false, "Debug mode")
	ver := flag.Bool("v", false, "Show version")
	flag.Parse()

	if *ver {
		fmt.Println(VERSION)
		PrintRepo()
		syscall.Exit(0)
	}

	config := common.Config{
		Listen: *listen,
		Debug: *debug,
		Age: *age,
	}

	Core.New(&config)
}