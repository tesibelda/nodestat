// nodestate is an external exe telegraf plugin that gather linux node stats
//
// Author: Tesifonte Belda
// License: The MIT License (MIT)

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tesibelda/nodestat/internal/fsproc"
	"github.com/tesibelda/nodestat/internal/fssys"
	"golang.org/x/exp/slices"
)

var showHelp = flag.Bool("help", false, "show help")
var showVersion = flag.Bool("version", false, "show version and exit")
var Version string = ""

func main() {
	var err error

	// parse command line options
	flag.Parse()
	if *showVersion {
		fmt.Println("nodestat", Version)
		os.Exit(0)
	}
	if *showHelp {
		fmt.Println("nodestat [--version] [--help] [collection]...")
		fmt.Println("  Possible collections are: fc_host net pressure. Default: all")
		os.Exit(0)
	}
	cols := flag.Args()

	// Get and display fibrechannels info
	if len(cols) == 0 || slices.Contains(cols, "fc_host") {
		if err = fssys.GatherSysFcHostInfo(); err != nil {
			fmt.Fprintln(os.Stderr, "Could not obtain fibrechannels info:", err)
			os.Exit(1)
		}
	}

	// Get and display network interfaces info
	if len(cols) == 0 || slices.Contains(cols, "net") {
		if err = fssys.GatherSysNetInfo(); err != nil {
			fmt.Fprintln(os.Stderr, "Could not obtain network interfaces info:", err)
			os.Exit(1)
		}
	}

	// Get and display pressure info if available
	if len(cols) == 0 || slices.Contains(cols, "fc_host") {
		if err = fsproc.GatherProcPressureInfo(); err != nil {
			fmt.Fprintln(os.Stderr, "Could not obtain pressure info:", err)
			os.Exit(1)
		}
	}
}
