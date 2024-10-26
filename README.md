# nodestat input plugin

nodestat is an exec input plugin for [telegraf](https://github.com/influxdata/telegraf) that gathers status and basic stats from Linux /proc and /sys pseudo filesystems. It is aimed to be small, fast and to gather metrics not (yet) provided by native telegraf's input. Inspired by [node_exporter](https://github.com/prometheus/node_exporter).

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

Command line options are:
```
nodestat [--version] [--help] [collector]...
```
Current collectors are:\
 fc_host - fibrechannel metrics from /sys/class/fc_host/\
 net - network interface metrics from /sys/class/net/\
 pressure - metrics from /proc/pressure/\
 userprocs - metrics for the number of processes and threads per user\
Metric timestamp precision will be 1s.

* Restart or reload Telegraf.

# Quick test in your environment

* Run nodestat without parameters
```
/path/to/nodestat
```

# Example output

```plain
nodestat_fc_host,nodename=20000025ff1bab79,type=NPort\ (fabric\ via\ point-to-point),fibrechannel=host12 port_state="Online",port_state_code=0u,loss_of_signal_count=0u,loss_of_sync_count=0u,nos_count=1u,link_failure_count=1u,seconds_since_last_reset=14476044u,error_frames=0u,rx_frames=181397571u,tx_frames=749365874u,fcp_packet_aborts=0u 1662965695000000000
nodestat_net,interface=eno1,protocol=ethernet carrier=1i,flag_running=true,flag_up=true,operstate_code=0i,flag_lower_up=true,dormant=0i,duplex="full",link_mode=0i,operstate="up"  1662965695000000000
nodestat_net,interface=enp3s0f0,protocol=ethernet flag_running=true,flag_up=true,carrier=0i,duplex="unknown",link_mode=0i,operstate="down",operstate_code=5i,flag_lower_up=true,dormant=0i 1662965695000000000
nodestat_pressure cpu_waiting_avg60=0,io_waiting_avg60=21.1,io_stalled_avg60=20.97 1662965695000000000
nodestat_userprocs,group=root,user=root threads=173i,processes=124i 1662965695000000000
nodestat_userprocs,group=postfix,user=postfix processes=3i,threads=3i 1662965695000000000
```

# Metrics
See [Metrics](https://github.com/tesibelda/nodestat/blob/master/METRICS.md)

# Build Instructions

Download the repo

    $ git clone https://github.com/tesibelda/nodestat.git

build the "nodestat" binary. If under Windows set GOOS environment variable to linux

    $ go build -o bin/nodestat cmd/main.go

 If you use [go-task](https://github.com/go-task/task) execute one of these
 
    $ task build

# Author

Tesifonte Belda (https://github.com/tesibelda)

# License

[The MIT License (MIT)](https://github.com/tesibelda/nodestat/blob/master/LICENSE)
