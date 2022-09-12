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

	"github.com/prometheus/procfs"
	"github.com/tesibelda/nodestat/pkg/simplemetric"
)

// GatherProcPressureInfo prints pressure metrics from /proc/pressure/
func GatherProcPressureInfo() error {
	var (
		psiResources = []string{"cpu", "io", "memory"}
		err          error
	)

	fs, err := procfs.NewDefaultFS()
	if err != nil {
		return fmt.Errorf("failed to open procfs: %w", err)
	}

	fields := make(map[string]interface{}, 5)
	m := simplemetric.New("nodestat_pressure", nil, fields)

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
			fields["cpu_waiting_avg60"] = float64(stats.Some.Avg60)
		case "io":
			fields["io_waiting_avg60"] = float64(stats.Some.Avg60)
			fields["io_stalled_avg60"] = float64(stats.Full.Avg60)
		case "mem":
			fields["memory_waiting_avg60"] = float64(stats.Some.Avg60)
			fields["memory_stalled_avg60"] = float64(stats.Full.Avg60)
		}
	}
	fmt.Fprintln(os.Stdout, m.String("influx"))

	return nil
}
