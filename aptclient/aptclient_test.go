package aptclient

import (
	"reflect"
	"testing"
)

func TestParseUpgradeOut(t *testing.T) {
	var parseUpgradeOutTests = []struct {
		name        string
		upgradeOut  string
		upgradeInfo *UpgradeInfo
	}{
		{
			"upgraded 4 packages",
			`Reading package lists...  Done
			Building dependency tree...
			Reading state information...  Done
			Calculating upgrade...  Done
			The following packages will be upgraded:
				google-cloud-sdk (332.0.0-0 => 333.0.0-0)
				isc-dhcp-client (4.4.1-2.1ubuntu5 => 4.4.1-2.1ubuntu5.20.04.1)
				isc-dhcp-common1.3-_3 332 (4.4.1-2.1ubuntu5(somne)) => 4.4.1-2.1ubuntu5.20.04.1)
				linux-libc-dev(linux) (5.4.0-67.75 => 5.4.0-70.78)
			4 upgraded, 1 newly installed, 2 to remove and 5 not upgraded.
			Need to get 70.9 MB of archives.
			After this operation, 8911 kB disk space will be freed.
			Do you want to continue? [Y/n]
		`,
			&UpgradeInfo{

				Upgraded:     4,
				NewInstalled: 1,
				ToRemove:     2,
				NotUpgraded:  5,
				UpgradedPackages: []UpgradedPackage{
					{
						Name:        "google-cloud-sdk",
						FromVersion: "332.0.0-0",
						ToVersion:   "333.0.0-0",
					},
					{
						Name:        "isc-dhcp-client",
						FromVersion: "4.4.1-2.1ubuntu5",
						ToVersion:   "4.4.1-2.1ubuntu5.20.04.1",
					},
					{
						Name:        "isc-dhcp-common1.3-_3 332",
						FromVersion: "4.4.1-2.1ubuntu5(somne))",
						ToVersion:   "4.4.1-2.1ubuntu5.20.04.1",
					},
					{
						Name:        "linux-libc-dev(linux)",
						FromVersion: "5.4.0-67.75",
						ToVersion:   "5.4.0-70.78",
					},
				},
			},
		},
		{
			"upgraded 0 packages",
			`Reading package lists... Done
			Building dependency tree
			Reading state information... Done
			Calculating upgrade... Done
			0 upgraded, 0 newly installed, 0 to remove and 0 not upgraded.
		`,
			&UpgradeInfo{
				Upgraded:         0,
				NewInstalled:     0,
				ToRemove:         0,
				NotUpgraded:      0,
				UpgradedPackages: []UpgradedPackage{},
			},
		},
		{
			"invalid upgraded count number",
			`Reading package lists... Done
			Building dependency tree
			Reading state information... Done
			Calculating upgrade... Done
			notnum upgraded, 0 newly installed, 0 to remove and 0 not upgraded.
		`,
			nil,
		},
	}

	for _, tt := range parseUpgradeOutTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s, err := parseUpgradeOut(tt.upgradeOut)
			if err != nil && tt.upgradeInfo != nil {
				t.Errorf("got error %+v, but wanted %+v", err, tt.upgradeInfo)
			}
			if !reflect.DeepEqual(tt.upgradeInfo, s) {
				t.Errorf("got %+v, want %+v", s, tt.upgradeInfo)
			}
		})
	}
}
