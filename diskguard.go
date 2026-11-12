//go:build linux || darwin || freebsd

// Package diskguard reports free space and free inodes for a filesystem path
// and decides whether they are below a threshold.
package diskguard

import "syscall"

// Usage describes the filesystem that holds a path.
type Usage struct {
	Path        string
	TotalBytes  uint64
	FreeBytes   uint64
	TotalInodes uint64
	FreeInodes  uint64
}

// FreePercent returns the percentage of space available to unprivileged users.
func (u Usage) FreePercent() float64 { return pct(u.FreeBytes, u.TotalBytes) }

// FreeInodePercent returns the percentage of free inodes. Filesystems without
// fixed inode tables report 100.
func (u Usage) FreeInodePercent() float64 {
	if u.TotalInodes == 0 {
		return 100
	}
	return pct(u.FreeInodes, u.TotalInodes)
}

func pct(free, total uint64) float64 {
	if total == 0 {
		return 0
	}
	return float64(free) * 100 / float64(total)
}

// Stat reads filesystem usage for path.
func Stat(path string) (Usage, error) {
	var s syscall.Statfs_t
	if err := syscall.Statfs(path, &s); err != nil {
		return Usage{Path: path}, err
	}
	bs := uint64(s.Bsize)
	return Usage{
		Path:        path,
		TotalBytes:  uint64(s.Blocks) * bs,
		FreeBytes:   uint64(s.Bavail) * bs,
		TotalInodes: uint64(s.Files),
		FreeInodes:  uint64(s.Ffree),
	}, nil
}

// Low reports whether free space or free inodes are below minPercent.
func (u Usage) Low(minPercent float64) bool {
	return u.FreePercent() < minPercent || u.FreeInodePercent() < minPercent
}
