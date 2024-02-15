// fsproc functions show metrics from linux /proc filesystem using influx line protocol
//
// Author: Tesifonte Belda
// License: The MIT License (MIT)
//
// References:
//  https://github.com/prometheus/node_exporter/tree/master/collector/pressure_linux.go

package fsproc

import (
	"fmt"
	"os"
	"os/user"
	"time"

	"github.com/prometheus/procfs"
	"github.com/tesibelda/lightmetric/metric"
)

type userInfo struct {
	processes int
	threads   int
}

// GatherProcUserProcsInfo prints number of process per user metrics from /proc/<PID>/status>
func GatherProcUserProcsInfo() error {
	fs, err := procfs.NewDefaultFS()
	if err != nil {
		return fmt.Errorf("failed to open procfs: %w", err)
	}

	p, err := fs.AllProcs()
	if err != nil {
		return fmt.Errorf("unable to list all processes: %w", err)
	}

	uidprocs := make(map[string]userInfo, 10)
	totalProcs, th, ok := 0, 0, false
	info := userInfo{}
	for _, pid := range p {
		status, err := pid.NewStatus()
		if err != nil {
			// PIDs can vanish between getting the list and getting stats.
			continue
		}

		th = 0
		stat, err := pid.Stat()
		if err == nil {
			th = stat.NumThreads
		}

		if info, ok = uidprocs[status.UIDs[0]]; !ok {
			info = userInfo{}
		}
		info.processes++
		info.threads += th
		uidprocs[status.UIDs[0]] = info
		totalProcs++
	}
	if totalProcs == 0 {
		return fmt.Errorf("unable to list any processes")
	}

	fields := make(map[string]interface{}, 2)
	tags := make(map[string]string, 2)
	var t time.Time
	var m metric.Metric
	for k, v := range uidprocs {
		usr, err := user.LookupId(k)
		if err != nil {
			continue
		}
		if len(usr.Username) > 0 {
			grp, err := user.LookupGroupId(usr.Gid)
			if err != nil {
				grp = &user.Group{}
			}

			fields["processes"] = v.processes
			fields["threads"] = v.threads
			tags["user"] = usr.Username
			tags["group"] = grp.Name
			t = metric.TimeWithPrecision(time.Now(), time.Second)
			m = metric.New("nodestat_userprocs", tags, fields, t)
			fmt.Fprint(os.Stdout, m.String(metric.InfluxLp))
		}
	}
	return nil
}
