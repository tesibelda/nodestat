// fsproc functions show metrics from linux /sys filesystem using influx line protocol
//
// Author: Tesifonte Belda
// License: The MIT License (MIT)
//
// References:
//  https://pkg.go.dev/github.com/prometheus/procfs@v0.15.1/sysfs#FibreChannelClass

package fssys

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/prometheus/procfs/sysfs"
	"github.com/tesibelda/lightmetric/metric"
)

// fc state code mapping
var fcState = map[string]int{
	"online":   0,
	"Online":   0,
	"unknown":  1,
	"Unknown":  1,
	"blocked":  2,
	"Blocked":  2,
	"linkdown": 3,
	"Linkdown": 3,
}

// GatherSysFcHostInfo prints fibrechannel metrics from /sys/class/fc_host/
func GatherSysFcHostInfo() error {
	var (
		m     metric.Metric
		state int
		ok    bool
	)

	fs, err := sysfs.NewDefaultFS()
	if err != nil {
		return fmt.Errorf("failed to open sysfs: %w", err)
	}

	fcDevices, err := fs.FibreChannelClass()
	if err != nil {
		var pErr *os.PathError
		if errors.Is(err, os.ErrNotExist) || errors.As(err, &pErr) {
			return nil
		}
		return fmt.Errorf("failed to retrieve fibrechannel stats: %w", err)
	}
	t := metric.TimeWithPrecision(time.Now(), time.Second)

	tags := make(map[string]string, 3)
	fields := make(map[string]interface{}, 10)

	for _, fcInfo := range fcDevices {
		if state, ok = fcState[*fcInfo.PortState]; !ok {
			state = 1
		}

		tags["fibrechannel"] = *fcInfo.Name
		tags["nodename"] = *fcInfo.NodeName
		tags["type"] = *fcInfo.PortType
		fields["port_state"] = *fcInfo.PortState
		fields["port_state_code"] = state
		fields["link_failure_count"] = *fcInfo.Counters.LinkFailureCount
		fields["seconds_since_last_reset"] = *fcInfo.Counters.SecondsSinceLastReset
		fields["loss_of_signal_count"] = *fcInfo.Counters.LossOfSignalCount
		fields["loss_of_sync_count"] = *fcInfo.Counters.LossOfSyncCount
		fields["nos_count"] = *fcInfo.Counters.NosCount
		fields["error_frames"] = *fcInfo.Counters.ErrorFrames
		fields["rx_frames"] = *fcInfo.Counters.RXFrames
		fields["tx_frames"] = *fcInfo.Counters.TXFrames

		m = metric.New("nodestat_fc_host", tags, fields, t)
		fmt.Fprint(os.Stdout, m.String(metric.InfluxLp))
	}

	return nil
}
