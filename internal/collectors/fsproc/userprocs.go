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
	"strconv"
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

	var (
		uidprocs       = make(map[uint64]userInfo, 10)
		info           userInfo
		status         procfs.ProcStatus
		stat           procfs.ProcStat
		totalProcs, th int
		ok             bool
	)
	for _, pid := range p {
		if status, err = pid.NewStatus(); err != nil {
			// PIDs can vanish between getting the list and getting stats.
			continue
		}

		th = 0
		if stat, err = pid.Stat(); err == nil {
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

	var (
		fields = make(map[string]interface{}, 2)
		tags   = make(map[string]string, 2)
		t      time.Time
		m      metric.Metric
		usr    *user.User
		grp    *user.Group
	)
	for k, v := range uidprocs {
		if usr, err = user.LookupId(strconv.FormatUint(k, 10)); err != nil {
			continue
		}
		if len(usr.Username) > 0 {
			if grp, err = user.LookupGroupId(usr.Gid); err != nil {
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
