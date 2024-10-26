// nodestate is an external exe telegraf plugin that gather linux node stats
//
// Author: Tesifonte Belda
// License: The MIT License (MIT)

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tesibelda/nodestat/internal/collectors"
	"golang.org/x/exp/slices"
)

var Version string

func main() {
	var (
		showHelp    = flag.Bool("help", false, "show help")
		showVersion = flag.Bool("version", false, "show version and exit")
		exclude     = flag.Bool("exclude", false, "exclude given collectors")
		err         error
	)

	// parse command line options
	flag.Parse()
	if *showVersion {
		fmt.Println("nodestat", Version)
		os.Exit(0)
	}

	// get available collectors
	colsav := collectors.GetInfo()

	if *showHelp {
		help(colsav)
		os.Exit(0)
	}
	cols := flag.Args()

	// parse given collector list
	for _, in := range cols {
		if !collectors.CollectorAvailable(in) {
			fmt.Fprintf(os.Stderr, "Collector %s not available\n", in)
			help(colsav)
			os.Exit(1)
		}
	}

	// run enabled collectors
	for _, col := range colsav {
		if len(cols) > 0 {
			if *exclude && slices.Contains(cols, col.Name) {
				continue
			}
			if !*exclude && !slices.Contains(cols, col.Name) {
				continue
			}
			if !*exclude && !col.IsDefault {
				continue
			}
		}
		if err = collectors.Gather(col.Name); err != nil {
			fmt.Fprintf(os.Stderr, "Could not get %s info: %s\n", col.What, err)
			os.Exit(1)
		}
	}
}

func help(colsav []collectors.CollectorInfo) {
	fmt.Println("nodestat [--version] [--help] [--exclude] [collector]...")
	fmt.Println("  Available collectors are:")
	for _, col := range colsav {
		fmt.Printf(
			"   %s: collects %s info. Enabled by default: %t\n",
			col.Name,
			col.What,
			col.IsDefault,
		)
	}
}
