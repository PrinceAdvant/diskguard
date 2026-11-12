//go:build linux || darwin || freebsd

package diskguard

import "testing"

func TestLow(t *testing.T) {
	u := Usage{TotalBytes: 100, FreeBytes: 5, TotalInodes: 100, FreeInodes: 50}
	if !u.Low(10) {
		t.Fatal("5 percent free should be low at a 10 percent threshold")
	}
	u.FreeBytes = 50
	if u.Low(10) {
		t.Fatal("50 percent free should not be low")
	}
	u.FreeInodes = 2
	if !u.Low(10) {
		t.Fatal("low inodes should count as low")
	}
}

func TestStatRoot(t *testing.T) {
	u, err := Stat("/")
	if err != nil {
		t.Fatal(err)
	}
	if u.TotalBytes == 0 {
		t.Fatal("expected a non-zero size for /")
	}
}
