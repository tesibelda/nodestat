// collectors contains a list of collector's information
//  New collectors should be added to collectInfos in init
//
// Author: Tesifonte Belda
// License: The MIT License (MIT)
//

package collectors

import (
	"fmt"

	"github.com/tesibelda/nodestat/internal/collectors/fsproc"
	"github.com/tesibelda/nodestat/internal/collectors/fssys"
)

type gatherf func() error
type CollectorInfo struct {
	Name      string
	IsDefault bool
	What      string
	gfunc     gatherf
}

var collectInfos []CollectorInfo

func init() {
	var c = make([]CollectorInfo, 0, 4)
	var ci CollectorInfo

	ci = CollectorInfo{"fc_host", true, "fibrechannels", fssys.GatherSysFcHostInfo}
	c = append(c, ci)
	ci = CollectorInfo{"net", true, "network interfaces", fssys.GatherSysNetInfo}
	c = append(c, ci)
	ci = CollectorInfo{"pressure", true, "pressure", fsproc.GatherProcPressureInfo}
	c = append(c, ci)
	ci = CollectorInfo{"userprocs", true, "processes per user", fsproc.GatherProcUserProcsInfo}
	c = append(c, ci)
	collectInfos = c
}

func GetInfo() []CollectorInfo {
	return collectInfos
}

func Gather(colname string) error {
	var f gatherf

	for _, v := range collectInfos {
		if v.Name == colname {
			f = v.gfunc
		}
	}

	if f != nil {
		return f()
	}

	return fmt.Errorf("collector not available")
}

func CollectorAvailable(in string) bool {
	for _, col := range collectInfos {
		if in == col.Name {
			return true
		}
	}
	return false
}
