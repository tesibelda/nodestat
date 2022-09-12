# Telegraf exec nodestat input

nodestat is an exec input plugin for [Telegraf](https://github.com/influxdata/telegraf) that gathers status and basic stats from Linux /proc and /sys pseudo filesystems. It is aimed to be small, fast and to gather metrics not (yet) provided by native telegraf's input. Inspired by [node_exporter](https://github.com/prometheus/node_exporter).

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/tesibelda/nodestat/raw/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/tesibelda/nodestat)](https://goreportcard.com/report/github.com/tesibelda/nodestat)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/tesibelda/nodestat?display_name=release)


# Configuration

* Download the [latest release package](https://github.com/tesibelda/nodestat/releases/latest).

* Edit telegraf's exec input configuration with nodestat binary and influx data format.

```
## Gather nodestat stats
[[inputs.exec]]
  commands = ["/path/to/nodestat_binary"]
  data_format = "influx"
```

You can optionally activate only certain collectors, for example:
```
## Gather nodestat stats
[[inputs.exec]]
  commands = ["/path/to/nodestat_binary net pressure"]
  data_format = "influx"
```
Metric timestamp precision will be 1s.

* Restart or reload Telegraf.

# Quick test in your environment

* Run nodestat without parameters
```
/path/to/nodestat
```

# Example output

```plain
nodestat_fc_host,fibrechannel=host12,nodename=20000025ff1bab79,type=NPort\ (fabric\ via\ point-to-point) port_state="Online",port_state_code=0i 1662965695000000000
nodestat_net,interface=eno1,protocol=ethernet link_mode=0i,flag_lower_up=true,flag_running=true,carrier=1i,dormant=0i,duplex="full",operstate="up",flag_up=true 1662965695000000000
nodestat_net,interface=enp3s0f0,protocol=ethernet operstate="down",flag_up=true,carrier=0i,dormant=0i,duplex="unknown",link_mode=0i,flag_lower_up=true,flag_running=true 1662965695000000000
nodestat_pressure io_waiting_avg60=16.67,io_stalled_avg60=16.49,cpu_waiting_avg60=0 1662965695000000000
```

# Metrics
See [Metrics](https://github.com/tesibelda/nodestat/blob/master/METRICS.md)

# Build Instructions

Download the repo

    $ git clone git@github.com:tesibelda/nodestat.git

build the "nodestat" binary

    $ go build -o bin/nodestat cmd/main.go

 If you use [go-task](https://github.com/go-task/task) execute one of these
 
    $ task linux:build

# Author

Tesifonte Belda (https://github.com/tesibelda)

# License

[The MIT License (MIT)](https://github.com/tesibelda/nodestat/blob/master/LICENSE)
