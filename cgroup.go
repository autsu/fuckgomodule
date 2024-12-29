package holmes

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/sys/unix"
)

type CGroupGetter interface {
	CGroupCPUCore() (float64, error)
	CGroupMemoryLimit() (uint64, error)
}

var (
	_ CGroupGetter = (*CGroupV1)(nil)
	_ CGroupGetter = (*CGroupV2)(nil)
)

type CGroupV1 struct {
}

func (C *CGroupV1) CGroupCPUCore() (float64, error) {
	return getCGroupCPUCore()
}

func (C *CGroupV1) CGroupMemoryLimit() (uint64, error) {
	return getCGroupMemoryLimit()
}

type CGroupV2 struct {
}

func (C *CGroupV2) CGroupCPUCore() (float64, error) {
	v, err := os.ReadFile(cgroupCpuMaxPathV2)
	if err != nil {
		return 0, err
	}
	fields := strings.Fields(string(v))
	if len(fields) != 2 {
		return 0, fmt.Errorf("invalid cgroup v2 format")
	}
	quota := fields[0]
	period := fields[1]
	if quota == "max" {
		return float64(runtime.GOMAXPROCS(-1)), nil
	} else {
		periodInt, err := parseUint(period, 10, 64)
		if err != nil {
			return 0, err
		}
		quotaInt, err := parseUint(quota, 10, 64)
		if err != nil {
			return 0, err
		}
		return float64(quotaInt) / float64(periodInt), nil
	}
}

func (C *CGroupV2) CGroupMemoryLimit() (uint64, error) {
	v, err := os.ReadFile(cgroupMemLimitPathV2)
	if err != nil {
		return 0, err
	}
	mem, err := parseUint(strings.TrimSpace(string(v)), 10, 64)
	if err != nil {
		return 0, err
	}
	return mem, nil
}

func NewCGroup() CGroupGetter {
	switch Mode() {
	case Legacy:
		return &CGroupV1{}
	case Unified:
		return &CGroupV2{}
	default:
		panic("invalid cgroup mode")
	}
}

// Mode returns the cgroups mode running on the host
func Mode() CGMode {
	checkMode.Do(func() {
		var st unix.Statfs_t
		if err := unix.Statfs(unifiedMountpoint, &st); err != nil {
			cgMode = Unavailable
			return
		}
		switch st.Type {
		case unix.CGROUP2_SUPER_MAGIC:
			cgMode = Unified
		default:
			cgMode = Legacy
			if err := unix.Statfs(filepath.Join(unifiedMountpoint, "unified"), &st); err != nil {
				return
			}
			if st.Type == unix.CGROUP2_SUPER_MAGIC {
				cgMode = Hybrid
			}
		}
	})
	return cgMode
}
