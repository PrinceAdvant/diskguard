// Command diskguard checks free disk space and inodes for one or more paths
// and exits non-zero when any is below the threshold.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/PrinceAdvant/diskguard"
)

func main() {
	min := flag.Float64("min", 15, "minimum free percent for space and inodes")
	flag.Parse()
	paths := flag.Args()
	if len(paths) == 0 {
		paths = []string{"/", "/var/lib/docker"}
	}
	bad := false
	for _, p := range paths {
		u, err := diskguard.Stat(p)
		if err != nil {
			fmt.Printf("%s error: %v\n", p, err)
			bad = true
			continue
		}
		status := "ok"
		if u.Low(*min) {
			status = "LOW"
			bad = true
		}
		fmt.Printf("%s %s space=%.1f%% free (%.1f GB) inodes=%.1f%% free\n", p, status, u.FreePercent(), float64(u.FreeBytes)/1e9, u.FreeInodePercent())
	}
	if bad {
		os.Exit(1)
	}
}
