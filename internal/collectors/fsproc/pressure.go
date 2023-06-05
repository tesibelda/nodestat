// fsproc functions show metrics from linux /proc filesystem using influx line protocol
//
// Author: Tesifonte Belda
// License: The MIT License (MIT)
//
// References:
//  https://github.com/prometheus/node_exporter/tree/master/collector/pressure_linux.go

package fsproc

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/prometheus/procfs"
	"github.com/tesibelda/lightmetric/metric"
)

// GatherProcPressureInfo prints pressure metrics from /proc/pressure/
func GatherProcPressureInfo() error {
	var psiResources = []string{"cpu", "io", "memory"}

	fs, err := procfs.NewDefaultFS()
	if err != nil {
		return fmt.Errorf("failed to open procfs: %w", err)
	}

	fields := make(map[string]interface{}, 5)

	for _, res := range psiResources {
		stats, err := fs.PSIStatsForResource(res)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) || errors.Is(err, syscall.ENOTSUP) {
				return nil
			}
			return fmt.Errorf("failed to retrieve pressure stats: %w", err)
		}

		switch res {
		case "cpu":
			if stats.Some != nil {
				fields["cpu_waiting_avg60"] = stats.Some.Avg60
			}
		case "io":
			if stats.Some != nil {
				fields["io_waiting_avg60"] = stats.Some.Avg60
			}
			if stats.Full != nil {
				fields["io_stalled_avg60"] = stats.Full.Avg60
			}
		case "mem":
			if stats.Some != nil {
				fields["memory_waiting_avg60"] = stats.Some.Avg60
			}
			if stats.Full != nil {
				fields["memory_stalled_avg60"] = stats.Full.Avg60
			}
		}
	}
	t := metric.TimeWithPrecision(time.Now(), time.Second)
	m := metric.New("nodestat_pressure", nil, fields, t)
	fmt.Fprint(os.Stdout, m.String(metric.InfluxLp))

	return nil
}
