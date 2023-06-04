// fsproc functions show metrics from linux /sys filesystem using influx line protocol
//
// Author: Tesifonte Belda
// License: The MIT License (MIT)
//
// References:
//  https://www.kernel.org/doc/html/latest/networking/operstates.html
//  https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/if_arp.h#L30
//  https://www.kernel.org/doc/Documentation/ABI/testing/sysfs-class-net
//  https://stackoverflow.com/questions/11679514/what-is-the-difference-between-iff-up-and-iff-running
//  https://pkg.go.dev/golang.org/x/sys/unix#section-readme

package fssys

import (
	"fmt"
	"os"

	"github.com/prometheus/procfs/sysfs"
	"github.com/tesibelda/nodestat/pkg/simplemetric"
)

// network interface type mapping
var iftype = map[int64]string{
	0:   "net/rom",
	1:   "ethernet",
	2:   "experimentalethernet",
	3:   "x25",
	6:   "ethernet",
	15:  "framerelay",
	19:  "atm",
	32:  "infiniband",
	271: "x25",
	272: "x25",
	768: "tunnel",
	769: "tunnel",
	770: "framerelay",
	772: "loopback",
	774: "fddi",
	784: "fibrechannel",
	785: "fibrechannel",
	786: "fibrechannel",
	787: "fibrechannel",
	801: "wireless",
	802: "wireless",
	803: "wireless",
	804: "wireless",
	805: "wireless",
}

// interface operstate code mapping
//
//	ref: https://www.kernel.org/doc/html/latest/networking/operstates.html
var ifOperState = map[string]int{
	"up":             0,
	"dormant":        1,
	"unknown":        2,
	"testing":        3,
	"lowerlayerdown": 4,
	"down":           5,
	"notpresent":     6,
}

// GatherSysNetInfo prints network interface metrics from /sys/class/net/
func GatherSysNetInfo() error {
	var (
		netClass         sysfs.NetClass
		err              error
		itype            string
		carrier, dormant int64
		linkmode, flags  int64
		state            int
		ok               bool
	)

	netClass, err = getNetClassInfo()
	if err != nil {
		return err
	}

	tags := make(map[string]string, 2)
	fields := make(map[string]interface{}, 10)
	m := simplemetric.New("nodestat_net", tags, fields)

	for _, ifaceInfo := range netClass {
		if ifaceInfo.Type == nil {
			continue
		}

		carrier, dormant, linkmode = -1, -1, -1
		flags = 0
		if ifaceInfo.Carrier != nil {
			carrier = *ifaceInfo.Carrier
		}
		if ifaceInfo.Dormant != nil {
			dormant = *ifaceInfo.Dormant
		}
		if ifaceInfo.LinkMode != nil {
			linkmode = *ifaceInfo.LinkMode
		}
		if ifaceInfo.Flags != nil {
			flags = *ifaceInfo.Flags
		}

		itype, ok = iftype[*ifaceInfo.Type]
		if !ok {
			itype = "other"
		}

		if itype != "loopback" {
			state, ok = ifOperState[ifaceInfo.OperState]
			if !ok {
				state = 1
			}

			tags["interface"] = ifaceInfo.Name
			tags["protocol"] = itype
			fields["carrier"] = carrier
			fields["dormant"] = dormant
			fields["duplex"] = ifaceInfo.Duplex
			fields["ifalias"] = ifaceInfo.IfAlias
			fields["link_mode"] = linkmode
			fields["operstate"] = ifaceInfo.OperState
			fields["operstate_code"] = state
			fields["flag_lower_up"] = ((flags>>16)%2 == 0)
			fields["flag_running"] = ((flags>>6)%2 == 0)
			fields["flag_up"] = !(flags%2 == 0)

			fmt.Fprintln(os.Stdout, m.String("influx"))
		}
	}

	return nil
}

// From https://github.com/prometheus/node_exporter/blob/master/collector/netclass_linux.go
func getNetClassInfo() (sysfs.NetClass, error) {
	fs, err := sysfs.NewDefaultFS()
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}

	netClass := sysfs.NetClass{}
	netDevices, err := fs.NetClassDevices()
	if err != nil {
		return netClass, err
	}

	for _, device := range netDevices {
		interfaceClass, err := fs.NetClassByIface(device)
		if err != nil {
			return netClass, err
		}
		netClass[device] = *interfaceClass
	}

	return netClass, nil
}
