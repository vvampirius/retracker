package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

const VERSION = `0.10.0`

var (
	ErrorLog = log.New(os.Stderr, `error#`, log.Lshortfile)
	DebugLog = log.New(os.Stdout, `debug#`, log.Lshortfile)

	//go:embed favicon.ico
	faviconIco []byte
)

func helpText() {
	fmt.Println("# https://github.com/vvampirius/retracker\n")
	flag.PrintDefaults()
}

func main() {
	listen := flag.String("l", ":80", "Listen address:port")
	age := flag.Float64("a", 180, "Keep 'n' minutes peer in memory")
	debug := flag.Bool("d", false, "Debug mode")
	xrealip := flag.Bool("x", false, "Get RemoteAddr from X-Real-IP header")
	forwards := flag.String("f", "", "Load forwards from YAML file")
	forwardTimeout := flag.Int("t", 2, "Timeout (sec) for forward requests (used with -f)")
	enablePrometheus := flag.Bool("p", false, "Enable Prometheus metrics")
	announceResponseInterval := flag.Int("i", 30, "Announce response interval (sec)")
	ver := flag.Bool("v", false, "Show version")
	help := flag.Bool("h", false, "print this help")
	flag.Parse()

	if *help {
		helpText()
		os.Exit(0)
	}

	if *ver {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	fmt.Printf("Starting version %s\n", VERSION)

	config := Config{
		AnnounceResponseInterval: *announceResponseInterval,
		Listen:                   *listen,
		Debug:                    *debug,
		Age:                      *age,
		XRealIP:                  *xrealip,
		ForwardTimeout:           *forwardTimeout,
	}

	if *forwards != `` {
		if err := config.ReloadForwards(*forwards); err != nil {
			ErrorLog.Fatalln(err.Error())
		}
	}

	tempStorage, err := NewTempStorage(``)
	if err != nil {
		os.Exit(1)
	}

	core := NewCore(&config, tempStorage)

	// https://github.com/vvampirius/retracker/issues/7
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/x-icon")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(faviconIco)))
		w.Write(faviconIco)
	})

	http.HandleFunc("/scrape", core.httpScrapeHandler)
	http.HandleFunc("/announce", core.Receiver.Announce.httpHandler)
	if *enablePrometheus {
		p, err := NewPrometheus()
		if err != nil {
			os.Exit(1)
		}
		core.Receiver.Announce.Prometheus = p
	}
	if err := http.ListenAndServe(config.Listen, nil); err != nil { // set listen port
		ErrorLog.Println(err)
	}
}
